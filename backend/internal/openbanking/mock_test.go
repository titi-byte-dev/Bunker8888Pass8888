package openbanking_test

import (
	"context"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/openbanking"
)

func TestMockProvider_ListTransactions(t *testing.T) {
	txs, err := openbanking.MockProvider{}.ListTransactions(context.Background(), "u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(txs) < 3 {
		t.Fatalf("esperava movimentos mock, obteve %d", len(txs))
	}
	if txs[0].Amount >= 0 {
		t.Fatal("débitos SaaS devem ser negativos")
	}
}
