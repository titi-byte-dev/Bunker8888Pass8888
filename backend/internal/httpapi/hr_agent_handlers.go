package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

type requestComplianceRequest struct {
	InvoiceID string `json:"invoice_id"`
	Reason    string `json:"reason"`
}

func handleRequestComplianceReport(bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bus == nil {
			writeError(w, http.StatusServiceUnavailable, "event bus indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		var req requestComplianceRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		_ = eventbus.PublishJSON(r.Context(), bus, eventbus.HRComplianceRequested, userID, "hr.compliance", map[string]any{
			"invoice_id": req.InvoiceID,
			"reason":     req.Reason,
		})
		writeJSON(w, http.StatusOK, map[string]any{"status": "requested"})
	}
}
