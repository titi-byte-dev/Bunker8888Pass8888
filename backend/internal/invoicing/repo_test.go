package invoicing

import (
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestInvoicing_Integration(t *testing.T) {
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
		{"inv-owner@test.local", &ownerID},
		{"inv-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'inv-%@test.local'`)
	})

	repo := NewRepo(pool)

	// Tipo invalido e rejeitado.
	if _, err := repo.Issue(ctx, ownerID, "xpto", "", []byte("b")); err != ErrInvalidType {
		t.Fatalf("Issue tipo invalido: esperado ErrInvalidType, got %v", err)
	}

	// Emite duas faturas: numeracao sequencial sem buracos.
	d1, err := repo.Issue(ctx, ownerID, DocInvoice, "", []byte("blob-ft-1"))
	if err != nil {
		t.Fatal(err)
	}
	d2, err := repo.Issue(ctx, ownerID, DocInvoice, "lead-xyz", []byte("blob-ft-2"))
	if err != nil {
		t.Fatal(err)
	}
	if d1.Seq != 1 || d2.Seq != 2 {
		t.Fatalf("seq esperado 1,2; got %d,%d", d1.Seq, d2.Seq)
	}
	if d2.SourceLeadID != "lead-xyz" {
		t.Fatalf("source_lead_id nao guardado: %q", d2.SourceLeadID)
	}

	// Series independentes por tipo: pro-forma comeca tambem em 1.
	pf, err := repo.Issue(ctx, ownerID, DocProforma, "", []byte("blob-pf-1"))
	if err != nil {
		t.Fatal(err)
	}
	if pf.Seq != 1 {
		t.Fatalf("proforma seq esperado 1; got %d", pf.Seq)
	}
	if pf.Number[:2] != "PF" || d1.Number[:2] != "FT" {
		t.Fatalf("prefixos errados: pf=%q ft=%q", pf.Number, d1.Number)
	}

	// Isolamento por dono.
	if got, _ := repo.List(ctx, ownerID); len(got) != 3 {
		t.Fatalf("List dono: esperava 3, got %d", len(got))
	}
	if got, _ := repo.List(ctx, strangerID); len(got) != 0 {
		t.Fatalf("List estranho: esperava 0, got %d", len(got))
	}

	// Estranho nao muda estado.
	if _, err := repo.UpdateStatus(ctx, strangerID, d1.ID, StatusPaid); err != ErrNotFound {
		t.Fatalf("UpdateStatus estranho: esperado ErrNotFound, got %v", err)
	}
	paid, err := repo.UpdateStatus(ctx, ownerID, d1.ID, StatusPaid)
	if err != nil || paid.Status != StatusPaid {
		t.Fatalf("UpdateStatus: paid=%+v err=%v", paid, err)
	}

	// Estado invalido rejeitado.
	if _, err := repo.UpdateStatus(ctx, ownerID, d1.ID, "deleted"); err != ErrInvalidStatus {
		t.Fatalf("estado invalido: esperado ErrInvalidStatus, got %v", err)
	}
}
