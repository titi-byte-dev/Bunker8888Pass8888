package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// DraftCandidateTool prepara rascunho de candidato com triagem às cegas (AGENT-007).
type DraftCandidateTool struct{}

func NewDraftCandidateTool() *DraftCandidateTool { return &DraftCandidateTool{} }

type draftCandidateInput struct {
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func (DraftCandidateTool) Descriptor() Descriptor {
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
		Name:        "draft_candidate_from_email",
		Description: "Cria rascunho de candidato com triagem às cegas a partir de e-mail.",
		InputSchema: schema,
		Permissions: []Permission{PermHRWriteCandidateDraft, PermMailReadMeta},
		Sensitive:   true,
	}
}

func (DraftCandidateTool) Validate(input json.RawMessage) error {
	var in draftCandidateInput
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
	if !IsRecruitmentEmail(in.Subject, in.Body) {
		return fmt.Errorf("%w: não parece e-mail de candidatura", ErrInvalidToolInput)
	}
	if err := RejectIfLooksLikeInstruction(in.Body); err != nil {
		return err
	}
	return nil
}

func (DraftCandidateTool) Execute(_ context.Context, req Request) (json.RawMessage, error) {
	var in draftCandidateInput
	_ = DecodeInput(req.Input, &in)
	blind := BlindScreenCV(in.Body)
	wrapped := WrapExternalData("email:"+in.FromEmail, blind)
	out := map[string]any{
		"status": "draft",
		"candidate": map[string]any{
			"email":    in.FromEmail,
			"subject":  in.Subject,
			"summary":  wrapped,
			"blind":    true,
			"source":   "email",
		},
		"meta": map[string]string{
			"agent_id": req.AgentID,
			"user_id":  req.UserID,
		},
	}
	return json.Marshal(out)
}
