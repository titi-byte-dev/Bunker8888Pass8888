// Package sessions gere tokens de sessão persistidos na base de dados.
package sessions

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrInvalidSession indica token inexistente ou expirado.
var ErrInvalidSession = errors.New("sessions: sessão inválida ou expirada")

// Repo dá acesso à tabela "sessions".
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria um repositório de sessões.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Create grava uma sessão. `tokenHash` é o SHA-256 do token (nunca o token cru).
func (r *Repo) Create(ctx context.Context, tokenHash []byte, userID string, ttlSeconds int) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sessions (token_hash, user_id, expires_at)
		VALUES ($1, $2, now() + make_interval(secs => $3))`,
		tokenHash, userID, ttlSeconds,
	)
	return err
}

// UserIDByToken devolve o utilizador dono de um token válido (não expirado).
func (r *Repo) UserIDByToken(ctx context.Context, tokenHash []byte) (string, error) {
	var userID string
	err := r.pool.QueryRow(ctx, `
		SELECT user_id::text FROM sessions
		WHERE token_hash = $1 AND expires_at > now()`, tokenHash,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidSession
		}
		return "", err
	}
	return userID, nil
}

// Delete remove uma sessão (logout).
func (r *Repo) Delete(ctx context.Context, tokenHash []byte) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}
