// Package httpapi monta o router HTTP da API do AegisPass.
//
// Didático: separamos a construção do router (aqui) do arranque do processo
// (cmd/server). Assim o router pode ser testado isoladamente, sem abrir portas.
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

// Deps são as dependências injetadas no router.
type Deps struct {
	Auth     *auth.Service
	Vault    *vault.Repo
	Hub      *realtime.Hub // nil desactiva WebSocket e notificações push
	Wipe     *security.WipeService
	Users    *users.Repo
	Shifts    *shifts.Repo
	Geofence  *geofence.Repo
	AdminKey string // vazio desactiva POST /api/admin/.../remote-wipe
}

// ctxKey é um tipo privado para chaves de context (evita colisões entre pacotes).
type ctxKey string

const userIDKey ctxKey = "userID"

// NewRouter devolve o http.Handler com todas as rotas registadas.
func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()
	ap := accessPolicyDeps{Shifts: deps.Shifts, Geofence: deps.Geofence}

	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /api/time", handleServerTime)

	if deps.Auth != nil {
		mux.HandleFunc("POST /api/auth/register", handleRegister(deps.Auth))
		mux.HandleFunc("POST /api/auth/login", handleLoginWithAccessPolicy(deps.Auth, ap))
		mux.HandleFunc("GET /api/auth/kdf", handleKDFParams(deps.Auth))
	}
	if deps.Auth != nil && deps.Vault != nil {
		vd := vaultDeps{repo: deps.Vault, hub: deps.Hub}
		mux.Handle("GET /api/vault", requireAuthWithAccessPolicy(deps.Auth, ap, handleListItems(deps.Vault)))
		mux.Handle("POST /api/vault", requireAuthWithAccessPolicy(deps.Auth, ap, handleCreateItem(vd)))
		mux.Handle("GET /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleGetItem(deps.Vault)))
		mux.Handle("PUT /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleUpdateItem(vd)))
		mux.Handle("DELETE /api/vault/{id}", requireAuthWithAccessPolicy(deps.Auth, ap, handleDeleteItem(vd)))
	}
	if deps.Auth != nil && deps.Hub != nil {
		mux.HandleFunc("GET /api/ws/vault", handleVaultWS(deps.Auth, deps.Hub, ap))
	}
	if deps.Shifts != nil && deps.Auth != nil {
		mux.Handle("GET /api/access/shift", requireAuth(deps.Auth, handleGetAccessShift(deps.Shifts)))
	}
	if deps.Geofence != nil && deps.Auth != nil {
		mux.Handle("GET /api/access/geofence", requireAuth(deps.Auth, handleGetAccessGeofence(deps.Geofence)))
	}
	if deps.Shifts != nil && deps.Users != nil && deps.AdminKey != "" {
		mux.HandleFunc(
			"PUT /api/admin/users/{id}/access-shift",
			handleAdminSetAccessShift(deps.AdminKey, deps.Users, deps.Shifts),
		)
	}
	if deps.Geofence != nil && deps.Users != nil && deps.AdminKey != "" {
		mux.HandleFunc(
			"PUT /api/admin/users/{id}/access-geofence",
			handleAdminSetAccessGeofence(deps.AdminKey, deps.Users, deps.Geofence),
		)
	}
	if deps.Wipe != nil && deps.Users != nil && deps.AdminKey != "" {
		mux.HandleFunc(
			"POST /api/admin/users/{id}/remote-wipe",
			handleAdminRemoteWipe(deps.AdminKey, deps.Users, deps.Wipe),
		)
	}
	if deps.Wipe != nil && deps.Auth != nil {
		mux.Handle(
			"POST /api/security/remote-wipe/self",
			requireAuth(deps.Auth, handleSelfRemoteWipe(deps.Wipe)),
		)
	}

	return mux
}

// requireAuth é um middleware: valida o token Bearer e, se válido, injeta o
// userID no context antes de chamar o handler seguinte.
func requireAuth(svc *auth.Service, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "token em falta")
			return
		}
		userID, err := svc.Authenticate(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "sessão inválida")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// bearerToken extrai o token do header "Authorization: Bearer <token>".
func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if strings.HasPrefix(h, prefix) {
		return strings.TrimPrefix(h, prefix)
	}
	return ""
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// writeJSON escreve uma resposta JSON com o status indicado.
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError padroniza as respostas de erro: { "error": "mensagem" }.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
