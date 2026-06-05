package agent

import (
	"context"
	"encoding/json"
	"fmt"
)

// Runner executa tools com validação + política de permissões.
type Runner struct {
	Registry *Registry
	Policy   Policy
}

// NewRunner cria o executor com registo e política.
func NewRunner(reg *Registry, policy Policy) *Runner {
	if reg == nil {
		reg = NewRegistry()
	}
	if policy == nil {
		policy = PermissivePolicy{}
	}
	return &Runner{Registry: reg, Policy: policy}
}

// Run valida permissões, valida input e executa a tool.
func (run *Runner) Run(ctx context.Context, name string, req Request) (json.RawMessage, error) {
	t, ok := run.Registry.Get(name)
	if !ok {
		return nil, ErrToolNotFound
	}
	desc := t.Descriptor()
	if err := run.Policy.Allows(ctx, req, desc.Permissions); err != nil {
		return nil, err
	}
	if err := t.Validate(req.Input); err != nil {
		return nil, err
	}
	out, err := t.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("agent: execução de %s: %w", name, err)
	}
	return out, nil
}
