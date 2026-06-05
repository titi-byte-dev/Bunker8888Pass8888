package mail

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInboxNotFound = errors.New("mail: mensagem de inbox não encontrada")

// InboxMessage espelha mail_inbox_messages (stub até MAIL-002).
type InboxMessage struct {
	ID          string
	OwnerID     string
	FromEmail   string
	Subject     string
	Body        string
	ReceivedAt  time.Time
	ProcessedAt *time.Time
	CreatedAt   time.Time
}

// InboxRepo acede à caixa de entrada simulada.
type InboxRepo struct {
	pool *pgxpool.Pool
}

func NewInboxRepo(pool *pgxpool.Pool) *InboxRepo {
	return &InboxRepo{pool: pool}
}

// CreateInboxMessage insere uma mensagem (dev/seed ou relay futuro).
func (r *InboxRepo) CreateInboxMessage(ctx context.Context, ownerID, fromEmail, subject, body string) (*InboxMessage, error) {
	fromEmail = strings.TrimSpace(fromEmail)
	subject = strings.TrimSpace(subject)
	body = strings.TrimSpace(body)
	if !looksLikeEmail(fromEmail) {
		return nil, ErrInvalidDest
	}
	if subject == "" || len(subject) > 500 {
		return nil, errors.New("mail: subject inválido")
	}
	if body == "" || len(body) > 20000 {
		return nil, errors.New("mail: body inválido")
	}
	m := &InboxMessage{OwnerID: ownerID, FromEmail: fromEmail, Subject: subject, Body: body}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO mail_inbox_messages (owner_id, from_email, subject, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, received_at, created_at`,
		ownerID, fromEmail, subject, body,
	).Scan(&m.ID, &m.ReceivedAt, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ListInbox devolve mensagens do utilizador (mais recentes primeiro).
func (r *InboxRepo) ListInbox(ctx context.Context, ownerID string, unprocessedOnly bool) ([]InboxMessage, error) {
	q := `
		SELECT id::text, owner_id::text, from_email, subject, body,
		       received_at, processed_at, created_at
		FROM mail_inbox_messages
		WHERE owner_id = $1`
	if unprocessedOnly {
		q += ` AND processed_at IS NULL`
	}
	q += ` ORDER BY received_at DESC`

	rows, err := r.pool.Query(ctx, q, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []InboxMessage
	for rows.Next() {
		var m InboxMessage
		if err := rows.Scan(
			&m.ID, &m.OwnerID, &m.FromEmail, &m.Subject, &m.Body,
			&m.ReceivedAt, &m.ProcessedAt, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// MarkProcessed marca uma mensagem como já tratada pela prospeção.
func (r *InboxRepo) MarkProcessed(ctx context.Context, ownerID, messageID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE mail_inbox_messages SET processed_at = now()
		WHERE id = $1 AND owner_id = $2 AND processed_at IS NULL`,
		messageID, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInboxNotFound
	}
	return nil
}

// GetInboxMessage devolve uma mensagem (testes/uso interno).
func (r *InboxRepo) GetInboxMessage(ctx context.Context, ownerID, messageID string) (*InboxMessage, error) {
	m := &InboxMessage{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, owner_id::text, from_email, subject, body,
		       received_at, processed_at, created_at
		FROM mail_inbox_messages WHERE id = $1 AND owner_id = $2`,
		messageID, ownerID,
	).Scan(
		&m.ID, &m.OwnerID, &m.FromEmail, &m.Subject, &m.Body,
		&m.ReceivedAt, &m.ProcessedAt, &m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInboxNotFound
		}
		return nil, err
	}
	return m, nil
}
