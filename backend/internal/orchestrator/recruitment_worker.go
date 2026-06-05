package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// RecruitmentWorker sugere triagem às cegas quando chega e-mail de candidatura.
type RecruitmentWorker struct {
	Bus *eventbus.Bus
}

func NewRecruitmentWorker(bus *eventbus.Bus) *RecruitmentWorker {
	return &RecruitmentWorker{Bus: bus}
}

func (RecruitmentWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "recruitment",
		Description: "Sugere triagem às cegas quando chega candidatura por e-mail",
		Handles:     []string{eventbus.MailInboxReceived},
	}
}

func (w *RecruitmentWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		InboxID   string `json:"inbox_id"`
		Subject   string `json:"subject"`
		FromEmail string `json:"from_email"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	if !agent.IsRecruitmentEmail(meta.Subject, "") {
		return nil
	}
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.recruitment", map[string]any{
		"action":     "screen_candidate",
		"reason":     "mail.inbox.received",
		"inbox_id":   meta.InboxID,
		"from_email": meta.FromEmail,
		"subject":    meta.Subject,
		"agent_id":   "recruitment",
		"auto_run":   false,
	})
}
