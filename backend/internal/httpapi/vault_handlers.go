package httpapi

import (
	"encoding/base64"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

// createItemRequest é o corpo de POST /api/vault. O blob já vem cifrado do cliente.
type createItemRequest struct {
	Type string `json:"type"`
	Blob string `json:"blob"` // base64 de nonce||ciphertext||tag
}

func handleCreateItem(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)

		var req createItemRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.Type == "" || req.Blob == "" {
			writeError(w, http.StatusBadRequest, "type e blob são obrigatórios")
			return
		}
		blob, err := base64.StdEncoding.DecodeString(req.Blob)
		if err != nil {
			writeError(w, http.StatusBadRequest, "blob base64 inválido")
			return
		}

		item, err := repo.Create(r.Context(), userID, req.Type, blob)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao gravar item")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":         item.ID,
			"type":       item.Type,
			"created_at": item.CreatedAt,
		})
	}
}

func handleListItems(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)

		items, err := repo.ListByUser(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar itens")
			return
		}

		out := make([]map[string]any, 0, len(items))
		for _, it := range items {
			out = append(out, map[string]any{
				"id":         it.ID,
				"type":       it.Type,
				"blob":       base64.StdEncoding.EncodeToString(it.Blob),
				"created_at": it.CreatedAt,
				"updated_at": it.UpdatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	}
}
