/**
 * Painel de saúde de segurança (DW-003).
 *
 * Didático: combinamos o score de higiene (VAULT-008) com o resultado de
 * fugas (DW-001) num único número 0–100. O histórico fica só no cliente
 * (localStorage) — o servidor nunca vê passwords.
 */
import type { HygieneSummary } from "$lib/vault/hygiene";
import type { BreachCheckResult } from "./breach";

export type SecurityHealthSnapshot = {
  at: string;
  hygieneScore: number;
  exposedCount: number;
  weakCount: number;
  reusedCount: number;
  /** Score composto 0–100 */
  compositeScore: number;
};

export type SecurityHealthReport = SecurityHealthSnapshot & {
  trend: "up" | "down" | "flat" | "unknown";
  trendDelta: number;
  totalLogins: number;
  scannedCount: number;
};

const HISTORY_KEY = "aegis:security-health-history";
const HISTORY_MAX = 12;

/** Score composto — penaliza exposições em fugas mais que fraqueza local. */
export function computeCompositeScore(
  hygieneScore: number,
  exposedCount: number,
  weakCount: number,
  reusedCount: number,
  totalLogins: number,
): number {
  if (totalLogins === 0) return 100;
  let score = hygieneScore;
  score -= Math.min(40, exposedCount * 20);
  score -= Math.min(15, weakCount * 3);
  score -= Math.min(10, reusedCount * 2);
  return clamp(Math.round(score), 0, 100);
}

export function buildHealthReport(
  summary: HygieneSummary,
  breachByItem: Map<string, BreachCheckResult>,
): SecurityHealthReport {
  const exposedCount = [...breachByItem.values()].filter((b) => b.breached).length;
  const compositeScore = computeCompositeScore(
    summary.overallScore,
    exposedCount,
    summary.weakCount,
    summary.reusedCount,
    summary.totalLogins,
  );

  const snapshot: SecurityHealthSnapshot = {
    at: new Date().toISOString(),
    hygieneScore: summary.overallScore,
    exposedCount,
    weakCount: summary.weakCount,
    reusedCount: summary.reusedCount,
    compositeScore,
  };

  const history = loadHealthHistory();
  const previous = history[0];
  let trend: SecurityHealthReport["trend"] = "unknown";
  let trendDelta = 0;
  if (previous) {
    trendDelta = compositeScore - previous.compositeScore;
    if (trendDelta > 0) trend = "up";
    else if (trendDelta < 0) trend = "down";
    else trend = "flat";
  }

  return {
    ...snapshot,
    trend,
    trendDelta,
    totalLogins: summary.totalLogins,
    scannedCount: breachByItem.size,
  };
}

export function loadHealthHistory(): SecurityHealthSnapshot[] {
  if (typeof localStorage === "undefined") return [];
  try {
    const raw = localStorage.getItem(HISTORY_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as SecurityHealthSnapshot[];
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

export function saveHealthSnapshot(snapshot: SecurityHealthSnapshot): void {
  if (typeof localStorage === "undefined") return;
  const history = loadHealthHistory().filter((h) => h.at !== snapshot.at);
  history.unshift(snapshot);
  localStorage.setItem(HISTORY_KEY, JSON.stringify(history.slice(0, HISTORY_MAX)));
}

export function healthGrade(score: number): "good" | "warn" | "bad" {
  if (score >= 75) return "good";
  if (score >= 50) return "warn";
  return "bad";
}

function clamp(n: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, n));
}
