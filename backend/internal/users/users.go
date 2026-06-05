// Package users guarda e lê utilizadores da base de dados.
package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound é devolvido quando um utilizador não existe.
var ErrNotFound = errors.New("users: utilizador não encontrado")

// ErrEmailTaken é devolvido quando o email já está registado.
var ErrEmailTaken = errors.New("users: email já registado")

// User espelha uma linha da tabela "users".
type User struct {
	ID           string
	Email        string
	Verifier     []byte
	VerifierSalt []byte
	KDFSalt      []byte
	KDFTime      int
	KDFMemory    int
	KDFThreads   int
}

// Repo dá acesso à tabela "users".
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria um repositório de utilizadores.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Create insere um novo utilizador e preenche u.ID com o id gerado.
func (r *Repo) Create(ctx context.Context, u *User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text`,
		u.Email, u.Verifier, u.VerifierSalt, u.KDFSalt, u.KDFTime, u.KDFMemory, u.KDFThreads,
	).Scan(&u.ID)
	if err != nil {
		// 23505 = unique_violation no PostgreSQL (email duplicado).
		if isUniqueViolation(err) {
			return ErrEmailTaken
		}
		return fmt.Errorf("inserir utilizador: %w", err)
	}
	return nil
}

// ByEmail procura um utilizador pelo email. Devolve ErrNotFound se não existir.
func (r *Repo) ByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads
		FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.Verifier, &u.VerifierSalt, &u.KDFSalt, &u.KDFTime, &u.KDFMemory, &u.KDFThreads)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

// Summary é a vista pública de um utilizador (sem segredos criptográficos).
type Summary struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// ListSummaries devolve todos os utilizadores (id + email) para painel admin.
func (r *Repo) ListSummaries(ctx context.Context) ([]Summary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, email, created_at
		FROM users
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("listar utilizadores: %w", err)
	}
	defer rows.Close()

	var out []Summary
	for rows.Next() {
		var s Summary
		var createdAt time.Time
		if err := rows.Scan(&s.ID, &s.Email, &createdAt); err != nil {
			return nil, err
		}
		s.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		out = append(out, s)
	}
	return out, rows.Err()
}

// ByID procura um utilizador pelo UUID. Devolve ErrNotFound se não existir.
func (r *Repo) ByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, email, verifier, verifier_salt, kdf_salt, kdf_time, kdf_memory, kdf_threads
		FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Verifier, &u.VerifierSalt, &u.KDFSalt, &u.KDFTime, &u.KDFMemory, &u.KDFThreads)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

// isUniqueViolation deteta o erro de chave única do PostgreSQL sem importar o
// pacote pgconn diretamente nos chamadores.
func isUniqueViolation(err error) bool {
	// O erro do pgx implementa um SQLState() string == "23505" em violações únicas.
	type sqlStater interface{ SQLState() string }
	var s sqlStater
	if errors.As(err, &s) {
		return s.SQLState() == "23505"
	}
	return false
}
