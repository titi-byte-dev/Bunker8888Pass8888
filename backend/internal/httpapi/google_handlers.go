package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/googleworkspace"
)

func handleGoogleWorkspaceStatus(svc *googleworkspace.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			writeError(w, http.StatusServiceUnavailable, "google workspace indisponível")
			return
		}
		st := svc.Status(r.Context())
		writeJSON(w, http.StatusOK, map[string]any{"status": st})
	}
}
