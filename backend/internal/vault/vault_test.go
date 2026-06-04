package vault

import (
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestCRUD_Integration(t *testing.T) {
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

	var userID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
		VALUES ('vault-crud@test.local', '\x00', '\x00', '\x00', 1, 8192, 1)
		RETURNING id::text`).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM vault_items WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	repo := NewRepo(pool)
	created, err := repo.Create(ctx, userID, TypeLogin, []byte("blob-cifrado"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, userID, created.ID); err != nil {
		t.Fatal(err)
	}
	updated, err := repo.Update(ctx, userID, created.ID, TypeNote, []byte("novo-blob"))
	if err != nil || updated.Type != TypeNote {
		t.Fatalf("Update: %v type=%s", err, updated.Type)
	}
	if err := repo.Delete(ctx, userID, created.ID); err != nil {
		t.Fatal(err)
	}
}
