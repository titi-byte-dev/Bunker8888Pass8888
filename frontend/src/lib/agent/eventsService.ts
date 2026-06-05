/**
 * Feed de eventos do Event Bus (AGENT-004).
 */
import { loadSessionToken } from "$lib/session";

export interface AgentEvent {
  id: string;
  type: string;
  source: string;
  payload: Record<string, unknown>;
  createdAt: string;
  label: string;
}

interface AgentEventDTO {
  id: string;
  type: string;
  source: string;
  payload: Record<string, unknown>;
  created_at: string;
}

function labelForType(type: string): string {
  switch (type) {
    case "mail.inbox.received":
      return "E-mail recebido no alias";
    case "crm.prospection.run":
      return "Prospeção executada";
    case "agent.tool.executed":
      return "Tool de agente executada";
    default:
      return type;
  }
}

function mapEvent(d: AgentEventDTO): AgentEvent {
  return {
    id: d.id,
    type: d.type,
    source: d.source,
    payload: d.payload ?? {},
    createdAt: d.created_at,
    label: labelForType(d.type),
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
