/**
 * Auditoria do Guardião (AGENT-002) — prova que tools correm sem Master Key.
 */
import { loadSessionToken } from "$lib/session";

export interface GuardianAuditEntry {
  id: string;
  agentId: string;
  toolName: string;
  success: boolean;
  errorMsg: string;
  createdAt: string;
}

export async function listGuardianAudit(): Promise<GuardianAuditEntry[]> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão expirada — inicia sessão de novo.");
  const res = await fetch("/api/agent/audit", {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  const j = (await res.json()) as {
    events: Array<{
      id: string;
      agent_id: string;
      tool_name: string;
      success: boolean;
      error_msg: string;
      created_at: string;
    }>;
  };
  return (j.events ?? []).map((e) => ({
    id: e.id,
    agentId: e.agent_id,
    toolName: e.tool_name,
    success: e.success,
    errorMsg: e.error_msg,
    createdAt: e.created_at,
  }));
}
