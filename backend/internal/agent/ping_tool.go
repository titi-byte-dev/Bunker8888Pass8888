package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// PingTool — diagnóstico; confirma que o pipeline de tools funciona.
type PingTool struct{}

func NewPingTool() *PingTool { return &PingTool{} }

type pingInput struct {
	Message string `json:"message"`
}

func (PingTool) Descriptor() Descriptor {
	schema := json.RawMessage(`{
		"type": "object",
		"required": ["message"],
		"properties": {
			"message": { "type": "string", "minLength": 1, "maxLength": 200 }
		}
	}`)
	return Descriptor{
		Name:        "ping",
		Description: "Eco de diagnóstico para validar o sistema de tools.",
		InputSchema: schema,
		Permissions: []Permission{PermNone},
		Sensitive:   false,
	}
}

func (PingTool) Validate(input json.RawMessage) error {
	var in pingInput
	if err := DecodeInput(input, &in); err != nil {
		return err
	}
	in.Message = strings.TrimSpace(in.Message)
	if in.Message == "" || len(in.Message) > 200 {
		return fmt.Errorf("%w: message inválida", ErrInvalidToolInput)
	}
	return nil
}

func (PingTool) Execute(_ context.Context, req Request) (json.RawMessage, error) {
	var in pingInput
	_ = DecodeInput(req.Input, &in)
	out := map[string]string{
		"pong":    in.Message,
		"agent":   req.AgentID,
		"user_id": req.UserID,
	}
	return json.Marshal(out)
}
