package agent_test

import (
	"context"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

type recruitInbox struct {
	msgs []mail.InboxMessage
}

func (r recruitInbox) ListInbox(_ context.Context, _ string, unprocessed bool) ([]mail.InboxMessage, error) {
	if !unprocessed {
		return r.msgs, nil
	}
	var out []mail.InboxMessage
	for _, m := range r.msgs {
		if m.ProcessedAt == nil {
			out = append(out, m)
		}
	}
	return out, nil
}

func TestRecruitment_RunDrafts(t *testing.T) {
	reg := agent.NewRegistry()
	reg.MustRegister(agent.NewDraftCandidateTool())
	run := agent.NewRunner(reg, agent.PermissivePolicy{})
	inbox := recruitInbox{msgs: []mail.InboxMessage{{
		ID: "m1", FromEmail: "cand@test.com",
		Subject: "Candidatura Dev",
		Body:    "Género: M\nSkills: Go",
	}}}
	svc := &agent.Recruitment{Runner: run, Inbox: inbox}
	drafts, err := svc.Run(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(drafts) != 1 || !drafts[0].Blind {
		t.Fatalf("drafts=%+v", drafts)
	}
}
