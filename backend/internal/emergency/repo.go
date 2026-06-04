package emergency

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo acede às tabelas emergency_*.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) UpsertConfig(ctx context.Context, ownerID, heirEmail string, waitDays int, blob []byte) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO emergency_configs (owner_user_id, heir_email, wait_days, encrypted_blob)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (owner_user_id) DO UPDATE SET
			heir_email = EXCLUDED.heir_email,
			wait_days = EXCLUDED.wait_days,
			encrypted_blob = COALESCE(EXCLUDED.encrypted_blob, emergency_configs.encrypted_blob),
			updated_at = now()`,
		ownerID, heirEmail, waitDays, blob,
	)
	return err
}

func (r *Repo) GetConfig(ctx context.Context, ownerID string) (*Config, error) {
	var cfg Config
	var blob []byte
	err := r.pool.QueryRow(ctx, `
		SELECT heir_email, wait_days, encrypted_blob, updated_at
		FROM emergency_configs WHERE owner_user_id = $1`, ownerID,
	).Scan(&cfg.HeirEmail, &cfg.WaitDays, &blob, &cfg.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotConfigured
		}
		return nil, err
	}
	cfg.HasBlob = len(blob) > 0
	return &cfg, nil
}

func (r *Repo) GetEncryptedBlob(ctx context.Context, ownerID string) ([]byte, error) {
	var blob []byte
	err := r.pool.QueryRow(ctx,
		`SELECT encrypted_blob FROM emergency_configs WHERE owner_user_id = $1`, ownerID,
	).Scan(&blob)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotConfigured
		}
		return nil, err
	}
	return blob, nil
}

func (r *Repo) DeleteConfig(ctx context.Context, ownerID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM emergency_configs WHERE owner_user_id = $1`, ownerID)
	return err
}

func (r *Repo) CreateRequest(ctx context.Context, ownerID, heirUserID string, unlocksAt time.Time) (*Request, error) {
	req := &Request{}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO emergency_requests (owner_user_id, heir_user_id, unlocks_at)
		VALUES ($1, $2, $3)
		RETURNING id::text, owner_user_id::text, heir_user_id::text, status,
			requested_at, unlocks_at`,
		ownerID, heirUserID, unlocksAt,
	).Scan(&req.ID, &req.OwnerUserID, &req.HeirUserID, &req.Status, &req.RequestedAt, &req.UnlocksAt)
	if err != nil {
		return nil, err
	}
	_ = r.pool.QueryRow(ctx, `SELECT email FROM users WHERE id = $1`, heirUserID).Scan(&req.HeirEmail)
	return req, nil
}

func (r *Repo) GetRequest(ctx context.Context, id string) (*Request, error) {
	req := &Request{}
	var rejected, consumed *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT r.id::text, r.owner_user_id::text, r.heir_user_id::text, u.email,
			r.status, r.requested_at, r.unlocks_at, r.rejected_at, r.consumed_at
		FROM emergency_requests r
		JOIN users u ON u.id = r.heir_user_id
		WHERE r.id = $1`, id,
	).Scan(
		&req.ID, &req.OwnerUserID, &req.HeirUserID, &req.HeirEmail,
		&req.Status, &req.RequestedAt, &req.UnlocksAt, &rejected, &consumed,
	)
	if err != nil {
		return nil, err
	}
	req.RejectedAt = rejected
	req.ConsumedAt = consumed
	return req, nil
}

func (r *Repo) GetActiveRequest(ctx context.Context, ownerID, heirUserID string) (*Request, error) {
	req := &Request{}
	var rejected, consumed *time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT r.id::text, r.owner_user_id::text, r.heir_user_id::text, u.email,
			r.status, r.requested_at, r.unlocks_at, r.rejected_at, r.consumed_at
		FROM emergency_requests r
		JOIN users u ON u.id = r.heir_user_id
		WHERE r.owner_user_id = $1 AND r.heir_user_id = $2
		  AND r.status IN ('waiting', 'ready')
		ORDER BY r.requested_at DESC LIMIT 1`, ownerID, heirUserID,
	).Scan(
		&req.ID, &req.OwnerUserID, &req.HeirUserID, &req.HeirEmail,
		&req.Status, &req.RequestedAt, &req.UnlocksAt, &rejected, &consumed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	req.RejectedAt = rejected
	req.ConsumedAt = consumed
	return req, nil
}

func (r *Repo) ListRequestsByOwner(ctx context.Context, ownerID string) ([]Request, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT r.id::text, r.owner_user_id::text, r.heir_user_id::text, u.email,
			r.status, r.requested_at, r.unlocks_at, r.rejected_at, r.consumed_at
		FROM emergency_requests r
		JOIN users u ON u.id = r.heir_user_id
		WHERE r.owner_user_id = $1
		ORDER BY r.requested_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Request
	for rows.Next() {
		req := Request{}
		var rejected, consumed *time.Time
		if err := rows.Scan(
			&req.ID, &req.OwnerUserID, &req.HeirUserID, &req.HeirEmail,
			&req.Status, &req.RequestedAt, &req.UnlocksAt, &rejected, &consumed,
		); err != nil {
			return nil, err
		}
		req.RejectedAt = rejected
		req.ConsumedAt = consumed
		out = append(out, req)
	}
	return out, rows.Err()
}

func (r *Repo) SetRejected(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE emergency_requests SET status = 'rejected', rejected_at = now()
		WHERE id = $1 AND status = 'waiting'`, id)
	return err
}

func (r *Repo) SetReady(ctx context.Context, id string, readyAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE emergency_requests SET status = 'ready', unlocks_at = $2
		WHERE id = $1 AND status IN ('waiting', 'ready')`, id, readyAt)
	return err
}

func (r *Repo) SetConsumed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE emergency_requests SET status = 'consumed', consumed_at = now()
		WHERE id = $1 AND status = 'ready'`, id)
	return err
}