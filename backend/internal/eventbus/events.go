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
	// HREmployeeCreated — ficha de empregado vazia criada (pedido de onboarding).
	HREmployeeCreated = "hr.employee.created"
	// OpsStockLow — inventário atingiu nível de reordenação (AGENT-008).
	OpsStockLow = "ops.stock.low"
	// HRRecruitmentRun — triagem às cegas executada (AGENT-007).
	HRRecruitmentRun = "hr.recruitment.run"
	// AgentToolExecuted — tool de agente executada (complementa guardian audit).
	AgentToolExecuted = "agent.tool.executed"
	// OrchestratorActionSuggested — orquestrador sugere acção (human-in-the-loop).
	OrchestratorActionSuggested = "orchestrator.action.suggested"
	// OrchestratorActionApproved — utilizador aprovou uma sugestão (AGENT-009).
	OrchestratorActionApproved = "orchestrator.action.approved"
	// OrchestratorActionRejected — utilizador rejeitou uma sugestão (AGENT-009).
	OrchestratorActionRejected = "orchestrator.action.rejected"
	// WildcardSubscribe recebe todos os tipos (só uso interno/debug).
	WildcardSubscribe = "*"
)
