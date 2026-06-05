// Package invoicing implementa faturacao com numeracao legal sequencial
// (FIN-006). O conteudo do documento e um blob cifrado com a Master Key (ZK);
// apenas o numero legal e o estado vivem em claro, geridos pelo servidor.
package invoicing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("invoicing: documento nao encontrado")
	ErrInvalidType   = errors.New("invoicing: tipo de documento invalido")
	ErrInvalidStatus = errors.New("invoicing: estado invalido")
)

// Tipos de documento e respectivos prefixos legais.
const (
	DocProforma = "proforma"
	DocInvoice  = "invoice"
	DocReceipt  = "receipt"
)

// Estados possiveis de um documento.
const (
	StatusIssued = "issued"
	StatusPaid   = "paid"
	StatusVoid   = "void"
)

var prefixes = map[string]string{
	DocProforma: "PF",
	DocInvoice:  "FT",
	DocReceipt:  "RC",
}

func validType(t string) bool   { _, ok := prefixes[t]; return ok }
func validStatus(s string) bool { return s == StatusIssued || s == StatusPaid || s == StatusVoid }

// Document espelha uma linha de invoices.
type Document struct {
	ID           string
	OwnerID      string
	DocType      string
	Year         int
	Seq          int64
	Number       string
	Status       string
	SourceLeadID string // "" quando nao nasce de um lead
	Blob         []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repo acede a invoices.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// formatNumber constroi o numero legal: "<prefixo> <ano>/<seq:0000>".
func formatNumber(docType string, year int, seq int64) string {
	return fmt.Sprintf("%s %d/%04d", prefixes[docType], year, seq)
}

// Issue emite um documento novo, atribuindo o proximo numero legal sequencial
// por (dono, tipo, ano). Usa um advisory lock por essa chave para serializar
// emissoes concorrentes — sem ele, dois pedidos podiam ler o mesmo "ultimo seq"
// e colidir (ou abrir buracos) na numeracao.
func (r *Repo) Issue(ctx context.Context, ownerID, docType, sourceLeadID string, blob []byte) (*Document, error) {
	if !validType(docType) {
		return nil, ErrInvalidType
	}
	if len(blob) == 0 {
		return nil, errors.New("invoicing: blob vazio")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	year := time.Now().UTC().Year()

	// Lock por dono+tipo+ano (a numeracao reinicia a cada ano/serie).
	lockKey := fmt.Sprintf("%s|%s|%d", ownerID, docType, year)
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, lockKey); err != nil {
		return nil, err
	}

	var lastSeq int64
	err = tx.QueryRow(ctx, `
		SELECT seq FROM invoices
		WHERE owner_id = $1 AND doc_type = $2 AND year = $3
		ORDER BY seq DESC LIMIT 1`,
		ownerID, docType, year,
	).Scan(&lastSeq)
	if errors.Is(err, pgx.ErrNoRows) {
		lastSeq = 0
	} else if err != nil {
		return nil, err
	}

	seq := lastSeq + 1
	number := formatNumber(docType, year, seq)

	d := &Document{
		OwnerID: ownerID, DocType: docType, Year: year, Seq: seq,
		Number: number, Status: StatusIssued, SourceLeadID: sourceLeadID, Blob: blob,
	}

	var lead any
	if sourceLeadID != "" {
		lead = sourceLeadID
	}
	if err := tx.QueryRow(ctx, `
		INSERT INTO invoices (owner_id, doc_type, year, seq, number, status, source_lead_id, blob)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id::text, created_at, updated_at`,
		ownerID, docType, year, seq, number, StatusIssued, lead, blob,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

func scanDoc(row pgx.Row) (*Document, error) {
	var d Document
	var lead *string
	if err := row.Scan(&d.ID, &d.OwnerID, &d.DocType, &d.Year, &d.Seq, &d.Number,
		&d.Status, &lead, &d.Blob, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, err
	}
	if lead != nil {
		d.SourceLeadID = *lead
	}
	return &d, nil
}

const selectCols = `id::text, owner_id::text, doc_type, year, seq, number, status,
	source_lead_id::text, blob, created_at, updated_at`

// List devolve os documentos do dono, mais recentes primeiro.
func (r *Repo) List(ctx context.Context, ownerID string) ([]Document, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+selectCols+` FROM invoices WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Document
	for rows.Next() {
		d, err := scanDoc(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return out, rows.Err()
}

// Get devolve um documento do dono.
func (r *Repo) Get(ctx context.Context, ownerID, id string) (*Document, error) {
	d, err := scanDoc(r.pool.QueryRow(ctx,
		`SELECT `+selectCols+` FROM invoices WHERE id = $1 AND owner_id = $2`, id, ownerID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}

// UpdateStatus muda o estado do documento (issued -> paid | void). NUNCA apaga:
// a tabela e fiscalmente append-only.
func (r *Repo) UpdateStatus(ctx context.Context, ownerID, id, status string) (*Document, error) {
	if !validStatus(status) {
		return nil, ErrInvalidStatus
	}
	tag, err := r.pool.Exec(ctx,
		`UPDATE invoices SET status = $1, updated_at = now() WHERE id = $2 AND owner_id = $3`,
		status, id, ownerID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.Get(ctx, ownerID, id)
}
