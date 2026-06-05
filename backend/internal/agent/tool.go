// Package agent define o sistema de tools para agentes de IA (AGENT-001).
//
// Didático: em vez do LLM tocar na base de dados, ele pede para chamar uma
// "tool" com argumentos JSON; o nosso Go valida e executa com segurança.
package agent

import (
	"context"
	"encoding/json"
	"errors"
)

// Permission descreve o que uma tool precisa de aceder (menor privilégio).
type Permission string

const (
	PermNone           Permission = "none"
	PermMailReadMeta   Permission = "mail:read_metadata"
	PermCRMWriteLead           Permission = "crm:write_lead_draft"
	PermHRWriteCandidateDraft  Permission = "hr:write_candidate_draft"
)

// Descriptor expõe metadados para function-calling / auditoria.
type Descriptor struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
	Permissions []Permission    `json:"permissions"`
	// Sensitive indica acções que no futuro exigem human-in-the-loop (AGENT-009).
	Sensitive bool `json:"sensitive"`
}

// Request contextualiza uma execução de tool (tenant + utilizador + agente).
type Request struct {
	TenantID string
	UserID   string
	AgentID  string
	Input    json.RawMessage
}

// Tool é o contrato reimplementado de raiz (ver epic-agents).
type Tool interface {
	Descriptor() Descriptor
	Validate(input json.RawMessage) error
	Execute(ctx context.Context, req Request) (json.RawMessage, error)
}

var (
	ErrToolNotFound      = errors.New("agent: tool não encontrada")
	ErrDuplicateTool     = errors.New("agent: tool duplicada")
	ErrPermissionDenied  = errors.New("agent: permissão negada")
	ErrInvalidToolInput  = errors.New("agent: input inválido")
	ErrExternalAsCommand = errors.New("agent: conteúdo externo parece instrução")
)
