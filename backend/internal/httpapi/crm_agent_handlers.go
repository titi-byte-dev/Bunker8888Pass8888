package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

type reportDealClosedRequest struct {
	LeadID string `json:"lead_id"`
}

func handleReportDealClosed(bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if bus == nil {
			writeError(w, http.StatusServiceUnavailable, "event bus indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		var req reportDealClosedRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.LeadID == "" {
			writeError(w, http.StatusBadRequest, "lead_id obrigatório")
			return
		}
		_ = eventbus.PublishJSON(r.Context(), bus, eventbus.CRMDealClosed, userID, "crm.deal", map[string]any{
			"lead_id": req.LeadID,
		})
		writeJSON(w, http.StatusOK, map[string]any{"status": "reported", "lead_id": req.LeadID})
	}
}
