package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/burnnotes"
)

// Notas auto-destrutivas (SHARE-005). Tal como os secret links, o ciphertext é
// opaco para o servidor (a chave vive no fragmento do URL). A diferenca: cada
// nota lê-se UMA vez (arde a seguir) e o autor recebe um burn_token para a
// destruir antes de ser lida.

type createBurnNoteRequest struct {
	Ciphertext string `json:"ciphertext"` // nota cifrada no cliente (base64)
	TTLSeconds int    `json:"ttl_seconds"`
}

type burnNoteRequest struct {
	BurnToken string `json:"burn_token"`
}

// handleCreateBurnNote guarda uma nota cifrada em RAM e devolve o id, o
// burn_token e a expiracao. Requer sessao — só utilizadores criam notas.
func handleCreateBurnNote(store *burnnotes.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createBurnNoteRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		ct, err := base64.StdEncoding.DecodeString(req.Ciphertext)
		if err != nil || len(ct) == 0 {
			writeError(w, http.StatusBadRequest, "ciphertext base64 inválido")
			return
		}
		id, burnToken, expiresAt, err := store.Create(ct, time.Duration(req.TTLSeconds)*time.Second)
		switch {
		case errors.Is(err, burnnotes.ErrTooLarge):
			writeError(w, http.StatusRequestEntityTooLarge, "nota demasiado grande")
			return
		case errors.Is(err, burnnotes.ErrFull):
			writeError(w, http.StatusServiceUnavailable, "limite de notas temporário atingido")
			return
		case err != nil:
			writeError(w, http.StatusBadRequest, "pedido inválido")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":         id,
			"burn_token": burnToken,
			"expires_at": expiresAt.UTC().Format(time.RFC3339),
		})
	}
}

// handleConsumeBurnNote devolve o ciphertext UMA vez e queima a nota. PÚBLICO: a
// chave de cifra está no fragmento do URL, que nunca chega aqui.
func handleConsumeBurnNote(store *burnnotes.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct, err := store.Consume(r.PathValue("id"))
		if errors.Is(err, burnnotes.ErrNotFound) {
			writeError(w, http.StatusNotFound, "nota inexistente, expirada ou ja destruida")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"ciphertext": base64.StdEncoding.EncodeToString(ct),
		})
	}
}

// handleBurnNote destroi uma nota antes de ser lida, mediante o burn_token. É um
// modelo de CAPACIDADE: quem tem o token pode destruir; o servidor compara-o em
// tempo constante. Público (a posse do token é a autorizacao).
func handleBurnNote(store *burnnotes.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnNoteRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		err := store.Burn(r.PathValue("id"), req.BurnToken)
		switch {
		case errors.Is(err, burnnotes.ErrNotFound):
			writeError(w, http.StatusNotFound, "nota inexistente, expirada ou ja destruida")
			return
		case errors.Is(err, burnnotes.ErrBadToken):
			writeError(w, http.StatusForbidden, "token de destruicao invalido")
			return
		case err != nil:
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
