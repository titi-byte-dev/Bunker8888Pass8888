package openbanking

import (
	"context"
	"time"
)

// Service orquestra consentimento e sync com o provider activo.
type Service struct {
	Repo     *Repo
	Provider Provider
}

func NewService(repo *Repo, provider Provider) *Service {
	return &Service{Repo: repo, Provider: provider}
}

// Status devolve a ligação actual (ou pending se ainda não existir).
func (s *Service) Status(ctx context.Context, ownerID string) (*Connection, error) {
	c, err := s.Repo.Get(ctx, ownerID, s.Provider.Name())
	if err != nil {
		return nil, err
	}
	if c == nil {
		return s.Repo.GetOrCreate(ctx, ownerID, s.Provider.Name())
	}
	if c.Status == StatusConnected && c.ConsentExpiresAt != nil && c.ConsentExpiresAt.Before(time.Now()) {
		c.Status = StatusExpired
	}
	return c, nil
}

// Connect simula consentimento PSD2 — em produção redireciona ao banco.
func (s *Service) Connect(ctx context.Context, ownerID string) (*Connection, error) {
	if _, err := s.Repo.GetOrCreate(ctx, ownerID, s.Provider.Name()); err != nil {
		return nil, err
	}
	expires := time.Now().UTC().Add(90 * 24 * time.Hour)
	return s.Repo.MarkConnected(ctx, ownerID, s.Provider.Name(), expires)
}

// Sync obtém movimentos do provider (mock em dev) após ligação activa.
func (s *Service) Sync(ctx context.Context, ownerID string) ([]Transaction, *Connection, error) {
	st, err := s.Status(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}
	if st.Status != StatusConnected {
		return nil, st, ErrNotConnected
	}
	txs, err := s.Provider.ListTransactions(ctx, ownerID)
	if err != nil {
		return nil, nil, err
	}
	if err := s.Repo.TouchSync(ctx, ownerID, s.Provider.Name()); err != nil {
		return nil, nil, err
	}
	st, _ = s.Status(ctx, ownerID)
	return txs, st, nil
}
