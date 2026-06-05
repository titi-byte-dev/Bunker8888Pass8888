// Package sessions gere tokens de sessão persistidos na base de dados.
package sessions

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

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

// SessionInfo metadados de uma sessão HTTP (sem expor o token).
type SessionInfo struct {
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// ListByUser devolve sessões activas do utilizador (ordenadas da mais recente).
func (r *Repo) ListByUser(ctx context.Context, userID string) ([]SessionInfo, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT token_hash, created_at, expires_at FROM sessions
		WHERE user_id = $1 AND expires_at > now()
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []SessionInfo
	for rows.Next() {
		var hash []byte
		var s SessionInfo
		if err := rows.Scan(&hash, &s.CreatedAt, &s.ExpiresAt); err != nil {
			return nil, err
		}
		s.ID = hex.EncodeToString(hash)
		out = append(out, s)
	}
	return out, rows.Err()
}

// DeleteByIDForUser revoga uma sessão pelo id opaco (hex do token_hash).
func (r *Repo) DeleteByIDForUser(ctx context.Context, userID, idHex string) error {
	hash, err := hex.DecodeString(idHex)
	if err != nil {
		return ErrInvalidSession
	}
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM sessions WHERE user_id = $1 AND token_hash = $2`, userID, hash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInvalidSession
	}
	return nil
}

// DeleteAllForUserExcept revoga todas as sessões excepto a indicada (logout remoto).
func (r *Repo) DeleteAllForUserExcept(ctx context.Context, userID string, keepHash []byte) (int64, error) {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM sessions WHERE user_id = $1 AND token_hash <> $2`, userID, keepHash)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// Delete remove uma sessão (logout).
func (r *Repo) Delete(ctx context.Context, tokenHash []byte) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}

// DeleteAllForUser revoga TODAS as sessões de um utilizador (remote wipe).
//
// ⚠️ Segurança: mesmo que o dispositivo não receba o push WebSocket, fica
// impossibilitado de voltar a autenticar-se com tokens antigos.
func (r *Repo) DeleteAllForUser(ctx context.Context, userID string) (int64, error) {
	tag, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
