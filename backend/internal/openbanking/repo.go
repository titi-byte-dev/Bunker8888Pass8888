package openbanking

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotConnected — sync pedido sem consentimento activo.
var ErrNotConnected = errors.New("openbanking: ligação inactiva")

// Repo persiste metadados de ligação bancária por utilizador.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) GetOrCreate(ctx context.Context, ownerID, provider string) (*Connection, error) {
	c := &Connection{OwnerID: ownerID, Provider: provider, Status: StatusPending}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO bank_connections (owner_id, provider, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (owner_id, provider) DO UPDATE SET updated_at = now()
		RETURNING id::text, status, consent_expires_at, last_sync_at, created_at, updated_at`,
		ownerID, provider, StatusPending,
	).Scan(&c.ID, &c.Status, &c.ConsentExpiresAt, &c.LastSyncAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *Repo) Get(ctx context.Context, ownerID, provider string) (*Connection, error) {
	c := &Connection{OwnerID: ownerID, Provider: provider}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, status, consent_expires_at, last_sync_at, created_at, updated_at
		FROM bank_connections WHERE owner_id = $1 AND provider = $2`,
		ownerID, provider,
	).Scan(&c.ID, &c.Status, &c.ConsentExpiresAt, &c.LastSyncAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *Repo) MarkConnected(ctx context.Context, ownerID, provider string, expiresAt time.Time) (*Connection, error) {
	c := &Connection{OwnerID: ownerID, Provider: provider}
	err := r.pool.QueryRow(ctx, `
		UPDATE bank_connections
		SET status = $1, consent_expires_at = $2, updated_at = now()
		WHERE owner_id = $3 AND provider = $4
		RETURNING id::text, status, consent_expires_at, last_sync_at, created_at, updated_at`,
		StatusConnected, expiresAt, ownerID, provider,
	).Scan(&c.ID, &c.Status, &c.ConsentExpiresAt, &c.LastSyncAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return r.GetOrCreate(ctx, ownerID, provider)
		}
		return nil, err
	}
	c.Status = StatusConnected
	return c, nil
}

func (r *Repo) TouchSync(ctx context.Context, ownerID, provider string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE bank_connections SET last_sync_at = now(), updated_at = now()
		WHERE owner_id = $1 AND provider = $2 AND status = $3`,
		ownerID, provider, StatusConnected)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotConnected
	}
	return nil
}
