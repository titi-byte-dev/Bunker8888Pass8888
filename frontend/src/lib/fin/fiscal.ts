/**
 * Categorização fiscal (FIN-005) — calculada no cliente após decifrar subscrições.
 *
 * O código fiscal vive no blob ZK da subscrição; o servidor nunca vê categorias
 * nem totais dedutíveis — só bytes opacos.
 */
import type { Subscription, SubscriptionPayload } from "./subscriptions";
import { monthlyCost } from "./alerts";

/** Códigos fiscais simplificados (base IRC — PME Portugal). */
export type FiscalCode =
  | "pendente"
  | "dedutivel_100"
  | "dedutivel_50"
  | "nao_dedutivel"
  | "investimento"
  | "iva_recuperavel";

export interface FiscalCategory {
  code: FiscalCode;
  label: string;
  /** Percentagem típica de dedução em sede de IRC (didático, não aconselhamento fiscal). */
  deductPct: number;
  hint: string;
}

export const FISCAL_CATEGORIES: FiscalCategory[] = [
  {
    code: "pendente",
    label: "Por classificar",
    deductPct: 0,
    hint: "Aguarda revisão humana ou sugestão do agente.",
  },
  {
    code: "dedutivel_100",
    label: "Despesa dedutível 100%",
    deductPct: 100,
    hint: "Software SaaS, ferramentas de produtividade, serviços cloud.",
  },
  {
    code: "dedutivel_50",
    label: "Despesa dedutível 50%",
    deductPct: 50,
    hint: "Despesas mistas (ex.: refeições de negócio) — raro em SaaS.",
  },
  {
    code: "nao_dedutivel",
    label: "Não dedutível",
    deductPct: 0,
    hint: "Multas, donativos sem benefício fiscal, uso pessoal.",
  },
  {
    code: "investimento",
    label: "Investimento / ativo",
    deductPct: 0,
    hint: "Hardware ou licenças perpétuas amortizáveis — consultar contabilista.",
  },
  {
    code: "iva_recuperavel",
    label: "IVA recuperável (23%)",
    deductPct: 100,
    hint: "Compra com IVA dedutível — base para mapa de IVA.",
  },
];

export function fiscalLabel(code: FiscalCode | undefined): string {
  const c = FISCAL_CATEGORIES.find((x) => x.code === (code ?? "pendente"));
  return c?.label ?? "Por classificar";
}

/** Regras heurísticas — substituível por agente LLM com Guardião (AGENT-002). */
export function suggestFiscalCode(payload: Pick<SubscriptionPayload, "name" | "category">): FiscalCode {
  const name = payload.name.toLowerCase();
  const cat = (payload.category ?? "").toLowerCase();

  if (/multa|donativo|pessoal/.test(name)) return "nao_dedutivel";
  if (/hardware|servidor|macbook|laptop|impressora/.test(name)) return "investimento";
  if (/refeição|restaurante|hotel|viagem/.test(name + cat)) return "dedutivel_50";

  if (
    /saas|software|cloud|design|dev|crm|erp|notion|figma|github|google|microsoft|slack|zoom/.test(
      name + cat,
    )
  ) {
    return "dedutivel_100";
  }
  return "pendente";
}

export interface FiscalLine {
  subscriptionId: string;
  name: string;
  monthly: number;
  currency: string;
  fiscalCode: FiscalCode;
  deductPct: number;
  deductibleMonthly: number;
}

export interface FiscalSummary {
  lines: FiscalLine[];
  totalMonthly: number;
  totalDeductibleMonthly: number;
  pendingCount: number;
  currency: string;
}

/** Agrega custos mensais por código fiscal (após decifrar). */
export function fiscalSummary(subs: Subscription[]): FiscalSummary {
  const active = subs.filter((s) => s.active);
  const currency = active[0]?.currency ?? "EUR";
  const lines: FiscalLine[] = active.map((s) => {
    const code = s.fiscalCode ?? "pendente";
    const cat = FISCAL_CATEGORIES.find((c) => c.code === code) ?? FISCAL_CATEGORIES[0];
    const monthly = monthlyCost(s);
    const deductibleMonthly = (monthly * cat.deductPct) / 100;
    return {
      subscriptionId: s.id,
      name: s.name,
      monthly,
      currency: s.currency,
      fiscalCode: code,
      deductPct: cat.deductPct,
      deductibleMonthly,
    };
  });
  const totalMonthly = lines.reduce((a, l) => a + l.monthly, 0);
  const totalDeductibleMonthly = lines.reduce((a, l) => a + l.deductibleMonthly, 0);
  const pendingCount = lines.filter((l) => l.fiscalCode === "pendente").length;
  return { lines, totalMonthly, totalDeductibleMonthly, pendingCount, currency };
}

/** Aplica sugestões automáticas onde ainda está «pendente». */
export function applyFiscalSuggestions(
  subs: Subscription[],
): { id: string; fiscalCode: FiscalCode }[] {
  return subs
    .filter((s) => s.active && (!s.fiscalCode || s.fiscalCode === "pendente"))
    .map((s) => ({ id: s.id, fiscalCode: suggestFiscalCode(s) }));
}
