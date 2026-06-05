// Package burnnotes implementa as "notas auto-destrutivas" (SHARE-005): notas
// cifradas que ardem APOS A PRIMEIRA LEITURA (burn-after-read).
//
// Assenta no mesmo padrao efemero da SHARE-003 (secretlinks) — ciphertext so em
// RAM, chave no fragmento do URL — mas acrescenta duas ideias proprias:
//
//  1. BURN-AFTER-READ ESTRITO: cada nota lê-se exatamente uma vez; não há opção
//     de varias visualizacoes. Ao ler, a nota é removida da RAM.
//  2. BURN MANUAL (capacidade): quem cria recebe um burn_token secreto. Com ele
//     pode destruir a nota ANTES de ser lida (revogacao proativa). O token é
//     comparado em tempo constante para nao dar pistas por timing.
//
// A passphrase opcional (2.ª camada) é tratada inteiramente no cliente: a nota é
// duplamente cifrada e o servidor nunca sabe se existe passphrase — esse facto
// viaja no fragmento do URL, nao aqui. O servidor so guarda bytes opacos.
//
//	   criar                                   ler (1x, depois arde)
//	┌────────┐  ct opaco + burn_token   ┌────────┐  POST /notes/{id}
//	│ Autor  │ ───────────────────────▶ │  RAM   │ ◀──────────────── Leitor
//	└───┬────┘                          └────────┘  devolve ct, depois apaga
//	    │ guarda burn_token
//	    └── POST /notes/{id}/burn (token) ─▶ destroi antes de ser lida
package burnnotes

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// Limites de defesa (memoria e abuso). Conservadores de proposito.
const (
	// MaxCiphertextBytes limita o tamanho de cada nota cifrada (64 KiB).
	MaxCiphertextBytes = 64 * 1024
	// MaxTTL é o tempo de vida maximo de uma nota por ler (7 dias).
	MaxTTL = 7 * 24 * time.Hour
	// MinTTL evita notas instantaneamente expiradas.
	MinTTL = 10 * time.Second
	// DefaultTTL aplica-se quando o pedido nao indica TTL.
	DefaultTTL = 24 * time.Hour
	// MaxEntries limita o total de notas vivas em RAM (anti-exaustao de memoria).
	MaxEntries = 10_000
)

// Erros expostos. ErrNotFound cobre inexistente, expirada e já lida/queimada —
// uma resposta unica, para nao dar um oraculo sobre que ids existiram.
var (
	ErrNotFound = errors.New("burnnotes: nota inexistente, expirada ou ja destruida")
	ErrTooLarge = errors.New("burnnotes: ciphertext demasiado grande")
	ErrEmpty    = errors.New("burnnotes: ciphertext vazio")
	ErrFull     = errors.New("burnnotes: limite de notas em memoria atingido")
	ErrBadToken = errors.New("burnnotes: token de destruicao invalido")
)

type note struct {
	ciphertext []byte
	burnToken  string
	expiresAt  time.Time
}

// Store guarda as notas em RAM. Seguro para uso concorrente.
type Store struct {
	mu    sync.Mutex
	notes map[string]*note
	// now permite injetar o relogio nos testes (tempo determinista).
	now func() time.Time
}

// NewStore cria um store vazio.
func NewStore() *Store {
	return &Store{
		notes: make(map[string]*note),
		now:   time.Now,
	}
}

// clampTTL corta o TTL pedido para [MinTTL, MaxTTL]; 0 usa o default.
func clampTTL(ttl time.Duration) time.Duration {
	switch {
	case ttl <= 0:
		return DefaultTTL
	case ttl < MinTTL:
		return MinTTL
	case ttl > MaxTTL:
		return MaxTTL
	default:
		return ttl
	}
}

// Create guarda uma nota cifrada e devolve o id, o burn_token (capacidade de
// destruicao) e a hora de expiracao.
func (s *Store) Create(ciphertext []byte, ttl time.Duration) (id, burnToken string, expiresAt time.Time, err error) {
	if len(ciphertext) == 0 {
		return "", "", time.Time{}, ErrEmpty
	}
	if len(ciphertext) > MaxCiphertextBytes {
		return "", "", time.Time{}, ErrTooLarge
	}
	ttl = clampTTL(ttl)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Aproveita para limpar expiradas antes de medir a lotacao.
	s.reapLocked(s.now())
	if len(s.notes) >= MaxEntries {
		return "", "", time.Time{}, ErrFull
	}

	id, err = randToken(16)
	if err != nil {
		return "", "", time.Time{}, err
	}
	burnToken, err = randToken(24)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Copia defensiva: nao guardamos a fatia do chamador.
	ct := make([]byte, len(ciphertext))
	copy(ct, ciphertext)

	expiresAt = s.now().Add(ttl)
	s.notes[id] = &note{ciphertext: ct, burnToken: burnToken, expiresAt: expiresAt}
	return id, burnToken, expiresAt, nil
}

// Consume devolve o ciphertext UMA vez e destroi imediatamente a nota (arde após
// a leitura). Notas expiradas sao tratadas como inexistentes.
func (s *Store) Consume(id string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, ok := s.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	if !s.now().Before(n.expiresAt) {
		delete(s.notes, id) // expirada: apaga e finge que nunca existiu
		return nil, ErrNotFound
	}
	delete(s.notes, id) // burn-after-read: sai sempre da RAM
	return n.ciphertext, nil
}

// Burn destroi uma nota antes de ser lida, mediante o burn_token correto. O
// token é comparado em tempo constante para nao revelar nada por timing.
func (s *Store) Burn(id, burnToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, ok := s.notes[id]
	if !ok {
		return ErrNotFound
	}
	if subtle.ConstantTimeCompare([]byte(n.burnToken), []byte(burnToken)) != 1 {
		return ErrBadToken
	}
	delete(s.notes, id)
	return nil
}

// Reap remove as notas expiradas e devolve quantas foram removidas.
func (s *Store) Reap(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.reapLocked(now)
}

func (s *Store) reapLocked(now time.Time) int {
	removed := 0
	for id, n := range s.notes {
		if !now.Before(n.expiresAt) {
			delete(s.notes, id)
			removed++
		}
	}
	return removed
}

// Count devolve o numero de notas vivas (util em testes/diagnostico).
func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.notes)
}

// randToken gera um token aleatorio URL-safe com n bytes de entropia.
func randToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
