package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
)

func handleRunRecruitment(svc *agent.Recruitment, audit *guardian.AuditRepo, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		drafts, err := svc.Run(r.Context(), userID)
		if audit != nil {
			_ = audit.Record(r.Context(), userID, "recruitment", "recruitment_run", err)
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha na triagem")
			return
		}
		_ = eventbus.PublishJSON(r.Context(), bus, eventbus.HRRecruitmentRun, userID, "agent.recruitment", map[string]any{
			"draft_count": len(drafts),
		})
		writeJSON(w, http.StatusOK, map[string]any{
			"agent_id": "recruitment",
			"drafts":   drafts,
			"count":    len(drafts),
		})
	}
}
