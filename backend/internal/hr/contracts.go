package hr

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Gestao de contratos (HR-005) + assinatura digital (HR-006).
//
// Contratos sao ficheiros cifrados ficheiro-a-ficheiro (file_key embrulhada com
// a Master Key), tal como os campos da HR-001. A assinatura digital (ECDSA
// P-256) e feita no cliente sobre o digest do ciphertext; o servidor so guarda e
// faz cumprir a posse — nunca decifra nem assina.

// MaxContractBytes limita o ciphertext de um contrato (5 MiB), como nos anexos.
const MaxContractBytes = 5 << 20

// ErrTooLargeContract: o contrato excede MaxContractBytes.
var ErrTooLargeContract = errors.New("hr: contrato demasiado grande")

// ErrNoSigningIdentity: o utilizador ainda nao tem identidade de assinatura.
var ErrNoSigningIdentity = errors.New("hr: sem identidade de assinatura")

// Contract espelha uma linha de "employee_contracts" (com ou sem os bytes).
type Contract struct {
	ID            string
	RecordID      string
	MetaBlob      []byte
	DataBlob      []byte // so presente no GET individual
	WrappedKey    []byte
	ByteSize      int64
	ContentDigest string // sha256(data_blob) hex, presente se assinado
	Signature     []byte // ECDSA, presente se assinado
	SignedBy      string
	SignedAt      *time.Time
	CreatedBy     string
	CreatedAt     time.Time
}

// SigningIdentity espelha "hr_signing_identities".
type SigningIdentity struct {
	OwnerID           string
	PublicKey         []byte
	WrappedPrivateKey []byte
	CreatedAt         time.Time
}

// PutSigningIdentity grava (ou substitui) a identidade de assinatura do dono.
func (r *Repo) PutSigningIdentity(ctx context.Context, ownerID string, publicKey, wrappedPrivate []byte) error {
	if len(publicKey) == 0 || len(wrappedPrivate) == 0 {
		return errors.New("hr: identidade de assinatura incompleta")
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO hr_signing_identities (owner_id, public_key, wrapped_private_key)
		VALUES ($1, $2, $3)
		ON CONFLICT (owner_id) DO UPDATE SET
			public_key          = EXCLUDED.public_key,
			wrapped_private_key = EXCLUDED.wrapped_private_key`,
		ownerID, publicKey, wrappedPrivate,
	)
	return err
}

// GetSigningIdentity devolve a identidade de assinatura do dono, ou
// ErrNoSigningIdentity se ainda nao existir.
func (r *Repo) GetSigningIdentity(ctx context.Context, ownerID string) (*SigningIdentity, error) {
	si := &SigningIdentity{OwnerID: ownerID}
	err := r.pool.QueryRow(ctx, `
		SELECT public_key, wrapped_private_key, created_at
		FROM hr_signing_identities WHERE owner_id = $1`, ownerID,
	).Scan(&si.PublicKey, &si.WrappedPrivateKey, &si.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoSigningIdentity
		}
		return nil, err
	}
	return si, nil
}

// GetSignerPublicKey devolve so a chave publica de um utilizador (para verificar
// assinaturas). Publica por definicao.
func (r *Repo) GetSignerPublicKey(ctx context.Context, signerID string) ([]byte, error) {
	var pk []byte
	err := r.pool.QueryRow(ctx,
		`SELECT public_key FROM hr_signing_identities WHERE owner_id = $1`, signerID,
	).Scan(&pk)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoSigningIdentity
		}
		return nil, err
	}
	return pk, nil
}

// AddContract guarda um contrato cifrado numa ficha. Exige posse da ficha.
func (r *Repo) AddContract(ctx context.Context, ownerID, recordID string, metaBlob, dataBlob, wrappedKey []byte) (*Contract, error) {
	if len(metaBlob) == 0 || len(dataBlob) == 0 || len(wrappedKey) == 0 {
		return nil, ErrInvalidField
	}
	if len(dataBlob) > MaxContractBytes {
		return nil, ErrTooLargeContract
	}
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	c := &Contract{RecordID: recordID, MetaBlob: metaBlob, WrappedKey: wrappedKey, ByteSize: int64(len(dataBlob)), CreatedBy: ownerID}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO employee_contracts (record_id, meta_blob, data_blob, wrapped_key, byte_size, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, created_at`,
		recordID, metaBlob, dataBlob, wrappedKey, len(dataBlob), ownerID,
	).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	if _, err := r.AppendAudit(ctx, ownerID, AuditContractAdd, c.ID); err != nil {
		return nil, err
	}
	return c, nil
}

// scanContractMeta lê os campos de listagem (sem data_blob).
func scanContractMeta(rows pgx.Rows) (Contract, error) {
	var c Contract
	var signedBy *string
	err := rows.Scan(
		&c.ID, &c.RecordID, &c.MetaBlob, &c.WrappedKey, &c.ByteSize,
		&c.ContentDigest, &c.Signature, &signedBy, &c.SignedAt,
		&c.CreatedBy, &c.CreatedAt,
	)
	if signedBy != nil {
		c.SignedBy = *signedBy
	}
	return c, err
}

const contractMetaCols = `id::text, record_id::text, meta_blob, wrapped_key, byte_size,
	COALESCE(content_digest, ''), signature, signed_by::text, signed_at,
	COALESCE(created_by::text, ''), created_at`

// ListContracts devolve os contratos de uma ficha (sem os bytes). Exige posse.
func (r *Repo) ListContracts(ctx context.Context, ownerID, recordID string) ([]Contract, error) {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx,
		`SELECT `+contractMetaCols+` FROM employee_contracts WHERE record_id = $1 ORDER BY created_at DESC`,
		recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Contract
	for rows.Next() {
		c, err := scanContractMeta(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// GetContract devolve um contrato completo (com os bytes cifrados). Exige posse.
func (r *Repo) GetContract(ctx context.Context, ownerID, recordID, contractID string) (*Contract, error) {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	var c Contract
	var signedBy *string
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, record_id::text, meta_blob, data_blob, wrapped_key, byte_size,
		       COALESCE(content_digest, ''), signature, signed_by::text, signed_at,
		       COALESCE(created_by::text, ''), created_at
		FROM employee_contracts WHERE id = $1 AND record_id = $2`, contractID, recordID,
	).Scan(
		&c.ID, &c.RecordID, &c.MetaBlob, &c.DataBlob, &c.WrappedKey, &c.ByteSize,
		&c.ContentDigest, &c.Signature, &signedBy, &c.SignedAt, &c.CreatedBy, &c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if signedBy != nil {
		c.SignedBy = *signedBy
	}
	return &c, nil
}

// SignContract grava a assinatura digital de um contrato (HR-006). A assinatura e
// a verificacao acontecem no cliente; o servidor so persiste e audita.
func (r *Repo) SignContract(ctx context.Context, ownerID, recordID, contractID, contentDigest string, signature []byte) (*Contract, error) {
	if contentDigest == "" || len(signature) == 0 {
		return nil, ErrInvalidField
	}
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	tag, err := r.pool.Exec(ctx, `
		UPDATE employee_contracts
		SET content_digest = $1, signature = $2, signed_by = $3, signed_at = now()
		WHERE id = $4 AND record_id = $5`,
		contentDigest, signature, ownerID, contractID, recordID,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	if _, err := r.AppendAudit(ctx, ownerID, AuditContractSign, contractID); err != nil {
		return nil, err
	}
	return r.GetContract(ctx, ownerID, recordID, contractID)
}

// DeleteContract remove um contrato. Exige posse da ficha.
func (r *Repo) DeleteContract(ctx context.Context, ownerID, recordID, contractID string) error {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM employee_contracts WHERE id = $1 AND record_id = $2`, contractID, recordID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	if _, err := r.AppendAudit(ctx, ownerID, AuditContractDelete, contractID); err != nil {
		return err
	}
	return nil
}
