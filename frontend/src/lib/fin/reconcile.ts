/**
 * Reconciliação local movimentos ↔ subscrições (FIN-003 / AGENT-006).
 * Funções puras — os dados bancários nunca são enviados ao servidor em claro.
 */
import { monthlyCost } from "./alerts";
import type { Subscription } from "./subscriptions";

export interface BankTransaction {
  id: string;
  amount: number;
  currency: string;
  bookedAt: string;
  description: string;
  merchantRef?: string;
}

export interface ReconcileMatch {
  transactionId: string;
  subscriptionId: string;
  subscriptionName: string;
  amount: number;
}

export interface ReconcileResult {
  matched: ReconcileMatch[];
  unmatched: BankTransaction[];
}

const AMOUNT_TOLERANCE = 0.02;

/** Cruza débitos com custos mensais das subscrições activas. */
export function reconcileTransactions(
  txs: BankTransaction[],
  subs: Subscription[],
): ReconcileResult {
  const active = subs.filter((s) => s.active);
  const matched: ReconcileMatch[] = [];
  const unmatched: BankTransaction[] = [];
  const usedSub = new Set<string>();

  for (const tx of txs) {
    const debit = Math.abs(tx.amount);
    let found: ReconcileMatch | null = null;
    for (const sub of active) {
      if (usedSub.has(sub.id)) continue;
      const expected = monthlyCost(sub);
      if (Math.abs(debit - expected) <= AMOUNT_TOLERANCE) {
        found = {
          transactionId: tx.id,
          subscriptionId: sub.id,
          subscriptionName: sub.name,
          amount: tx.amount,
        };
        usedSub.add(sub.id);
        break;
      }
    }
    if (found) matched.push(found);
    else unmatched.push(tx);
  }
  return { matched, unmatched };
}
