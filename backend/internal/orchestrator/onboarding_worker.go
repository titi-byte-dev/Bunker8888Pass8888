package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// OnboardingWorker sugere completar onboarding quando uma ficha vazia é criada.
// Não cifra nem preenche campos — Zero-Knowledge exige Master Key no cliente.
type OnboardingWorker struct {
	Bus *eventbus.Bus
}

func NewOnboardingWorker(bus *eventbus.Bus) *OnboardingWorker {
	return &OnboardingWorker{Bus: bus}
}

func (OnboardingWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "onboarding",
		Description: "Sugere completar onboarding RH quando uma ficha vazia é criada",
		Handles:     []string{eventbus.HREmployeeCreated},
	}
}

func (w *OnboardingWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		RecordID string `json:"record_id"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.onboarding", map[string]any{
		"action":    "run_onboarding",
		"reason":    "hr.employee.created",
		"record_id": meta.RecordID,
		"agent_id":  "onboarding",
		"auto_run":  false,
	})
}
