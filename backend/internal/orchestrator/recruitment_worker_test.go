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

func TestRecruitmentWorker_SuggestsScreening(t *testing.T) {
	store := &memStore{}
	bus := eventbus.New(store, nil, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Run(ctx)

	orc := orchestrator.New(bus, nil, orchestrator.NewRecruitmentWorker(bus))
	orc.Start()

	var wg sync.WaitGroup
	var suggested eventbus.Event
	wg.Add(1)
	bus.Subscribe(eventbus.OrchestratorActionSuggested, func(_ context.Context, ev eventbus.Event) error {
		suggested = ev
		wg.Done()
		return nil
	})

	payload, _ := json.Marshal(map[string]string{
		"inbox_id": "in-1", "subject": "Candidatura Dev", "from_email": "c@test.com",
	})
	_ = bus.Publish(ctx, eventbus.Event{
		Type: eventbus.MailInboxReceived, UserID: "u1", Payload: payload,
	})

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout")
	}
	var body map[string]any
	_ = json.Unmarshal(suggested.Payload, &body)
	if body["action"] != "screen_candidate" {
		t.Fatalf("payload=%v", body)
	}
}
