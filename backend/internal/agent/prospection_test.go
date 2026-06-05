package agent_test

import (
	"context"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

func TestProspection_RunDrafts(t *testing.T) {
	inbox := &fakeInbox{msgs: []mail.InboxMessage{{
		ID: "m1", OwnerID: "u1", FromEmail: "maria@empresa.pt",
		Subject: "Pedido de demo", Body: "Gostávamos de agendar uma demonstração.",
		ReceivedAt: time.Now(),
	}}}
	reg := agent.NewRegistry()
	reg.MustRegister(agent.NewDraftLeadTool())
	run := agent.NewRunner(reg, guardian.Policy{})
	svc := &agent.Prospection{Runner: run, Inbox: inbox}

	drafts, err := svc.Run(context.Background(), "u1")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(drafts) != 1 {
		t.Fatalf("drafts=%d", len(drafts))
	}
	if drafts[0].Email != "maria@empresa.pt" || drafts[0].MessageID != "m1" {
		t.Fatalf("draft=%+v", drafts[0])
	}
}
