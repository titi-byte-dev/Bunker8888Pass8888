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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/emergency"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/passkeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/recovery"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharekeys"
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
	Recovery  *recovery.Repo
	Devices   *clidevices.Repo
	CLIca     *clidevices.CA
	CLICertTTL time.Duration
	Passkeys   *passkeys.Service
	Emergency  *emergency.Service
	Sentinel   *sentinel.Service
	ShareKeys  *sharekeys.Repo // nil desactiva endpoints /api/share/*
	AdminKey string // vazio desactiva endpoints /api/admin/*
	Pool     *pgxpool.Pool
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
		mux.HandleFunc("POST /api/auth/login", handleLoginWithAccessPolicy(deps.Auth, ap, deps.Sentinel, deps.Users))
		mux.HandleFunc("GET /api/auth/kdf", handleKDFParams(deps.Auth))
	}
	if deps.Auth != nil && deps.Users != nil {
		mux.Handle("GET /api/auth/session", requireAuth(deps.Auth, handleAuthSession(deps.Users)))
		mux.Handle("GET /api/auth/sessions", requireAuth(deps.Auth, handleListSessions(deps.Auth)))
		mux.Handle("DELETE /api/auth/sessions/{id}", requireAuth(deps.Auth, handleRevokeSession(deps.Auth)))
		mux.Handle("POST /api/auth/sessions/revoke-others", requireAuth(deps.Auth, handleRevokeOtherSessions(deps.Auth)))
		mux.Handle("POST /api/auth/logout", requireAuth(deps.Auth, handleLogout(deps.Auth)))
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
	registerAdminRoutes(mux, deps)
	if deps.Wipe != nil && deps.Auth != nil {
		mux.Handle(
			"POST /api/security/remote-wipe/self",
			requireAuth(deps.Auth, handleSelfRemoteWipe(deps.Wipe)),
		)
	}
	if deps.Recovery != nil && deps.Auth != nil {
		mux.Handle("PUT /api/vault/recovery-backup", requireAuth(deps.Auth, handlePutRecoveryBackup(deps.Recovery)))
		mux.Handle("GET /api/vault/recovery-backup", requireAuth(deps.Auth, handleGetRecoveryBackupSelf(deps.Recovery)))
		mux.Handle("GET /api/vault/recovery-backup/status", requireAuth(deps.Auth, handleRecoveryBackupStatus(deps.Recovery)))
	}
	if deps.Recovery != nil {
		mux.HandleFunc("GET /api/vault/recovery-backup/lookup", handleGetRecoveryBackupByEmail(deps.Recovery))
	}
	if deps.Emergency != nil && deps.Auth != nil {
		mux.Handle("PUT /api/emergency/config", requireAuth(deps.Auth, handlePutEmergencyConfig(deps.Emergency)))
		mux.Handle("GET /api/emergency/config", requireAuth(deps.Auth, handleGetEmergencyConfig(deps.Emergency)))
		mux.Handle("DELETE /api/emergency/config", requireAuth(deps.Auth, handleDeleteEmergencyConfig(deps.Emergency)))
		mux.Handle("GET /api/emergency/requests", requireAuth(deps.Auth, handleListEmergencyRequests(deps.Emergency)))
		mux.Handle("POST /api/emergency/requests/{id}/reject", requireAuth(deps.Auth, handleRejectEmergencyRequest(deps.Emergency)))
		mux.Handle("POST /api/emergency/requests/{id}/approve", requireAuth(deps.Auth, handleApproveEmergencyRequest(deps.Emergency)))
		mux.Handle("POST /api/emergency/request", requireAuth(deps.Auth, handleCreateEmergencyRequest(deps.Emergency)))
		mux.Handle("GET /api/emergency/request/status", requireAuth(deps.Auth, handleGetEmergencyRequestStatus(deps.Emergency)))
		mux.Handle("GET /api/emergency/access", requireAuth(deps.Auth, handleFetchEmergencyAccess(deps.Emergency)))
	}
	registerCLIRoutes(mux, deps)
	registerPasskeyRoutes(mux, deps, ap)
	if deps.Sentinel != nil && deps.Auth != nil {
		mux.Handle("GET /api/security/sentinel/events", requireAuth(deps.Auth, handleListSentinelEvents(deps.Sentinel)))
	}
	if deps.Sentinel != nil && deps.Passkeys != nil && deps.Users != nil && deps.Auth != nil {
		mux.HandleFunc("POST /api/auth/sentinel/step-up/begin", handleSentinelStepUpBegin(deps.Passkeys, deps.Sentinel, deps.Users))
		mux.HandleFunc("POST /api/auth/sentinel/step-up/finish", handleSentinelStepUpFinish(deps.Auth, deps.Passkeys, deps.Sentinel, deps.Users, ap))
	}
	if deps.ShareKeys != nil && deps.Auth != nil {
		mux.Handle("PUT /api/share/keypair", requireAuth(deps.Auth, handlePutShareKeypair(deps.ShareKeys)))
		mux.Handle("GET /api/share/keypair", requireAuth(deps.Auth, handleGetShareKeypair(deps.ShareKeys)))
		mux.Handle("GET /api/share/keypair/status", requireAuth(deps.Auth, handleShareKeypairStatus(deps.ShareKeys)))
		mux.Handle("GET /api/share/public-key", requireAuth(deps.Auth, handleGetSharePublicKey(deps.ShareKeys)))
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
