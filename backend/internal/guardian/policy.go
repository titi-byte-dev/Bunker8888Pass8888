// Package guardian implementa o Guardião (AGENT-002): autorização por agente
// e auditoria — o servidor nunca decifra; nega por omissão (zero-trust).
package guardian

import (
	"context"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
)

// grants mapeia agent_id → permissões permitidas.
var grants = map[string][]agent.Permission{
	// Utilizador humano via API (testes, palette futura).
	"manual": {
		agent.PermNone,
		agent.PermCRMWriteLead,
		agent.PermMailReadMeta,
	},
	// Agente de prospeção (AGENT-003): só metadados de mail + rascunho CRM.
	"prospection": {
		agent.PermMailReadMeta,
		agent.PermCRMWriteLead,
	},
}

// Policy aplica menor privilégio por agent_id.
type Policy struct{}

func (Policy) Allows(_ context.Context, req agent.Request, perms []agent.Permission) error {
	allowed, ok := grants[req.AgentID]
	if !ok {
		return agent.ErrPermissionDenied
	}
	allowedSet := make(map[agent.Permission]struct{}, len(allowed))
	for _, p := range allowed {
		allowedSet[p] = struct{}{}
	}
	for _, need := range perms {
		if need == agent.PermNone {
			continue
		}
		if _, ok := allowedSet[need]; !ok {
			return agent.ErrPermissionDenied
		}
	}
	return nil
}
