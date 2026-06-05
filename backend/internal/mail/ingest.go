package mail

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrWebhookIgnored = errors.New("mail: mensagem ignorada (sem alias correspondente)")
	ErrInvalidWebhook = errors.New("mail: payload de webhook inválido")
)

// IngestService encaminha e-mail recebido (Mailpit) para a inbox do dono do alias.
type IngestService struct {
	Aliases *Repo
	Inbox   *InboxRepo
	Mailpit *MailpitClient
	Relay   *RelayService // MAIL-004: reencaminha para destination (opcional)
	Limiter *RateLimiter  // MAIL-005: anti-abuso
}

// IngestResult resume o que aconteceu com uma mensagem recebida.
type IngestResult struct {
	MessageID   string `json:"message_id"`
	AliasUsed   string `json:"alias_used"`
	InboxID     string `json:"inbox_id,omitempty"`
	OwnerID     string `json:"owner_id,omitempty"`
	FromEmail   string `json:"from_email,omitempty"`
	Subject     string `json:"subject,omitempty"`
	BodyPreview string `json:"body_preview,omitempty"`
	Status      string `json:"status"`
	Relayed     bool   `json:"relayed,omitempty"`
	RelayTo     string `json:"relay_to,omitempty"`
}

// ProcessMailpitWebhook trata o POST do Mailpit quando chega SMTP para um alias.
func (s *IngestService) ProcessMailpitWebhook(ctx context.Context, raw []byte) (*IngestResult, error) {
	if s == nil || s.Aliases == nil || s.Inbox == nil {
		return nil, errors.New("mail: ingest não configurado")
	}
	summary, err := ParseMailpitSummary(raw)
	if err != nil || summary.ID == "" {
		return nil, ErrInvalidWebhook
	}
	alias, usedAddr, err := s.resolveAliasFromRecipients(ctx, summary.RecipientAddresses())
	if err != nil {
		return nil, err
	}
	if alias == nil {
		return &IngestResult{MessageID: summary.ID, Status: "ignored"}, ErrWebhookIgnored
	}
	body, err := s.fetchBody(ctx, summary)
	if err != nil {
		return nil, err
	}
	return s.ingestAliasMessage(ctx, alias, usedAddr, summary.ID, summary.FromAddress(), summary.Subject, body)
}

func (s *IngestService) fetchBody(ctx context.Context, summary *mailpitSummary) (string, error) {
	if s.Mailpit != nil {
		if body, err := s.Mailpit.FetchBody(ctx, summary.ID); err == nil && body != "" {
			return body, nil
		}
	}
	if snippet := strings.TrimSpace(summary.Snippet); snippet != "" {
		return snippet, nil
	}
	return "(corpo vazio)", nil
}
