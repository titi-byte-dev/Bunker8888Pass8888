package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// CRMDealWorker sugere emitir pro-forma quando um negócio fecha (CRM-003).
type CRMDealWorker struct {
	Bus *eventbus.Bus
}

func NewCRMDealWorker(bus *eventbus.Bus) *CRMDealWorker {
	return &CRMDealWorker{Bus: bus}
}

func (CRMDealWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "crm_deal",
		Description: "Sugere emitir pro-forma quando um lead passa a ganho",
		Handles:     []string{eventbus.CRMDealClosed},
	}
}

func (w *CRMDealWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		LeadID string `json:"lead_id"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.crm", map[string]any{
		"action":   "issue_proforma",
		"reason":   "crm.deal_closed",
		"lead_id":  meta.LeadID,
		"agent_id": "crm_deal",
		"auto_run": false,
	})
}
