/**
 * Human-in-the-loop (AGENT-009) — aprovar ou rejeitar sugestões do orquestrador.
 */
import { loadSessionToken } from "$lib/session";

export type ApprovalStatus = "pending" | "approved" | "rejected";

export interface ApprovalResult {
  status: ApprovalStatus;
  action: string;
  suggestion: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

async function decide(suggestionId: string, path: "approve" | "reject"): Promise<ApprovalResult> {
  const res = await fetch(`/api/agent/orchestrator/actions/${encodeURIComponent(suggestionId)}/${path}`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token()}` },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return (await res.json()) as ApprovalResult;
}

/** Regista aprovação humana — o cliente executa a acção ZK em seguida. */
export async function approveSuggestion(suggestionId: string): Promise<ApprovalResult> {
  return decide(suggestionId, "approve");
}

/** Rejeita sugestão — fica auditada no feed. */
export async function rejectSuggestion(suggestionId: string): Promise<ApprovalResult> {
  return decide(suggestionId, "reject");
}
