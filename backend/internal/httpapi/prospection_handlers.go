package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
)

func handleRunProspection(svc *agent.Prospection, audit *guardian.AuditRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		drafts, err := svc.Run(r.Context(), userID)
		if audit != nil {
			_ = audit.Record(r.Context(), userID, "prospection", "prospection_run", err)
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha na prospeção")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"agent_id": "prospection",
			"drafts":   drafts,
			"count":    len(drafts),
		})
	}
}
