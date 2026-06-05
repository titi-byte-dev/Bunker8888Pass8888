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

func TestFinanceWorker_SuggestsReview(t *testing.T) {
	store := &memStore{}
	bus := eventbus.New(store, nil, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Run(ctx)

	orc := orchestrator.New(bus, nil, orchestrator.NewFinanceWorker(bus))
	orc.Start()

	var wg sync.WaitGroup
	var suggested eventbus.Event
	wg.Add(1)
	bus.Subscribe(eventbus.OrchestratorActionSuggested, func(_ context.Context, ev eventbus.Event) error {
		suggested = ev
		wg.Done()
		return nil
	})

	payload, _ := json.Marshal(map[string]any{
		"alert_count": 2, "subscription_ids": []string{"s1", "s2"}, "monthly_saving": 49.0,
	})
	_ = bus.Publish(ctx, eventbus.Event{
		Type: eventbus.FinSubscriptionStale, UserID: "u1", Payload: payload,
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
	if body["action"] != "review_saas_licenses" {
		t.Fatalf("payload=%v", body)
	}
}

func TestFinanceWorker_SuggestsReconcile(t *testing.T) {
	store := &memStore{}
	bus := eventbus.New(store, nil, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Run(ctx)

	orc := orchestrator.New(bus, nil, orchestrator.NewFinanceWorker(bus))
	orc.Start()

	var wg sync.WaitGroup
	var suggested eventbus.Event
	wg.Add(1)
	bus.Subscribe(eventbus.OrchestratorActionSuggested, func(_ context.Context, ev eventbus.Event) error {
		suggested = ev
		wg.Done()
		return nil
	})

	payload, _ := json.Marshal(map[string]any{
		"transaction_count": 4, "matched_count": 2, "unmatched_count": 2,
	})
	_ = bus.Publish(ctx, eventbus.Event{
		Type: eventbus.FinTransactionsSynced, UserID: "u1", Payload: payload,
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
	if body["action"] != "reconcile_payments" {
		t.Fatalf("payload=%v", body)
	}
}
