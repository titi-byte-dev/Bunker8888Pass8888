package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// registerRequest é o corpo esperado em POST /api/auth/register.
// Os bytes (auth_hash, salt) viajam codificados em base64 no JSON.
type registerRequest struct {
	Email    string `json:"email"`
	AuthHash string `json:"auth_hash"`
	KDF      struct {
		Salt    string `json:"salt"`
		Time    int    `json:"time"`
		Memory  int    `json:"memory"`
		Threads int    `json:"threads"`
	} `json:"kdf"`
}

func handleRegister(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.Email == "" || req.AuthHash == "" {
			writeError(w, http.StatusBadRequest, "email e auth_hash são obrigatórios")
			return
		}
		authHash, err1 := base64.StdEncoding.DecodeString(req.AuthHash)
		salt, err2 := base64.StdEncoding.DecodeString(req.KDF.Salt)
		if err1 != nil || err2 != nil {
			writeError(w, http.StatusBadRequest, "base64 inválido")
			return
		}

		err := svc.Register(r.Context(), req.Email, authHash, auth.ClientKDF{
			Salt: salt, Time: req.KDF.Time, Memory: req.KDF.Memory, Threads: req.KDF.Threads,
		})
		if errors.Is(err, users.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email já registado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao registar")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"status": "registado"})
	}
}

// loginRequest é o corpo esperado em POST /api/auth/login.
type loginRequest struct {
	Email    string `json:"email"`
	AuthHash string `json:"auth_hash"`
}

func handleLogin(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		authHash, err := base64.StdEncoding.DecodeString(req.AuthHash)
		if err != nil {
			writeError(w, http.StatusBadRequest, "base64 inválido")
			return
		}

		token, err := svc.Login(r.Context(), req.Email, authHash)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			// 401 genérico: não dizemos se foi o email ou o auth hash.
			writeError(w, http.StatusUnauthorized, "credenciais inválidas")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha no login")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}

// handleKDFParams devolve o salt/parâmetros KDF do cliente para um email, para o
// cliente conseguir re-derivar o Auth Hash antes de fazer login.
func handleKDFParams(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			writeError(w, http.StatusBadRequest, "email em falta")
			return
		}
		kdf, err := svc.KDFParamsFor(r.Context(), email)
		if err != nil {
			// Não distinguimos "não existe" para limitar enumeração.
			writeError(w, http.StatusNotFound, "não encontrado")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"salt":    base64.StdEncoding.EncodeToString(kdf.Salt),
			"time":    kdf.Time,
			"memory":  kdf.Memory,
			"threads": kdf.Threads,
		})
	}
}

// decodeJSON lê e valida o corpo JSON do pedido, rejeitando campos desconhecidos.
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
