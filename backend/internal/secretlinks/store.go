// Package secretlinks implementa os "secret links" efémeros (SHARE-003):
// ligações de uso único / com TTL para partilhar um segredo fora do sistema.
//
// ⚠️ SERVIDO VIA RAM — DECISÃO DE DESIGN CENTRAL:
//
// O ciphertext vive SÓ em memória deste processo (um mapa protegido por mutex).
// Nunca toca o disco: não há tabela, não há migração, não há log do conteúdo.
// Quando o link expira (tempo) ou esgota as visualizações, a entrada é removida
// do mapa e o segredo deixa de existir — não fica rasto em lado nenhum.
//
// A chave de cifra NUNCA chega ao servidor: viaja no fragmento (#) do URL, que o
// browser não envia. O servidor guarda bytes que não consegue decifrar (ZK).
//
//	     criar                              consumir (1x)
//	┌────────┐  ct = AES-GCM(k, segredo)  ┌────────┐  POST /links/{id}
//	│ Autor  │ ─────────────────────────▶ │  RAM   │ ◀──────────────── Destinatário
//	└────────┘  k vai no #fragmento        └────────┘  devolve ct, depois apaga
//	   link = /s/{id}#k=...                  TTL + max_views; reaper limpa expirados
package secretlinks

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// Limites de defesa (memória e abuso). Conservadores de propósito.
const (
	// MaxCiphertextBytes limita o tamanho de cada segredo cifrado (64 KiB).
	MaxCiphertextBytes = 64 * 1024
	// MaxTTL é o tempo de vida máximo de um link (24h). TTLs maiores são cortados.
	MaxTTL = 24 * time.Hour
	// MinTTL evita links instantaneamente expirados.
	MinTTL = 10 * time.Second
	// DefaultTTL aplica-se quando o pedido não indica TTL.
	DefaultTTL = time.Hour
	// MaxViewsCap limita as visualizações por link (a maioria será 1).
	MaxViewsCap = 100
	// MaxEntries limita o total de links vivos em RAM (anti-exaustão de memória).
	MaxEntries = 10_000
)

// Erros expostos. ErrNotFound cobre inexistente, expirado e já consumido — uma
// resposta única, para não dar um oráculo sobre que ids existiram.
var (
	ErrNotFound = errors.New("secretlinks: link inexistente, expirado ou ja utilizado")
	ErrTooLarge = errors.New("secretlinks: ciphertext demasiado grande")
	ErrEmpty    = errors.New("secretlinks: ciphertext vazio")
	ErrFull     = errors.New("secretlinks: limite de links em memoria atingido")
)

type entry struct {
	ciphertext []byte
	expiresAt  time.Time
	maxViews   int
	views      int
}

// Store guarda os links em RAM. Seguro para uso concorrente.
type Store struct {
	mu      sync.Mutex
	entries map[string]*entry
	// now permite injetar o relógio nos testes (tempo determinista).
	now func() time.Time
}

// NewStore cria um store vazio.
func NewStore() *Store {
	return &Store{
		entries: make(map[string]*entry),
		now:     time.Now,
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

// clampViews corta as visualizações para [1, MaxViewsCap]; 0 vira 1 (uso único).
func clampViews(v int) int {
	switch {
	case v <= 0:
		return 1
	case v > MaxViewsCap:
		return MaxViewsCap
	default:
		return v
	}
}

// Create guarda um ciphertext e devolve o id e a hora de expiração.
func (s *Store) Create(ciphertext []byte, ttl time.Duration, maxViews int) (string, time.Time, error) {
	if len(ciphertext) == 0 {
		return "", time.Time{}, ErrEmpty
	}
	if len(ciphertext) > MaxCiphertextBytes {
		return "", time.Time{}, ErrTooLarge
	}
	ttl = clampTTL(ttl)
	maxViews = clampViews(maxViews)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Aproveita para limpar expirados antes de medir a lotação.
	s.reapLocked(s.now())
	if len(s.entries) >= MaxEntries {
		return "", time.Time{}, ErrFull
	}

	id, err := newID()
	if err != nil {
		return "", time.Time{}, err
	}

	// Cópia defensiva: não guardamos a fatia do chamador.
	ct := make([]byte, len(ciphertext))
	copy(ct, ciphertext)

	expiresAt := s.now().Add(ttl)
	s.entries[id] = &entry{ciphertext: ct, expiresAt: expiresAt, maxViews: maxViews}
	return id, expiresAt, nil
}

// Consume devolve o ciphertext uma vez e contabiliza a visualização. Quando as
// visualizações se esgotam (ou se já expirou), a entrada é removida da RAM.
func (s *Store) Consume(id string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[id]
	if !ok {
		return nil, ErrNotFound
	}
	now := s.now()
	if !now.Before(e.expiresAt) {
		delete(s.entries, id) // expirado: apaga e finge que nunca existiu
		return nil, ErrNotFound
	}

	e.views++
	ct := e.ciphertext
	if e.views >= e.maxViews {
		delete(s.entries, id) // esgotou as visualizações: apaga já
	}
	return ct, nil
}

// Reap remove as entradas expiradas e devolve quantas foram removidas.
func (s *Store) Reap(now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.reapLocked(now)
}

func (s *Store) reapLocked(now time.Time) int {
	removed := 0
	for id, e := range s.entries {
		if !now.Before(e.expiresAt) {
			delete(s.entries, id)
			removed++
		}
	}
	return removed
}

// Count devolve o número de links vivos (útil em testes/diagnóstico).
func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

// newID gera um id aleatório URL-safe (16 bytes => 22 chars base64url).
func newID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
