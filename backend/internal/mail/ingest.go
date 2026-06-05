package mail

import (
	"context"
	"errors"
	"fmt"
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
	var alias *Alias
	for _, addr := range summary.RecipientAddresses() {
		if !strings.HasSuffix(addr, "@"+AliasDomain) {
			continue
		}
		a, err := s.Aliases.GetActiveAliasByAddress(ctx, addr)
		if err == nil {
			alias = a
			break
		}
	}
	if alias == nil {
		return &IngestResult{MessageID: summary.ID, Status: "ignored"}, ErrWebhookIgnored
	}
	if s.Limiter != nil {
		if err := s.Limiter.AllowInbound(ctx, alias.OwnerID); err != nil {
			return nil, err
		}
	}
	body, err := s.fetchBody(ctx, summary)
	if err != nil {
		return nil, err
	}
	from := summary.FromAddress()
	if from == "" {
		from = "unknown@unknown"
	}
	subject := strings.TrimSpace(summary.Subject)
	if subject == "" {
		subject = "(sem assunto)"
	}
	msg, err := s.Inbox.CreateInboxMessage(ctx, alias.OwnerID, from, subject, body)
	if err != nil {
		return nil, fmt.Errorf("mail: criar inbox: %w", err)
	}
	preview := body
	if len(preview) > 200 {
		preview = preview[:200]
	}
	result := &IngestResult{
		MessageID:   summary.ID,
		AliasUsed:   alias.AliasAddress,
		InboxID:     msg.ID,
		OwnerID:     alias.OwnerID,
		FromEmail:   from,
		Subject:     subject,
		BodyPreview: preview,
		Status:      "ingested",
	}
	if s.Limiter != nil {
		_ = s.Limiter.Record(ctx, alias.OwnerID, alias.ID, DirInbound, from, alias.AliasAddress)
	}
	if s.Relay != nil {
		relayOK := s.Limiter == nil
		if s.Limiter != nil {
			relayOK = s.Limiter.AllowRelay(ctx, alias.ID) == nil
		}
		if relayOK {
			if relayOut, err := s.Relay.ForwardInbound(ctx, alias, from, subject, body); err == nil && relayOut != nil {
				result.Relayed = true
				result.RelayTo = relayOut.To
				if s.Limiter != nil {
					_ = s.Limiter.Record(ctx, alias.OwnerID, alias.ID, DirOutboundRelay, alias.AliasAddress, relayOut.To)
				}
			}
		}
		// Falha de relay não bloqueia ingestão na inbox — registo local primeiro.
	}
	return result, nil
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
