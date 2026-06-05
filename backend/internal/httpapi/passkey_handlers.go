package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/passkeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

type passkeyBeginRequest struct {
	Email string `json:"email"`
}

type passkeyFinishRequest struct {
	SessionID  string          `json:"session_id"`
	Name       string          `json:"name"`
	Credential json.RawMessage `json:"credential"`
}

func handlePasskeyRegisterBegin(svc *passkeys.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		options, sessionID, err := svc.BeginRegistration(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao iniciar registo passkey")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"options":    options,
			"session_id": sessionID,
		})
	}
}

func handlePasskeyRegisterFinish(pk *passkeys.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req passkeyFinishRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.SessionID == "" || req.Name == "" || len(req.Credential) == 0 {
			writeError(w, http.StatusBadRequest, "session_id, name e credential são obrigatórios")
			return
		}
		if err := pk.FinishRegistration(r.Context(), userID, req.SessionID, req.Name, req.Credential); err != nil {
			if errors.Is(err, passkeys.ErrSessionExpired) || errors.Is(err, passkeys.ErrNotFound) {
				writeError(w, http.StatusBadRequest, "sessão WebAuthn inválida ou expirada")
				return
			}
			writeError(w, http.StatusBadRequest, "registo passkey falhou")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"status": "registada"})
	}
}

func handlePasskeyLoginBegin(pk *passkeys.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req passkeyBeginRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.Email == "" {
			writeError(w, http.StatusBadRequest, "email em falta")
			return
		}
		options, sessionID, err := pk.BeginLogin(r.Context(), req.Email)
		if errors.Is(err, users.ErrNotFound) || errors.Is(err, passkeys.ErrNotFound) {
			writeError(w, http.StatusNotFound, "não encontrado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao iniciar login passkey")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"options":    options,
			"session_id": sessionID,
		})
	}
}

func handlePasskeyLoginFinish(authSvc *auth.Service, pk *passkeys.Service, ap accessPolicyDeps, sent *sentinel.Service, userRepo *users.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req passkeyFinishRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.SessionID == "" || len(req.Credential) == 0 {
			writeError(w, http.StatusBadRequest, "session_id e credential são obrigatórios")
			return
		}
		userID, err := pk.FinishLogin(r.Context(), req.SessionID, req.Credential)
		if errors.Is(err, passkeys.ErrSessionExpired) || errors.Is(err, passkeys.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "passkey inválida ou sessão expirada")
			return
		}
		if err != nil {
			writeError(w, http.StatusUnauthorized, "login passkey falhou")
			return
		}
		if err := writeAccessPolicyError(w, assertUserAccessPolicy(r.Context(), r, userID, ap)); err != nil {
			return
		}
		email := ""
		if userRepo != nil {
			if u, e := userRepo.ByID(r.Context(), userID); e == nil {
				email = u.Email
			}
		}
		finishLoginWithSentinel(w, r, userID, email, authSvc, sent)
	}
}

func handlePasskeyList(pk *passkeys.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		list, err := pk.ListMeta(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar passkeys")
			return
		}
		out := make([]map[string]string, 0, len(list))
		for _, c := range list {
			out = append(out, map[string]string{
				"id":         c.ID,
				"name":       c.Name,
				"created_at": c.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"passkeys": out})
	}
}

func registerPasskeyRoutes(mux *http.ServeMux, deps Deps, ap accessPolicyDeps) {
	if deps.Passkeys == nil || deps.Auth == nil {
		return
	}
	mux.Handle("POST /api/auth/passkey/register/begin", requireAuth(deps.Auth, handlePasskeyRegisterBegin(deps.Passkeys)))
	mux.Handle("POST /api/auth/passkey/register/finish", requireAuth(deps.Auth, handlePasskeyRegisterFinish(deps.Passkeys)))
	mux.HandleFunc("POST /api/auth/passkey/login/begin", handlePasskeyLoginBegin(deps.Passkeys))
	mux.HandleFunc("POST /api/auth/passkey/login/finish", handlePasskeyLoginFinish(deps.Auth, deps.Passkeys, ap, deps.Sentinel, deps.Users))
	mux.Handle("GET /api/auth/passkey", requireAuth(deps.Auth, handlePasskeyList(deps.Passkeys)))
}
