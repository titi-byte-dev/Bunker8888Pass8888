package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// FinanceWorker sugere rever licenças SaaS sem uso (AGENT-006 / FIN-002).
type FinanceWorker struct {
	Bus *eventbus.Bus
}

func NewFinanceWorker(bus *eventbus.Bus) *FinanceWorker {
	return &FinanceWorker{Bus: bus}
}

func (FinanceWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "finance",
		Description: "Sugere rever licenças SaaS sem uso ou sem login no cofre",
		Handles:     []string{eventbus.FinSubscriptionStale},
	}
}

func (w *FinanceWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		AlertCount     int      `json:"alert_count"`
		SubscriptionIDs []string `json:"subscription_ids"`
		MonthlySaving  float64  `json:"monthly_saving"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.finance", map[string]any{
		"action":           "review_saas_licenses",
		"reason":           "fin.subscription.stale",
		"alert_count":      meta.AlertCount,
		"subscription_ids": meta.SubscriptionIDs,
		"monthly_saving":   meta.MonthlySaving,
		"agent_id":         "finance",
		"auto_run":         false,
	})
}
