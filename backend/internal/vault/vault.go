// Package vault guarda e lê itens cifrados do cofre.
//
// ⚠️ O servidor trata o campo Blob como bytes opacos: foi cifrado no cliente e
// nunca é decifrado aqui (modelo Zero-Knowledge).
package vault

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica item inexistente ou de outro utilizador (mesma resposta).
var ErrNotFound = errors.New("vault: item não encontrado")

// Item espelha uma linha da tabela "vault_items".
type Item struct {
	ID        string
	UserID    string
	Type      string
	Blob      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repo dá acesso à tabela "vault_items".
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria um repositório de itens do cofre.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Create insere um item cifrado pertencente a um utilizador.
func (r *Repo) Create(ctx context.Context, userID, itemType string, blob []byte) (*Item, error) {
	if err := ValidateType(itemType); err != nil {
		return nil, err
	}
	if len(blob) == 0 {
		return nil, fmt.Errorf("vault: blob vazio")
	}
	it := &Item{UserID: userID, Type: itemType, Blob: blob}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO vault_items (user_id, item_type, blob)
		VALUES ($1, $2, $3)
		RETURNING id::text, created_at, updated_at`,
		userID, itemType, blob,
	).Scan(&it.ID, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// ListByUser devolve itens do utilizador; itemType opcional filtra por tipo.
func (r *Repo) ListByUser(ctx context.Context, userID, itemType string) ([]Item, error) {
	if itemType != "" {
		if err := ValidateType(itemType); err != nil {
			return nil, err
		}
	}
	query := `
		SELECT id::text, user_id::text, item_type, blob, created_at, updated_at
		FROM vault_items WHERE user_id = $1`
	args := []any{userID}
	if itemType != "" {
		query += ` AND item_type = $2`
		args = append(args, itemType)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.UserID, &it.Type, &it.Blob, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// GetByID devolve um item se pertencer ao utilizador.
func (r *Repo) GetByID(ctx context.Context, userID, id string) (*Item, error) {
	it := &Item{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, user_id::text, item_type, blob, created_at, updated_at
		FROM vault_items WHERE id = $1 AND user_id = $2`, id, userID,
	).Scan(&it.ID, &it.UserID, &it.Type, &it.Blob, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return it, nil
}

// Update substitui blob e tipo de um item existente.
func (r *Repo) Update(ctx context.Context, userID, id, itemType string, blob []byte) (*Item, error) {
	if err := ValidateType(itemType); err != nil {
		return nil, err
	}
	if len(blob) == 0 {
		return nil, fmt.Errorf("vault: blob vazio")
	}
	it := &Item{}
	err := r.pool.QueryRow(ctx, `
		UPDATE vault_items SET item_type = $3, blob = $4, updated_at = now()
		WHERE id = $1 AND user_id = $2
		RETURNING id::text, user_id::text, item_type, blob, created_at, updated_at`,
		id, userID, itemType, blob,
	).Scan(&it.ID, &it.UserID, &it.Type, &it.Blob, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return it, nil
}

// Delete remove um item (hard delete).
func (r *Repo) Delete(ctx context.Context, userID, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM vault_items WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
