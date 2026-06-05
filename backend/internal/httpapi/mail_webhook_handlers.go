package httpapi

import (
	"errors"
	"io"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

func handleMailpitWebhook(svc *mail.IngestService, secret string, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if secret == "" || r.URL.Query().Get("secret") != secret {
			writeError(w, http.StatusUnauthorized, "webhook não autorizado")
			return
		}
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
		if err != nil {
			writeError(w, http.StatusBadRequest, "corpo inválido")
			return
		}
		result, err := svc.ProcessMailpitWebhook(r.Context(), raw)
		if errors.Is(err, mail.ErrWebhookIgnored) {
			writeJSON(w, http.StatusOK, map[string]any{"status": "ignored", "result": result})
			return
		}
		if errors.Is(err, mail.ErrInvalidWebhook) {
			writeError(w, http.StatusBadRequest, "payload inválido")
			return
		}
		if mapMailRateError(w, err) {
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao ingerir e-mail")
			return
		}
		if result != nil && result.OwnerID != "" {
			_ = eventbus.PublishJSON(r.Context(), bus, eventbus.MailInboxReceived, result.OwnerID, "mail.ingest", map[string]any{
				"inbox_id":     result.InboxID,
				"alias":        result.AliasUsed,
				"from_email":   result.FromEmail,
				"subject":      result.Subject,
				"body_preview": result.BodyPreview,
				"relayed":      result.Relayed,
				"relay_to":     result.RelayTo,
				"mailpit_id":   result.MessageID,
			})
		}
		writeJSON(w, http.StatusCreated, map[string]any{"status": "ingested", "result": result})
	}
}
