package eventbus

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PGStore persiste eventos em agent_events.
type PGStore struct {
	pool *pgxpool.Pool
}

func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool}
}

// Append grava o evento e devolve o id gerado.
func (s *PGStore) Append(ctx context.Context, ev Event) (string, error) {
	if s == nil || s.pool == nil {
		return "", nil
	}
	payload := ev.Payload
	if len(payload) == 0 {
		payload = []byte("{}")
	}
	var id string
	err := s.pool.QueryRow(ctx, `
		INSERT INTO agent_events (user_id, event_type, source, payload)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text`,
		ev.UserID, ev.Type, ev.Source, payload,
	).Scan(&id)
	return id, err
}

// Record é uma linha lida de agent_events.
type Record struct {
	ID        string
	UserID    string
	Type      string
	Source    string
	Payload   []byte
	CreatedAt time.Time
}

// ListRecent devolve eventos do utilizador (feed de acções).
func (s *PGStore) ListRecent(ctx context.Context, userID string, limit int) ([]Record, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id::text, user_id::text, event_type, source, payload, created_at
		FROM agent_events WHERE user_id = $1
		ORDER BY created_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.UserID, &r.Type, &r.Source, &r.Payload, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
