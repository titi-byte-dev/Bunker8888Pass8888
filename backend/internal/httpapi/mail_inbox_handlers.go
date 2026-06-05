package httpapi

import (
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

type inboxCreateRequest struct {
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func mapInboxError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, mail.ErrInboxNotFound):
		writeError(w, http.StatusNotFound, "mensagem não encontrada")
	case errors.Is(err, mail.ErrInvalidDest):
		writeError(w, http.StatusBadRequest, "from_email inválido")
	default:
		writeError(w, http.StatusInternalServerError, "falha na caixa de entrada")
	}
	return true
}

func inboxMessageJSON(m *mail.InboxMessage) map[string]any {
	out := map[string]any{
		"id":          m.ID,
		"from_email":  m.FromEmail,
		"subject":     m.Subject,
		"body":        m.Body,
		"received_at": m.ReceivedAt,
		"created_at":  m.CreatedAt,
	}
	if m.ProcessedAt != nil {
		out["processed_at"] = *m.ProcessedAt
	}
	return out
}

func handleListInbox(inbox *mail.InboxRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		unprocessed := r.URL.Query().Get("unprocessed") == "1"
		msgs, err := inbox.ListInbox(r.Context(), userID, unprocessed)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar inbox")
			return
		}
		items := make([]map[string]any, 0, len(msgs))
		for i := range msgs {
			items = append(items, inboxMessageJSON(&msgs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"messages": items})
	}
}

func handleCreateInboxMessage(inbox *mail.InboxRepo, bus *eventbus.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req inboxCreateRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		m, err := inbox.CreateInboxMessage(r.Context(), userID, req.FromEmail, req.Subject, req.Body)
		if mapInboxError(w, err) {
			return
		}
		if bus != nil {
			preview := req.Body
			if len(preview) > 200 {
				preview = preview[:200]
			}
			_ = eventbus.PublishJSON(r.Context(), bus, eventbus.MailInboxReceived, userID, "mail.inbox.simulate", map[string]any{
				"inbox_id":     m.ID,
				"from_email":   m.FromEmail,
				"subject":      m.Subject,
				"body_preview": preview,
			})
		}
		writeJSON(w, http.StatusCreated, inboxMessageJSON(m))
	}
}

func handleMarkInboxProcessed(inbox *mail.InboxRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		if err := inbox.MarkProcessed(r.Context(), userID, id); mapInboxError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "processed"})
	}
}
