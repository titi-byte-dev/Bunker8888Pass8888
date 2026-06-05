package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

// CandidateDraft é um rascunho de candidato com triagem às cegas (AGENT-007).
type CandidateDraft struct {
	MessageID string `json:"message_id"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Summary   string `json:"summary"`
	Blind     bool   `json:"blind"`
	Source    string `json:"source"`
}

// Recruitment orquestra inbox → draft_candidate_from_email.
type Recruitment struct {
	Runner *Runner
	Inbox  MailInboxLister
}

// Run processa e-mails de candidatura pendentes.
func (r *Recruitment) Run(ctx context.Context, userID string) ([]CandidateDraft, error) {
	if r.Runner == nil || r.Inbox == nil {
		return nil, fmt.Errorf("agent: recrutamento não configurado")
	}
	msgs, err := r.Inbox.ListInbox(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	drafts := make([]CandidateDraft, 0)
	for _, m := range msgs {
		if !IsRecruitmentEmail(m.Subject, m.Body) {
			continue
		}
		inRaw, _ := json.Marshal(map[string]string{
			"from_email": m.FromEmail,
			"subject":    m.Subject,
			"body":       m.Body,
		})
		out, err := r.Runner.Run(ctx, "draft_candidate_from_email", Request{
			UserID:  userID,
			AgentID: "recruitment",
			Input:   inRaw,
		})
		if err != nil {
			continue
		}
		var parsed struct {
			Candidate struct {
				Email   string `json:"email"`
				Subject string `json:"subject"`
				Summary string `json:"summary"`
				Blind   bool   `json:"blind"`
				Source  string `json:"source"`
			} `json:"candidate"`
		}
		if err := json.Unmarshal(out, &parsed); err != nil {
			continue
		}
		drafts = append(drafts, CandidateDraft{
			MessageID: m.ID,
			Email:     parsed.Candidate.Email,
			Subject:   parsed.Candidate.Subject,
			Summary:   parsed.Candidate.Summary,
			Blind:     parsed.Candidate.Blind,
			Source:    parsed.Candidate.Source,
		})
	}
	return drafts, nil
}
