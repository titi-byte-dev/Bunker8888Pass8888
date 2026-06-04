package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/recovery"
)

type recoveryBackupRequest struct {
	Blob string `json:"blob"`
}

func handlePutRecoveryBackup(repo *recovery.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req recoveryBackupRequest
		if err := decodeJSON(r, &req); err != nil || req.Blob == "" {
			writeError(w, http.StatusBadRequest, "blob em falta")
			return
		}
		blob, err := base64.StdEncoding.DecodeString(req.Blob)
		if err != nil {
			writeError(w, http.StatusBadRequest, "base64 inválido")
			return
		}
		if err := repo.Upsert(r.Context(), userID, blob); err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao gravar backup")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleGetRecoveryBackupSelf(repo *recovery.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, err := repo.GetByUserID(r.Context(), userID)
		if errors.Is(err, recovery.ErrNotFound) {
			writeError(w, http.StatusNotFound, "backup não configurado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"blob": base64.StdEncoding.EncodeToString(blob),
		})
	}
}

func handleGetRecoveryBackupByEmail(repo *recovery.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			writeError(w, http.StatusBadRequest, "email em falta")
			return
		}
		blob, err := repo.GetByEmail(r.Context(), email)
		if errors.Is(err, recovery.ErrNotFound) {
			// Resposta genérica — não revelar se o email existe.
			writeError(w, http.StatusNotFound, "não encontrado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"blob": base64.StdEncoding.EncodeToString(blob),
		})
	}
}

func handleRecoveryBackupStatus(repo *recovery.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		ok, err := repo.HasBackup(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"configured": ok})
	}
}
