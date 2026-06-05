package orchestrator

import (
	"context"
	"encoding/json"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// OperationsWorker sugere ordem de compra quando o stock está baixo.
type OperationsWorker struct {
	Bus *eventbus.Bus
}

func NewOperationsWorker(bus *eventbus.Bus) *OperationsWorker {
	return &OperationsWorker{Bus: bus}
}

func (OperationsWorker) Descriptor() Descriptor {
	return Descriptor{
		ID:          "operations",
		Description: "Sugere ordem de compra quando inventário atinge nível de reordenação",
		Handles:     []string{eventbus.OpsStockLow},
	}
}

func (w *OperationsWorker) Handle(ctx context.Context, ev eventbus.Event) error {
	var meta struct {
		ItemID       string `json:"item_id"`
		Name         string `json:"name"`
		SKU          string `json:"sku"`
		Quantity     int    `json:"quantity"`
		ReorderLevel int    `json:"reorder_level"`
		Unit         string `json:"unit"`
	}
	_ = json.Unmarshal(ev.Payload, &meta)
	return eventbus.PublishJSON(ctx, w.Bus, eventbus.OrchestratorActionSuggested, ev.UserID, "orchestrator.operations", map[string]any{
		"action":        "create_purchase_order",
		"reason":        "ops.stock.low",
		"item_id":       meta.ItemID,
		"name":          meta.Name,
		"sku":           meta.SKU,
		"quantity":      meta.Quantity,
		"reorder_level": meta.ReorderLevel,
		"unit":          meta.Unit,
		"agent_id":      "operations",
		"auto_run":      false,
	})
}
