package eventbus

import (
	"context"
	"encoding/json"
	"errors"
)

// ErrSuggestionNotFound — evento de sugestão inexistente ou de outro utilizador.
var ErrSuggestionNotFound = errors.New("sugestão não encontrada")

// ErrAlreadyDecided — a sugestão já foi aprovada ou rejeitada.
var ErrAlreadyDecided = errors.New("sugestão já decidida")

// Decision é approve ou reject (contrato da API AGENT-009).
type Decision string

const (
	DecisionApprove Decision = "approve"
	DecisionReject  Decision = "reject"
)

// suggestionPayload extrai o campo action de uma sugestão do orquestrador.
type suggestionPayload struct {
	Action string `json:"action"`
}

// decisionPayload liga uma decisão ao evento de sugestão original.
type decisionPayload struct {
	SuggestionID string `json:"suggestion_id"`
	Action       string `json:"action"`
}

// ApprovalStore lê sugestões e decisões já tomadas (PGStore em produção).
type ApprovalStore interface {
	GetByID(ctx context.Context, userID, id string) (Record, error)
	DecisionMap(ctx context.Context, userID string) (map[string]string, error)
}

// Decide regista aprovação ou rejeição de uma sugestão e publica no bus.
//
// Didático: o servidor nunca executa acções que exijam Master Key — só regista
// a intenção humana; o cliente corre tools ZK após approve.
func Decide(ctx context.Context, store ApprovalStore, bus *Bus, userID, suggestionID string, d Decision) (string, map[string]any, error) {
	if store == nil {
		return "", nil, ErrSuggestionNotFound
	}
	rec, err := store.GetByID(ctx, userID, suggestionID)
	if err != nil {
		return "", nil, err
	}
	if rec.Type != OrchestratorActionSuggested {
		return "", nil, ErrSuggestionNotFound
	}
	decisions, err := store.DecisionMap(ctx, userID)
	if err != nil {
		return "", nil, err
	}
	if decisions[suggestionID] != "" {
		return "", nil, ErrAlreadyDecided
	}
	var sp suggestionPayload
	_ = json.Unmarshal(rec.Payload, &sp)
	eventType := OrchestratorActionApproved
	if d == DecisionReject {
		eventType = OrchestratorActionRejected
	}
	payload := map[string]any{
		"suggestion_id": suggestionID,
		"action":        sp.Action,
		"decision":      string(d),
	}
	if err := PublishJSON(ctx, bus, eventType, userID, "human.approval", payload); err != nil {
		return "", nil, err
	}
	return sp.Action, payload, nil
}
