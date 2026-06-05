package mail

import (
	"context"
	"encoding/json"
	"testing"
)

func TestProcessSMTPIngest_InvalidPayload(t *testing.T) {
	svc := &IngestService{Aliases: &Repo{}, Inbox: &InboxRepo{}}
	_, err := svc.ProcessSMTPIngest(context.Background(), []byte(`{}`))
	if err != ErrInvalidWebhook {
		t.Fatalf("err=%v", err)
	}
}

func TestSMTPIngestPayload_JSON(t *testing.T) {
	raw := `{"message_id":"abc","from":"a@b.com","to":["x@aegis.email"],"subject":"Hi","body":"text"}`
	var p SMTPIngestPayload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatal(err)
	}
	if p.MessageID != "abc" || len(p.To) != 1 {
		t.Fatalf("payload=%+v", p)
	}
}
