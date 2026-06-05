package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharekeys"
)

type keypairRequest struct {
	PublicKey         string `json:"public_key"`
	WrappedPrivateKey string `json:"wrapped_private_key"`
	Algorithm         string `json:"algorithm"`
}

// handlePutShareKeypair grava o par de chaves de partilha do utilizador.
// A chave privada chega já cifrada com a Master Key (o servidor nunca a vê em claro).
func handlePutShareKeypair(repo *sharekeys.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req keypairRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		pub, err := base64.StdEncoding.DecodeString(req.PublicKey)
		if err != nil || len(pub) == 0 {
			writeError(w, http.StatusBadRequest, "public_key base64 inválida")
			return
		}
		priv, err := base64.StdEncoding.DecodeString(req.WrappedPrivateKey)
		if err != nil || len(priv) == 0 {
			writeError(w, http.StatusBadRequest, "wrapped_private_key base64 inválida")
			return
		}
		algorithm := strings.TrimSpace(req.Algorithm)
		if algorithm == "" {
			writeError(w, http.StatusBadRequest, "algorithm em falta")
			return
		}
		if err := repo.Upsert(r.Context(), userID, pub, priv, algorithm); err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao gravar par de chaves")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// handleGetShareKeypair devolve o par do dono (inclui a chave privada cifrada,
// que só o próprio consegue abrir com a Master Key).
func handleGetShareKeypair(repo *sharekeys.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		kp, err := repo.GetByUserID(r.Context(), userID)
		if errors.Is(err, sharekeys.ErrNotFound) {
			writeError(w, http.StatusNotFound, "par de chaves não configurado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"public_key":          base64.StdEncoding.EncodeToString(kp.PublicKey),
			"wrapped_private_key": base64.StdEncoding.EncodeToString(kp.WrappedPrivateKey),
			"algorithm":           kp.Algorithm,
		})
	}
}

// handleShareKeypairStatus indica se o utilizador já activou a partilha.
func handleShareKeypairStatus(repo *sharekeys.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		ok, err := repo.HasKeypair(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"configured": ok})
	}
}

// handleGetSharePublicKey devolve a chave pública de um colega (por email), para
// lhe partilhar um segredo. Requer sessão — só utilizadores autenticados podem
// procurar chaves públicas. Nunca devolve a chave privada.
func handleGetSharePublicKey(repo *sharekeys.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("email")))
		if email == "" {
			writeError(w, http.StatusBadRequest, "email em falta")
			return
		}
		pk, err := repo.GetPublicKeyByEmail(r.Context(), email)
		if errors.Is(err, sharekeys.ErrNotFound) {
			// Resposta genérica — não revelar se o email existe nem se tem partilha.
			writeError(w, http.StatusNotFound, "não encontrado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"email":      email,
			"user_id":    pk.UserID,
			"public_key": base64.StdEncoding.EncodeToString(pk.PublicKey),
			"algorithm":  pk.Algorithm,
		})
	}
}
