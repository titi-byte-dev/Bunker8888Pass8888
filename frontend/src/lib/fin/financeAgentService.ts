/**
 * Agente financeiro — licenças SaaS (AGENT-006 / FIN-002).
 */
import { loadSessionToken } from "$lib/session";
import type { Alert } from "./alerts";
import { costSummary } from "./alerts";
import type { Subscription } from "./subscriptions";

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

/** Reporta alertas ao orquestrador — só envia IDs (sem nomes/custos em claro). */
export async function reportStaleAlerts(alerts: Alert[], subs: Subscription[]): Promise<void> {
  const ids = [...new Set(alerts.map((a) => a.subscriptionId))];
  const summary = costSummary(subs);
  const res = await fetch("/api/agent/finance/report-stale", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      subscription_ids: ids,
      alert_count: alerts.length,
      monthly_saving: summary.potentialMonthlySaving,
    }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
}
