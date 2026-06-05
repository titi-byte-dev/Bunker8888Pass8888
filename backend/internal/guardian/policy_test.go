package guardian_test

import (
	"context"
	"errors"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/guardian"
)

func TestPolicy_AllowsManualPing(t *testing.T) {
	p := guardian.Policy{}
	err := p.Allows(context.Background(), agent.Request{AgentID: "manual"}, []agent.Permission{agent.PermNone})
	if err != nil {
		t.Fatalf("manual ping: %v", err)
	}
}

func TestPolicy_DeniesUnknownAgent(t *testing.T) {
	p := guardian.Policy{}
	err := p.Allows(context.Background(), agent.Request{AgentID: "evil"}, []agent.Permission{agent.PermNone})
	if !errors.Is(err, agent.ErrPermissionDenied) {
		t.Fatalf("esperava negado, got %v", err)
	}
}

func TestPolicy_ProspectionNeedsCRM(t *testing.T) {
	p := guardian.Policy{}
	err := p.Allows(context.Background(), agent.Request{AgentID: "prospection"}, []agent.Permission{agent.PermCRMWriteLead})
	if err != nil {
		t.Fatalf("prospection crm: %v", err)
	}
}
