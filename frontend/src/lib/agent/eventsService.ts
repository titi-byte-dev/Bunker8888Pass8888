/**
 * Feed de eventos do Event Bus (AGENT-004).
 */
import { loadSessionToken } from "$lib/session";

import type { ApprovalStatus } from "./approvalService";

export interface AgentEvent {
  id: string;
  type: string;
  source: string;
  payload: Record<string, unknown>;
  createdAt: string;
  label: string;
  /** AGENT-009 — pending até o utilizador decidir. */
  approvalStatus?: ApprovalStatus;
}

interface AgentEventDTO {
  id: string;
  type: string;
  source: string;
  payload: Record<string, unknown>;
  created_at: string;
  approval_status?: ApprovalStatus;
}

function labelForType(type: string, payload?: Record<string, unknown>): string {
  switch (type) {
    case "mail.inbox.received":
      return "E-mail recebido no alias";
    case "crm.prospection.run":
      return "Prospeção executada";
    case "agent.tool.executed":
      return "Tool de agente executada";
    case "orchestrator.action.suggested": {
      const action = payload?.action;
      if (action === "run_prospection") return "Sugestão: correr prospeção";
      if (action === "run_onboarding") return "Sugestão: completar onboarding";
      if (action === "create_purchase_order") return "Sugestão: ordem de compra";
      if (action === "screen_candidate") return "Sugestão: triagem às cegas";
      if (action === "review_saas_licenses") return "Sugestão: rever licenças SaaS";
      if (action === "reconcile_payments") return "Sugestão: reconciliar pagamentos";
      if (action === "issue_proforma") return "Sugestão: emitir pro-forma";
      if (action === "calculate_commission") return "Sugestão: calcular comissão";
      if (action === "generate_rgpd_report") return "Sugestão: relatório RGPD";
      return "Sugestão do orquestrador";
    }
    case "fin.subscription.stale":
      return "Licenças SaaS sem uso reportadas";
    case "fin.transactions.synced":
      return "Movimentos bancários sincronizados";
    case "crm.deal_closed":
      return "Negócio fechado no CRM";
    case "fin.invoice.paid":
      return "Fatura marcada como paga";
    case "hr.compliance.requested":
      return "Pedido de relatório RGPD";
    case "hr.recruitment.run":
      return "Triagem de candidatos executada";
    case "hr.employee.created":
      return "Ficha de empregado criada";
    case "ops.stock.low":
      return "Stock baixo no inventário";
    case "orchestrator.action.approved":
      return "Sugestão aprovada";
    case "orchestrator.action.rejected":
      return "Sugestão rejeitada";
    default:
      return type;
  }
}

function mapEvent(d: AgentEventDTO): AgentEvent {
  const approvalStatus = d.approval_status;
  let label = labelForType(d.type, d.payload);
  if (d.type === "orchestrator.action.suggested" && approvalStatus === "approved") {
    label = "Sugestão aprovada";
  } else if (d.type === "orchestrator.action.suggested" && approvalStatus === "rejected") {
    label = "Sugestão rejeitada";
  }
  return {
    id: d.id,
    type: d.type,
    source: d.source,
    payload: d.payload ?? {},
    createdAt: d.created_at,
    label,
    approvalStatus,
  };
}

export async function listAgentEvents(): Promise<AgentEvent[]> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão expirada — inicia sessão de novo.");
  const res = await fetch("/api/agent/events", {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  const j = (await res.json()) as { events: AgentEventDTO[] };
  return (j.events ?? []).map(mapEvent);
}
