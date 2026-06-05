package guardian

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditEntry regista uma execução de tool (sem input/output sensível).
type AuditEntry struct {
	ID        string
	UserID    string
	AgentID   string
	ToolName  string
	Success   bool
	ErrorMsg  string
	CreatedAt time.Time
}

// AuditRepo persiste eventos em agent_tool_audit.
type AuditRepo struct {
	pool *pgxpool.Pool
}

func NewAuditRepo(pool *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{pool: pool}
}

// Record grava o resultado de uma execução.
func (r *AuditRepo) Record(ctx context.Context, userID, agentID, toolName string, runErr error) error {
	if r == nil || r.pool == nil {
		return nil
	}
	success := runErr == nil
	var errMsg *string
	if runErr != nil {
		s := runErr.Error()
		if len(s) > 500 {
			s = s[:500]
		}
		errMsg = &s
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO agent_tool_audit (user_id, agent_id, tool_name, success, error_msg)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, agentID, toolName, success, errMsg,
	)
	return err
}

// ListRecent devolve os últimos eventos do utilizador (painel de auditoria).
func (r *AuditRepo) ListRecent(ctx context.Context, userID string, limit int) ([]AuditEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, user_id::text, agent_id, tool_name, success,
		       COALESCE(error_msg, ''), created_at
		FROM agent_tool_audit
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.AgentID, &e.ToolName, &e.Success, &e.ErrorMsg, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
