package emergency

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// Service orquestra regras de negócio de emergência.
type Service struct {
	repo  *Repo
	users *users.Repo
}

// NewService constrói o serviço.
func NewService(repo *Repo, users *users.Repo) *Service {
	return &Service{repo: repo, users: users}
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

// SetConfig grava ou actualiza a configuração do herdeiro.
func (s *Service) SetConfig(ctx context.Context, ownerID, heirEmail string, waitDays int, blob []byte) error {
	heirEmail = normalizeEmail(heirEmail)
	if heirEmail == "" {
		return errors.New("emergency: email do herdeiro em falta")
	}
	if waitDays < 1 || waitDays > 90 {
		return errors.New("emergency: wait_days inválido (1–90)")
	}
	owner, err := s.users.ByID(ctx, ownerID)
	if err != nil {
		return err
	}
	if normalizeEmail(owner.Email) == heirEmail {
		return ErrSelfHeir
	}
	return s.repo.UpsertConfig(ctx, ownerID, heirEmail, waitDays, blob)
}

// GetConfig devolve a configuração do titular.
func (s *Service) GetConfig(ctx context.Context, ownerID string) (*Config, error) {
	return s.repo.GetConfig(ctx, ownerID)
}

// DeleteConfig remove herdeiro e blob.
func (s *Service) DeleteConfig(ctx context.Context, ownerID string) error {
	return s.repo.DeleteConfig(ctx, ownerID)
}

// CreateRequest inicia um pedido de acesso (herdeiro autenticado).
func (s *Service) CreateRequest(ctx context.Context, heirUserID, ownerEmail string) (*Request, error) {
	owner, err := s.users.ByEmail(ctx, normalizeEmail(ownerEmail))
	if err != nil {
		return nil, err
	}
	if owner.ID == heirUserID {
		return nil, ErrSelfHeir
	}

	heir, err := s.users.ByID(ctx, heirUserID)
	if err != nil {
		return nil, err
	}
	heirEmail := normalizeEmail(heir.Email)

	cfg, err := s.repo.GetConfig(ctx, owner.ID)
	if err != nil {
		return nil, err
	}
	if cfg.HeirEmail != heirEmail {
		return nil, ErrHeirMismatch
	}

	if active, err := s.repo.GetActiveRequest(ctx, owner.ID, heirUserID); err != nil {
		return nil, err
	} else if active != nil {
		return nil, ErrActiveRequest
	}

	unlocksAt := time.Now().UTC().Add(time.Duration(cfg.WaitDays) * 24 * time.Hour)
	return s.repo.CreateRequest(ctx, owner.ID, heirUserID, unlocksAt)
}

// ListRequestsForOwner lista pedidos recebidos pelo titular.
func (s *Service) ListRequestsForOwner(ctx context.Context, ownerID string) ([]Request, error) {
	reqs, err := s.repo.ListRequestsByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for i := range reqs {
		s.applyPromotion(ctx, &reqs[i], now)
	}
	return reqs, nil
}

// RejectRequest rejeita um pedido pendente.
func (s *Service) RejectRequest(ctx context.Context, ownerID, requestID string) error {
	req, err := s.repo.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}
	if req.OwnerUserID != ownerID {
		return ErrForbidden
	}
	if req.Status != StatusWaiting {
		return ErrForbidden
	}
	return s.repo.SetRejected(ctx, requestID)
}

// ApproveEarly torna o pedido imediatamente disponível.
func (s *Service) ApproveEarly(ctx context.Context, ownerID, requestID string) error {
	req, err := s.repo.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}
	if req.OwnerUserID != ownerID {
		return ErrForbidden
	}
	if req.Status != StatusWaiting {
		return ErrForbidden
	}
	return s.repo.SetReady(ctx, requestID, time.Now().UTC())
}

// GetHeirRequest devolve o pedido activo do herdeiro para um titular.
func (s *Service) GetHeirRequest(ctx context.Context, heirUserID, ownerEmail string) (*Request, error) {
	owner, err := s.users.ByEmail(ctx, normalizeEmail(ownerEmail))
	if err != nil {
		return nil, err
	}
	req, err := s.repo.GetActiveRequest(ctx, owner.ID, heirUserID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	now := time.Now().UTC()
	s.applyPromotion(ctx, req, now)
	return req, nil
}

// FetchAccessBlob devolve o blob cifrado quando o pedido está ready.
func (s *Service) FetchAccessBlob(ctx context.Context, heirUserID, ownerEmail string) ([]byte, error) {
	owner, err := s.users.ByEmail(ctx, normalizeEmail(ownerEmail))
	if err != nil {
		return nil, err
	}
	req, err := s.repo.GetActiveRequest(ctx, owner.ID, heirUserID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}
	now := time.Now().UTC()
	s.applyPromotion(ctx, req, now)
	if req.Status != StatusReady {
		return nil, ErrNotReady
	}

	blob, err := s.repo.GetEncryptedBlob(ctx, owner.ID)
	if err != nil {
		return nil, err
	}
	if len(blob) == 0 {
		return nil, ErrNoEncryptedBlob
	}
	if err := s.repo.SetConsumed(ctx, req.ID); err != nil {
		return nil, err
	}
	return blob, nil
}

func (s *Service) applyPromotion(ctx context.Context, req *Request, now time.Time) {
	next := PromoteIfExpired(now, req.Status, req.UnlocksAt)
	if next == req.Status {
		return
	}
	if err := s.repo.SetReady(ctx, req.ID, req.UnlocksAt); err == nil {
		req.Status = StatusReady
	}
}
