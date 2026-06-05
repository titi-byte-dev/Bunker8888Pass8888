package openbanking

import (
	"context"
	"time"
)

// MockProvider simula um banco para desenvolvimento (FIN-003a).
// Os valores espelham subscrições SaaS típicas para testar reconciliação.
type MockProvider struct{}

func (MockProvider) Name() string { return "mock" }

func (MockProvider) ListTransactions(_ context.Context, _ string) ([]Transaction, error) {
	now := time.Now().UTC()
	return []Transaction{
		{
			ID: "mock-tx-1", Amount: -12.99, Currency: "EUR",
			BookedAt: now.AddDate(0, 0, -3), Description: "NETFLIX.COM",
			MerchantRef: "netflix",
		},
		{
			ID: "mock-tx-2", Amount: -15.00, Currency: "EUR",
			BookedAt: now.AddDate(0, 0, -5), Description: "SPOTIFY AB",
			MerchantRef: "spotify",
		},
		{
			ID: "mock-tx-3", Amount: -49.00, Currency: "EUR",
			BookedAt: now.AddDate(0, 0, -1), Description: "AWS EMEA",
			MerchantRef: "aws",
		},
		{
			ID: "mock-tx-4", Amount: -8.50, Currency: "EUR",
			BookedAt: now.AddDate(0, 0, -2), Description: "CAFE LOCAL LISBOA",
			MerchantRef: "cafe",
		},
	}, nil
}
