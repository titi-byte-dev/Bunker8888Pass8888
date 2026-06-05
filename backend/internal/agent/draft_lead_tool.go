package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// DraftLeadTool prepara um rascunho de lead a partir de metadados de e-mail.
// AGENT-003 persistirá em CRM; aqui só estruturamos dados já sanitizados.
type DraftLeadTool struct{}

func NewDraftLeadTool() *DraftLeadTool { return &DraftLeadTool{} }

type draftLeadInput struct {
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (DraftLeadTool) Descriptor() Descriptor {
	schema := json.RawMessage(`{
		"type": "object",
		"required": ["from_email", "subject", "body"],
		"properties": {
			"from_email": { "type": "string", "format": "email" },
			"subject": { "type": "string", "maxLength": 500 },
			"body": { "type": "string", "maxLength": 20000 }
		}
	}`)
	return Descriptor{
		Name:        "draft_lead_from_email",
		Description: "Cria rascunho de lead CRM a partir de e-mail (conteúdo tratado como dados externos).",
		InputSchema: schema,
		Permissions: []Permission{PermCRMWriteLead, PermMailReadMeta},
		Sensitive:   true,
	}
}

func looksLikeEmail(s string) bool {
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-1
}

func (DraftLeadTool) Validate(input json.RawMessage) error {
	var in draftLeadInput
	if err := DecodeInput(input, &in); err != nil {
		return err
	}
	in.FromEmail = strings.TrimSpace(in.FromEmail)
	in.Subject = strings.TrimSpace(in.Subject)
	if !looksLikeEmail(in.FromEmail) {
		return fmt.Errorf("%w: from_email inválido", ErrInvalidToolInput)
	}
	if in.Subject == "" || len(in.Subject) > 500 {
		return fmt.Errorf("%w: subject inválido", ErrInvalidToolInput)
	}
	if len(in.Body) == 0 || len(in.Body) > 20000 {
		return fmt.Errorf("%w: body inválido", ErrInvalidToolInput)
	}
	if err := RejectIfLooksLikeInstruction(in.Body); err != nil {
		return err
	}
	if err := RejectIfLooksLikeInstruction(in.Subject); err != nil {
		return err
	}
	return nil
}

func (DraftLeadTool) Execute(_ context.Context, req Request) (json.RawMessage, error) {
	var in draftLeadInput
	_ = DecodeInput(req.Input, &in)
	wrapped := WrapExternalData("email:"+in.FromEmail, in.Body)
	out := map[string]any{
		"status": "draft",
		"lead": map[string]any{
			"email":           in.FromEmail,
			"subject":         in.Subject,
			"notes":           wrapped,
			"suggested_stage": "new",
			"source":          "email",
		},
		"meta": map[string]string{
			"agent_id": req.AgentID,
			"user_id":  req.UserID,
		},
	}
	return json.Marshal(out)
}
