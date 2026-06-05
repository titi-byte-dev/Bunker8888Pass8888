package agent_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

type fakeInbox struct {
	msgs []mail.InboxMessage
}

func (f *fakeInbox) ListInbox(_ context.Context, ownerID string, unprocessedOnly bool) ([]mail.InboxMessage, error) {
	out := make([]mail.InboxMessage, 0)
	for _, m := range f.msgs {
		if m.OwnerID != ownerID {
			continue
		}
		if unprocessedOnly && m.ProcessedAt != nil {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

func TestListMailInboxTool(t *testing.T) {
	inbox := &fakeInbox{msgs: []mail.InboxMessage{{
		ID: "m1", OwnerID: "u1", FromEmail: "lead@empresa.pt",
		Subject: "Orçamento", Body: "Olá, precisamos de proposta.",
		ReceivedAt: time.Now(),
	}}}
	reg := agent.NewRegistry()
	reg.MustRegister(agent.NewListMailInboxTool(inbox))
	run := agent.NewRunner(reg, guardian.Policy{})

	raw, _ := json.Marshal(map[string]bool{"unprocessed_only": true})
	out, err := run.Run(context.Background(), "list_mail_inbox", agent.Request{
		UserID:  "u1",
		AgentID: "prospection",
		Input:   raw,
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	var parsed struct {
		Count    int `json:"count"`
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.Count != 1 || parsed.Messages[0].ID != "m1" {
		t.Fatalf("output inesperado: %+v", parsed)
	}
}
