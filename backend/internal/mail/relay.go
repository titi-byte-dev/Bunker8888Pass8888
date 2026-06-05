package mail

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// RelayService reencaminha e-mail recebido no alias para o destino real (MAIL-004).
// Didático: em dev usamos Mailpit SMTP; em produção aponta para Postfix na VPS.
type RelayService struct {
	SMTPHost string // ex.: "mailpit:1025"
}

// RelayResult resume um envio SMTP.
type RelayResult struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ForwardInbound envia cópia do e-mail recebido para alias.Destination.
// O remetente visível é o alias — o Reply-To mantém o contacto original.
func (r *RelayService) ForwardInbound(_ context.Context, alias *Alias, originalFrom, subject, body string) (*RelayResult, error) {
	if r == nil || strings.TrimSpace(r.SMTPHost) == "" {
		return nil, nil
	}
	if alias == nil || !alias.Active {
		return nil, fmt.Errorf("mail: alias inactivo")
	}
	originalFrom = strings.TrimSpace(originalFrom)
	if originalFrom == "" {
		originalFrom = "unknown@unknown"
	}
	subject = strings.TrimSpace(subject)
	if subject == "" {
		subject = "(sem assunto)"
	}

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", alias.AliasAddress)
	fmt.Fprintf(&msg, "To: %s\r\n", alias.Destination)
	fmt.Fprintf(&msg, "Reply-To: %s\r\n", originalFrom)
	fmt.Fprintf(&msg, "Subject: [via %s] %s\r\n", alias.AliasAddress, subject)
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
	fmt.Fprintf(&msg, "--- Reencaminhado pelo AegisPass ---\r\n")
	fmt.Fprintf(&msg, "Alias: %s\r\n", alias.AliasAddress)
	fmt.Fprintf(&msg, "Remetente original: %s\r\n\r\n", originalFrom)
	msg.WriteString(body)

	host, port, err := net.SplitHostPort(r.SMTPHost)
	if err != nil {
		// Sem porta — assume 25 (didático; dev usa host:port explícito).
		host = r.SMTPHost
		port = "25"
	}
	addr := net.JoinHostPort(host, port)
	if err := smtp.SendMail(addr, nil, alias.AliasAddress, []string{alias.Destination}, msg.Bytes()); err != nil {
		return nil, fmt.Errorf("mail: relay SMTP: %w", err)
	}
	return &RelayResult{From: alias.AliasAddress, To: alias.Destination}, nil
}
