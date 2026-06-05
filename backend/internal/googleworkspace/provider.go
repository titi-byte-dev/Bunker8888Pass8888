package googleworkspace

import "context"

// Provider abstrai Drive/Sheets/Gmail (GOOGLE-002–004).
// Em dev usamos MockProvider; em produção ServiceAccountProvider com delegação.
type Provider interface {
	Name() string
	// Ping confirma que as credenciais e scopes estão acessíveis.
	Ping(ctx context.Context) error
}
