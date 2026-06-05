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

	// --- Crypto-shredding (HR-003) + certificado (HR-004) -------------------
	// Repor o salario para o eliminar a seguir via shred.
	if _, err := repo.PutField(ctx, ownerID, rec.ID, "salary", []byte("v-sal3"), []byte("k-sal3")); err != nil {
		t.Fatal(err)
	}
	cert, err := repo.ShredField(ctx, ownerID, rec.ID, "salary")
	if err != nil {
		t.Fatal(err)
	}
	if cert.CertHash == "" || cert.ValueDigest == "" || cert.RecordID != rec.ID {
		t.Fatalf("certificado inválido: %+v", cert)
	}
	// O cert_hash tem de bater certo com o canónico reconstruído.
	want := sha256hex([]byte(canonicalCert(rec.ID, "salary", cert.ValueDigest, cert.ShreddedAt, ownerID)))
	if cert.CertHash != want {
		t.Fatalf("cert_hash não verifica: got %s want %s", cert.CertHash, want)
	}
	// O campo ficou sem chave (shredded) mas o value_blob mantém-se.
	_, fields, _ = repo.GetRecord(ctx, ownerID, rec.ID)
	var salary *Field
	for i := range fields {
		if fields[i].FieldName == "salary" {
			salary = &fields[i]
		}
	}
	if salary == nil || salary.WrappedKey != nil || salary.ShreddedAt == nil || len(salary.ValueBlob) == 0 {
		t.Fatalf("campo eliminado em estado inesperado: %+v", salary)
	}
	// Segundo shred é idempotente (ErrAlreadyShredded), sem novo certificado.
	if _, err := repo.ShredField(ctx, ownerID, rec.ID, "salary"); err != ErrAlreadyShredded {
		t.Fatalf("re-shred: esperado ErrAlreadyShredded, got %v", err)
	}
	// Estranho não pode emitir shred.
	if _, err := repo.ShredField(ctx, strangerID, rec.ID, "full_name"); err != ErrNotFound {
		t.Fatalf("shred estranho: esperado ErrNotFound, got %v", err)
	}
	// Certificado aparece na listagem do dono.
	certs, err := repo.ListCertificates(ctx, ownerID)
	if err != nil || len(certs) != 1 {
		t.Fatalf("ListCertificates: certs=%d err=%v", len(certs), err)
	}
	// Shred da ficha inteira emite cert por cada campo ainda com chave (full_name).
	more, err := repo.ShredRecord(ctx, ownerID, rec.ID)
	if err != nil || len(more) != 1 {
		t.Fatalf("ShredRecord: emitidos=%d err=%v", len(more), err)
	}
	if certs2, _ := repo.ListCertificates(ctx, ownerID); len(certs2) != 2 {
		t.Fatalf("após ShredRecord esperava 2 certificados, got %d", len(certs2))
	}

	// --- Logs imutaveis encadeados (HR-002) ---------------------------------
	// As accoes acima (create, puts, shred, shred-record) ja escreveram na cadeia.
	ok, broken, err := repo.VerifyAudit(ctx, ownerID)
	if err != nil || !ok {
		t.Fatalf("VerifyAudit: ok=%v broken=%d err=%v", ok, broken, err)
	}
	entries, err := repo.ListAudit(ctx, ownerID)
	if err != nil || len(entries) == 0 {
		t.Fatalf("ListAudit: entries=%d err=%v", len(entries), err)
	}
	// Seq comeca em 1 e o primeiro encadeia ao GENESIS.
	if entries[0].Seq != 1 || entries[0].PrevHash != "GENESIS" {
		t.Fatalf("primeira entrada inesperada: %+v", entries[0])
	}
	// Adulterar uma entrada antiga tem de partir a cadeia.
	if _, err := pool.Exec(ctx,
		`UPDATE audit_log SET detail = 'forjado' WHERE owner_id = $1 AND seq = 1`, ownerID,
	); err != nil {
		t.Fatal(err)
	}
	ok, broken, err = repo.VerifyAudit(ctx, ownerID)
	if err != nil || ok || broken != 1 {
		t.Fatalf("cadeia adulterada: ok=%v broken=%d err=%v (esperado ok=false broken=1)", ok, broken, err)
	}

	// --- Relatorio de conformidade RGPD (HR-008) ----------------------------
	rep, err := repo.ComplianceReportFor(ctx, ownerID)
	if err != nil {
		t.Fatal(err)
	}
	if rep.RecordCount != 1 || rep.CertificateCount != 2 {
		t.Fatalf("relatorio com contagens inesperadas: %+v", rep)
	}
	// A cadeia foi adulterada acima, logo o relatorio reflecte invalidade.
	if rep.AuditChainValid || rep.AuditBrokenSeq != 1 {
		t.Fatalf("relatorio deveria reportar cadeia partida: %+v", rep)
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
