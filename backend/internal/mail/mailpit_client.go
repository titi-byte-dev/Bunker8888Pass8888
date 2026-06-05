package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MailpitClient busca o corpo completo de mensagens via API REST do Mailpit.
// Didático: o webhook só envia um "snippet"; para prospeção precisamos do Text.
type MailpitClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewMailpitClient(baseURL string) *MailpitClient {
	return &MailpitClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}
}

type mailpitAddress struct {
	Address string `json:"Address"`
}

type mailpitSummary struct {
	ID      string           `json:"ID"`
	From    mailpitAddress   `json:"From"`
	To      []mailpitAddress `json:"To"`
	Subject string           `json:"Subject"`
	Snippet string           `json:"Snippet"`
}

type mailpitMessage struct {
	Text    string `json:"Text"`
	HTML    string `json:"HTML"`
	Snippet string `json:"Snippet"`
}

// FetchBody obtém texto plano (ou snippet) de uma mensagem Mailpit.
func (c *MailpitClient) FetchBody(ctx context.Context, messageID string) (string, error) {
	if c == nil || c.BaseURL == "" {
		return "", fmt.Errorf("mail: mailpit não configurado")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/api/v1/message/"+messageID, nil)
	if err != nil {
		return "", err
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mail: mailpit HTTP %d", res.StatusCode)
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var msg mailpitMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return "", err
	}
	if t := strings.TrimSpace(msg.Text); t != "" {
		return t, nil
	}
	if h := strings.TrimSpace(msg.HTML); h != "" {
		return h, nil
	}
	return strings.TrimSpace(msg.Snippet), nil
}

// ParseMailpitSummary decodifica o JSON do webhook Mailpit.
func ParseMailpitSummary(raw []byte) (*mailpitSummary, error) {
	var s mailpitSummary
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// RecipientAddresses devolve todos os endereços To (normalizados em minúsculas).
func (s *mailpitSummary) RecipientAddresses() []string {
	out := make([]string, 0, len(s.To))
	for _, t := range s.To {
		addr := strings.ToLower(strings.TrimSpace(t.Address))
		if addr != "" {
			out = append(out, addr)
		}
	}
	return out
}

func (s *mailpitSummary) FromAddress() string {
	return strings.ToLower(strings.TrimSpace(s.From.Address))
}
