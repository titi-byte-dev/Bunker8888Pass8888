package sharedvaults

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestSharedVaultAttachments_Integration(t *testing.T) {
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

	// Dona (owner) + colega que entra como viewer.
	var ownerID, mateID string
	for _, u := range []struct {
		email string
		dst   *string
	}{
		{"att-owner@test.local", &ownerID},
		{"att-mate@test.local", &mateID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'att-%@test.local'`)
	})

	repo := NewRepo(pool)
	v, err := repo.CreateVault(ctx, ownerID, []byte("name-blob"), "", []byte("k-owner"))
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.AddMember(ctx, v.ID, ownerID, mateID, RoleViewer, []byte("k-mate")); err != nil {
		t.Fatal(err)
	}

	meta := []byte("meta-cifrado")
	data := []byte("data-cifrado-do-ficheiro")

	// Viewer não pode carregar anexos.
	if _, err := repo.AddAttachment(ctx, v.ID, mateID, meta, data); err != ErrForbidden {
		t.Fatalf("viewer AddAttachment: esperado ErrForbidden, got %v", err)
	}
	// Não-membro nem sequer vê o cofre.
	if _, err := repo.AddAttachment(ctx, v.ID, "00000000-0000-0000-0000-000000000000", meta, data); err != ErrNotFound {
		t.Fatalf("não-membro AddAttachment: esperado ErrNotFound, got %v", err)
	}
	// Limite de tamanho do ciphertext.
	if _, err := repo.AddAttachment(ctx, v.ID, ownerID, meta, make([]byte, MaxAttachmentBytes+1)); err != ErrTooLarge {
		t.Fatalf("anexo grande: esperado ErrTooLarge, got %v", err)
	}

	// Owner carrega; byte_size = tamanho do data_blob.
	a, err := repo.AddAttachment(ctx, v.ID, ownerID, meta, data)
	if err != nil {
		t.Fatal(err)
	}
	if a.ByteSize != int64(len(data)) {
		t.Fatalf("byte_size: esperado %d, got %d", len(data), a.ByteSize)
	}

	// A listagem traz metadados, mas nunca os bytes do ficheiro.
	list, err := repo.ListAttachments(ctx, v.ID, mateID)
	if err != nil || len(list) != 1 || list[0].ID != a.ID {
		t.Fatalf("ListAttachments viewer: list=%+v err=%v", list, err)
	}

	// O viewer pode descarregar (GET traz os bytes intactos).
	got, err := repo.GetAttachment(ctx, v.ID, mateID, a.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.DataBlob, data) || !bytes.Equal(got.MetaBlob, meta) {
		t.Fatalf("GetAttachment devolveu blobs alterados")
	}

	// Viewer não apaga; owner apaga.
	if err := repo.DeleteAttachment(ctx, v.ID, mateID, a.ID); err != ErrForbidden {
		t.Fatalf("viewer DeleteAttachment: esperado ErrForbidden, got %v", err)
	}
	if err := repo.DeleteAttachment(ctx, v.ID, ownerID, a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetAttachment(ctx, v.ID, ownerID, a.ID); err != ErrNotFound {
		t.Fatalf("após apagar: esperado ErrNotFound, got %v", err)
	}

	// Revogar o colega corta o acesso aos anexos.
	if err := repo.RemoveMember(ctx, v.ID, ownerID, mateID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.ListAttachments(ctx, v.ID, mateID); err != ErrNotFound {
		t.Fatalf("ex-membro ListAttachments: esperado ErrNotFound, got %v", err)
	}
}
