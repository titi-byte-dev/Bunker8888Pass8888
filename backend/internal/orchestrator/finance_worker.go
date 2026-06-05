package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// FinanceWorker sugere acções financeiras: licenças SaaS (FIN-002) e reconciliação (FIN-003).
type FinanceWorker struct {
	Bus *eventbus.Bus
}

func NewFinanceWorker(bus *eventbus.Bus) *FinanceWorker {
	return &FinanceWorker{Bus: bus}
}

func (FinanceWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "finance",
		Description: "Sugere rever licenças SaaS e reconciliar pagamentos bancários",
		Handles: []string{
			eventbus.FinSubscriptionStale,
			eventbus.FinTransactionsSynced,
		},
	}
}

func (w *FinanceWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	switch ev.Type {
	case eventbus.FinSubscriptionStale:
		return w.suggestReviewSaaS(ctx, ev)
	case eventbus.FinTransactionsSynced:
		return w.suggestReconcile(ctx, ev)
	default:
		return nil
	}
}

func (w *FinanceWorker) suggestReviewSaaS(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		AlertCount      int      `json:"alert_count"`
		SubscriptionIDs []string `json:"subscription_ids"`
		MonthlySaving   float64  `json:"monthly_saving"`
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

func (w *FinanceWorker) suggestReconcile(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		TransactionCount int `json:"transaction_count"`
		MatchedCount     int `json:"matched_count"`
		UnmatchedCount   int `json:"unmatched_count"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.finance", map[string]any{
		"action":            "reconcile_payments",
		"reason":            "fin.transactions.synced",
		"transaction_count": meta.TransactionCount,
		"matched_count":     meta.MatchedCount,
		"unmatched_count":   meta.UnmatchedCount,
		"agent_id":          "finance",
		"auto_run":          false,
	})
}
