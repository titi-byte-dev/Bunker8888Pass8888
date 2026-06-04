package passkeys

import (
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

func TestSessionStore_PutGet(t *testing.T) {
	s := NewSessionStore(time.Minute)
	id := s.Put(webauthn.SessionData{Challenge: "abc"})
	got, err := s.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if got.Challenge != "abc" {
		t.Fatalf("challenge=%q", got.Challenge)
	}
	_, err = s.Get(id)
	if err != ErrNotFound {
		t.Fatalf("reuse devia falhar, got %v", err)
	}
}

func TestSessionStore_Expired(t *testing.T) {
	s := NewSessionStore(-time.Second)
	id := s.Put(webauthn.SessionData{Challenge: "x"})
	_, err := s.Get(id)
	if err != ErrSessionExpired {
		t.Fatalf("got %v", err)
	}
}
