package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

type reportStaleRequest struct {
	SubscriptionIDs []string `json:"subscription_ids"`
	AlertCount      int      `json:"alert_count"`
	MonthlySaving   float64  `json:"monthly_saving"`
}

// handleReportStaleSubscriptions publica fin.subscription.stale (AGENT-006).
// Didático: o cliente envia só IDs — os nomes/custos nunca saem do browser em claro.
func handleReportStaleSubscriptions(bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bus == nil {
			writeError(w, http.StatusServiceUnavailable, "event bus indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		var req reportStaleRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.AlertCount <= 0 || len(req.SubscriptionIDs) == 0 {
			writeError(w, http.StatusBadRequest, "sem alertas para reportar")
			return
		}
		_ = eventbus.PublishJSON(r.Context(), bus, eventbus.FinSubscriptionStale, userID, "fin.alerts", map[string]any{
			"subscription_ids": req.SubscriptionIDs,
			"alert_count":      req.AlertCount,
			"monthly_saving":   req.MonthlySaving,
		})
		writeJSON(w, http.StatusOK, map[string]any{"status": "reported", "alert_count": req.AlertCount})
	}
}

type reportSyncRequest struct {
	TransactionCount int `json:"transaction_count"`
	MatchedCount     int `json:"matched_count"`
	UnmatchedCount   int `json:"unmatched_count"`
}

// handleReportTransactionsSynced publica fin.transactions.synced (FIN-003 / AGENT-006).
// Didático: o servidor recebe apenas contagens — os movimentos ficam no cliente.
func handleReportTransactionsSynced(bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bus == nil {
			writeError(w, http.StatusServiceUnavailable, "event bus indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		var req reportSyncRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.TransactionCount <= 0 {
			writeError(w, http.StatusBadRequest, "sem movimentos para reportar")
			return
		}
		_ = eventbus.PublishJSON(r.Context(), bus, eventbus.FinTransactionsSynced, userID, "fin.banking", map[string]any{
			"transaction_count": req.TransactionCount,
			"matched_count":     req.MatchedCount,
			"unmatched_count":   req.UnmatchedCount,
		})
		writeJSON(w, http.StatusOK, map[string]any{
			"status":            "reported",
			"transaction_count": req.TransactionCount,
			"unmatched_count":   req.UnmatchedCount,
		})
	}
}
