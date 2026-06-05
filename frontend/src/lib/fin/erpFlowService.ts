/**
 * Eventos do fluxo ERP dev (DoD Fase 3) — só metadados, sem PII.
 */
import { loadSessionToken } from "$lib/session";

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

async function post(path: string, body: Record<string, unknown>): Promise<void> {
  const res = await fetch(path, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
}

export async function reportInvoicePaid(invoiceId: string, invoiceNumber: string): Promise<void> {
  await post("/api/agent/finance/report-invoice-paid", { invoice_id: invoiceId, invoice_number: invoiceNumber });
}

export async function reportCommissionRecorded(invoiceId: string): Promise<void> {
  await post("/api/agent/hr/request-compliance", { invoice_id: invoiceId, reason: "erp_flow_complete" });
}
