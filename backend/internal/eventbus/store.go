package eventbus

import (
	"context"
	"encoding/json"
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

// GetByID devolve um evento do utilizador (isolamento multi-tenant).
func (s *PGStore) GetByID(ctx context.Context, userID, id string) (Record, error) {
	if s == nil || s.pool == nil {
		return Record{}, ErrSuggestionNotFound
	}
	var r Record
	err := s.pool.QueryRow(ctx, `
		SELECT id::text, user_id::text, event_type, source, payload, created_at
		FROM agent_events WHERE id = $1 AND user_id = $2`, id, userID,
	).Scan(&r.ID, &r.UserID, &r.Type, &r.Source, &r.Payload, &r.CreatedAt)
	if err != nil {
		return Record{}, ErrSuggestionNotFound
	}
	return r, nil
}

// DecisionMap mapeia suggestion_id → approved|rejected a partir do feed recente.
func (s *PGStore) DecisionMap(ctx context.Context, userID string) (map[string]string, error) {
	if s == nil || s.pool == nil {
		return map[string]string{}, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT event_type, payload
		FROM agent_events
		WHERE user_id = $1
		  AND event_type IN ($2, $3)
		ORDER BY created_at DESC
		LIMIT 200`, userID, OrchestratorActionApproved, OrchestratorActionRejected)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]string)
	for rows.Next() {
		var eventType string
		var payload []byte
		if err := rows.Scan(&eventType, &payload); err != nil {
			return nil, err
		}
		var dp struct {
			SuggestionID string `json:"suggestion_id"`
		}
		_ = json.Unmarshal(payload, &dp)
		if dp.SuggestionID == "" || out[dp.SuggestionID] != "" {
			continue
		}
		switch eventType {
		case OrchestratorActionApproved:
			out[dp.SuggestionID] = "approved"
		case OrchestratorActionRejected:
			out[dp.SuggestionID] = "rejected"
		}
	}
	return out, rows.Err()
}
