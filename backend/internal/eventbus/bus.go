package eventbus

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

// Event é uma mensagem publicada no bus.
type Event struct {
	Type      string          `json:"type"`
	UserID    string          `json:"user_id"`
	Source    string          `json:"source,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// Handler processa um evento. Erros são registados; não bloqueiam outros handlers.
type Handler func(ctx context.Context, ev Event) error

// Store persiste eventos (opcional — testes usam nil).
type Store interface {
	Append(ctx context.Context, ev Event) (string, error)
}

// Bus distribui eventos via channel Go para subscritores registados.
type Bus struct {
	ch      chan Event
	subs    map[string][]Handler
	store   Store
	logger  *slog.Logger
	mu      sync.RWMutex
	bufSize int
}

// New cria o bus com buffer e store opcional.
func New(store Store, logger *slog.Logger, bufSize int) *Bus {
	if bufSize <= 0 {
		bufSize = 256
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Bus{
		ch:      make(chan Event, bufSize),
		subs:    make(map[string][]Handler),
		store:   store,
		logger:  logger,
		bufSize: bufSize,
	}
}

// Subscribe regista um handler para um tipo (ou WildcardSubscribe).
func (b *Bus) Subscribe(eventType string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[eventType] = append(b.subs[eventType], h)
}

// Publish envia um evento de forma não bloqueante (descarta se buffer cheio).
func (b *Bus) Publish(ctx context.Context, ev Event) error {
	if b == nil {
		return nil
	}
	if ev.CreatedAt.IsZero() {
		ev.CreatedAt = time.Now().UTC()
	}
	if b.store != nil {
		if _, err := b.store.Append(ctx, ev); err != nil {
			b.logger.Warn("eventbus: persistência falhou", "type", ev.Type, "err", err)
		}
	}
	select {
	case b.ch <- ev:
	default:
		b.logger.Warn("eventbus: buffer cheio, evento descartado", "type", ev.Type)
	}
	return nil
}

// Run inicia a goroutine de dispatch (chamar uma vez no arranque do servidor).
func (b *Bus) Run(ctx context.Context) {
	if b == nil {
		return
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-b.ch:
				b.dispatch(ctx, ev)
			}
		}
	}()
}

func (b *Bus) dispatch(ctx context.Context, ev Event) {
	b.mu.RLock()
	handlers := append([]Handler{}, b.subs[ev.Type]...)
	handlers = append(handlers, b.subs[WildcardSubscribe]...)
	b.mu.RUnlock()
	for _, h := range handlers {
		if err := h(ctx, ev); err != nil {
			b.logger.Warn("eventbus: handler falhou", "type", ev.Type, "err", err)
		}
	}
}

// PublishJSON helper para construir payload.
func PublishJSON(ctx context.Context, b *Bus, eventType, userID, source string, payload any) error {
	if b == nil {
		return nil
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return b.Publish(ctx, Event{
		Type:    eventType,
		UserID:  userID,
		Source:  source,
		Payload: raw,
	})
}
