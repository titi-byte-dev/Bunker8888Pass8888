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
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

// Deps são as dependências injetadas no router. Quando Auth/Vault são nil (ex:
// sem base de dados configurada), só o /healthz fica disponível.
type Deps struct {
	Auth  *auth.Service
	Vault *vault.Repo
}

// ctxKey é um tipo privado para chaves de context (evita colisões entre pacotes).
type ctxKey string

const userIDKey ctxKey = "userID"

// NewRouter devolve o http.Handler com todas as rotas registadas.
func NewRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handleHealth)

	if deps.Auth != nil {
		mux.HandleFunc("POST /api/auth/register", handleRegister(deps.Auth))
		mux.HandleFunc("POST /api/auth/login", handleLogin(deps.Auth))
		mux.HandleFunc("GET /api/auth/kdf", handleKDFParams(deps.Auth))
	}
	if deps.Auth != nil && deps.Vault != nil {
		// Rotas protegidas: passam pelo middleware de sessão.
		mux.Handle("GET /api/vault", requireAuth(deps.Auth, handleListItems(deps.Vault)))
		mux.Handle("POST /api/vault", requireAuth(deps.Auth, handleCreateItem(deps.Vault)))
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
