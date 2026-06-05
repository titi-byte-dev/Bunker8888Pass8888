package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/fin"
)

// Monitorizacao de custos SaaS (FIN-001). Tudo opaco: o cliente cifra a
// subscricao (nome, custo, ciclo...) com a Master Key e envia o blob em base64.

type subscriptionRequest struct {
	Blob string `json:"blob"`
}

func mapFinError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, fin.ErrNotFound):
		writeError(w, http.StatusNotFound, "subscrição não encontrada")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação de subscrições")
	}
	return true
}

func subscriptionJSON(s *fin.Subscription) map[string]any {
	return map[string]any{
		"id":         s.ID,
		"blob":       base64.StdEncoding.EncodeToString(s.Blob),
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
	}
}

func decodeBlob(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	var req subscriptionRequest
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

func handleListSubscriptions(repo *fin.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		subs, err := repo.List(r.Context(), userID)
		if mapFinError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(subs))
		for i := range subs {
			out = append(out, subscriptionJSON(&subs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"subscriptions": out})
	}
}

func handleCreateSubscription(repo *fin.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, ok := decodeBlob(w, r)
		if !ok {
			return
		}
		s, err := repo.Create(r.Context(), userID, blob)
		if mapFinError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, subscriptionJSON(s))
	}
}

func handleUpdateSubscription(repo *fin.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, ok := decodeBlob(w, r)
		if !ok {
			return
		}
		s, err := repo.Update(r.Context(), userID, r.PathValue("id"), blob)
		if mapFinError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, subscriptionJSON(s))
	}
}

func handleDeleteSubscription(repo *fin.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.Delete(r.Context(), userID, r.PathValue("id")); mapFinError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
