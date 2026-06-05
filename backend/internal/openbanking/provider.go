package openbanking

import "context"

// Provider abstrai um banco ou agregador Open Banking (FIN-003).
// Implementações reais usam mTLS + consentimento PSD2; em dev usamos MockProvider.
type Provider interface {
	Name() string
	// ListTransactions devolve movimentos recentes após consentimento activo.
	ListTransactions(ctx context.Context, ownerID string) ([]Transaction, error)
}
