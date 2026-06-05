// Package eventbus implementa o Event Bus in-process (AGENT-004).
//
// Didático: Event-Driven Architecture desacopla agentes — em vez de A chamar B
// directamente, A publica um evento e os subscritores reagem. Começamos com
// channels Go; NATS fica para escala multi-nó (futuro).
package eventbus

// Tipos de evento do domínio agentes (contrato estável para subscritores).
const (
	// MailInboxReceived — e-mail ingerido para um alias do utilizador.
	MailInboxReceived = "mail.inbox.received"
	// CRMProspectionRun — utilizador ou sistema pediu rascunhos de prospeção.
	CRMProspectionRun = "crm.prospection.run"
	// AgentToolExecuted — tool de agente executada (complementa guardian audit).
	AgentToolExecuted = "agent.tool.executed"
	// OrchestratorActionSuggested — orquestrador sugere acção (human-in-the-loop).
	OrchestratorActionSuggested = "orchestrator.action.suggested"
	// WildcardSubscribe recebe todos os tipos (só uso interno/debug).
	WildcardSubscribe = "*"
)
