// Package recovery guarda backups cifrados da Master Key (VAULT-018).
package recovery

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica que o utilizador não tem backup de recuperação.
var ErrNotFound = errors.New("recovery: backup não encontrado")

// Repo acede à tabela recovery_backups.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// GetByUserID devolve o blob cifrado do utilizador autenticado.
func (r *Repo) GetByUserID(ctx context.Context, userID string) ([]byte, error) {
	var blob []byte
	err := r.pool.QueryRow(ctx,
		`SELECT blob FROM recovery_backups WHERE user_id = $1`, userID,
	).Scan(&blob)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return blob, nil
}

// GetByEmail devolve o blob para recuperação (lookup por email).
func (r *Repo) GetByEmail(ctx context.Context, email string) ([]byte, error) {
	var blob []byte
	err := r.pool.QueryRow(ctx, `
		SELECT rb.blob FROM recovery_backups rb
		JOIN users u ON u.id = rb.user_id
		WHERE u.email = $1`, email,
	).Scan(&blob)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return blob, nil
}

// Upsert grava ou actualiza o backup cifrado.
func (r *Repo) Upsert(ctx context.Context, userID string, blob []byte) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO recovery_backups (user_id, blob)
		VALUES ($1, $2)
	 ON CONFLICT (user_id) DO UPDATE SET blob = EXCLUDED.blob, updated_at = now()`,
		userID, blob,
	)
	return err
}

// HasBackup indica se existe backup para o utilizador.
func (r *Repo) HasBackup(ctx context.Context, userID string) (bool, error) {
	var ok bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM recovery_backups WHERE user_id = $1)`, userID,
	).Scan(&ok)
	return ok, err
}
