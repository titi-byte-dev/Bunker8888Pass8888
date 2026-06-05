/**
 * Open Banking scaffold (FIN-003) — mock provider em dev.
 */
import { loadSessionToken } from "$lib/session";
import type { BankTransaction } from "./reconcile";

export type ConnectionStatus = "pending" | "connected" | "expired";

export interface BankConnection {
  id: string;
  provider: string;
  status: ConnectionStatus;
  consentExpiresAt?: string;
  lastSyncAt?: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

async function authed(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return res;
}

export async function getBankingStatus(): Promise<BankConnection> {
  const j = (await (await authed("/api/fin/banking/status")).json()) as {
    connection: BankConnection;
  };
  return j.connection;
}

export async function connectBank(): Promise<BankConnection> {
  const j = (await (await authed("/api/fin/banking/connect", { method: "POST" })).json()) as {
    connection: BankConnection;
  };
  return j.connection;
}

export async function syncTransactions(): Promise<BankTransaction[]> {
  const j = (await (await authed("/api/fin/banking/sync", { method: "POST" })).json()) as {
    transactions: Array<{
      id: string;
      amount: number;
      currency: string;
      booked_at: string;
      description: string;
      merchant_ref?: string;
    }>;
  };
  return (j.transactions ?? []).map((t) => ({
    id: t.id,
    amount: t.amount,
    currency: t.currency,
    bookedAt: t.booked_at,
    description: t.description,
    merchantRef: t.merchant_ref,
  }));
}

/** Reporta contagens de reconciliação ao orquestrador — sem movimentos em claro. */
export async function reportTransactionsSynced(
  transactionCount: number,
  matchedCount: number,
  unmatchedCount: number,
): Promise<void> {
  await authed("/api/agent/finance/report-sync", {
    method: "POST",
    body: JSON.stringify({
      transaction_count: transactionCount,
      matched_count: matchedCount,
      unmatched_count: unmatchedCount,
    }),
  });
}
