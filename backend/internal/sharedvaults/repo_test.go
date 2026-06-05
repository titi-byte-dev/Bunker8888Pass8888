package sharedvaults

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestSharedVaults_Integration(t *testing.T) {
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

	// Três utilizadores: dona, colega (member) e estranho (não-membro).
	var ownerID, mateID, strangerID string
	for _, u := range []struct {
		email string
		dst   *string
	}{
		{"sv-owner@test.local", &ownerID},
		{"sv-mate@test.local", &mateID},
		{"sv-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'sv-%@test.local'`)
	})

	repo := NewRepo(pool)

	// --- Criar cofre: a dona fica owner com a sua cópia da chave -------------
	v, err := repo.CreateVault(ctx, ownerID, []byte("name-blob"), "", []byte("k-wrapped-owner"))
	if err != nil {
		t.Fatal(err)
	}
	if v.ID == "" || v.OwnerID != ownerID {
		t.Fatalf("cofre criado inválido: %+v", v)
	}

	// A dona vê o cofre; o estranho não.
	if vm, err := repo.GetForUser(ctx, ownerID, v.ID); err != nil || vm.Role != RoleOwner {
		t.Fatalf("owner GetForUser: vm=%+v err=%v", vm, err)
	}
	if _, err := repo.GetForUser(ctx, strangerID, v.ID); err != ErrNotFound {
		t.Fatalf("stranger GetForUser: esperado ErrNotFound, got %v", err)
	}

	// --- Permissões: estranho não pode convidar nem escrever ----------------
	if err := repo.AddMember(ctx, v.ID, strangerID, mateID, RoleMember, []byte("k")); err != ErrNotFound {
		t.Fatalf("AddMember por não-membro: esperado ErrNotFound, got %v", err)
	}
	if err := repo.AddMember(ctx, v.ID, ownerID, mateID, RoleOwner, []byte("k")); err != ErrInvalidRole {
		t.Fatalf("AddMember role=owner: esperado ErrInvalidRole, got %v", err)
	}

	// --- A dona convida o colega como viewer (só leitura) -------------------
	if err := repo.AddMember(ctx, v.ID, ownerID, mateID, RoleViewer, []byte("k-wrapped-mate")); err != nil {
		t.Fatal(err)
	}
	vmMate, err := repo.GetForUser(ctx, mateID, v.ID)
	if err != nil || vmMate.Role != RoleViewer || !bytes.Equal(vmMate.WrappedVaultKey, []byte("k-wrapped-mate")) {
		t.Fatalf("colega após convite: vm=%+v err=%v", vmMate, err)
	}

	// Viewer pode ler itens, mas não escrever.
	if _, err := repo.CreateItem(ctx, v.ID, mateID, "note", []byte("blob")); err != ErrForbidden {
		t.Fatalf("viewer CreateItem: esperado ErrForbidden, got %v", err)
	}
	it, err := repo.CreateItem(ctx, v.ID, ownerID, "note", []byte("blob-cifrado"))
	if err != nil {
		t.Fatal(err)
	}
	items, err := repo.ListItems(ctx, v.ID, mateID)
	if err != nil || len(items) != 1 || items[0].ID != it.ID {
		t.Fatalf("viewer ListItems: items=%+v err=%v", items, err)
	}

	// --- Promover a member: já pode escrever --------------------------------
	if err := repo.AddMember(ctx, v.ID, ownerID, mateID, RoleMember, []byte("k-wrapped-mate")); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateItem(ctx, v.ID, mateID, "note", []byte("do-colega")); err != nil {
		t.Fatalf("member CreateItem deveria passar: %v", err)
	}

	// --- Owner é imutável; revogação imediata do colega ---------------------
	if err := repo.RemoveMember(ctx, v.ID, ownerID, ownerID); err != ErrOwnerImmutable {
		t.Fatalf("remover owner: esperado ErrOwnerImmutable, got %v", err)
	}
	if err := repo.RemoveMember(ctx, v.ID, ownerID, mateID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetForUser(ctx, mateID, v.ID); err != ErrNotFound {
		t.Fatalf("após revogação: esperado ErrNotFound, got %v", err)
	}
	if _, err := repo.ListItems(ctx, v.ID, mateID); err != ErrNotFound {
		t.Fatalf("ex-membro ListItems: esperado ErrNotFound, got %v", err)
	}

	// --- Listar cofres do owner; apagar (só owner) --------------------------
	if vaults, err := repo.ListForUser(ctx, ownerID); err != nil || len(vaults) != 1 {
		t.Fatalf("ListForUser owner: vaults=%d err=%v", len(vaults), err)
	}
	if err := repo.DeleteVault(ctx, v.ID, mateID); err != ErrNotFound {
		t.Fatalf("apagar por não-membro: esperado ErrNotFound, got %v", err)
	}
	if err := repo.DeleteVault(ctx, v.ID, ownerID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetForUser(ctx, ownerID, v.ID); err != ErrNotFound {
		t.Fatalf("após apagar: esperado ErrNotFound, got %v", err)
	}
}
