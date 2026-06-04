package shifts

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica que o utilizador não tem linha em access_shifts.
var ErrNotFound = errors.New("shifts: política não encontrada")

// Repo acede à tabela access_shifts.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório de turnos.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Get devolve a política de turno. Se não existir linha, devolve política
// desactivada (sem restrição).
func (r *Repo) Get(ctx context.Context, userID string) (Policy, error) {
	var p Policy
	var raw []byte
	err := r.pool.QueryRow(ctx, `
		SELECT user_id::text, timezone, schedule, enabled, max_clock_skew_seconds
		FROM access_shifts WHERE user_id = $1`, userID,
	).Scan(&p.UserID, &p.Timezone, &raw, &p.Enabled, &p.MaxClockSkewSecs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Policy{UserID: userID, Enabled: false, Timezone: "UTC"}, nil
		}
		return Policy{}, err
	}
	p.Schedule, err = ParseSchedule(raw)
	if err != nil {
		return Policy{}, err
	}
	return p, nil
}

// Upsert grava ou actualiza a política de turno de um utilizador.
func (r *Repo) Upsert(ctx context.Context, p Policy) error {
	raw, err := MarshalSchedule(p.Schedule)
	if err != nil {
		return err
	}
	if p.Timezone == "" {
		p.Timezone = "UTC"
	}
	if p.MaxClockSkewSecs <= 0 {
		p.MaxClockSkewSecs = 300
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO access_shifts (user_id, timezone, schedule, enabled, max_clock_skew_seconds)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			timezone = EXCLUDED.timezone,
			schedule = EXCLUDED.schedule,
			enabled = EXCLUDED.enabled,
			max_clock_skew_seconds = EXCLUDED.max_clock_skew_seconds,
			updated_at = now()`,
		p.UserID, p.Timezone, raw, p.Enabled, p.MaxClockSkewSecs,
	)
	return err
}
