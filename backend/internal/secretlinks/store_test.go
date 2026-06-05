package secretlinks

import (
	"bytes"
	"testing"
	"time"
)

// newClockedStore devolve um store com um relógio controlado pelo teste.
func newClockedStore(now *time.Time) *Store {
	s := NewStore()
	s.now = func() time.Time { return *now }
	return s
}

func TestCreateAndConsumeOnce(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := newClockedStore(&now)

	id, expiresAt, err := s.Create([]byte("ct-cifrado"), time.Hour, 1)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("id vazio")
	}
	if !expiresAt.Equal(now.Add(time.Hour)) {
		t.Fatalf("expiresAt=%v, esperado %v", expiresAt, now.Add(time.Hour))
	}

	got, err := s.Consume(id)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte("ct-cifrado")) {
		t.Fatalf("ciphertext difere: %q", got)
	}

	// Segunda leitura: uso único esgotado -> apagado da RAM.
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("segundo consumo: esperado ErrNotFound, got %v", err)
	}
	if s.Count() != 0 {
		t.Fatalf("Count=%d, esperado 0 (sem rasto)", s.Count())
	}
}

func TestMultiViewThenExhausted(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := newClockedStore(&now)

	id, _, err := s.Create([]byte("x"), time.Hour, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if _, err := s.Consume(id); err != nil {
			t.Fatalf("consumo %d: %v", i+1, err)
		}
	}
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("4o consumo: esperado ErrNotFound, got %v", err)
	}
}

func TestExpiryByTime(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := newClockedStore(&now)

	id, _, err := s.Create([]byte("x"), time.Minute, 5)
	if err != nil {
		t.Fatal(err)
	}
	// Avança o relógio para lá da expiração.
	now = now.Add(2 * time.Minute)
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("consumo após expirar: esperado ErrNotFound, got %v", err)
	}
	if s.Count() != 0 {
		t.Fatalf("Count=%d, esperado 0 após expirar", s.Count())
	}
}

func TestReapRemovesExpired(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := newClockedStore(&now)

	if _, _, err := s.Create([]byte("curto"), MinTTL, 1); err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.Create([]byte("longo"), time.Hour, 1); err != nil {
		t.Fatal(err)
	}
	now = now.Add(30 * time.Second) // passa o MinTTL (10s), não a 1h
	if removed := s.Reap(now); removed != 1 {
		t.Fatalf("Reap removeu %d, esperado 1", removed)
	}
	if s.Count() != 1 {
		t.Fatalf("Count=%d, esperado 1", s.Count())
	}
}

func TestClampingAndValidation(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	s := newClockedStore(&now)

	// Vazio e demasiado grande são rejeitados.
	if _, _, err := s.Create(nil, time.Hour, 1); err != ErrEmpty {
		t.Fatalf("vazio: esperado ErrEmpty, got %v", err)
	}
	if _, _, err := s.Create(make([]byte, MaxCiphertextBytes+1), time.Hour, 1); err != ErrTooLarge {
		t.Fatalf("grande: esperado ErrTooLarge, got %v", err)
	}

	// TTL acima do máximo é cortado para MaxTTL; views<=0 vira 1.
	_, expiresAt, err := s.Create([]byte("x"), 72*time.Hour, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !expiresAt.Equal(now.Add(MaxTTL)) {
		t.Fatalf("TTL não cortado: expiresAt=%v, esperado %v", expiresAt, now.Add(MaxTTL))
	}
}

func TestConsumeUnknown(t *testing.T) {
	s := NewStore()
	if _, err := s.Consume("nao-existe"); err != ErrNotFound {
		t.Fatalf("id desconhecido: esperado ErrNotFound, got %v", err)
	}
}
