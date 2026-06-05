package mail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRateLimited = errors.New("mail: limite de taxa excedido")

// Direction do evento de mail (auditoria + rate limit).
type Direction string

const (
	DirInbound       Direction = "inbound"
	DirOutboundRelay Direction = "outbound_relay"
	DirCompose       Direction = "compose"
)

// RateConfig limites por hora (MAIL-005). Zero = omissão segura abaixo.
type RateConfig struct {
	InboundPerHour  int
	RelayPerHour    int
	ComposePerHour  int
}

func (c RateConfig) withDefaults() RateConfig {
	out := c
	if out.InboundPerHour <= 0 {
		out.InboundPerHour = 120
	}
	if out.RelayPerHour <= 0 {
		out.RelayPerHour = 60
	}
	if out.ComposePerHour <= 0 {
		out.ComposePerHour = 30
	}
	return out
}

// RateLimiter conta eventos na BD e nega abuso.
type RateLimiter struct {
	pool *pgxpool.Pool
	cfg  RateConfig
}

func NewRateLimiter(pool *pgxpool.Pool, cfg RateConfig) *RateLimiter {
	return &RateLimiter{pool: pool, cfg: cfg.withDefaults()}
}

func (rl *RateLimiter) countSince(ctx context.Context, q string, args ...any) (int, error) {
	if rl == nil || rl.pool == nil {
		return 0, nil
	}
	since := time.Now().Add(-time.Hour)
	all := append([]any{since}, args...)
	var n int
	err := rl.pool.QueryRow(ctx, q, all...).Scan(&n)
	return n, err
}

// AllowInbound verifica limite de ingestão por owner.
func (rl *RateLimiter) AllowInbound(ctx context.Context, ownerID string) error {
	n, err := rl.countSince(ctx, `
		SELECT count(*) FROM mail_relay_log
		WHERE owner_id = $2 AND direction = 'inbound' AND created_at >= $1`, ownerID)
	if err != nil {
		return err
	}
	if n >= rl.cfg.InboundPerHour {
		return fmt.Errorf("%w: inbound %d/h", ErrRateLimited, rl.cfg.InboundPerHour)
	}
	return nil
}

// AllowRelay verifica limite de reencaminhamento por alias.
func (rl *RateLimiter) AllowRelay(ctx context.Context, aliasID string) error {
	if aliasID == "" {
		return nil
	}
	n, err := rl.countSince(ctx, `
		SELECT count(*) FROM mail_relay_log
		WHERE alias_id = $2 AND direction = 'outbound_relay' AND created_at >= $1`, aliasID)
	if err != nil {
		return err
	}
	if n >= rl.cfg.RelayPerHour {
		return fmt.Errorf("%w: relay %d/h", ErrRateLimited, rl.cfg.RelayPerHour)
	}
	return nil
}

// AllowCompose verifica limite de envio manual por owner.
func (rl *RateLimiter) AllowCompose(ctx context.Context, ownerID string) error {
	n, err := rl.countSince(ctx, `
		SELECT count(*) FROM mail_relay_log
		WHERE owner_id = $2 AND direction = 'compose' AND created_at >= $1`, ownerID)
	if err != nil {
		return err
	}
	if n >= rl.cfg.ComposePerHour {
		return fmt.Errorf("%w: compose %d/h", ErrRateLimited, rl.cfg.ComposePerHour)
	}
	return nil
}

// Record regista um evento após operação bem-sucedida.
func (rl *RateLimiter) Record(ctx context.Context, ownerID, aliasID string, dir Direction, from, to string) error {
	if rl == nil || rl.pool == nil {
		return nil
	}
	var aid any = nil
	if aliasID != "" {
		aid = aliasID
	}
	_, err := rl.pool.Exec(ctx, `
		INSERT INTO mail_relay_log (owner_id, alias_id, direction, from_email, to_email)
		VALUES ($1, $2, $3, $4, $5)`,
		ownerID, aid, string(dir), from, to)
	return err
}
