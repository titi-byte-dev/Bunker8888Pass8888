package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// ComplianceWorker sugere gerar relatório RGPD após ciclo ERP (DoD Fase 3).
type ComplianceWorker struct {
	Bus *eventbus.Bus
}

func NewComplianceWorker(bus *eventbus.Bus) *ComplianceWorker {
	return &ComplianceWorker{Bus: bus}
}

func (ComplianceWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "compliance",
		Description: "Sugere relatório de conformidade RGPD após fecho operacional",
		Handles:     []string{eventbus.HRComplianceRequested},
	}
}

func (w *ComplianceWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		InvoiceID string `json:"invoice_id"`
		Reason    string `json:"reason"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.compliance", map[string]any{
		"action":     "generate_rgpd_report",
		"reason":     "hr.compliance.requested",
		"invoice_id": meta.InvoiceID,
		"agent_id":   "compliance",
		"auto_run":   false,
	})
}
