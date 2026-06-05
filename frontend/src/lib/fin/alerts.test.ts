import { describe, expect, it } from "vitest";
import {
  costSummary,
  detectAlerts,
  monthlyCost,
  totalMonthly,
  totalYearly,
  UNUSED_THRESHOLD_DAYS,
} from "./alerts";
import type { Subscription } from "./subscriptions";

const NOW = Date.parse("2026-06-05T00:00:00Z");
const DAY = 24 * 60 * 60 * 1000;

function sub(over: Partial<Subscription> = {}): Subscription {
  return {
    id: over.id ?? "s1",
    name: over.name ?? "Serviço",
    cost: over.cost ?? 10,
    currency: over.currency ?? "EUR",
    cycle: over.cycle ?? "monthly",
    active: over.active ?? true,
    vaultItemId: over.vaultItemId,
    lastUsedAt: over.lastUsedAt,
    category: over.category,
    vaultItemTitle: over.vaultItemTitle,
  };
}

describe("FIN-002 custos e alertas", () => {
  it("normaliza custo anual para mensal", () => {
    expect(monthlyCost(sub({ cost: 120, cycle: "yearly" }))).toBe(10);
    expect(monthlyCost(sub({ cost: 10, cycle: "monthly" }))).toBe(10);
  });

  it("soma só subscrições activas", () => {
    const subs = [
      sub({ id: "a", cost: 10 }),
      sub({ id: "b", cost: 120, cycle: "yearly" }), // 10/mes
      sub({ id: "c", cost: 99, active: false }), // ignorada
    ];
    expect(totalMonthly(subs)).toBe(20);
    expect(totalYearly(subs)).toBe(240);
  });

  it("marca como stale uma subscrição activa sem uso recente", () => {
    const old = new Date(NOW - (UNUSED_THRESHOLD_DAYS + 5) * DAY).toISOString();
    const alerts = detectAlerts([sub({ id: "a", vaultItemId: "v1", lastUsedAt: old })], NOW);
    expect(alerts).toHaveLength(1);
    expect(alerts[0].kind).toBe("stale");
    expect(alerts[0].monthlySaving).toBe(10);
  });

  it("não marca stale uma subscrição usada recentemente", () => {
    const recent = new Date(NOW - 3 * DAY).toISOString();
    const alerts = detectAlerts([sub({ id: "a", vaultItemId: "v1", lastUsedAt: recent })], NOW);
    expect(alerts).toHaveLength(0);
  });

  it("marca orphan uma subscrição sem login no cofre", () => {
    const recent = new Date(NOW - 1 * DAY).toISOString();
    const alerts = detectAlerts([sub({ id: "a", lastUsedAt: recent })], NOW);
    expect(alerts.map((a) => a.kind)).toEqual(["orphan"]);
  });

  it("uma subscrição nunca usada e sem cofre dispara stale + orphan", () => {
    const alerts = detectAlerts([sub({ id: "a" })], NOW);
    expect(alerts.map((a) => a.kind).sort()).toEqual(["orphan", "stale"]);
  });

  it("ignora subscrições inactivas nos alertas", () => {
    expect(detectAlerts([sub({ id: "a", active: false })], NOW)).toHaveLength(0);
  });

  it("resume custos e poupança potencial sem duplicar stale", () => {
    const old = new Date(NOW - 200 * DAY).toISOString();
    const subs = [
      sub({ id: "a", cost: 30, vaultItemId: "v", lastUsedAt: old }), // stale 30
      sub({ id: "b", cost: 10, vaultItemId: "v", lastUsedAt: new Date(NOW).toISOString() }),
    ];
    const sum = costSummary(subs, NOW);
    expect(sum.monthly).toBe(40);
    expect(sum.activeCount).toBe(2);
    expect(sum.potentialMonthlySaving).toBe(30);
  });
});
