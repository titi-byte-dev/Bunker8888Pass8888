package agent

import "context"

// Policy decide se um pedido pode usar as permissões declaradas pela tool.
type Policy interface {
	Allows(ctx context.Context, req Request, perms []Permission) error
}

// PermissivePolicy permite tudo em desenvolvimento; substituir pelo Guardião (AGENT-002).
type PermissivePolicy struct{}

func (PermissivePolicy) Allows(_ context.Context, _ Request, _ []Permission) error {
	return nil
}
