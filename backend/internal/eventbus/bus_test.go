package eventbus_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

type memStore struct {
	events []eventbus.Event
}

func (m *memStore) Append(_ context.Context, ev eventbus.Event) (string, error) {
	m.events = append(m.events, ev)
	return "id-1", nil
}

func TestBus_PublishSubscribe(t *testing.T) {
	store := &memStore{}
	bus := eventbus.New(store, nil, 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Run(ctx)

	var wg sync.WaitGroup
	var got eventbus.Event
	wg.Add(1)
	bus.Subscribe(eventbus.MailInboxReceived, func(_ context.Context, ev eventbus.Event) error {
		got = ev
		wg.Done()
		return nil
	})

	payload, _ := json.Marshal(map[string]string{"inbox_id": "x"})
	_ = bus.Publish(ctx, eventbus.Event{
		Type:    eventbus.MailInboxReceived,
		UserID:  "u1",
		Source:  "mail.ingest",
		Payload: payload,
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout à espera do handler")
	}
	if got.UserID != "u1" || len(store.events) != 1 {
		t.Fatalf("got=%+v store=%d", got, len(store.events))
	}
}
