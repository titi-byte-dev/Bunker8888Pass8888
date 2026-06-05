package sharekeys

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestKeypair_Integration(t *testing.T) {
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

	const email = "sharekeys-it@test.local"
	var userID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
		VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
		RETURNING id::text`, email).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM user_keypairs WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	repo := NewRepo(pool)

	if ok, err := repo.HasKeypair(ctx, userID); err != nil || ok {
		t.Fatalf("HasKeypair inicial: ok=%v err=%v", ok, err)
	}
	if _, err := repo.GetByUserID(ctx, userID); err != ErrNotFound {
		t.Fatalf("GetByUserID inexistente: esperado ErrNotFound, got %v", err)
	}

	pub := []byte("public-spki-der")
	priv := []byte("wrapped-pkcs8")
	const alg = "RSA-OAEP-3072-SHA256"
	if err := repo.Upsert(ctx, userID, pub, priv, alg); err != nil {
		t.Fatal(err)
	}

	kp, err := repo.GetByUserID(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(kp.PublicKey, pub) || !bytes.Equal(kp.WrappedPrivateKey, priv) || kp.Algorithm != alg {
		t.Fatalf("keypair gravado difere: %+v", kp)
	}

	byEmail, err := repo.GetPublicKeyByEmail(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(byEmail.PublicKey, pub) || byEmail.UserID != userID {
		t.Fatalf("GetPublicKeyByEmail difere: %+v", byEmail)
	}

	// Upsert é idempotente / actualiza.
	pub2 := []byte("public-spki-der-rotated")
	if err := repo.Upsert(ctx, userID, pub2, priv, alg); err != nil {
		t.Fatal(err)
	}
	kp2, err := repo.GetByUserID(ctx, userID)
	if err != nil || !bytes.Equal(kp2.PublicKey, pub2) {
		t.Fatalf("Upsert rotativo falhou: %+v err=%v", kp2, err)
	}

	if _, err := repo.GetPublicKeyByEmail(ctx, "ninguem@test.local"); err != ErrNotFound {
		t.Fatalf("GetPublicKeyByEmail inexistente: esperado ErrNotFound, got %v", err)
	}
}
