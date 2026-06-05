package eventbus_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

type approvalMemStore struct {
	records []eventbus.Record
	events  []eventbus.Event
}

func (m *approvalMemStore) GetByID(_ context.Context, userID, id string) (eventbus.Record, error) {
	for _, r := range m.records {
		if r.ID == id && r.UserID == userID {
			return r, nil
		}
	}
	return eventbus.Record{}, eventbus.ErrSuggestionNotFound
}

func (m *approvalMemStore) DecisionMap(_ context.Context, _ string) (map[string]string, error) {
	out := make(map[string]string)
	for _, ev := range m.events {
		var dp struct {
			SuggestionID string `json:"suggestion_id"`
		}
		_ = json.Unmarshal(ev.Payload, &dp)
		if dp.SuggestionID == "" || out[dp.SuggestionID] != "" {
			continue
		}
		switch ev.Type {
		case eventbus.OrchestratorActionApproved:
			out[dp.SuggestionID] = "approved"
		case eventbus.OrchestratorActionRejected:
			out[dp.SuggestionID] = "rejected"
		}
	}
	return out, nil
}

func TestDecide_ApprovePublishesEvent(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{"action": "run_prospection"})
	store := &approvalMemStore{records: []eventbus.Record{{
		ID: "sug-1", UserID: "u1", Type: eventbus.OrchestratorActionSuggested, Payload: payload,
	}}}
	bus := eventbus.New(nil, nil, 8)
	ctx := context.Background()
	bus.Run(ctx)

	var approved eventbus.Event
	done := make(chan struct{})
	bus.Subscribe(eventbus.OrchestratorActionApproved, func(_ context.Context, ev eventbus.Event) error {
		approved = ev
		close(done)
		return nil
	})

	action, _, err := eventbus.Decide(ctx, store, bus, "u1", "sug-1", eventbus.DecisionApprove)
	if err != nil {
		t.Fatalf("decide: %v", err)
	}
	if action != "run_prospection" {
		t.Fatalf("action=%q", action)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	if approved.Type != eventbus.OrchestratorActionApproved {
		t.Fatalf("type=%s", approved.Type)
	}
}

func TestDecide_AlreadyDecided(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{"action": "run_prospection"})
	decPayload, _ := json.Marshal(map[string]string{"suggestion_id": "sug-1"})
	store := &approvalMemStore{
		records: []eventbus.Record{{
			ID: "sug-1", UserID: "u1", Type: eventbus.OrchestratorActionSuggested, Payload: payload,
		}},
		events: []eventbus.Event{{
			Type: eventbus.OrchestratorActionApproved, Payload: decPayload,
		}},
	}
	_, _, err := eventbus.Decide(context.Background(), store, nil, "u1", "sug-1", eventbus.DecisionApprove)
	if err != eventbus.ErrAlreadyDecided {
		t.Fatalf("err=%v", err)
	}
}
