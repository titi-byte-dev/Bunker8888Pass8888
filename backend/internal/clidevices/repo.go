package clidevices

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica dispositivo não registado ou revogado.
var ErrNotFound = errors.New("clidevices: dispositivo não encontrado")

// Device é um registo de certificado CLI.
type Device struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt string
}

// Repo acede à tabela cli_devices.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Register associa um certificado (fingerprint) a um utilizador.
func (r *Repo) Register(ctx context.Context, userID, name string, fingerprint []byte) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO cli_devices (user_id, name, cert_fingerprint)
		VALUES ($1, $2, $3)`,
		userID, name, fingerprint,
	)
	return err
}

// LookupActiveByFingerprint devolve o user_id do dispositivo activo.
func (r *Repo) LookupActiveByFingerprint(ctx context.Context, fingerprint []byte) (string, error) {
	var userID string
	err := r.pool.QueryRow(ctx, `
		SELECT user_id FROM cli_devices
		WHERE cert_fingerprint = $1 AND revoked_at IS NULL`,
		fingerprint,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return userID, nil
}

// ListByUser devolve os dispositivos activos do utilizador.
func (r *Repo) ListByUser(ctx context.Context, userID string) ([]Device, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, user_id::text, name, created_at::text
		FROM cli_devices
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Device
	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// Revoke marca um dispositivo como revogado (só se pertencer ao utilizador).
func (r *Repo) Revoke(ctx context.Context, userID, deviceID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE cli_devices SET revoked_at = now()
		WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL`,
		deviceID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
