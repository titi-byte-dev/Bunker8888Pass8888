// Package passkeys implementa WebAuthn (passkeys) para autenticação ao servidor.
//
// > 💡 **Conceito — Passkey:** credencial WebAuthn ligada a um dispositivo
// (biometria ou PIN). Substitui o envio do Auth Hash no login, mas NÃO substitui
// a Master Password — a Master Key continua a ser derivada localmente (Zero-Knowledge).
package passkeys

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// ErrNotFound indica credencial ou sessão inexistente.
var ErrNotFound = errors.New("passkeys: não encontrado")

// ErrSessionExpired indica que a sessão WebAuthn expirou (challenge TTL).
var ErrSessionExpired = errors.New("passkeys: sessão WebAuthn expirou")

// Credential registo na BD.
type Credential struct {
	ID              string
	UserID          string
	Name            string
	CredentialID    []byte
	PublicKey       []byte
	SignCount       uint32
	BackupEligible  bool
	BackupState     bool
	Transports      []string
	CreatedAt       string
}

// Repo acede à tabela webauthn_credentials.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Save insere uma credencial recém-registada.
func (r *Repo) Save(ctx context.Context, userID, name string, cred webauthn.Credential) error {
	transports := ""
	if len(cred.Transport) > 0 {
		b, _ := json.Marshal(cred.Transport)
		transports = string(b)
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO webauthn_credentials
			(user_id, credential_id, public_key, sign_count, backup_eligible, backup_state, transports, name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, cred.ID, cred.PublicKey, cred.Authenticator.SignCount,
		cred.Flags.BackupEligible, cred.Flags.BackupState, transports, name,
	)
	return err
}

// ListByUserID devolve credenciais WebAuthn do utilizador.
func (r *Repo) ListByUserID(ctx context.Context, userID string) ([]webauthn.Credential, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT credential_id, public_key, sign_count, backup_eligible, backup_state, transports
		FROM webauthn_credentials WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []webauthn.Credential
	for rows.Next() {
		var credID, pubKey []byte
		var signCount int64
		var backupEligible, backupState bool
		var transportsJSON *string
		if err := rows.Scan(&credID, &pubKey, &signCount, &backupEligible, &backupState, &transportsJSON); err != nil {
			return nil, err
		}
		var transports []protocol.AuthenticatorTransport
		if transportsJSON != nil && *transportsJSON != "" {
			var ts []string
			_ = json.Unmarshal([]byte(*transportsJSON), &ts)
			for _, t := range ts {
				transports = append(transports, protocol.AuthenticatorTransport(t))
			}
		}
		out = append(out, webauthn.Credential{
			ID:        credID,
			PublicKey: pubKey,
			Transport: transports,
			Flags: webauthn.CredentialFlags{
				BackupEligible: backupEligible,
				BackupState:    backupState,
			},
			Authenticator: webauthn.Authenticator{
				SignCount: uint32(signCount),
			},
		})
	}
	return out, rows.Err()
}

// ListMeta devolve metadados para a UI (sem chaves).
func (r *Repo) ListMeta(ctx context.Context, userID string) ([]Credential, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, user_id::text, name, created_at::text
		FROM webauthn_credentials WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Credential
	for rows.Next() {
		var c Credential
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpdateSignCount actualiza o contador anti-replay após login bem-sucedido.
func (r *Repo) UpdateSignCount(ctx context.Context, credentialID []byte, count uint32) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE webauthn_credentials SET sign_count = $2 WHERE credential_id = $1`,
		credentialID, count,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UserByCredentialID resolve o utilizador a partir do credential id (login).
func (r *Repo) UserByCredentialID(ctx context.Context, credentialID []byte) (string, error) {
	var userID string
	err := r.pool.QueryRow(ctx, `
		SELECT user_id::text FROM webauthn_credentials WHERE credential_id = $1`, credentialID,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return userID, nil
}

// webUser implementa webauthn.User para a biblioteca go-webauthn.
type webUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u webUser) WebAuthnID() []byte                          { return u.id }
func (u webUser) WebAuthnName() string                        { return u.name }
func (u webUser) WebAuthnDisplayName() string                 { return u.displayName }
func (u webUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

func newWebUser(u *users.User, creds []webauthn.Credential) (webUser, error) {
	parsed, err := uuid.Parse(u.ID)
	if err != nil {
		return webUser{}, err
	}
	return webUser{
		id:          parsed[:],
		name:        u.Email,
		displayName: u.Email,
		credentials: creds,
	}, nil
}

// sessionEntry guarda SessionData entre begin e finish.
type sessionEntry struct {
	data      webauthn.SessionData
	expiresAt time.Time
}

// SessionStore memória para desafios WebAuthn (TTL curto).
type SessionStore struct {
	mu    sync.Mutex
	items map[string]sessionEntry
	ttl   time.Duration
}

// NewSessionStore cria store com TTL (ex: 2 minutos).
func NewSessionStore(ttl time.Duration) *SessionStore {
	return &SessionStore{items: make(map[string]sessionEntry), ttl: ttl}
}

// Put guarda session data e devolve id opaco.
func (s *SessionStore) Put(data webauthn.SessionData) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.NewString()
	s.items[id] = sessionEntry{data: data, expiresAt: time.Now().Add(s.ttl)}
	return id
}

// Get recupera e remove session data (one-time use).
func (s *SessionStore) Get(id string) (webauthn.SessionData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.items[id]
	if !ok {
		return webauthn.SessionData{}, ErrNotFound
	}
	delete(s.items, id)
	if time.Now().After(entry.expiresAt) {
		return webauthn.SessionData{}, ErrSessionExpired
	}
	return entry.data, nil
}

// Service orquestra WebAuthn com users + repo.
type Service struct {
	wa     *webauthn.WebAuthn
	repo   *Repo
	users  *users.Repo
	store  *SessionStore
}

// Config para criar o serviço WebAuthn.
type Config struct {
	RPDisplayName string
	RPID          string
	RPOrigins     []string
}

// NewService constrói o serviço passkeys.
func NewService(cfg Config, repo *Repo, users *users.Repo) (*Service, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.RPDisplayName,
		RPID:          cfg.RPID,
		RPOrigins:     cfg.RPOrigins,
	})
	if err != nil {
		return nil, err
	}
	return &Service{
		wa:    wa,
		repo:  repo,
		users: users,
		store: NewSessionStore(2 * time.Minute),
	}, nil
}

// BeginRegistration inicia registo de passkey (utilizador já autenticado).
func (s *Service) BeginRegistration(ctx context.Context, userID string) (protocol.CredentialCreation, string, error) {
	_, _, wu, err := s.loadWebUser(ctx, userID)
	if err != nil {
		return protocol.CredentialCreation{}, "", err
	}
	options, session, err := s.wa.BeginRegistration(wu)
	if err != nil {
		return protocol.CredentialCreation{}, "", err
	}
	return *options, s.store.Put(*session), nil
}

// FinishRegistration conclui registo e persiste credencial.
func (s *Service) FinishRegistration(ctx context.Context, userID, sessionID, name string, body []byte) error {
	session, err := s.store.Get(sessionID)
	if err != nil {
		return err
	}
	_, _, wu, err := s.loadWebUser(ctx, userID)
	if err != nil {
		return err
	}
	parsed, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(body))
	if err != nil {
		return err
	}
	cred, err := s.wa.CreateCredential(wu, session, parsed)
	if err != nil {
		return err
	}
	return s.repo.Save(ctx, userID, name, *cred)
}

// BeginLogin inicia autenticação passkey para um email.
func (s *Service) BeginLogin(ctx context.Context, email string) (protocol.CredentialAssertion, string, error) {
	u, err := s.users.ByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			// Resposta genérica — não revelar se email existe.
			return protocol.CredentialAssertion{}, "", users.ErrNotFound
		}
		return protocol.CredentialAssertion{}, "", err
	}
	creds, err := s.repo.ListByUserID(ctx, u.ID)
	if err != nil {
		return protocol.CredentialAssertion{}, "", err
	}
	if len(creds) == 0 {
		return protocol.CredentialAssertion{}, "", ErrNotFound
	}
	wu, err := newWebUser(u, creds)
	if err != nil {
		return protocol.CredentialAssertion{}, "", err
	}
	options, session, err := s.wa.BeginLogin(wu)
	if err != nil {
		return protocol.CredentialAssertion{}, "", err
	}
	return *options, s.store.Put(*session), nil
}

// FinishLogin valida asserção e devolve user_id.
func (s *Service) FinishLogin(ctx context.Context, sessionID string, body []byte) (string, error) {
	session, err := s.store.Get(sessionID)
	if err != nil {
		return "", err
	}
	parsed, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	rawID := parsed.RawID
	userID, err := s.repo.UserByCredentialID(ctx, rawID)
	if err != nil {
		return "", err
	}
	u, err := s.users.ByID(ctx, userID)
	if err != nil {
		return "", err
	}
	creds, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	wu, err := newWebUser(u, creds)
	if err != nil {
		return "", err
	}
	cred, err := s.wa.ValidateLogin(wu, session, parsed)
	if err != nil {
		return "", err
	}
	if err := s.repo.UpdateSignCount(ctx, cred.ID, cred.Authenticator.SignCount); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Service) loadWebUser(ctx context.Context, userID string) (*users.User, []webauthn.Credential, webUser, error) {
	u, err := s.users.ByID(ctx, userID)
	if err != nil {
		return nil, nil, webUser{}, err
	}
	creds, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, nil, webUser{}, err
	}
	wu, err := newWebUser(u, creds)
	if err != nil {
		return nil, nil, webUser{}, err
	}
	return u, creds, wu, nil
}

// ListMeta devolve passkeys registadas do utilizador.
func (s *Service) ListMeta(ctx context.Context, userID string) ([]Credential, error) {
	return s.repo.ListMeta(ctx, userID)
}
