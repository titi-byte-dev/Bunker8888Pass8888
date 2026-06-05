package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

type shiftPolicyRequest struct {
	Timezone         string                `json:"timezone"`
	Schedule         shifts.WeeklySchedule `json:"schedule"`
	Enabled          bool                  `json:"enabled"`
	MaxClockSkewSecs int                   `json:"max_clock_skew_seconds"`
}

func handleServerTime(w http.ResponseWriter, _ *http.Request) {
	now := time.Now().UTC()
	writeJSON(w, http.StatusOK, map[string]any{
		"server_time": now.Format(time.RFC3339),
		"unix_ms":     now.UnixMilli(),
	})
}

func handleGetAccessShift(shiftsRepo *shifts.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		p, err := shiftsRepo.Get(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		within, err := shifts.IsWithinShift(time.Now().UTC(), p)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"enabled":                p.Enabled,
			"timezone":               p.Timezone,
			"schedule":               p.Schedule,
			"max_clock_skew_seconds": p.MaxClockSkewSecs,
			"within_shift":           within,
			"server_time":            time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func handleAdminSetAccessShift(userRepo *users.Repo, shiftsRepo *shifts.Repo) http.HandlerFunc {
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

		var req shiftPolicyRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		p := shifts.Policy{
			UserID:           targetID,
			Timezone:         req.Timezone,
			Schedule:         req.Schedule,
			Enabled:          req.Enabled,
			MaxClockSkewSecs: req.MaxClockSkewSecs,
		}
		if err := shiftsRepo.Upsert(r.Context(), p); err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao gravar turno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func assertUserWithinShift(ctx context.Context, shiftsRepo *shifts.Repo, userID string) error {
	p, err := shiftsRepo.Get(ctx, userID)
	if err != nil {
		return err
	}
	return shifts.AssertWithinShift(time.Now().UTC(), p)
}
