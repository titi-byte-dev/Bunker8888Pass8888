package hr

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Crypto-shredding (HR-003) + certificado de eliminação (HR-004).
//
// Shred = destruir a chave do campo (wrapped_key := NULL) e carimbar shredded_at.
// O value_blob fica, mas sem a chave é lixo indecifrável. Cada shred emite um
// certificado verificável: uma impressão digital (sha256) sobre os factos da
// eliminação, que o cliente pode recalcular para confirmar a prova.

// ErrAlreadyShredded: o campo já tinha sido eliminado (sem chave).
var ErrAlreadyShredded = errors.New("hr: campo já eliminado")

// Certificate é a prova criptográfica de uma eliminação (HR-004).
type Certificate struct {
	ID          string
	OwnerID     string
	RecordID    string // pode estar vazio se a ficha foi apagada depois
	FieldName   string
	ValueDigest string // sha256(value_blob) em hex
	ShreddedAt  string // RFC3339Nano — a mesma string usada no cálculo do hash
	CertHash    string // sha256(canonical) em hex
	IssuedAt    time.Time
}

// canonicalCert constrói a string canónica que é hasheada para o cert_hash.
// O formato é estável e versionado (v1) para o cliente reproduzir exactamente.
func canonicalCert(recordID, fieldName, valueDigest, shreddedAt, ownerID string) string {
	return "v1|" + recordID + "|" + fieldName + "|" + valueDigest + "|" + shreddedAt + "|" + ownerID
}

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// ShredField destrói a chave de um campo (crypto-shredding) e emite o respectivo
// certificado, tudo numa transação. Idempotência: se já estava eliminado,
// devolve ErrAlreadyShredded sem criar um segundo certificado.
func (r *Repo) ShredField(ctx context.Context, ownerID, recordID, fieldName string) (*Certificate, error) {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Lê o campo e bloqueia a linha. Precisamos do value_blob para o digest.
	var valueBlob []byte
	var alreadyNull bool
	err = tx.QueryRow(ctx, `
		SELECT value_blob, wrapped_key IS NULL
		FROM employee_fields
		WHERE record_id = $1 AND field_name = $2
		FOR UPDATE`, recordID, fieldName,
	).Scan(&valueBlob, &alreadyNull)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if alreadyNull {
		return nil, ErrAlreadyShredded
	}

	shreddedAt := time.Now().UTC().Format(time.RFC3339Nano)
	valueDigest := sha256hex(valueBlob)
	certHash := sha256hex([]byte(canonicalCert(recordID, fieldName, valueDigest, shreddedAt, ownerID)))

	// Destrói a chave: o campo torna-se permanentemente indecifrável.
	if _, err := tx.Exec(ctx, `
		UPDATE employee_fields
		SET wrapped_key = NULL, shredded_at = now(), updated_at = now()
		WHERE record_id = $1 AND field_name = $2`, recordID, fieldName,
	); err != nil {
		return nil, err
	}

	cert := &Certificate{
		OwnerID: ownerID, RecordID: recordID, FieldName: fieldName,
		ValueDigest: valueDigest, ShreddedAt: shreddedAt, CertHash: certHash,
	}
	if err := tx.QueryRow(ctx, `
		INSERT INTO erasure_certificates
			(owner_id, record_id, field_name, value_digest, shredded_at, cert_hash)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, issued_at`,
		ownerID, recordID, fieldName, valueDigest, shreddedAt, certHash,
	).Scan(&cert.ID, &cert.IssuedAt); err != nil {
		return nil, err
	}

	// Regista a eliminacao na cadeia de auditoria, na mesma transacao.
	if _, err := r.appendAuditTx(ctx, tx, ownerID, AuditFieldShred, fieldName); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return cert, nil
}

// ShredRecord elimina (crypto-shred) TODOS os campos ainda com chave de uma
// ficha, emitindo um certificado por cada um. Erasure RGPD da ficha inteira.
func (r *Repo) ShredRecord(ctx context.Context, ownerID, recordID string) ([]Certificate, error) {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	// Nomes dos campos ainda não eliminados.
	rows, err := r.pool.Query(ctx, `
		SELECT field_name FROM employee_fields
		WHERE record_id = $1 AND wrapped_key IS NOT NULL
		ORDER BY field_name ASC`, recordID)
	if err != nil {
		return nil, err
	}
	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			rows.Close()
			return nil, err
		}
		names = append(names, n)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var certs []Certificate
	for _, n := range names {
		cert, err := r.ShredField(ctx, ownerID, recordID, n)
		if err != nil {
			if errors.Is(err, ErrAlreadyShredded) {
				continue // corrida benigna: outro shred chegou primeiro
			}
			return nil, err
		}
		certs = append(certs, *cert)
	}
	return certs, nil
}

// ListCertificates devolve os certificados de eliminação do utilizador, recentes
// primeiro. Sobrevivem à remoção das fichas (record_id pode estar a NULL).
func (r *Repo) ListCertificates(ctx context.Context, ownerID string) ([]Certificate, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, COALESCE(record_id::text, ''),
		       field_name, value_digest, shredded_at, cert_hash, issued_at
		FROM erasure_certificates
		WHERE owner_id = $1
		ORDER BY issued_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Certificate
	for rows.Next() {
		var c Certificate
		if err := rows.Scan(
			&c.ID, &c.OwnerID, &c.RecordID, &c.FieldName,
			&c.ValueDigest, &c.ShreddedAt, &c.CertHash, &c.IssuedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
