import { describe, expect, it } from "vitest";
import { reconcileTransactions, type BankTransaction } from "./reconcile";
import type { Subscription } from "./subscriptions";

const subs: Subscription[] = [
  {
    id: "s1",
    name: "Netflix",
    cost: 12.99,
    currency: "EUR",
    cycle: "monthly",
    category: "",
    active: true,
    createdAt: "",
    updatedAt: "",
  },
  {
    id: "s2",
    name: "Spotify",
    cost: 15,
    currency: "EUR",
    cycle: "monthly",
    category: "",
    active: true,
    createdAt: "",
    updatedAt: "",
  },
];

describe("reconcileTransactions (FIN-003)", () => {
  it("associa débitos a subscrições pelo valor mensal", () => {
    const txs: BankTransaction[] = [
      { id: "t1", amount: -12.99, currency: "EUR", bookedAt: "2026-01-01", description: "NETFLIX" },
      { id: "t2", amount: -15, currency: "EUR", bookedAt: "2026-01-02", description: "SPOTIFY" },
      { id: "t3", amount: -8.5, currency: "EUR", bookedAt: "2026-01-03", description: "CAFE" },
    ];
    const r = reconcileTransactions(txs, subs);
    expect(r.matched).toHaveLength(2);
    expect(r.unmatched).toHaveLength(1);
    expect(r.matched[0].subscriptionName).toBe("Netflix");
  });
});
