// Package commissions implementa o registo de comissoes sobre faturas pagas
// (FIN-007). O conteudo (beneficiario, percentagem, valor) e um blob cifrado
// com a Master Key (ZK); apenas a ligacao a fatura e o estado vivem em claro.
package commissions

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("commissions: comissao nao encontrada")
	ErrInvalidStatus = errors.New("commissions: estado invalido")
)

// Estados de liquidacao ao beneficiario.
const (
	StatusPending = "pending"
	StatusPaid    = "paid"
	StatusVoid    = "void"
)

func validStatus(s string) bool {
	return s == StatusPending || s == StatusPaid || s == StatusVoid
}

// Commission espelha uma linha de commissions.
type Commission struct {
	ID        string
	OwnerID   string
	InvoiceID string // "" quando a fatura foi removida ou nao foi indicada
	Status    string
	Blob      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repo acede a commissions.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// Create regista uma comissao nova (estado inicial pending). invoiceID e
// opcional; quando indicado, e validado contra uma fatura do mesmo dono para
// evitar referencias cruzadas entre contas.
func (r *Repo) Create(ctx context.Context, ownerID, invoiceID string, blob []byte) (*Commission, error) {
	if len(blob) == 0 {
		return nil, errors.New("commissions: blob vazio")
	}

	var invoice any
	if invoiceID != "" {
		// A fatura tem de pertencer ao dono. Caso contrario, ErrNotFound.
		var exists bool
		if err := r.pool.QueryRow(ctx,
			`SELECT EXISTS(SELECT 1 FROM invoices WHERE id = $1 AND owner_id = $2)`,
			invoiceID, ownerID,
		).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrNotFound
		}
		invoice = invoiceID
	}

	c := &Commission{OwnerID: ownerID, InvoiceID: invoiceID, Status: StatusPending, Blob: blob}
	if err := r.pool.QueryRow(ctx, `
		INSERT INTO commissions (owner_id, invoice_id, status, blob)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, created_at, updated_at`,
		ownerID, invoice, StatusPending, blob,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return c, nil
}

func scanCommission(row pgx.Row) (*Commission, error) {
	var c Commission
	var invoice *string
	if err := row.Scan(&c.ID, &c.OwnerID, &invoice, &c.Status,
		&c.Blob, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	if invoice != nil {
		c.InvoiceID = *invoice
	}
	return &c, nil
}

const selectCols = `id::text, owner_id::text, invoice_id::text, status, blob, created_at, updated_at`

// List devolve as comissoes do dono, mais recentes primeiro.
func (r *Repo) List(ctx context.Context, ownerID string) ([]Commission, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+selectCols+` FROM commissions WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Commission
	for rows.Next() {
		c, err := scanCommission(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, rows.Err()
}

// Get devolve uma comissao do dono.
func (r *Repo) Get(ctx context.Context, ownerID, id string) (*Commission, error) {
	c, err := scanCommission(r.pool.QueryRow(ctx,
		`SELECT `+selectCols+` FROM commissions WHERE id = $1 AND owner_id = $2`, id, ownerID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

// UpdateStatus muda o estado de liquidacao (pending -> paid | void).
func (r *Repo) UpdateStatus(ctx context.Context, ownerID, id, status string) (*Commission, error) {
	if !validStatus(status) {
		return nil, ErrInvalidStatus
	}
	tag, err := r.pool.Exec(ctx,
		`UPDATE commissions SET status = $1, updated_at = now() WHERE id = $2 AND owner_id = $3`,
		status, id, ownerID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return r.Get(ctx, ownerID, id)
}
