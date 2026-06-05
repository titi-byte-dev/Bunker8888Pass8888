package sentinel

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrChallengeNotFound = errors.New("sentinel: desafio não encontrado ou expirado")

// LoginEvent é uma tentativa de login registada para auditoria.
type LoginEvent struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Email           string     `json:"email"`
	ClientIP        string     `json:"client_ip"`
	GeoLat          *float64   `json:"geo_lat,omitempty"`
	GeoLon          *float64   `json:"geo_lon,omitempty"`
	Success         bool       `json:"success"`
	Suspicious      bool       `json:"suspicious"`
	Reason          string     `json:"reason"`
	StepUpRequired  bool       `json:"step_up_required"`
	CreatedAt       time.Time  `json:"created_at"`
}

// Challenge é um desafio de step-up pendente (passkey).
type Challenge struct {
	ID        string
	UserID    string
	Reason    string
	Detail    string
	ClientIP  string
	ExpiresAt time.Time
}

// Repo persiste eventos e desafios Sentinel.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// LastSuccessfulGeo devolve o último login bem-sucedido com coordenadas GPS.
func (r *Repo) LastSuccessfulGeo(ctx context.Context, userID string) (Point, bool, error) {
	var lat, lon *float64
	var at time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT geo_lat, geo_lon, created_at
		FROM login_events
		WHERE user_id = $1 AND success = true AND geo_lat IS NOT NULL AND geo_lon IS NOT NULL
		ORDER BY created_at DESC
		LIMIT 1`, userID,
	).Scan(&lat, &lon, &at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Point{}, false, nil
		}
		return Point{}, false, err
	}
	if lat == nil || lon == nil {
		return Point{}, false, nil
	}
	return Point{Lat: *lat, Lon: *lon, At: at.UTC()}, true, nil
}

// InsertEvent regista uma tentativa de login.
func (r *Repo) InsertEvent(ctx context.Context, e LoginEvent) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO login_events
			(user_id, email, client_ip, geo_lat, geo_lon, success, suspicious, reason, step_up_required)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		nullUUID(e.UserID), e.Email, e.ClientIP, e.GeoLat, e.GeoLon,
		e.Success, e.Suspicious, e.Reason, e.StepUpRequired,
	)
	return err
}

func nullUUID(id string) *string {
	if id == "" {
		return nil
	}
	return &id
}

// CreateChallenge grava um desafio de step-up.
func (r *Repo) CreateChallenge(ctx context.Context, c Challenge) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx, `
		INSERT INTO sentinel_challenges (user_id, reason, detail, client_ip, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text`,
		c.UserID, c.Reason, c.Detail, c.ClientIP, c.ExpiresAt,
	).Scan(&id)
	return id, err
}

// GetChallenge devolve um desafio válido (não expirado, não verificado).
func (r *Repo) GetChallenge(ctx context.Context, challengeID string) (Challenge, error) {
	var c Challenge
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, user_id::text, reason, detail, client_ip, expires_at
		FROM sentinel_challenges
		WHERE id = $1 AND verified_at IS NULL AND expires_at > now()`, challengeID,
	).Scan(&c.ID, &c.UserID, &c.Reason, &c.Detail, &c.ClientIP, &c.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Challenge{}, ErrChallengeNotFound
		}
		return Challenge{}, err
	}
	return c, nil
}

// MarkChallengeVerified marca o desafio como concluído.
func (r *Repo) MarkChallengeVerified(ctx context.Context, challengeID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE sentinel_challenges SET verified_at = now()
		WHERE id = $1 AND verified_at IS NULL`, challengeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrChallengeNotFound
	}
	return nil
}

// ListEvents devolve eventos recentes do utilizador.
func (r *Repo) ListEvents(ctx context.Context, userID string, limit int) ([]LoginEvent, error) {
	if limit <= 0 || limit > 100 {
		limit = 30
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, COALESCE(user_id::text, ''), email, client_ip, geo_lat, geo_lon,
		       success, suspicious, reason, step_up_required, created_at
		FROM login_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []LoginEvent
	for rows.Next() {
		var e LoginEvent
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.Email, &e.ClientIP, &e.GeoLat, &e.GeoLon,
			&e.Success, &e.Suspicious, &e.Reason, &e.StepUpRequired, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// CountRecentSuspicious conta alertas suspeitos nas últimas 24h.
func (r *Repo) CountRecentSuspicious(ctx context.Context, userID string) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)::int FROM login_events
		WHERE user_id = $1 AND suspicious = true AND created_at > now() - interval '24 hours'`,
		userID,
	).Scan(&n)
	return n, err
}
