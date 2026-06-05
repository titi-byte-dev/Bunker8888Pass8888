// Package sharekeys guarda o par de chaves assimétricas de cada utilizador
// (SHARE-001), a fundação da partilha em Zero-Knowledge.
//
// ⚠️ O servidor trata ambos os campos como bytes opacos:
//   - PublicKey é partilhável (SPKI DER, em claro);
//   - WrappedPrivateKey foi cifrada no cliente com a Master Key do dono e nunca
//     é decifrada aqui (modelo Zero-Knowledge).
package sharekeys

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound indica que o utilizador ainda não tem par de chaves.
var ErrNotFound = errors.New("sharekeys: par de chaves não encontrado")

// Keypair espelha uma linha de "user_keypairs".
type Keypair struct {
	UserID            string
	PublicKey         []byte
	WrappedPrivateKey []byte
	Algorithm         string
}

// PublicKey é o subconjunto partilhável (sem a chave privada cifrada).
type PublicKey struct {
	UserID    string
	PublicKey []byte
	Algorithm string
}

// Repo acede à tabela user_keypairs.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Upsert grava ou actualiza o par de chaves do utilizador.
func (r *Repo) Upsert(ctx context.Context, userID string, publicKey, wrappedPrivateKey []byte, algorithm string) error {
	if len(publicKey) == 0 || len(wrappedPrivateKey) == 0 {
		return errors.New("sharekeys: chaves vazias")
	}
	if algorithm == "" {
		return errors.New("sharekeys: algoritmo em falta")
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_keypairs (user_id, public_key, wrapped_private_key, algorithm)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			public_key          = EXCLUDED.public_key,
			wrapped_private_key = EXCLUDED.wrapped_private_key,
			algorithm           = EXCLUDED.algorithm,
			updated_at          = now()`,
		userID, publicKey, wrappedPrivateKey, algorithm,
	)
	return err
}

// GetByUserID devolve o par completo do dono (inclui a chave privada cifrada).
func (r *Repo) GetByUserID(ctx context.Context, userID string) (*Keypair, error) {
	kp := &Keypair{UserID: userID}
	err := r.pool.QueryRow(ctx, `
		SELECT public_key, wrapped_private_key, algorithm
		FROM user_keypairs WHERE user_id = $1`, userID,
	).Scan(&kp.PublicKey, &kp.WrappedPrivateKey, &kp.Algorithm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return kp, nil
}

// GetPublicKeyByEmail devolve apenas a chave pública de um colega, para lhe
// partilhar um segredo. Nunca expõe a chave privada cifrada.
func (r *Repo) GetPublicKeyByEmail(ctx context.Context, email string) (*PublicKey, error) {
	pk := &PublicKey{}
	err := r.pool.QueryRow(ctx, `
		SELECT kp.user_id::text, kp.public_key, kp.algorithm
		FROM user_keypairs kp
		JOIN users u ON u.id = kp.user_id
		WHERE u.email = $1`, email,
	).Scan(&pk.UserID, &pk.PublicKey, &pk.Algorithm)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return pk, nil
}

// HasKeypair indica se o utilizador já configurou a partilha.
func (r *Repo) HasKeypair(ctx context.Context, userID string) (bool, error) {
	var ok bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM user_keypairs WHERE user_id = $1)`, userID,
	).Scan(&ok)
	return ok, err
}
