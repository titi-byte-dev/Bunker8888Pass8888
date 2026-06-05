// Package crm implementa leads do funil de vendas (CRM-001): blobs cifrados
// com a Master Key — o servidor nunca vê PII em claro.
package crm

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("crm: lead não encontrado")

// Lead espelha uma linha de crm_leads.
type Lead struct {
	ID        string
	OwnerID   string
	Blob      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repo acede a crm_leads.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, ownerID string, blob []byte) (*Lead, error) {
	if len(blob) == 0 {
		return nil, errors.New("crm: blob vazio")
	}
	l := &Lead{OwnerID: ownerID, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO crm_leads (owner_id, blob)
		VALUES ($1, $2)
		RETURNING id::text, created_at, updated_at`,
		ownerID, blob,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *Repo) List(ctx context.Context, ownerID string) ([]Lead, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, blob, created_at, updated_at
		FROM crm_leads WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Lead
	for rows.Next() {
		var l Lead
		if err := rows.Scan(&l.ID, &l.OwnerID, &l.Blob, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (r *Repo) Update(ctx context.Context, ownerID, id string, blob []byte) (*Lead, error) {
	if len(blob) == 0 {
		return nil, errors.New("crm: blob vazio")
	}
	l := &Lead{ID: id, OwnerID: ownerID, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		UPDATE crm_leads SET blob = $1, updated_at = now()
		WHERE id = $2 AND owner_id = $3
		RETURNING created_at, updated_at`,
		blob, id, ownerID,
	).Scan(&l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return l, nil
}

func (r *Repo) Delete(ctx context.Context, ownerID, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM crm_leads WHERE id = $1 AND owner_id = $2`, id, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
