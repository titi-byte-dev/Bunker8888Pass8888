/**
 * Alertas de licenças sem uso + agregados de custo (FIN-002).
 *
 * Funções puras sobre subscrições já decifradas — fáceis de testar e sem rede.
 * O "valor" do FIN-002 é cruzar custo com utilização: uma subscrição cara que
 * não se usa (ou sem login no cofre para confirmar uso) é dinheiro a escorrer.
 */
import type { Subscription } from "./subscriptions";

/** Dias sem uso a partir dos quais uma subscrição activa é "esquecida". */
export const UNUSED_THRESHOLD_DAYS = 60;

const DAY_MS = 24 * 60 * 60 * 1000;

/** Normaliza o custo de uma subscrição para o equivalente MENSAL. */
export function monthlyCost(sub: Subscription): number {
  return sub.cycle === "yearly" ? sub.cost / 12 : sub.cost;
}

/** Soma o custo mensal de todas as subscrições activas. */
export function totalMonthly(subs: Subscription[]): number {
  return subs.filter((s) => s.active).reduce((sum, s) => sum + monthlyCost(s), 0);
}

/** Soma o custo anual das subscrições activas. */
export function totalYearly(subs: Subscription[]): number {
  return totalMonthly(subs) * 12;
}

export type AlertKind = "stale" | "orphan";

export interface Alert {
  subscriptionId: string;
  name: string;
  kind: AlertKind;
  /** Quanto se poupa por mês se cancelar (custo mensal da subscrição). */
  monthlySaving: number;
  reason: string;
}

/** Dias decorridos desde lastUsedAt (Infinity se nunca usada). */
function daysSinceUse(sub: Subscription, now: number): number {
  if (!sub.lastUsedAt) return Infinity;
  const t = Date.parse(sub.lastUsedAt);
  if (Number.isNaN(t)) return Infinity;
  return (now - t) / DAY_MS;
}

/**
 * Gera alertas para subscrições ACTIVAS:
 *  - "stale"  → sem uso há mais de UNUSED_THRESHOLD_DAYS (ou nunca usada);
 *  - "orphan" → sem login associado no cofre (não dá para confirmar uso).
 * Uma subscrição pode disparar os dois.
 */
export function detectAlerts(
  subs: Subscription[],
  now: number = Date.now(),
  thresholdDays: number = UNUSED_THRESHOLD_DAYS,
): Alert[] {
  const out: Alert[] = [];
  for (const sub of subs) {
    if (!sub.active) continue;
    const saving = monthlyCost(sub);
    if (daysSinceUse(sub, now) > thresholdDays) {
      out.push({
        subscriptionId: sub.id,
        name: sub.name,
        kind: "stale",
        monthlySaving: saving,
        reason: sub.lastUsedAt
          ? `Sem uso há mais de ${thresholdDays} dias.`
          : "Nunca foi registada utilização.",
      });
    }
    if (!sub.vaultItemId) {
      out.push({
        subscriptionId: sub.id,
        name: sub.name,
        kind: "orphan",
        monthlySaving: saving,
        reason: "Sem login associado no cofre — utilização não confirmável.",
      });
    }
  }
  return out;
}

export interface CostSummary {
  monthly: number;
  yearly: number;
  activeCount: number;
  inactiveCount: number;
  /** Poupança potencial/mês somando subscrições "stale" (sem duplicar). */
  potentialMonthlySaving: number;
}

/** Resumo de custos + poupança potencial para o dashboard. */
export function costSummary(
  subs: Subscription[],
  now: number = Date.now(),
  thresholdDays: number = UNUSED_THRESHOLD_DAYS,
): CostSummary {
  const alerts = detectAlerts(subs, now, thresholdDays);
  const staleIds = new Set(alerts.filter((a) => a.kind === "stale").map((a) => a.subscriptionId));
  const potential = subs
    .filter((s) => staleIds.has(s.id))
    .reduce((sum, s) => sum + monthlyCost(s), 0);
  return {
    monthly: totalMonthly(subs),
    yearly: totalYearly(subs),
    activeCount: subs.filter((s) => s.active).length,
    inactiveCount: subs.filter((s) => !s.active).length,
    potentialMonthlySaving: potential,
  };
}
