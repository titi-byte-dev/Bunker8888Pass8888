/**
 * Acções de remediação quando passwords estão expostas (DW-002).
 */
import type { HygieneSummary } from "$lib/vault/hygiene";
import type { BreachCheckResult } from "./breach";

export type RemediationItem = {
  itemId: string;
  title: string;
  reason: "breach" | "weak_and_breach";
};

/** Itens que exigem alteração urgente de password. */
export function itemsRequiringPasswordChange(
  summary: HygieneSummary,
  breachByItem: Map<string, BreachCheckResult>,
): RemediationItem[] {
  const out: RemediationItem[] = [];
  for (const item of summary.items) {
    const breach = breachByItem.get(item.itemId);
    if (!breach?.breached) continue;
    out.push({
      itemId: item.itemId,
      title: item.title,
      reason: item.issues.includes("weak") ? "weak_and_breach" : "breach",
    });
  }
  return out;
}

/** URL de edição com contexto de remediação (DW-002). */
export function remediationEditUrl(itemId: string, reason: RemediationItem["reason"] = "breach"): string {
  return `/vault/${itemId}/edit?remediate=${reason}`;
}
