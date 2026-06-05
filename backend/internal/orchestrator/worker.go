// Package orchestrator coordena agentes via Event Bus (AGENT-005).
//
// Didático: o orquestrador não executa tools directamente — regista *workers*
// que reagem a eventos e publicam *sugestões* ou acções seguras. Acções que
// exigem Master Key ou aprovação humana ficam como sugestão (AGENT-009).
package orchestrator

import (
	"context"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

// Worker é um agente registado que reage a tipos de evento.
type Worker interface {
	// Descriptor expõe metadados para API e documentação.
	Descriptor() Descriptor
	// Handle processa um evento; deve ser idempotente quando possível.
	Handle(ctx context.Context, ev eventbus.Event) error
}

// Descriptor descreve um worker do orquestrador.
type Descriptor struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Handles     []string `json:"handles"`
}
