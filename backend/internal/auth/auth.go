// Package auth implementa o registo e o login por Auth Hash + sessões.
//
// Recordatório do modelo Zero-Knowledge:
//   - O cliente deriva a Master Key (fica no cliente) e um Auth Hash a partir da
//     Master Password.
//   - O cliente envia apenas o Auth Hash. O servidor NUNCA vê a password nem a
//     Master Key.
//   - O servidor guarda um hash do Auth Hash (verifier), nunca o Auth Hash cru.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sessions"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

// ErrInvalidCredentials é devolvido tanto para email inexistente como para
// Auth Hash errado — de propósito, para não revelar quais emails existem.
var ErrInvalidCredentials = errors.New("auth: credenciais inválidas")

// ClientKDF são os parâmetros/salt da KDF escolhidos pelo cliente.
type ClientKDF struct {
	Salt    []byte
	Time    int
	Memory  int
	Threads int
}

// Service junta as dependências necessárias à autenticação.
type Service struct {
	users      *users.Repo
	sessions   *sessions.Repo
	verifyKDF  crypto.KDFParams // KDF server-side aplicada ao Auth Hash recebido
	sessionTTL int              // segundos
}

// NewService constrói o serviço de autenticação.
func NewService(u *users.Repo, s *sessions.Repo, sessionTTLSeconds int) *Service {
	return &Service{
		users:      u,
		sessions:   s,
		verifyKDF:  crypto.DefaultKDFParams(),
		sessionTTL: sessionTTLSeconds,
	}
}

// Register cria um utilizador. `authHash` é o valor derivado no cliente.
func (s *Service) Register(ctx context.Context, email string, authHash []byte, kdf ClientKDF) error {
	// Geramos um salt próprio do servidor e guardamos Argon2id(authHash, salt).
	verifierSalt, err := crypto.GenerateSalt()
	if err != nil {
		return err
	}
	verifier := crypto.DeriveMasterKey(authHash, verifierSalt, s.verifyKDF)

	return s.users.Create(ctx, &users.User{
		Email:        email,
		Verifier:     verifier,
		VerifierSalt: verifierSalt,
		KDFSalt:      kdf.Salt,
		KDFTime:      kdf.Time,
		KDFMemory:    kdf.Memory,
		KDFThreads:   kdf.Threads,
	})
}

// KDFParamsFor devolve os parâmetros KDF do cliente para um email, necessários
// para o cliente re-derivar o Auth Hash antes do login.
func (s *Service) KDFParamsFor(ctx context.Context, email string) (ClientKDF, error) {
	u, err := s.users.ByEmail(ctx, email)
	if err != nil {
		return ClientKDF{}, err
	}
	return ClientKDF{Salt: u.KDFSalt, Time: u.KDFTime, Memory: u.KDFMemory, Threads: u.KDFThreads}, nil
}

// Login valida o Auth Hash e, em caso de sucesso, cria uma sessão e devolve o
// token em claro (que só o cliente vê; a BD guarda apenas o seu hash).
func (s *Service) Login(ctx context.Context, email string, authHash []byte) (string, error) {
	userID, err := s.ValidateCredentials(ctx, email, authHash)
	if err != nil {
		return "", err
	}
	return s.CreateSessionForUser(ctx, userID)
}

// ValidateCredentials confirma email + Auth Hash e devolve o id do utilizador.
func (s *Service) ValidateCredentials(ctx context.Context, email string, authHash []byte) (string, error) {
	u, err := s.users.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			dummySalt, _ := crypto.GenerateSalt()
			_ = crypto.DeriveMasterKey(authHash, dummySalt, s.verifyKDF)
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	candidate := crypto.DeriveMasterKey(authHash, u.VerifierSalt, s.verifyKDF)
	if !crypto.ConstantTimeEqual(candidate, u.Verifier) {
		return "", ErrInvalidCredentials
	}
	return u.ID, nil
}

// CreateSessionForUser emite um token de sessão para um utilizador já autenticado.
func (s *Service) CreateSessionForUser(ctx context.Context, userID string) (string, error) {
	token, err := newToken()
	if err != nil {
		return "", err
	}
	if err := s.sessions.Create(ctx, hashToken(token), userID, s.sessionTTL); err != nil {
		return "", err
	}
	return token, nil
}

// Authenticate valida um token de sessão e devolve o id do utilizador.
func (s *Service) Authenticate(ctx context.Context, token string) (string, error) {
	return s.sessions.UserIDByToken(ctx, hashToken(token))
}

// Logout invalida a sessão associada ao token.
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessions.Delete(ctx, hashToken(token))
}

// ListSessions devolve metadados das sessões activas do utilizador.
func (s *Service) ListSessions(ctx context.Context, userID string) ([]sessions.SessionInfo, error) {
	return s.sessions.ListByUser(ctx, userID)
}

// RevokeSession invalida uma sessão concreta (id opaco em hex).
func (s *Service) RevokeSession(ctx context.Context, userID, sessionIDHex string) error {
	return s.sessions.DeleteByIDForUser(ctx, userID, sessionIDHex)
}

// RevokeOtherSessions mantém só o token actual; devolve quantas sessões foram revogadas.
func (s *Service) RevokeOtherSessions(ctx context.Context, userID, currentToken string) (int64, error) {
	return s.sessions.DeleteAllForUserExcept(ctx, userID, hashToken(currentToken))
}

// SessionIDForToken devolve o identificador opaco da sessão actual (hex do hash).
func (s *Service) SessionIDForToken(token string) string {
	return hex.EncodeToString(hashToken(token))
}

// newToken gera um token de sessão aleatório (256 bits) em hexadecimal.
func newToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

// hashToken devolve o SHA-256 do token. Guardamos só o hash: rápido de verificar
// e, se a BD vazar, os tokens não são reutilizáveis.
func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}
