package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/passkeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// loginContextFromRequest extrai IP e GPS do pedido HTTP.
func loginContextFromRequest(r *http.Request, userID, email string) sentinel.LoginContext {
	lc := sentinel.LoginContext{
		UserID:   userID,
		Email:    email,
		ClientIP: clientIP(r),
		At:       time.Now().UTC(),
	}
	geo, err := geofence.ParseClientGeo(r.Header.Get("X-Geo-Latitude"), r.Header.Get("X-Geo-Longitude"))
	if err == nil && geo.Ok {
		lc.GeoLat = &geo.Lat
		lc.GeoLon = &geo.Lon
	}
	return lc
}

// finishLoginWithSentinel avalia Sentinel e emite token ou pede step-up passkey.
func finishLoginWithSentinel(
	w http.ResponseWriter,
	r *http.Request,
	userID, email string,
	authSvc *auth.Service,
	sent *sentinel.Service,
) {
	if sent == nil {
		token, err := authSvc.CreateSessionForUser(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha no login")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"token": token})
		return
	}

	lc := loginContextFromRequest(r, userID, email)
	assessment, err := sent.Evaluate(r.Context(), lc)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha na avaliação sentinel")
		return
	}

	if assessment.Suspicious {
		step, err := sent.CreateStepUpChallenge(r.Context(), lc, assessment)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao criar desafio")
			return
		}
		writeJSON(w, http.StatusForbidden, map[string]any{
			"code":         "sentinel_step_up",
			"error":        "verificação adicional necessária (Sentinel Mode)",
			"challenge_id": step.ChallengeID,
			"reason":       step.Reason,
			"detail":       step.Detail,
		})
		return
	}

	token, err := authSvc.CreateSessionForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "falha no login")
		return
	}
	_ = sent.CompleteSuccess(r.Context(), lc, assessment, false)
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func handleSentinelStepUpBegin(pk *passkeys.Service, sent *sentinel.Service, userRepo *users.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ChallengeID string `json:"challenge_id"`
		}
		if err := decodeJSON(r, &req); err != nil || req.ChallengeID == "" {
		 writeError(w, http.StatusBadRequest, "challenge_id em falta")
			return
		}
		ch, err := sent.GetChallenge(r.Context(), req.ChallengeID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "desafio inválido ou expirado")
			return
		}
		u, err := userRepo.ByID(r.Context(), ch.UserID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "utilizador inválido")
			return
		}
		options, sessionID, err := pk.BeginLogin(r.Context(), u.Email)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao iniciar passkey")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"options":      options,
			"session_id":   sessionID,
			"challenge_id": ch.ID,
			"email":        u.Email,
		})
	}
}

func handleSentinelStepUpFinish(
	authSvc *auth.Service,
	pk *passkeys.Service,
	sent *sentinel.Service,
	userRepo *users.Repo,
	ap accessPolicyDeps,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ChallengeID string          `json:"challenge_id"`
			SessionID   string          `json:"session_id"`
			Credential  json.RawMessage `json:"credential"`
		}
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.ChallengeID == "" || req.SessionID == "" || len(req.Credential) == 0 {
			writeError(w, http.StatusBadRequest, "campos em falta")
			return
		}

		ch, err := sent.GetChallenge(r.Context(), req.ChallengeID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "desafio inválido ou expirado")
			return
		}

		userID, err := pk.FinishLogin(r.Context(), req.SessionID, req.Credential)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "passkey inválida")
			return
		}
		if userID != ch.UserID {
			writeError(w, http.StatusForbidden, "passkey não corresponde ao desafio")
			return
		}
		if err := writeAccessPolicyError(w, assertUserAccessPolicy(r.Context(), r, userID, ap)); err != nil {
			return
		}
		if err := sent.VerifyChallenge(r.Context(), req.ChallengeID, userID); err != nil {
			writeError(w, http.StatusBadRequest, "desafio inválido")
			return
		}

		token, err := authSvc.CreateSessionForUser(r.Context(), userID)
		if err != nil {
		 writeError(w, http.StatusInternalServerError, "falha ao criar sessão")
			return
		}

		email := ""
		if u, err := userRepo.ByID(r.Context(), userID); err == nil {
			email = u.Email
		}
		lc := loginContextFromRequest(r, userID, email)
		_ = sent.CompleteSuccess(r.Context(), lc, sentinel.Assessment{Reason: ch.Reason, Suspicious: true}, true)

		writeJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}

func handleListSentinelEvents(sent *sentinel.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		events, alerts, err := sent.ListEvents(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar eventos")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"events":              events,
			"suspicious_last_24h": alerts,
		})
	}
}
