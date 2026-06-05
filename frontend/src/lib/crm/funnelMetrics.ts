/**
 * Métricas básicas do funil CRM (DoD Fase 2).
 * Funções puras sobre leads já decifrados — sem rede.
 */
import { LEAD_STAGES, type Lead, type LeadStage } from "./leads";

export interface FunnelMetrics {
  total: number;
  byStage: Record<LeadStage, number>;
  won: number;
  lost: number;
  open: number;
  conversionPct: number;
}

/** Conta leads por estágio e calcula conversão won / (won+lost). */
export function computeFunnelMetrics(leads: Lead[]): FunnelMetrics {
  const byStage = Object.fromEntries(LEAD_STAGES.map((s) => [s.id, 0])) as Record<LeadStage, number>;
  for (const l of leads) {
    byStage[l.stage] = (byStage[l.stage] ?? 0) + 1;
  }
  const won = byStage.won ?? 0;
  const lost = byStage.lost ?? 0;
  const closed = won + lost;
  return {
    total: leads.length,
    byStage,
    won,
    lost,
    open: leads.length - closed,
    conversionPct: closed > 0 ? Math.round((won / closed) * 100) : 0,
  };
}
