// Package hr implementa as Fichas de Empregado com cifragem CAMPO-A-CAMPO
// (HR-001): cada campo de uma ficha é cifrado de forma independente, com a sua
// própria chave, do lado do cliente.
//
// ⚠️ Modelo Zero-Knowledge: o servidor trata value_blob e wrapped_key como bytes
// opacos. Só o field_name (chave de esquema, ex.: "salary") é visível, tal como
// o item_type nos cofres. O servidor nunca decifra nada.
//
// Esta granularidade por campo é a fundação do crypto-shredding (HR-003): apagar
// o wrapped_key de um campo torna o seu valor permanentemente indecifrável, sem
// tocar nos restantes campos.
package hr

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Erros de domínio (mapeados para HTTP nos handlers).
var (
	// ErrNotFound: ficha inexistente OU não pertence ao requerente (mesma
	// resposta, para não revelar a existência de fichas alheias).
	ErrNotFound = errors.New("hr: ficha não encontrada")
	// ErrInvalidField: nome de campo vazio ou blob em falta.
	ErrInvalidField = errors.New("hr: campo inválido")
)

// Record espelha uma linha de "employee_records" (sem PII: tudo vive nos campos).
type Record struct {
	ID        string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Field espelha uma linha de "employee_fields".
type Field struct {
	ID         string
	RecordID   string
	FieldName  string
	ValueBlob  []byte
	WrappedKey []byte // nil quando o campo foi crypto-shredded (HR-003)
	ShreddedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Repo acede às tabelas employee_*.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// ownsRecord confirma que a ficha existe E pertence ao requerente.
func (r *Repo) ownsRecord(ctx context.Context, ownerID, recordID string) error {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM employee_records WHERE id = $1 AND owner_id = $2)`,
		recordID, ownerID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

// CreateRecord cria uma ficha vazia para o requerente. Os campos entram depois,
// um a um, via PutField — cada um com a sua própria chave.
func (r *Repo) CreateRecord(ctx context.Context, ownerID string) (*Record, error) {
	rec := &Record{OwnerID: ownerID}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO employee_records (owner_id)
		VALUES ($1)
		RETURNING id::text, created_at, updated_at`,
		ownerID,
	).Scan(&rec.ID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// ListRecords devolve as fichas do requerente (sem campos), mais recentes primeiro.
func (r *Repo) ListRecords(ctx context.Context, ownerID string) ([]Record, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, created_at, updated_at
		FROM employee_records
		WHERE owner_id = $1
		ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ID, &rec.OwnerID, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// GetRecord devolve uma ficha do requerente com TODOS os seus campos cifrados.
func (r *Repo) GetRecord(ctx context.Context, ownerID, recordID string) (*Record, []Field, error) {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, nil, err
	}
	rec := &Record{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, owner_id::text, created_at, updated_at
		FROM employee_records WHERE id = $1`, recordID,
	).Scan(&rec.ID, &rec.OwnerID, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	fields, err := r.listFields(ctx, recordID)
	if err != nil {
		return nil, nil, err
	}
	return rec, fields, nil
}

// listFields devolve os campos de uma ficha (ordenados por nome para estabilidade).
func (r *Repo) listFields(ctx context.Context, recordID string) ([]Field, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, record_id::text, field_name, value_blob, wrapped_key,
		       shredded_at, created_at, updated_at
		FROM employee_fields
		WHERE record_id = $1
		ORDER BY field_name ASC`, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Field
	for rows.Next() {
		var f Field
		if err := rows.Scan(
			&f.ID, &f.RecordID, &f.FieldName, &f.ValueBlob, &f.WrappedKey,
			&f.ShreddedAt, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// PutField cria ou actualiza um campo da ficha (upsert por field_name). value e
// wrappedKey já vêm cifrados do cliente. Repor um campo limpa o estado de shred.
func (r *Repo) PutField(ctx context.Context, ownerID, recordID, fieldName string, valueBlob, wrappedKey []byte) (*Field, error) {
	if fieldName == "" || len(valueBlob) == 0 || len(wrappedKey) == 0 {
		return nil, ErrInvalidField
	}
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return nil, err
	}
	f := &Field{RecordID: recordID, FieldName: fieldName, ValueBlob: valueBlob, WrappedKey: wrappedKey}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO employee_fields (record_id, field_name, value_blob, wrapped_key)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (record_id, field_name) DO UPDATE SET
			value_blob  = EXCLUDED.value_blob,
			wrapped_key = EXCLUDED.wrapped_key,
			shredded_at = NULL,
			updated_at  = now()
		RETURNING id::text, shredded_at, created_at, updated_at`,
		recordID, fieldName, valueBlob, wrappedKey,
	).Scan(&f.ID, &f.ShreddedAt, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// Após o upsert, refrescamos o timestamp da ficha (toque de auditoria).
	_, _ = r.pool.Exec(ctx, `UPDATE employee_records SET updated_at = now() WHERE id = $1`, recordID)
	return f, nil
}

// DeleteField remove por completo um campo da ficha (linha inteira). Para apagar
// só a chave mantendo o histórico cifrado, HR-003 usará um shred dedicado.
func (r *Repo) DeleteField(ctx context.Context, ownerID, recordID, fieldName string) error {
	if err := r.ownsRecord(ctx, ownerID, recordID); err != nil {
		return err
	}
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM employee_fields WHERE record_id = $1 AND field_name = $2`,
		recordID, fieldName)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteRecord apaga a ficha inteira (campos em cascata). Só o dono o pode fazer.
func (r *Repo) DeleteRecord(ctx context.Context, ownerID, recordID string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM employee_records WHERE id = $1 AND owner_id = $2`,
		recordID, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
