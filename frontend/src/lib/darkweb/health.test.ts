import { describe, expect, it } from "vitest";
import { computeCompositeScore, healthGrade } from "./health";
import { itemsRequiringPasswordChange, remediationEditUrl } from "./remediation";
import type { HygieneSummary } from "$lib/vault/hygiene";

describe("computeCompositeScore (DW-003)", () => {
  it("penaliza exposições em fugas", () => {
    expect(computeCompositeScore(80, 2, 0, 0, 5)).toBeLessThan(80);
    expect(computeCompositeScore(80, 0, 0, 0, 5)).toBe(80);
  });

  it("healthGrade classifica níveis", () => {
    expect(healthGrade(90)).toBe("good");
    expect(healthGrade(60)).toBe("warn");
    expect(healthGrade(30)).toBe("bad");
  });
});

describe("remediation (DW-002)", () => {
  const summary: HygieneSummary = {
    overallScore: 70,
    totalLogins: 2,
    weakCount: 1,
    reusedCount: 0,
    items: [
      { itemId: "a", title: "GitHub", score: 40, issues: ["weak"], reusedWith: [] },
      { itemId: "b", title: "Gmail", score: 90, issues: [], reusedWith: [] },
    ],
  };

  it("lista itens expostos em fugas", () => {
    const breaches = new Map([
      ["a", { breached: true, exposureCount: 100 }],
      ["b", { breached: false, exposureCount: 0 }],
    ]);
    const items = itemsRequiringPasswordChange(summary, breaches);
    expect(items).toHaveLength(1);
    expect(items[0]!.itemId).toBe("a");
    expect(items[0]!.reason).toBe("weak_and_breach");
  });

  it("remediationEditUrl inclui query remediate", () => {
    expect(remediationEditUrl("x", "weak_and_breach")).toContain("remediate=weak_and_breach");
  });
});
