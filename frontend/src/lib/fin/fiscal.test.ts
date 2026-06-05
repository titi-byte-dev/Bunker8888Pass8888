import { describe, expect, it } from "vitest";
import {
  applyFiscalSuggestions,
  fiscalSummary,
  suggestFiscalCode,
  type FiscalCode,
} from "./fiscal";
import type { Subscription } from "./subscriptions";

function sub(over: Partial<Subscription> = {}): Subscription {
  return {
    id: "s1",
    name: "Figma",
    cost: 30,
    currency: "EUR",
    cycle: "monthly",
    category: "Design",
    active: true,
    ...over,
  };
}

describe("FIN-005 suggestFiscalCode", () => {
  it("SaaS design → dedutivel_100", () => {
    expect(suggestFiscalCode({ name: "Figma", category: "Design" })).toBe("dedutivel_100");
  });

  it("hardware → investimento", () => {
    expect(suggestFiscalCode({ name: "MacBook Pro", category: "" })).toBe("investimento");
  });

  it("multa → nao_dedutivel", () => {
    expect(suggestFiscalCode({ name: "Multa estacionamento", category: "" })).toBe("nao_dedutivel");
  });
});

describe("FIN-005 fiscalSummary", () => {
  it("soma dedutível mensal", () => {
    const summary = fiscalSummary([
      sub({ fiscalCode: "dedutivel_100" as FiscalCode, cost: 100 }),
      sub({ id: "s2", name: "X", fiscalCode: "nao_dedutivel" as FiscalCode, cost: 50 }),
    ]);
    expect(summary.totalMonthly).toBe(150);
    expect(summary.totalDeductibleMonthly).toBe(100);
    expect(summary.pendingCount).toBe(0);
  });
});

describe("FIN-005 applyFiscalSuggestions", () => {
  it("só pendentes activas", () => {
    const out = applyFiscalSuggestions([
      sub({ fiscalCode: undefined }),
      sub({ id: "s2", active: false, fiscalCode: undefined }),
      sub({ id: "s3", fiscalCode: "dedutivel_100" as FiscalCode }),
    ]);
    expect(out).toHaveLength(1);
    expect(out[0].fiscalCode).toBe("dedutivel_100");
  });
});
