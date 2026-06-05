package sharedvaults

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// MaxAttachmentBytes limita o tamanho do CIPHERTEXT de um anexo (5 MiB). É um
// teto didático para guardar ficheiros na própria base de dados sem object
// storage; o cliente deve recusar ficheiros maiores antes de cifrar.
const MaxAttachmentBytes = 5 << 20 // 5 MiB

// ErrTooLarge: o anexo excede MaxAttachmentBytes.
var ErrTooLarge = errors.New("sharedvaults: anexo demasiado grande")

// Attachment é um anexo completo, incluindo os bytes cifrados do ficheiro.
type Attachment struct {
	ID        string
	VaultID   string
	MetaBlob  []byte // AES-GCM(chave_cofre, JSON{name,mime,size})
	DataBlob  []byte // AES-GCM(chave_cofre, bytes) — só presente no GET individual
	ByteSize  int64
	CreatedBy string
	CreatedAt time.Time
}

// AttachmentMeta é a vista "leve" de um anexo: tudo menos os bytes do ficheiro.
// Usada na listagem, para não arrastar megabytes só para mostrar nomes.
type AttachmentMeta struct {
	ID        string
	VaultID   string
	MetaBlob  []byte
	ByteSize  int64
	CreatedBy string
	CreatedAt time.Time
}

// AddAttachment guarda um anexo cifrado. Exige papel com escrita. metaBlob e
// dataBlob já vêm cifrados com a chave do cofre (o servidor nunca os abre).
func (r *Repo) AddAttachment(ctx context.Context, vaultID, userID string, metaBlob, dataBlob []byte) (*AttachmentMeta, error) {
	if len(metaBlob) == 0 || len(dataBlob) == 0 {
		return nil, errors.New("sharedvaults: anexo sem metadados ou dados")
	}
	if len(dataBlob) > MaxAttachmentBytes {
		return nil, ErrTooLarge
	}
	role, err := r.roleOf(ctx, vaultID, userID)
	if err != nil {
		return nil, err
	}
	if !canWriteItems(role) {
		return nil, ErrForbidden
	}
	a := &AttachmentMeta{VaultID: vaultID, MetaBlob: metaBlob, ByteSize: int64(len(dataBlob)), CreatedBy: userID}
	err = r.pool.QueryRow(ctx, `
		INSERT INTO shared_vault_attachments (vault_id, meta_blob, data_blob, byte_size, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, created_at`,
		vaultID, metaBlob, dataBlob, len(dataBlob), userID,
	).Scan(&a.ID, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListAttachments devolve só os metadados (sem os bytes). Qualquer membro lê.
func (r *Repo) ListAttachments(ctx context.Context, vaultID, userID string) ([]AttachmentMeta, error) {
	if _, err := r.roleOf(ctx, vaultID, userID); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, vault_id::text, meta_blob, byte_size,
		       COALESCE(created_by::text, ''), created_at
		FROM shared_vault_attachments
		WHERE vault_id = $1
		ORDER BY created_at DESC`, vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AttachmentMeta
	for rows.Next() {
		var a AttachmentMeta
		if err := rows.Scan(&a.ID, &a.VaultID, &a.MetaBlob, &a.ByteSize, &a.CreatedBy, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// GetAttachment devolve o anexo completo (com os bytes cifrados) para download.
// Qualquer membro (incl. viewer) pode descarregar.
func (r *Repo) GetAttachment(ctx context.Context, vaultID, userID, attID string) (*Attachment, error) {
	if _, err := r.roleOf(ctx, vaultID, userID); err != nil {
		return nil, err
	}
	a := &Attachment{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, vault_id::text, meta_blob, data_blob, byte_size,
		       COALESCE(created_by::text, ''), created_at
		FROM shared_vault_attachments
		WHERE id = $1 AND vault_id = $2`, attID, vaultID,
	).Scan(&a.ID, &a.VaultID, &a.MetaBlob, &a.DataBlob, &a.ByteSize, &a.CreatedBy, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

// DeleteAttachment remove um anexo. Exige papel com escrita (tudo menos viewer).
func (r *Repo) DeleteAttachment(ctx context.Context, vaultID, userID, attID string) error {
	role, err := r.roleOf(ctx, vaultID, userID)
	if err != nil {
		return err
	}
	if !canWriteItems(role) {
		return ErrForbidden
	}
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM shared_vault_attachments WHERE id = $1 AND vault_id = $2`, attID, vaultID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
