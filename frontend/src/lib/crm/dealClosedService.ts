/**
 * Reporta deal_closed ao orquestrador (CRM-003 / DoD Fase 3).
 */
import { loadSessionToken } from "$lib/session";

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

export async function reportDealClosed(leadId: string): Promise<void> {
  const res = await fetch("/api/agent/crm/report-deal-closed", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ lead_id: leadId }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
}
