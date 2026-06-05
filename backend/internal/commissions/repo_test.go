package commissions

import (
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestCommissions_Integration(t *testing.T) {
	url := os.Getenv("AEGIS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("AEGIS_TEST_DATABASE_URL não definido")
	}
	ctx := context.Background()
	pool, err := db.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}

	var ownerID, strangerID string
	for _, u := range []struct {
		email string
		dst   *string
	}{
		{"com-owner@test.local", &ownerID},
		{"com-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'com-%@test.local'`)
	})

	// Uma fatura do dono para ligar a comissao.
	var invoiceID string
	if err := pool.QueryRow(ctx, `
		INSERT INTO invoices (owner_id, doc_type, year, seq, number, status, blob)
		VALUES ($1, 'invoice', 2026, 1, 'FT 2026/0001', 'paid', '\x00')
		RETURNING id::text`, ownerID).Scan(&invoiceID); err != nil {
		t.Fatal(err)
	}

	repo := NewRepo(pool)

	// Blob vazio e rejeitado.
	if _, err := repo.Create(ctx, ownerID, "", nil); err == nil {
		t.Fatal("Create blob vazio: esperava erro")
	}

	// Comissao ligada a fatura do dono.
	c1, err := repo.Create(ctx, ownerID, invoiceID, []byte("blob-com-1"))
	if err != nil {
		t.Fatal(err)
	}
	if c1.Status != StatusPending {
		t.Fatalf("estado inicial esperado pending; got %q", c1.Status)
	}
	if c1.InvoiceID != invoiceID {
		t.Fatalf("invoice_id nao guardado: %q", c1.InvoiceID)
	}

	// Comissao sem fatura (avulsa) tambem e valida.
	if _, err := repo.Create(ctx, ownerID, "", []byte("blob-com-2")); err != nil {
		t.Fatal(err)
	}

	// Nao se pode ligar a uma fatura de outro dono.
	if _, err := repo.Create(ctx, strangerID, invoiceID, []byte("x")); err != ErrNotFound {
		t.Fatalf("Create fatura alheia: esperado ErrNotFound, got %v", err)
	}

	// Isolamento por dono.
	if got, _ := repo.List(ctx, ownerID); len(got) != 2 {
		t.Fatalf("List dono: esperava 2, got %d", len(got))
	}
	if got, _ := repo.List(ctx, strangerID); len(got) != 0 {
		t.Fatalf("List estranho: esperava 0, got %d", len(got))
	}

	// Estranho nao muda estado.
	if _, err := repo.UpdateStatus(ctx, strangerID, c1.ID, StatusPaid); err != ErrNotFound {
		t.Fatalf("UpdateStatus estranho: esperado ErrNotFound, got %v", err)
	}
	paid, err := repo.UpdateStatus(ctx, ownerID, c1.ID, StatusPaid)
	if err != nil || paid.Status != StatusPaid {
		t.Fatalf("UpdateStatus: paid=%+v err=%v", paid, err)
	}

	// Estado invalido rejeitado.
	if _, err := repo.UpdateStatus(ctx, ownerID, c1.ID, "deleted"); err != ErrInvalidStatus {
		t.Fatalf("estado invalido: esperado ErrInvalidStatus, got %v", err)
	}
}
