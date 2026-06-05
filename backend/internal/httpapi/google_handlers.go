package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/googledrive"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/googleworkspace"
)

type driveBlobRequest struct {
	Blob string `json:"blob"`
}

func driveFileJSON(f *googledrive.File) map[string]any {
	out := map[string]any{
		"id":         f.ID,
		"blob":       base64.StdEncoding.EncodeToString(f.Blob),
		"created_at": f.CreatedAt,
		"updated_at": f.UpdatedAt,
	}
	if f.ExternalID != nil {
		out["external_id"] = *f.ExternalID
	}
	return out
}

func decodeDriveBlob(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	var req driveBlobRequest
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

func mapDriveError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, googledrive.ErrNotFound) {
		writeError(w, http.StatusNotFound, "ficheiro não encontrado")
		return true
	}
	writeError(w, http.StatusInternalServerError, "falha na operação Drive")
	return true
}

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

func handleListDriveFiles(repo *googledrive.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		files, err := repo.List(r.Context(), userID)
		if mapDriveError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(files))
		for i := range files {
			out = append(out, driveFileJSON(&files[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"files": out})
	}
}

func handleCreateDriveFile(repo *googledrive.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		blob, ok := decodeDriveBlob(w, r)
		if !ok {
			return
		}
		f, err := repo.Create(r.Context(), userID, blob)
		if mapDriveError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, driveFileJSON(f))
	}
}

func handleDeleteDriveFile(repo *googledrive.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		if err := repo.Delete(r.Context(), userID, id); mapDriveError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
