package hr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Logs imutaveis com hashing encadeado (HR-002).
//
// Cada accao sobre as fichas acrescenta uma entrada a uma cadeia append-only,
// encadeada ao hash da anterior. Adulterar uma entrada antiga parte a cadeia de
// todas as seguintes — a fraude fica visivel ao recalcular os hashes.

// auditGenesis e o prev_hash da primeira entrada de cada dono.
const auditGenesis = "GENESIS"

// Accoes registadas (constantes para evitar strings soltas).
const (
	AuditRecordCreate = "record.create"
	AuditRecordDelete = "record.delete"
	AuditFieldPut     = "field.put"
	AuditFieldDelete  = "field.delete"
	AuditFieldShred   = "field.shred"
)

// AuditEntry espelha uma linha de "audit_log".
type AuditEntry struct {
	ID         string
	OwnerID    string
	Seq        int64
	Action     string
	Detail     string
	OccurredAt string // RFC3339Nano — a mesma string usada no entry_hash
	PrevHash   string
	EntryHash  string
}

// canonicalAudit constroi a string canonica versionada que e hasheada.
func canonicalAudit(seq int64, owner, action, detail, occurredAt, prevHash string) string {
	return fmt.Sprintf("v1|%d|%s|%s|%s|%s|%s", seq, owner, action, detail, occurredAt, prevHash)
}

// appendAuditTx acrescenta uma entrada a cadeia do dono dentro de uma transacao.
// Usa um advisory lock por dono para serializar appends concorrentes (sem ele,
// dois appends simultaneos podiam ler o mesmo "ultimo" e colidir no seq).
func (r *Repo) appendAuditTx(ctx context.Context, tx pgx.Tx, owner, action, detail string) (*AuditEntry, error) {
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, owner); err != nil {
		return nil, err
	}

	var lastSeq int64
	var prevHash string
	err := tx.QueryRow(ctx,
		`SELECT seq, entry_hash FROM audit_log WHERE owner_id = $1 ORDER BY seq DESC LIMIT 1`,
		owner,
	).Scan(&lastSeq, &prevHash)
	if errors.Is(err, pgx.ErrNoRows) {
		lastSeq = 0
		prevHash = auditGenesis
	} else if err != nil {
		return nil, err
	}

	seq := lastSeq + 1
	occurredAt := time.Now().UTC().Format(time.RFC3339Nano)
	entryHash := sha256hex([]byte(canonicalAudit(seq, owner, action, detail, occurredAt, prevHash)))

	e := &AuditEntry{
		OwnerID: owner, Seq: seq, Action: action, Detail: detail,
		OccurredAt: occurredAt, PrevHash: prevHash, EntryHash: entryHash,
	}
	if err := tx.QueryRow(ctx, `
		INSERT INTO audit_log (owner_id, seq, action, detail, occurred_at, prev_hash, entry_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text`,
		owner, seq, action, detail, occurredAt, prevHash, entryHash,
	).Scan(&e.ID); err != nil {
		return nil, err
	}
	return e, nil
}

// AppendAudit acrescenta uma entrada numa transacao propria (para chamadas fora
// de uma transacao existente).
func (r *Repo) AppendAudit(ctx context.Context, owner, action, detail string) (*AuditEntry, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	e, err := r.appendAuditTx(ctx, tx, owner, action, detail)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return e, nil
}

// ListAudit devolve a cadeia do dono por ordem cronologica (seq crescente).
func (r *Repo) ListAudit(ctx context.Context, owner string) ([]AuditEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, seq, action, detail, occurred_at, prev_hash, entry_hash
		FROM audit_log WHERE owner_id = $1 ORDER BY seq ASC`, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(
			&e.ID, &e.OwnerID, &e.Seq, &e.Action, &e.Detail,
			&e.OccurredAt, &e.PrevHash, &e.EntryHash,
		); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// VerifyAudit recalcula a cadeia do dono. Devolve ok=true se intacta; caso
// contrario ok=false e brokenSeq aponta a primeira entrada inconsistente.
func (r *Repo) VerifyAudit(ctx context.Context, owner string) (bool, int64, error) {
	entries, err := r.ListAudit(ctx, owner)
	if err != nil {
		return false, 0, err
	}
	prev := auditGenesis
	var expectedSeq int64 = 1
	for _, e := range entries {
		if e.Seq != expectedSeq || e.PrevHash != prev {
			return false, e.Seq, nil
		}
		want := sha256hex([]byte(canonicalAudit(e.Seq, owner, e.Action, e.Detail, e.OccurredAt, prev)))
		if want != e.EntryHash {
			return false, e.Seq, nil
		}
		prev = e.EntryHash
		expectedSeq++
	}
	return true, 0, nil
}
