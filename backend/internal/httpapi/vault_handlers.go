package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

type vaultDeps struct {
	repo *vault.Repo
	hub  *realtime.Hub
}

type itemRequest struct {
	Type string `json:"type"`
	Blob string `json:"blob"`
}

func handleCreateItem(d vaultDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req itemRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		blob, err := decodeItemRequest(req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		item, err := d.repo.Create(r.Context(), userID, req.Type, blob)
		if mapVaultError(w, err) {
			return
		}
		notifyVaultChange(d.hub, userID, realtime.EventCreated, item)
		writeJSON(w, http.StatusCreated, itemJSON(item, false))
	}
}

func handleListItems(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		items, err := repo.ListByUser(r.Context(), userID, r.URL.Query().Get("type"))
		if mapVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(items))
		for _, it := range items {
			out = append(out, itemJSON(&it, true))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	}
}

func handleGetItem(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		item, err := repo.GetByID(r.Context(), userID, r.PathValue("id"))
		if mapVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, itemJSON(item, true))
	}
}

func handleUpdateItem(d vaultDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		var req itemRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		blob, err := decodeItemRequest(req)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		item, err := d.repo.Update(r.Context(), userID, id, req.Type, blob)
		if mapVaultError(w, err) {
			return
		}
		notifyVaultChange(d.hub, userID, realtime.EventUpdated, item)
		writeJSON(w, http.StatusOK, itemJSON(item, false))
	}
}

func handleDeleteItem(d vaultDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		item, err := d.repo.GetByID(r.Context(), userID, id)
		if mapVaultError(w, err) {
			return
		}
		if err := d.repo.Delete(r.Context(), userID, id); mapVaultError(w, err) {
			return
		}
		notifyVaultDelete(d.hub, userID, item)
		w.WriteHeader(http.StatusNoContent)
	}
}

func notifyVaultChange(hub *realtime.Hub, userID, evType string, item *vault.Item) {
	if hub == nil {
		return
	}
	hub.Notify(userID, realtime.Event{
		Type:      evType,
		ItemID:    item.ID,
		ItemType:  item.Type,
		UpdatedAt: item.UpdatedAt,
	})
}

func notifyVaultDelete(hub *realtime.Hub, userID string, item *vault.Item) {
	if hub == nil {
		return
	}
	hub.Notify(userID, realtime.Event{
		Type:     realtime.EventDeleted,
		ItemID:   item.ID,
		ItemType: item.Type,
	})
}

func decodeItemRequest(req itemRequest) ([]byte, error) {
	if req.Type == "" || req.Blob == "" {
		return nil, errBadRequest("type e blob são obrigatórios")
	}
	blob, err := base64.StdEncoding.DecodeString(req.Blob)
	if err != nil {
		return nil, errBadRequest("blob base64 inválido")
	}
	return blob, nil
}

func itemJSON(it *vault.Item, includeBlob bool) map[string]any {
	out := map[string]any{
		"id":         it.ID,
		"type":       it.Type,
		"created_at": it.CreatedAt,
		"updated_at": it.UpdatedAt,
	}
	if includeBlob {
		out["blob"] = base64.StdEncoding.EncodeToString(it.Blob)
	}
	return out
}

func mapVaultError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, vault.ErrNotFound):
		writeError(w, http.StatusNotFound, "item não encontrado")
	case errors.Is(err, vault.ErrInvalidType):
		writeError(w, http.StatusBadRequest, "tipo de item inválido (permitidos: login, note, card)")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação do cofre")
	}
	return true
}

type errBadRequest string

func (e errBadRequest) Error() string { return string(e) }
