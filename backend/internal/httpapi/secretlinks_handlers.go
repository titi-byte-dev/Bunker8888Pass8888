package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/secretlinks"
)

type createLinkRequest struct {
	Ciphertext string `json:"ciphertext"` // segredo cifrado no cliente (base64)
	TTLSeconds int    `json:"ttl_seconds"`
	MaxViews   int    `json:"max_views"`
}

// handleCreateSecretLink guarda um segredo cifrado em RAM e devolve o id + a
// hora de expiração. Requer sessão — só utilizadores autenticados criam links.
func handleCreateSecretLink(store *secretlinks.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createLinkRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		ct, err := base64.StdEncoding.DecodeString(req.Ciphertext)
		if err != nil || len(ct) == 0 {
			writeError(w, http.StatusBadRequest, "ciphertext base64 inválido")
			return
		}
		id, expiresAt, err := store.Create(ct, time.Duration(req.TTLSeconds)*time.Second, req.MaxViews)
		switch {
		case errors.Is(err, secretlinks.ErrTooLarge):
			writeError(w, http.StatusRequestEntityTooLarge, "segredo demasiado grande")
			return
		case errors.Is(err, secretlinks.ErrFull):
			writeError(w, http.StatusServiceUnavailable, "limite de links temporário atingido")
			return
		case err != nil:
			writeError(w, http.StatusBadRequest, "pedido inválido")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":         id,
			"expires_at": expiresAt.UTC().Format(time.RFC3339),
		})
	}
}

// handleConsumeSecretLink devolve o ciphertext UMA vez e remove-o da RAM quando
// esgota as visualizações. É PÚBLICO: qualquer pessoa com o link o pode abrir
// (a chave de cifra está no fragmento do URL, que nunca chega aqui).
func handleConsumeSecretLink(store *secretlinks.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ct, err := store.Consume(r.PathValue("id"))
		if errors.Is(err, secretlinks.ErrNotFound) {
			writeError(w, http.StatusNotFound, "link inexistente, expirado ou ja utilizado")
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
