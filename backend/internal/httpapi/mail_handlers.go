package httpapi

import (
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

// Aliases de e-mail (MAIL-001). O destino e visivel ao servidor (relay futuro);
// o resto da app continua Zero-Knowledge.

type createAliasRequest struct {
	Destination string `json:"destination"`
	Label       string `json:"label"`
}

type setAliasActiveRequest struct {
	Active bool `json:"active"`
}

func mapMailError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, mail.ErrNotFound):
		writeError(w, http.StatusNotFound, "alias não encontrado")
	case errors.Is(err, mail.ErrInvalidDest):
		writeError(w, http.StatusBadRequest, "destino de e-mail inválido")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação de aliases")
	}
	return true
}

func aliasJSON(a *mail.Alias) map[string]any {
	return map[string]any{
		"id":            a.ID,
		"alias_address": a.AliasAddress,
		"destination":   a.Destination,
		"label":         a.Label,
		"active":        a.Active,
		"created_at":    a.CreatedAt,
	}
}

func handleListAliases(repo *mail.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		aliases, err := repo.ListAliases(r.Context(), userID)
		if mapMailError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(aliases))
		for i := range aliases {
			out = append(out, aliasJSON(&aliases[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"aliases": out})
	}
}

func handleCreateAlias(repo *mail.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req createAliasRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		a, err := repo.CreateAlias(r.Context(), userID, req.Destination, req.Label)
		if mapMailError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, aliasJSON(a))
	}
}

func handleSetAliasActive(repo *mail.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req setAliasActiveRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if err := repo.SetActive(r.Context(), userID, r.PathValue("id"), req.Active); mapMailError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleDeleteAlias(repo *mail.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteAlias(r.Context(), userID, r.PathValue("id")); mapMailError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
