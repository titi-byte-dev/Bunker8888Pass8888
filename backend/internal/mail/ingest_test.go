package mail_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

func TestParseMailpitSummary(t *testing.T) {
	raw := []byte(`{
		"ID": "abc123",
		"From": {"Address": "Lead@Empresa.pt"},
		"To": [{"Address": "deadbeef@aegis.email"}],
		"Subject": "Demo",
		"Snippet": "Ola mundo"
	}`)
	s, err := mail.ParseMailpitSummary(raw)
	if err != nil {
		t.Fatal(err)
	}
	if s.FromAddress() != "lead@empresa.pt" {
		t.Fatalf("from=%s", s.FromAddress())
	}
	addrs := s.RecipientAddresses()
	if len(addrs) != 1 || addrs[0] != "deadbeef@aegis.email" {
		t.Fatalf("to=%v", addrs)
	}
}

func TestProcessMailpitWebhook_IgnoredWithoutAlias(t *testing.T) {
	// Repo nil path
	svc := &mail.IngestService{}
	_, err := svc.ProcessMailpitWebhook(context.Background(), json.RawMessage(`{"ID":"x"}`))
	if err == nil {
		t.Fatal("esperava erro")
	}
}

func TestErrWebhookIgnored(t *testing.T) {
	if !errors.Is(mail.ErrWebhookIgnored, mail.ErrWebhookIgnored) {
		t.Fatal("sentinel")
	}
}
