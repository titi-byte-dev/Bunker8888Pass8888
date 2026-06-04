// Package vault guarda e lê itens cifrados do cofre.
//
// ⚠️ O servidor trata o campo Blob como bytes opacos: foi cifrado no cliente e
// nunca é decifrado aqui (modelo Zero-Knowledge).
package vault

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

// ListByUser devolve todos os itens de um utilizador (mais recentes primeiro).
func (r *Repo) ListByUser(ctx context.Context, userID string) ([]Item, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, user_id::text, item_type, blob, created_at, updated_at
		FROM vault_items WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// pgx.CollectRows + RowToStructByName seria uma opção; aqui fazemos o scan
	// manual para deixar explícito o mapeamento coluna -> campo.
	var items []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.UserID, &it.Type, &it.Blob, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	// rows.Err() reporta erros ocorridos DURANTE a iteração (ex: ligação caiu).
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
