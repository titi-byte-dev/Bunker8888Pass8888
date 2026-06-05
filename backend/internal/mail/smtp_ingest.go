package mail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// SMTPIngestPayload é o JSON enviado pelo pipe Postfix (MAIL-002 prod).
type SMTPIngestPayload struct {
	MessageID string   `json:"message_id"`
	From      string   `json:"from"`
	To        []string `json:"to"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
}

// ProcessSMTPIngest trata e-mail recebido via Postfix (sem API Mailpit).
func (s *IngestService) ProcessSMTPIngest(ctx context.Context, raw []byte) (*IngestResult, error) {
	if s == nil || s.Aliases == nil || s.Inbox == nil {
		return nil, errors.New("mail: ingest não configurado")
	}
	var p SMTPIngestPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, ErrInvalidWebhook
	}
	if p.MessageID == "" {
		return nil, ErrInvalidWebhook
	}
	alias, usedAddr, err := s.resolveAliasFromRecipients(ctx, p.To)
	if err != nil {
		return nil, err
	}
	if alias == nil {
		return &IngestResult{MessageID: p.MessageID, Status: "ignored"}, ErrWebhookIgnored
	}
	return s.ingestAliasMessage(ctx, alias, usedAddr, p.MessageID, p.From, p.Subject, p.Body)
}

func (s *IngestService) resolveAliasFromRecipients(ctx context.Context, to []string) (*Alias, string, error) {
	for _, addr := range to {
		addr = strings.TrimSpace(strings.ToLower(addr))
		if !strings.HasSuffix(addr, "@"+AliasDomain) {
			continue
		}
		a, err := s.Aliases.GetActiveAliasByAddress(ctx, addr)
		if err == nil {
			return a, addr, nil
		}
	}
	return nil, "", nil
}

func (s *IngestService) ingestAliasMessage(
	ctx context.Context,
	alias *Alias,
	aliasAddr string,
	messageID string,
	from string,
	subject string,
	body string,
) (*IngestResult, error) {
	if s.Limiter != nil {
		if err := s.Limiter.AllowInbound(ctx, alias.OwnerID); err != nil {
			return nil, err
		}
	}
	from = strings.TrimSpace(from)
	if from == "" {
		from = "unknown@unknown"
	}
	subject = strings.TrimSpace(subject)
	if subject == "" {
		subject = "(sem assunto)"
	}
	body = strings.TrimSpace(body)
	if body == "" {
		body = "(corpo vazio)"
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
		MessageID:   messageID,
		AliasUsed:   aliasAddr,
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
	}
	return result, nil
}
