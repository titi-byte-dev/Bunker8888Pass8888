package security

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// WipeAuditEvent é uma entrada do log de remote wipe (append-only).
type WipeAuditEvent struct {
	ID               string    `json:"id"`
	TargetUserID     string    `json:"target_user_id"`
	TargetEmail      string    `json:"target_email"`
	InitiatedBy      string    `json:"initiated_by"`
	Reason           string    `json:"reason"`
	DevicesNotified  int       `json:"devices_notified"`
	CreatedAt        time.Time `json:"created_at"`
}

// ListWipeEvents devolve eventos de auditoria mais recentes primeiro.
func ListWipeEvents(ctx context.Context, pool *pgxpool.Pool, limit int) ([]WipeAuditEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := pool.Query(ctx, `
		SELECT e.id::text, e.target_user_id::text, COALESCE(u.email, ''), e.initiated_by,
		       e.reason, e.devices_notified, e.created_at
		FROM remote_wipe_events e
		LEFT JOIN users u ON u.id = e.target_user_id
		ORDER BY e.created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("listar auditoria wipe: %w", err)
	}
	defer rows.Close()

	var out []WipeAuditEvent
	for rows.Next() {
		var e WipeAuditEvent
		if err := rows.Scan(
			&e.ID, &e.TargetUserID, &e.TargetEmail, &e.InitiatedBy,
			&e.Reason, &e.DevicesNotified, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
