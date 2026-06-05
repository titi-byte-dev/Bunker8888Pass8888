package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/orchestrator"
)

func handleOrchestratorStatus(orc *orchestrator.Orchestrator) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if orc == nil {
			writeJSON(w, http.StatusOK, map[string]any{"agents": []any{}})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"agents": orc.Descriptors(),
		})
	}
}
