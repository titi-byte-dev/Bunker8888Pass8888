package httpapi

import (
	"encoding/base64"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharedvaults"
)

// Anexos cifrados por ficheiro (SHARE-004). Tal como os itens, tudo é opaco: o
// cliente cifra os metadados e os bytes do ficheiro com a chave do cofre e
// envia em base64. O servidor só faz cumprir as permissões e o limite de tamanho.

type vaultAttachmentRequest struct {
	MetaBlob string `json:"meta_blob"` // nome/tipo/tamanho cifrados
	DataBlob string `json:"data_blob"` // bytes do ficheiro cifrados
}

func attachmentMetaJSON(a *sharedvaults.AttachmentMeta) map[string]any {
	return map[string]any{
		"id":         a.ID,
		"meta_blob":  base64.StdEncoding.EncodeToString(a.MetaBlob),
		"byte_size":  a.ByteSize,
		"created_by": a.CreatedBy,
		"created_at": a.CreatedAt,
	}
}

func handleListVaultAttachments(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		atts, err := repo.ListAttachments(r.Context(), r.PathValue("id"), userID)
		if mapSharedVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(atts))
		for i := range atts {
			out = append(out, attachmentMetaJSON(&atts[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"attachments": out})
	}
}

func handleAddVaultAttachment(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req vaultAttachmentRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		meta, err := base64.StdEncoding.DecodeString(req.MetaBlob)
		if err != nil || len(meta) == 0 {
			writeError(w, http.StatusBadRequest, "meta_blob base64 inválido")
			return
		}
		data, err := base64.StdEncoding.DecodeString(req.DataBlob)
		if err != nil || len(data) == 0 {
			writeError(w, http.StatusBadRequest, "data_blob base64 inválido")
			return
		}
		a, err := repo.AddAttachment(r.Context(), r.PathValue("id"), userID, meta, data)
		if mapSharedVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, attachmentMetaJSON(a))
	}
}

func handleGetVaultAttachment(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		a, err := repo.GetAttachment(r.Context(), r.PathValue("id"), userID, r.PathValue("attId"))
		if mapSharedVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id":         a.ID,
			"meta_blob":  base64.StdEncoding.EncodeToString(a.MetaBlob),
			"data_blob":  base64.StdEncoding.EncodeToString(a.DataBlob),
			"byte_size":  a.ByteSize,
			"created_by": a.CreatedBy,
			"created_at": a.CreatedAt,
		})
	}
}

func handleDeleteVaultAttachment(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteAttachment(r.Context(), r.PathValue("id"), userID, r.PathValue("attId")); mapSharedVaultError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
