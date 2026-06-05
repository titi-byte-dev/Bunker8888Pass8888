package fin

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestSaaSSubscriptions_Integration(t *testing.T) {
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
		{"fin-owner@test.local", &ownerID},
		{"fin-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'fin-%@test.local'`)
	})

	repo := NewRepo(pool)

	// Criar duas subscricoes (blobs opacos).
	s1, err := repo.Create(ctx, ownerID, []byte("blob-netflix"))
	if err != nil || s1.ID == "" {
		t.Fatalf("Create: s1=%+v err=%v", s1, err)
	}
	if _, err := repo.Create(ctx, ownerID, []byte("blob-figma")); err != nil {
		t.Fatal(err)
	}

	// Dono ve 2; estranho ve 0.
	if subs, _ := repo.List(ctx, ownerID); len(subs) != 2 {
		t.Fatalf("List dono: esperava 2, got %d", len(subs))
	}
	if subs, _ := repo.List(ctx, strangerID); len(subs) != 0 {
		t.Fatalf("List estranho: esperava 0, got %d", len(subs))
	}

	// Atualizar o blob; estranho nao consegue.
	if _, err := repo.Update(ctx, strangerID, s1.ID, []byte("x")); err != ErrNotFound {
		t.Fatalf("Update estranho: esperado ErrNotFound, got %v", err)
	}
	upd, err := repo.Update(ctx, ownerID, s1.ID, []byte("blob-netflix-v2"))
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := repo.List(ctx, ownerID); !containsBlob(got, []byte("blob-netflix-v2")) {
		t.Fatalf("update nao refletido; upd=%+v", upd)
	}

	// Apagar.
	if err := repo.Delete(ctx, ownerID, s1.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.Delete(ctx, ownerID, s1.ID); err != ErrNotFound {
		t.Fatalf("Delete repetido: esperado ErrNotFound, got %v", err)
	}
}

func containsBlob(subs []Subscription, want []byte) bool {
	for _, s := range subs {
		if bytes.Equal(s.Blob, want) {
			return true
		}
	}
	return false
}
