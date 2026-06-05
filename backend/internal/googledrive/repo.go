// Package googledrive guarda blobs cifrados de ficheiros «Drive» (GOOGLE-002).
// O servidor nunca decifra — em produção o mesmo blob pode ser replicado na API Google.
package googledrive

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("googledrive: ficheiro não encontrado")

// File é uma entrada opaca na BD.
type File struct {
	ID         string
	OwnerID    string
	Blob       []byte
	ExternalID *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Repo acede a google_drive_files.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, ownerID string, blob []byte) (*File, error) {
	if len(blob) == 0 {
		return nil, errors.New("googledrive: blob vazio")
	}
	f := &File{OwnerID: ownerID, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO google_drive_files (owner_id, blob)
		VALUES ($1, $2)
		RETURNING id::text, created_at, updated_at`,
		ownerID, blob,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repo) List(ctx context.Context, ownerID string) ([]File, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, blob, external_id, created_at, updated_at
		FROM google_drive_files WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []File
	for rows.Next() {
		var f File
		if err := rows.Scan(&f.ID, &f.OwnerID, &f.Blob, &f.ExternalID, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func (r *Repo) Get(ctx context.Context, ownerID, id string) (*File, error) {
	f := &File{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, owner_id::text, blob, external_id, created_at, updated_at
		FROM google_drive_files WHERE id = $1 AND owner_id = $2`,
		id, ownerID,
	).Scan(&f.ID, &f.OwnerID, &f.Blob, &f.ExternalID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return f, nil
}

func (r *Repo) Delete(ctx context.Context, ownerID, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM google_drive_files WHERE id = $1 AND owner_id = $2`, id, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
