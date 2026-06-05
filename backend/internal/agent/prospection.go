package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

// LeadDraft é um rascunho de lead pronto para o cliente cifrar (Zero-Knowledge).
type LeadDraft struct {
	MessageID       string `json:"message_id"`
	Email           string `json:"email"`
	Subject         string `json:"subject"`
	Notes           string `json:"notes"`
	SuggestedStage  string `json:"suggested_stage"`
	Source          string `json:"source"`
}

// Prospection orquestra inbox → draft_lead_from_email (AGENT-003).
// Didático: o servidor NÃO grava leads em claro — devolve rascunhos ao browser
// para o utilizador rever e cifrar com a Master Key.
type Prospection struct {
	Runner *Runner
	Inbox  MailInboxLister
}

// Run lê e-mails pendentes e gera rascunhos de lead por mensagem.
func (p *Prospection) Run(ctx context.Context, userID string) ([]LeadDraft, error) {
	if p.Runner == nil || p.Inbox == nil {
		return nil, fmt.Errorf("agent: prospeção não configurada")
	}
	msgs, err := p.Inbox.ListInbox(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	drafts := make([]LeadDraft, 0, len(msgs))
	for _, m := range msgs {
		if IsRecruitmentEmail(m.Subject, m.Body) {
			continue
		}
		inRaw, _ := json.Marshal(map[string]string{
			"from_email": m.FromEmail,
			"subject":    m.Subject,
			"body":       m.Body,
		})
		out, err := p.Runner.Run(ctx, "draft_lead_from_email", Request{
			UserID:  userID,
			AgentID: "prospection",
			Input:   inRaw,
		})
		if err != nil {
			// Ignora mensagens com injection ou input inválido — continua o lote.
			continue
		}
		var parsed struct {
			Lead struct {
				Email           string `json:"email"`
				Subject         string `json:"subject"`
				Notes           string `json:"notes"`
				SuggestedStage  string `json:"suggested_stage"`
				Source          string `json:"source"`
			} `json:"lead"`
		}
		if err := json.Unmarshal(out, &parsed); err != nil {
			continue
		}
		drafts = append(drafts, LeadDraft{
			MessageID:      m.ID,
			Email:          parsed.Lead.Email,
			Subject:        parsed.Lead.Subject,
			Notes:          parsed.Lead.Notes,
			SuggestedStage: parsed.Lead.SuggestedStage,
			Source:         parsed.Lead.Source,
		})
	}
	return drafts, nil
}
