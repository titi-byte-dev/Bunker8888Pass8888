package burnnotes

import (
	"testing"
	"time"
)

// newTestStore devolve um store com relogio controlado a partir de t0.
func newTestStore(t0 time.Time) (*Store, *time.Time) {
	now := t0
	s := NewStore()
	s.now = func() time.Time { return now }
	return s, &now
}

func TestCreateAndBurnAfterRead(t *testing.T) {
	s, _ := newTestStore(time.Unix(0, 0))
	id, _, _, err := s.Create([]byte("nota-cifrada"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	ct, err := s.Consume(id)
	if err != nil || string(ct) != "nota-cifrada" {
		t.Fatalf("primeira leitura: ct=%q err=%v", ct, err)
	}
	// Arde apos a leitura: a segunda vez nao encontra nada.
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("segunda leitura: esperado ErrNotFound, got %v", err)
	}
	if s.Count() != 0 {
		t.Fatalf("store devia estar vazio, tem %d", s.Count())
	}
}

func TestManualBurnRequiresToken(t *testing.T) {
	s, _ := newTestStore(time.Unix(0, 0))
	id, token, _, err := s.Create([]byte("segredo"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// Token errado nao destroi.
	if err := s.Burn(id, "token-errado"); err != ErrBadToken {
		t.Fatalf("burn token errado: esperado ErrBadToken, got %v", err)
	}
	if s.Count() != 1 {
		t.Fatalf("nota nao devia ter sido destruida")
	}
	// Token certo destroi antes de qualquer leitura.
	if err := s.Burn(id, token); err != nil {
		t.Fatalf("burn token certo: %v", err)
	}
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("apos burn manual: esperado ErrNotFound, got %v", err)
	}
}

func TestBurnUnknownNote(t *testing.T) {
	s, _ := newTestStore(time.Unix(0, 0))
	if err := s.Burn("nao-existe", "x"); err != ErrNotFound {
		t.Fatalf("burn de id inexistente: esperado ErrNotFound, got %v", err)
	}
}

func TestExpiryByTime(t *testing.T) {
	s, now := newTestStore(time.Unix(0, 0))
	id, _, _, err := s.Create([]byte("efemera"), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	*now = now.Add(2 * time.Minute) // passou a expiracao
	if _, err := s.Consume(id); err != ErrNotFound {
		t.Fatalf("apos expirar: esperado ErrNotFound, got %v", err)
	}
}

func TestReapRemovesExpired(t *testing.T) {
	s, now := newTestStore(time.Unix(0, 0))
	if _, _, _, err := s.Create([]byte("a"), time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := s.Create([]byte("b"), time.Hour); err != nil {
		t.Fatal(err)
	}
	*now = now.Add(10 * time.Minute)
	if removed := s.Reap(*now); removed != 1 {
		t.Fatalf("reap devia remover 1, removeu %d", removed)
	}
	if s.Count() != 1 {
		t.Fatalf("devia sobrar 1 nota, ha %d", s.Count())
	}
}

func TestClampingAndValidation(t *testing.T) {
	s, _ := newTestStore(time.Unix(0, 0))

	if _, _, _, err := s.Create(nil, time.Hour); err != ErrEmpty {
		t.Fatalf("ciphertext vazio: esperado ErrEmpty, got %v", err)
	}
	if _, _, _, err := s.Create(make([]byte, MaxCiphertextBytes+1), time.Hour); err != ErrTooLarge {
		t.Fatalf("ciphertext grande: esperado ErrTooLarge, got %v", err)
	}

	// TTL 0 vira DefaultTTL; TTL gigante é cortado a MaxTTL.
	_, _, exp, err := s.Create([]byte("x"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if got := exp.Sub(time.Unix(0, 0)); got != DefaultTTL {
		t.Fatalf("TTL default: esperado %v, got %v", DefaultTTL, got)
	}
	_, _, exp2, err := s.Create([]byte("y"), 1000*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if got := exp2.Sub(time.Unix(0, 0)); got != MaxTTL {
		t.Fatalf("TTL maximo: esperado %v, got %v", MaxTTL, got)
	}
}
