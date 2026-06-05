package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// ProspectionWorker sugere prospeção quando chega e-mail à inbox.
// Não executa tools automaticamente — Zero-Knowledge exige unlock humano.
type ProspectionWorker struct {
	Bus *eventbus.Bus
}

func NewProspectionWorker(bus *eventbus.Bus) *ProspectionWorker {
	return &ProspectionWorker{Bus: bus}
}

func (ProspectionWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "prospection",
		Description: "Sugere correr prospeção CRM quando chega e-mail ao alias",
		Handles:     []string{eventbus.MailInboxReceived},
	}
}

func (w *ProspectionWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		InboxID string `json:"inbox_id"`
		Alias   string `json:"alias"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.prospection", map[string]any{
		"action":    "run_prospection",
		"reason":    "mail.inbox.received",
		"inbox_id":  meta.InboxID,
		"alias":     meta.Alias,
		"agent_id":  "prospection",
		"auto_run":  false,
	})
}
