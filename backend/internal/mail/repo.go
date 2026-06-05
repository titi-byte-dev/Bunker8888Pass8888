// Package mail implementa os Aliases de e-mail (MAIL-001): enderecos
// descartaveis que reencaminham para o e-mail real do utilizador.
//
// Nota: o envio efectivo (relay SMTP) fica para o MAIL-002. Aqui geramos e
// geримos apenas a configuracao dos aliases. O destino e visivel ao servidor
// (e ele que, no futuro, fara o reencaminhamento) — uma excecao consciente ao
// modelo Zero-Knowledge, documentada na migracao 0017.
package mail

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AliasDomain e o dominio dos aliases gerados (didactico; sem MX real ainda).
const AliasDomain = "aegis.email"

// Erros de dominio.
var (
	ErrNotFound       = errors.New("mail: alias nao encontrado")
	ErrInvalidDest    = errors.New("mail: destino invalido")
	errAliasCollision = errors.New("mail: colisao de alias")
)

// Alias espelha uma linha de "email_aliases".
type Alias struct {
	ID           string
	OwnerID      string
	AliasAddress string
	Destination  string
	Label        string
	Active       bool
	CreatedAt    time.Time
}

// Repo acede a tabela email_aliases.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositorio.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// randomLocalPart gera a parte local do alias (10 hex = 5 bytes aleatorios).
func randomLocalPart() (string, error) {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// looksLikeEmail faz uma validacao minima do destino (didactica, nao RFC).
func looksLikeEmail(s string) bool {
	at := strings.IndexByte(s, '@')
	return at > 0 && at < len(s)-1 && !strings.ContainsAny(s, " \t\r\n")
}

// CreateAlias gera um alias unico que reencaminha para destination.
func (r *Repo) CreateAlias(ctx context.Context, ownerID, destination, label string) (*Alias, error) {
	destination = strings.TrimSpace(destination)
	if !looksLikeEmail(destination) {
		return nil, ErrInvalidDest
	}
	// Tenta algumas vezes em caso (raro) de colisao do endereco gerado.
	for attempt := 0; attempt < 5; attempt++ {
		local, err := randomLocalPart()
		if err != nil {
			return nil, err
		}
		address := local + "@" + AliasDomain
		a := &Alias{OwnerID: ownerID, AliasAddress: address, Destination: destination, Label: label, Active: true}
		err = r.pool.QueryRow(ctx, `
			INSERT INTO email_aliases (owner_id, alias_address, destination, label)
			VALUES ($1, $2, $3, $4)
			RETURNING id::text, created_at`,
			ownerID, address, destination, label,
		).Scan(&a.ID, &a.CreatedAt)
		if err == nil {
			return a, nil
		}
		if isUniqueViolation(err) {
			continue
		}
		return nil, err
	}
	return nil, errAliasCollision
}

// isUniqueViolation deteta a violacao de UNIQUE do Postgres (SQLSTATE 23505).
func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "23505")
}

// ListAliases devolve os aliases do utilizador, mais recentes primeiro.
func (r *Repo) ListAliases(ctx context.Context, ownerID string) ([]Alias, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, alias_address, destination, label, active, created_at
		FROM email_aliases WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Alias
	for rows.Next() {
		var a Alias
		if err := rows.Scan(
			&a.ID, &a.OwnerID, &a.AliasAddress, &a.Destination, &a.Label, &a.Active, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// SetActive liga/desliga o reencaminhamento de um alias.
func (r *Repo) SetActive(ctx context.Context, ownerID, aliasID string, active bool) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE email_aliases SET active = $1 WHERE id = $2 AND owner_id = $3`,
		active, aliasID, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteAlias remove um alias.
func (r *Repo) DeleteAlias(ctx context.Context, ownerID, aliasID string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM email_aliases WHERE id = $1 AND owner_id = $2`, aliasID, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// GetActiveAliasByAddress resolve o dono de um alias activo pelo endereço completo.
// Usado pelo ingest MAIL-002 quando SMTP chega para *@aegis.email.
func (r *Repo) GetActiveAliasByAddress(ctx context.Context, aliasAddress string) (*Alias, error) {
	aliasAddress = strings.ToLower(strings.TrimSpace(aliasAddress))
	a := &Alias{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, owner_id::text, alias_address, destination, label, active, created_at
		FROM email_aliases WHERE lower(alias_address) = $1 AND active = true`,
		aliasAddress,
	).Scan(&a.ID, &a.OwnerID, &a.AliasAddress, &a.Destination, &a.Label, &a.Active, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

// GetAlias devolve um alias do utilizador (usado em testes/uso interno).
func (r *Repo) GetAlias(ctx context.Context, ownerID, aliasID string) (*Alias, error) {
	a := &Alias{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, owner_id::text, alias_address, destination, label, active, created_at
		FROM email_aliases WHERE id = $1 AND owner_id = $2`, aliasID, ownerID,
	).Scan(&a.ID, &a.OwnerID, &a.AliasAddress, &a.Destination, &a.Label, &a.Active, &a.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}
