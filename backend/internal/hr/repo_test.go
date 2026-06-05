package hr

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
)

func TestEmployeeRecords_Integration(t *testing.T) {
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

	// Dois gestores de RH: a dona da ficha e um estranho.
	var ownerID, strangerID string
	for _, u := range []struct {
		email string
		dst   *string
	}{
		{"hr-owner@test.local", &ownerID},
		{"hr-stranger@test.local", &strangerID},
	} {
		if err := pool.QueryRow(ctx, `
			INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
			VALUES ($1, '\x00', '\x00', '\x00', 1, 8192, 1)
			RETURNING id::text`, u.email).Scan(u.dst); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email LIKE 'hr-%@test.local'`)
	})

	repo := NewRepo(pool)

	// --- Criar ficha vazia --------------------------------------------------
	rec, err := repo.CreateRecord(ctx, ownerID)
	if err != nil || rec.ID == "" {
		t.Fatalf("CreateRecord: rec=%+v err=%v", rec, err)
	}

	// O estranho não vê a ficha (isolamento por dono).
	if _, _, err := repo.GetRecord(ctx, strangerID, rec.ID); err != ErrNotFound {
		t.Fatalf("GetRecord estranho: esperado ErrNotFound, got %v", err)
	}

	// --- Campo a campo: cada campo tem o seu value_blob e wrapped_key --------
	if _, err := repo.PutField(ctx, ownerID, rec.ID, "full_name", []byte("v-name"), []byte("k-name")); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.PutField(ctx, ownerID, rec.ID, "salary", []byte("v-sal"), []byte("k-sal")); err != nil {
		t.Fatal(err)
	}
	// Campo inválido (sem chave) é recusado.
	if _, err := repo.PutField(ctx, ownerID, rec.ID, "nif", []byte("v"), nil); err != ErrInvalidField {
		t.Fatalf("PutField sem chave: esperado ErrInvalidField, got %v", err)
	}
	// Estranho não pode escrever.
	if _, err := repo.PutField(ctx, strangerID, rec.ID, "x", []byte("v"), []byte("k")); err != ErrNotFound {
		t.Fatalf("PutField estranho: esperado ErrNotFound, got %v", err)
	}

	// --- Ler a ficha com os campos cifrados ---------------------------------
	_, fields, err := repo.GetRecord(ctx, ownerID, rec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 {
		t.Fatalf("esperava 2 campos, got %d", len(fields))
	}
	// Ordenados por nome: full_name, salary.
	if fields[0].FieldName != "full_name" || !bytes.Equal(fields[0].ValueBlob, []byte("v-name")) {
		t.Fatalf("campo 0 inesperado: %+v", fields[0])
	}
	if !bytes.Equal(fields[1].WrappedKey, []byte("k-sal")) {
		t.Fatalf("wrapped_key do salario inesperado: %+v", fields[1])
	}

	// --- Upsert: reescrever um campo substitui valor e chave ----------------
	if _, err := repo.PutField(ctx, ownerID, rec.ID, "salary", []byte("v-sal2"), []byte("k-sal2")); err != nil {
		t.Fatal(err)
	}
	_, fields, _ = repo.GetRecord(ctx, ownerID, rec.ID)
	if len(fields) != 2 || !bytes.Equal(fields[1].ValueBlob, []byte("v-sal2")) {
		t.Fatalf("upsert do salario falhou: %+v", fields)
	}

	// --- Apagar um campo ----------------------------------------------------
	if err := repo.DeleteField(ctx, ownerID, rec.ID, "salary"); err != nil {
		t.Fatal(err)
	}
	if err := repo.DeleteField(ctx, ownerID, rec.ID, "salary"); err != ErrNotFound {
		t.Fatalf("DeleteField repetido: esperado ErrNotFound, got %v", err)
	}
	_, fields, _ = repo.GetRecord(ctx, ownerID, rec.ID)
	if len(fields) != 1 {
		t.Fatalf("após apagar, esperava 1 campo, got %d", len(fields))
	}

	// --- Listagem e remoção da ficha ----------------------------------------
	recs, err := repo.ListRecords(ctx, ownerID)
	if err != nil || len(recs) != 1 {
		t.Fatalf("ListRecords: recs=%+v err=%v", recs, err)
	}
	if err := repo.DeleteRecord(ctx, strangerID, rec.ID); err != ErrNotFound {
		t.Fatalf("DeleteRecord estranho: esperado ErrNotFound, got %v", err)
	}
	if err := repo.DeleteRecord(ctx, ownerID, rec.ID); err != nil {
		t.Fatal(err)
	}
	if recs, _ := repo.ListRecords(ctx, ownerID); len(recs) != 0 {
		t.Fatalf("após apagar ficha, esperava 0, got %d", len(recs))
	}
}
