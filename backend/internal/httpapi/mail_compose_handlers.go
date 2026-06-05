package httpapi

import (
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

type composeMailRequest struct {
	AliasID string `json:"alias_id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func handleComposeMail(aliases *mail.Repo, relay *mail.RelayService, limiter *mail.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req composeMailRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.AliasID == "" || req.To == "" {
			writeError(w, http.StatusBadRequest, "alias_id e to são obrigatórios")
			return
		}
		alias, err := aliases.GetAlias(r.Context(), userID, req.AliasID)
		if mapMailError(w, err) {
			return
		}
		if !alias.Active {
			writeError(w, http.StatusBadRequest, "alias desligado")
			return
		}
		if limiter != nil {
			if err := limiter.AllowCompose(r.Context(), userID); err != nil {
				mapMailRateError(w, err)
				return
			}
		}
		if relay == nil {
			writeError(w, http.StatusServiceUnavailable, "relay SMTP não configurado")
			return
		}
		out, err := relay.SendFromAlias(r.Context(), alias, req.To, req.Subject, req.Body)
		if mapMailError(w, err) {
			return
		}
		if limiter != nil {
			_ = limiter.Record(r.Context(), userID, alias.ID, mail.DirCompose, alias.AliasAddress, out.To)
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"from": out.From,
			"to":   out.To,
		})
	}
}

func mapMailRateError(w http.ResponseWriter, err error) bool {
	if errors.Is(err, mail.ErrRateLimited) {
		writeError(w, http.StatusTooManyRequests, "limite de e-mail excedido — tenta mais tarde")
		return true
	}
	return false
}
