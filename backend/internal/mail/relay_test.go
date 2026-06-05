package mail_test

import (
	"context"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

func TestRelayService_NilDisabled(t *testing.T) {
	var r *mail.RelayService
	out, err := r.ForwardInbound(context.Background(), &mail.Alias{
		AliasAddress: "a@aegis.email",
		Destination:  "u@gmail.com",
		Active:       true,
	}, "x@y.pt", "Hi", "body")
	if err != nil || out != nil {
		t.Fatalf("relay desactivado: out=%v err=%v", out, err)
	}
}

func TestRelayService_InactiveAlias(t *testing.T) {
	r := &mail.RelayService{SMTPHost: "localhost:1025"}
	_, err := r.ForwardInbound(context.Background(), &mail.Alias{Active: false}, "a@b.pt", "s", "b")
	if err == nil {
		t.Fatal("esperava erro para alias inactivo")
	}
}
