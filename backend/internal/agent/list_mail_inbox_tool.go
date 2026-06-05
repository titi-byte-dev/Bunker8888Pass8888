package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/mail"
)

// MailInboxLister abstrai a leitura da caixa (testável sem BD).
type MailInboxLister interface {
	ListInbox(ctx context.Context, ownerID string, unprocessedOnly bool) ([]mail.InboxMessage, error)
}

// ListMailInboxTool expõe metadados de e-mail ao agente de prospeção (AGENT-003).
type ListMailInboxTool struct {
	Inbox MailInboxLister
}

func NewListMailInboxTool(inbox MailInboxLister) *ListMailInboxTool {
	return &ListMailInboxTool{Inbox: inbox}
}

type listMailInboxInput struct {
	UnprocessedOnly bool `json:"unprocessed_only"`
}

func (ListMailInboxTool) Descriptor() Descriptor {
	schema := json.RawMessage(`{
		"type": "object",
		"properties": {
			"unprocessed_only": { "type": "boolean", "default": true }
		}
	}`)
	return Descriptor{
		Name:        "list_mail_inbox",
		Description: "Lista mensagens da caixa de entrada (metadados + corpo para prospeção).",
		InputSchema: schema,
		Permissions: []Permission{PermMailReadMeta},
		Sensitive:   false,
	}
}

func (ListMailInboxTool) Validate(input json.RawMessage) error {
	var in listMailInboxInput
	return DecodeInput(input, &in)
}

func (t ListMailInboxTool) Execute(ctx context.Context, req Request) (json.RawMessage, error) {
	if t.Inbox == nil {
		return nil, fmt.Errorf("agent: inbox não configurada")
	}
	var in listMailInboxInput
	_ = DecodeInput(req.Input, &in)
	// Por omissão só pendentes — evita reprocessar leads já importados.
	unprocessed := true
	if len(req.Input) > 0 {
		unprocessed = in.UnprocessedOnly
	}
	msgs, err := t.Inbox.ListInbox(ctx, req.UserID, unprocessed)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		items = append(items, map[string]any{
			"id":          m.ID,
			"from_email":  m.FromEmail,
			"subject":     m.Subject,
			"body":        WrapExternalData("inbox:"+m.ID, m.Body),
			"received_at": m.ReceivedAt,
		})
	}
	out := map[string]any{
		"messages": items,
		"count":    len(items),
	}
	return json.Marshal(out)
}
