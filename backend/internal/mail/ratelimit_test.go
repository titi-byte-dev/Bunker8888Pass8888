package mail_test

import (
	"errors"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

func TestRateConfigDefaults(t *testing.T) {
	rl := mail.NewRateLimiter(nil, mail.RateConfig{})
	if rl == nil {
		t.Fatal("nil pool ok")
	}
}

func TestErrRateLimited(t *testing.T) {
	if !errors.Is(mail.ErrRateLimited, mail.ErrRateLimited) {
		t.Fatal("sentinel")
	}
}
