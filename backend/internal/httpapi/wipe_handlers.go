package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

type wipeRequest struct {
	Reason string `json:"reason"`
}

// handleAdminRemoteWipe permite a um operador (chave admin) apagar dados locais
// nos dispositivos de um utilizador e revogar todas as sessões.
func handleAdminRemoteWipe(userRepo *users.Repo, wipe *security.WipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.PathValue("id")
		if targetID == "" {
			writeError(w, http.StatusBadRequest, "user id em falta")
			return
		}

		if _, err := userRepo.ByID(r.Context(), targetID); err != nil {
			if err == users.ErrNotFound {
				writeError(w, http.StatusNotFound, "utilizador não encontrado")
				return
			}
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}

		var req wipeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		result, err := wipe.ExecuteRemoteWipe(r.Context(), targetID, "admin", req.Reason)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "remote wipe falhou")
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

// handleSelfRemoteWipe permite ao utilizador autenticado apagar os seus próprios
// dispositivos (ex.: telemóvel perdido).
func handleSelfRemoteWipe(wipe *security.WipeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req wipeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		result, err := wipe.ExecuteRemoteWipe(r.Context(), userID, "self", req.Reason)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "remote wipe falhou")
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}
