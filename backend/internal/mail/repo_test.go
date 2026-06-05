package mail

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestEmailAliases_Integration(t *testing.T) {
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
		{"mail-owner@test.local", &ownerID},
		{"mail-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'mail-%@test.local'`)
	})

	repo := NewRepo(pool)

	// Destino invalido e recusado.
	if _, err := repo.CreateAlias(ctx, ownerID, "nao-e-email", "x"); err != ErrInvalidDest {
		t.Fatalf("CreateAlias destino invalido: esperado ErrInvalidDest, got %v", err)
	}

	// Criar dois aliases: enderecos gerados unicos e no dominio esperado.
	a1, err := repo.CreateAlias(ctx, ownerID, "real@gmail.com", "Netflix")
	if err != nil {
		t.Fatal(err)
	}
	a2, err := repo.CreateAlias(ctx, ownerID, "real@gmail.com", "Amazon")
	if err != nil {
		t.Fatal(err)
	}
	if a1.AliasAddress == a2.AliasAddress {
		t.Fatalf("aliases deviam ser unicos: %s", a1.AliasAddress)
	}
	if !strings.HasSuffix(a1.AliasAddress, "@"+AliasDomain) {
		t.Fatalf("alias com dominio inesperado: %s", a1.AliasAddress)
	}
	if !a1.Active {
		t.Fatalf("alias deveria nascer activo")
	}

	// Listagem do dono tem 2; o estranho nao ve nada.
	if as, _ := repo.ListAliases(ctx, ownerID); len(as) != 2 {
		t.Fatalf("ListAliases dono: esperava 2, got %d", len(as))
	}
	if as, _ := repo.ListAliases(ctx, strangerID); len(as) != 0 {
		t.Fatalf("ListAliases estranho: esperava 0, got %d", len(as))
	}

	// Desligar o reencaminhamento.
	if err := repo.SetActive(ctx, ownerID, a1.ID, false); err != nil {
		t.Fatal(err)
	}
	if got, _ := repo.GetAlias(ctx, ownerID, a1.ID); got.Active {
		t.Fatalf("alias deveria estar inactivo")
	}
	// Estranho nao mexe nos aliases alheios.
	if err := repo.SetActive(ctx, strangerID, a1.ID, true); err != ErrNotFound {
		t.Fatalf("SetActive estranho: esperado ErrNotFound, got %v", err)
	}

	// Apagar.
	if err := repo.DeleteAlias(ctx, ownerID, a1.ID); err != nil {
		t.Fatal(err)
	}
	if err := repo.DeleteAlias(ctx, ownerID, a1.ID); err != ErrNotFound {
		t.Fatalf("DeleteAlias repetido: esperado ErrNotFound, got %v", err)
	}
}
