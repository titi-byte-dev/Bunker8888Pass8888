// Package fin implementa a monitorizacao de custos SaaS (FIN-001): subscricoes
// guardadas como blobs cifrados com a Master Key. O servidor trata tudo como
// bytes opacos; os agregados de custo e os alertas de licencas sem uso (FIN-002)
// sao calculados no cliente, depois de decifrar.
package fin

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound: subscricao inexistente ou de outro utilizador.
var ErrNotFound = errors.New("fin: subscricao nao encontrada")

// Subscription espelha uma linha de "saas_subscriptions".
type Subscription struct {
	ID        string
	OwnerID   string
	Blob      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repo acede a tabela saas_subscriptions.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositorio.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Create guarda uma subscricao cifrada.
func (r *Repo) Create(ctx context.Context, ownerID string, blob []byte) (*Subscription, error) {
	if len(blob) == 0 {
		return nil, errors.New("fin: blob vazio")
	}
	s := &Subscription{OwnerID: ownerID, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO saas_subscriptions (owner_id, blob)
		VALUES ($1, $2)
		RETURNING id::text, created_at, updated_at`,
		ownerID, blob,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// List devolve as subscricoes do utilizador, mais recentes primeiro.
func (r *Repo) List(ctx context.Context, ownerID string) ([]Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, blob, created_at, updated_at
		FROM saas_subscriptions WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Subscription
	for rows.Next() {
		var s Subscription
		if err := rows.Scan(&s.ID, &s.OwnerID, &s.Blob, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

// Update substitui o blob de uma subscricao do utilizador.
func (r *Repo) Update(ctx context.Context, ownerID, id string, blob []byte) (*Subscription, error) {
	if len(blob) == 0 {
		return nil, errors.New("fin: blob vazio")
	}
	s := &Subscription{ID: id, OwnerID: ownerID, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		UPDATE saas_subscriptions SET blob = $1, updated_at = now()
		WHERE id = $2 AND owner_id = $3
		RETURNING created_at, updated_at`,
		blob, id, ownerID,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

// Delete remove uma subscricao do utilizador.
func (r *Repo) Delete(ctx context.Context, ownerID, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM saas_subscriptions WHERE id = $1 AND owner_id = $2`, id, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
