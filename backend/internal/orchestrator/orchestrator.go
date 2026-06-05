package orchestrator

import (
	"context"
	"log/slog"
	"sync"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// Orchestrator regista workers e liga-os ao Event Bus.
type Orchestrator struct {
	bus     *eventbus.Bus
	workers []Worker
	logger  *slog.Logger
	started sync.Once
}

// New cria o orquestrador com workers opcionais.
func New(bus *eventbus.Bus, logger *slog.Logger, workers ...Worker) *Orchestrator {
	if logger == nil {
		logger = slog.Default()
	}
	return &Orchestrator{bus: bus, workers: workers, logger: logger}
}

// Start subscreve cada worker aos tipos que declara.
func (o *Orchestrator) Start() {
	if o == nil || o.bus == nil {
		return
	}
	o.started.Do(func() {
		for _, w := range o.workers {
			desc := w.Descriptor()
			worker := w
			for _, evType := range desc.Handles {
				o.bus.Subscribe(evType, func(ctx context.Context, ev eventbus.Event) error {
					if err := worker.Handle(ctx, ev); err != nil {
						o.logger.Warn("orchestrator: worker falhou",
							"worker", desc.ID, "event", ev.Type, "err", err)
					}
					return nil
				})
			}
			o.logger.Info("orchestrator: worker registado", "id", desc.ID, "handles", desc.Handles)
		}
	})
}

// Descriptors devolve metadados de todos os workers.
func (o *Orchestrator) Descriptors() []Descriptor {
	if o == nil {
		return nil
	}
	out := make([]Descriptor, 0, len(o.workers))
	for _, w := range o.workers {
		out = append(out, w.Descriptor())
	}
	return out
}
