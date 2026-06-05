package orchestrator_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/orchestrator"
)

func TestOnboardingWorker_SuggestsAction(t *testing.T) {
	store := &memStore{}
	bus := eventbus.New(store, nil, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Run(ctx)

	orc := orchestrator.New(bus, nil, orchestrator.NewOnboardingWorker(bus))
	orc.Start()

	var wg sync.WaitGroup
	var suggested eventbus.Event
	wg.Add(1)
	bus.Subscribe(eventbus.OrchestratorActionSuggested, func(_ context.Context, ev eventbus.Event) error {
		suggested = ev
		wg.Done()
		return nil
	})

	payload, _ := json.Marshal(map[string]string{"record_id": "rec-1"})
	_ = bus.Publish(ctx, eventbus.Event{
		Type:    eventbus.HREmployeeCreated,
		UserID:  "u1",
		Payload: payload,
	})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	if suggested.Type != eventbus.OrchestratorActionSuggested || suggested.UserID != "u1" {
		t.Fatalf("suggested=%+v", suggested)
	}
	var body map[string]any
	_ = json.Unmarshal(suggested.Payload, &body)
	if body["action"] != "run_onboarding" || body["record_id"] != "rec-1" {
		t.Fatalf("payload=%v", body)
	}
}
