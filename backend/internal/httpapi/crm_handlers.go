package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/crm"
)

type leadRequest struct {
	Blob string `json:"blob"`
}

func mapCRMError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, crm.ErrNotFound):
		writeError(w, http.StatusNotFound, "lead não encontrado")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação CRM")
	}
	return true
}

func leadJSON(l *crm.Lead) map[string]any {
	return map[string]any{
		"id":         l.ID,
		"blob":       base64.StdEncoding.EncodeToString(l.Blob),
		"created_at": l.CreatedAt,
		"updated_at": l.UpdatedAt,
	}
}

func decodeLeadBlob(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	var req leadRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return nil, false
	}
	blob, err := base64.StdEncoding.DecodeString(req.Blob)
	if err != nil || len(blob) == 0 {
		writeError(w, http.StatusBadRequest, "blob base64 inválido")
		return nil, false
	}
	return blob, true
}

func handleListLeads(repo *crm.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		leads, err := repo.List(r.Context(), userID)
		if mapCRMError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(leads))
		for i := range leads {
			out = append(out, leadJSON(&leads[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"leads": out})
	}
}

func handleCreateLead(repo *crm.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, ok := decodeLeadBlob(w, r)
		if !ok {
			return
		}
		l, err := repo.Create(r.Context(), userID, blob)
		if mapCRMError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, leadJSON(l))
	}
}

func handleUpdateLead(repo *crm.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, ok := decodeLeadBlob(w, r)
		if !ok {
			return
		}
		l, err := repo.Update(r.Context(), userID, r.PathValue("id"), blob)
		if mapCRMError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, leadJSON(l))
	}
}

func handleDeleteLead(repo *crm.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.Delete(r.Context(), userID, r.PathValue("id")); mapCRMError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
