package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/ops"
)

type inventoryRequest struct {
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Quantity     int    `json:"quantity"`
	ReorderLevel int    `json:"reorder_level"`
	Unit         string `json:"unit"`
}

type adjustRequest struct {
	Delta int `json:"delta"`
}

func mapOpsError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, ops.ErrNotFound):
		writeError(w, http.StatusNotFound, "artigo não encontrado")
	case err.Error() == "ops: nome obrigatório" || err.Error() == "ops: quantidades inválidas":
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação de inventário")
	}
	return true
}

func inventoryItemJSON(it *ops.Item) map[string]any {
	return map[string]any{
		"id":            it.ID,
		"name":          it.Name,
		"sku":           it.SKU,
		"quantity":      it.Quantity,
		"reorder_level": it.ReorderLevel,
		"unit":          it.Unit,
		"low_stock":     ops.IsLowStock(it.Quantity, it.ReorderLevel),
		"created_at":    it.CreatedAt,
		"updated_at":    it.UpdatedAt,
	}
}

func publishLowStock(ctx context.Context, bus *eventbus.Bus, userID string, it *ops.Item) {
	if bus == nil || !ops.IsLowStock(it.Quantity, it.ReorderLevel) {
		return
	}
	_ = eventbus.PublishJSON(ctx, bus, eventbus.OpsStockLow, userID, "ops.inventory", map[string]any{
		"item_id":       it.ID,
		"name":          it.Name,
		"sku":           it.SKU,
		"quantity":      it.Quantity,
		"reorder_level": it.ReorderLevel,
		"unit":          it.Unit,
	})
}

func handleListInventory(repo *ops.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		items, err := repo.List(r.Context(), userID)
		if mapOpsError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(items))
		for i := range items {
			out = append(out, inventoryItemJSON(&items[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	}
}

func handleCreateInventory(repo *ops.Repo, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req inventoryRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		it, err := repo.Create(r.Context(), userID, ops.CreateInput{
			Name: req.Name, SKU: req.SKU, Quantity: req.Quantity,
			ReorderLevel: req.ReorderLevel, Unit: req.Unit,
		})
		if mapOpsError(w, err) {
			return
		}
		publishLowStock(r.Context(), bus, userID, it)
		writeJSON(w, http.StatusCreated, inventoryItemJSON(it))
	}
}

func handleUpdateInventory(repo *ops.Repo, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		var req inventoryRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		it, err := repo.Update(r.Context(), userID, id, ops.UpdateInput{
			Name: req.Name, SKU: req.SKU, Quantity: req.Quantity,
			ReorderLevel: req.ReorderLevel, Unit: req.Unit,
		})
		if mapOpsError(w, err) {
			return
		}
		publishLowStock(r.Context(), bus, userID, it)
		writeJSON(w, http.StatusOK, inventoryItemJSON(it))
	}
}

func handleAdjustInventory(repo *ops.Repo, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		var req adjustRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		it, err := repo.Adjust(r.Context(), userID, id, req.Delta)
		if mapOpsError(w, err) {
			return
		}
		publishLowStock(r.Context(), bus, userID, it)
		writeJSON(w, http.StatusOK, inventoryItemJSON(it))
	}
}

func handleDeleteInventory(repo *ops.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		if err := repo.Delete(r.Context(), userID, id); mapOpsError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
