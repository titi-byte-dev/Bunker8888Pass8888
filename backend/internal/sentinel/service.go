package sentinel

import (
	"context"
	"time"
)

const challengeTTL = 10 * time.Minute

// Service coordena avaliação, desafios e auditoria Sentinel.
type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// LoginContext traz metadados do pedido de login actual.
type LoginContext struct {
	UserID   string
	Email    string
	ClientIP string
	GeoLat   *float64
	GeoLon   *float64
	At       time.Time
}

// StepUpRequired indica que o login precisa de passkey antes de emitir token.
type StepUpRequired struct {
	ChallengeID string
	Reason      string
	Detail      string
}

// Evaluate verifica se o login actual parece impossível.
func (s *Service) Evaluate(ctx context.Context, lc LoginContext) (Assessment, error) {
	if lc.GeoLat == nil || lc.GeoLon == nil || lc.UserID == "" {
		return Assessment{}, nil
	}
	prev, ok, err := s.repo.LastSuccessfulGeo(ctx, lc.UserID)
	if err != nil {
		return Assessment{}, err
	}
	if !ok {
		return Assessment{}, nil
	}
	curr := Point{Lat: *lc.GeoLat, Lon: *lc.GeoLon, At: lc.At.UTC()}
	return AssessTravel(prev, curr), nil
}

// RecordFailure regista tentativa falhada (sem step-up).
func (s *Service) RecordFailure(ctx context.Context, lc LoginContext) error {
	return s.repo.InsertEvent(ctx, LoginEvent{
		UserID:   lc.UserID,
		Email:    lc.Email,
		ClientIP: lc.ClientIP,
		GeoLat:   lc.GeoLat,
		GeoLon:   lc.GeoLon,
		Success:  false,
	})
}

// CompleteSuccess regista login bem-sucedido (após token emitido).
func (s *Service) CompleteSuccess(ctx context.Context, lc LoginContext, a Assessment, stepUp bool) error {
	return s.repo.InsertEvent(ctx, LoginEvent{
		UserID:         lc.UserID,
		Email:          lc.Email,
		ClientIP:       lc.ClientIP,
		GeoLat:         lc.GeoLat,
		GeoLon:         lc.GeoLon,
		Success:        true,
		Suspicious:     a.Suspicious,
		Reason:         a.Reason,
		StepUpRequired: stepUp,
	})
}

// CreateStepUpChallenge cria desafio passkey pendente.
func (s *Service) CreateStepUpChallenge(ctx context.Context, lc LoginContext, a Assessment) (StepUpRequired, error) {
	id, err := s.repo.CreateChallenge(ctx, Challenge{
		UserID:    lc.UserID,
		Reason:    a.Reason,
		Detail:    a.Detail,
		ClientIP:  lc.ClientIP,
		ExpiresAt: lc.At.UTC().Add(challengeTTL),
	})
	if err != nil {
		return StepUpRequired{}, err
	}
	_ = s.repo.InsertEvent(ctx, LoginEvent{
		UserID:         lc.UserID,
		Email:          lc.Email,
		ClientIP:       lc.ClientIP,
		GeoLat:         lc.GeoLat,
		GeoLon:         lc.GeoLon,
		Success:        false,
		Suspicious:     true,
		Reason:         a.Reason,
		StepUpRequired: true,
	})
	return StepUpRequired{ChallengeID: id, Reason: a.Reason, Detail: a.Detail}, nil
}

// GetChallenge devolve desafio válido.
func (s *Service) GetChallenge(ctx context.Context, id string) (Challenge, error) {
	return s.repo.GetChallenge(ctx, id)
}

// VerifyChallenge marca desafio como verificado após passkey OK.
func (s *Service) VerifyChallenge(ctx context.Context, challengeID, userID string) error {
	c, err := s.repo.GetChallenge(ctx, challengeID)
	if err != nil {
		return err
	}
	if c.UserID != userID {
		return ErrChallengeNotFound
	}
	return s.repo.MarkChallengeVerified(ctx, challengeID)
}

// ListEvents devolve histórico para o utilizador autenticado.
func (s *Service) ListEvents(ctx context.Context, userID string) ([]LoginEvent, int, error) {
	events, err := s.repo.ListEvents(ctx, userID, 30)
	if err != nil {
		return nil, 0, err
	}
	n, err := s.repo.CountRecentSuspicious(ctx, userID)
	if err != nil {
		return events, 0, err
	}
	return events, n, nil
}
