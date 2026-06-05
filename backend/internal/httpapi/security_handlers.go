package httpapi

import (
	"net/http"
	"strings"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
)

// handleListSessions devolve sessões HTTP activas do utilizador (metadados).
func handleListSessions(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		currentID := svc.SessionIDForToken(bearerToken(r))

		list, err := svc.ListSessions(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar sessões")
			return
		}
		out := make([]map[string]any, 0, len(list))
		for _, s := range list {
			out = append(out, map[string]any{
				"id":         s.ID,
				"created_at": s.CreatedAt.UTC().Format(time.RFC3339),
				"expires_at": s.ExpiresAt.UTC().Format(time.RFC3339),
				"current":    s.ID == currentID,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"sessions": out})
	}
}

// handleRevokeSession apaga uma sessão concreta (id opaco).
func handleRevokeSession(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeError(w, http.StatusBadRequest, "id em falta")
			return
		}
		if err := svc.RevokeSession(r.Context(), userID, id); err != nil {
			writeError(w, http.StatusNotFound, "sessão não encontrada")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "revogada"})
	}
}

// handleRevokeOtherSessions termina todas as sessões excepto a actual.
func handleRevokeOtherSessions(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		token := bearerToken(r)
		n, err := svc.RevokeOtherSessions(r.Context(), userID, token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao revogar sessões")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"revoked": n})
	}
}

// handleLogout invalida a sessão Bearer actual.
func handleLogout(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "token em falta")
			return
		}
		_ = svc.Logout(r.Context(), token)
		writeJSON(w, http.StatusOK, map[string]string{"status": "terminada"})
	}
}

// handleRevokeCLIDevice revoga um certificado CLI do utilizador.
func handleRevokeCLIDevice(devices *clidevices.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeError(w, http.StatusBadRequest, "id em falta")
			return
		}
		if err := devices.Revoke(r.Context(), userID, id); err != nil {
			writeError(w, http.StatusNotFound, "dispositivo não encontrado")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "revogado"})
	}
}
