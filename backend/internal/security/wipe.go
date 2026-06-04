// Package security implementa acções de segurança do servidor (VAULT-012).
package security

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sessions"
)

// WipeService coordena remote wipe: push WebSocket + revogação de sessões + audit.
type WipeService struct {
	sessions *sessions.Repo
	hub      *realtime.Hub
	pool     *pgxpool.Pool
}

// NewWipeService constrói o serviço de remote wipe.
func NewWipeService(sessions *sessions.Repo, hub *realtime.Hub, pool *pgxpool.Pool) *WipeService {
	return &WipeService{sessions: sessions, hub: hub, pool: pool}
}

// RemoteWipeResult resume o que aconteceu num wipe.
type RemoteWipeResult struct {
	DevicesNotified int   `json:"devices_notified"`
	SessionsRevoked int64 `json:"sessions_revoked"`
}

// ExecuteRemoteWipe invalida sessões e envia push wipe aos dispositivos online.
//
// Ordem intencional:
//  1. Contar/notificar via WebSocket (dispositivos ainda autenticados na ligação WS)
//  2. Apagar sessões (tokens HTTP deixam de funcionar)
//  3. Registar auditoria
func (s *WipeService) ExecuteRemoteWipe(
	ctx context.Context,
	targetUserID, initiatedBy, reason string,
) (*RemoteWipeResult, error) {
	if s.hub == nil {
		return nil, fmt.Errorf("security: hub WebSocket indisponível")
	}

	devices := s.hub.ClientCount(targetUserID)
	s.hub.NotifyWipe(targetUserID, reason)

	revoked, err := s.sessions.DeleteAllForUser(ctx, targetUserID)
	if err != nil {
		return nil, fmt.Errorf("revogar sessões: %w", err)
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO remote_wipe_events (target_user_id, initiated_by, reason, devices_notified)
		VALUES ($1, $2, $3, $4)`,
		targetUserID, initiatedBy, reason, devices,
	)
	if err != nil {
		return nil, fmt.Errorf("registar auditoria: %w", err)
	}

	return &RemoteWipeResult{
		DevicesNotified: devices,
		SessionsRevoked: revoked,
	}, nil
}
