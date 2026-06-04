import { describe, it, expect } from "vitest";
import {
  passwordStrengthScore,
  analyzeVaultHygiene,
  toServerSummary,
} from "./hygiene";

describe("passwordStrengthScore (VAULT-008)", () => {
  it("password vazia tem score 0", () => {
    expect(passwordStrengthScore("")).toBe(0);
  });

  it("password comum tem score baixo", () => {
    expect(passwordStrengthScore("password")).toBeLessThan(50);
    expect(passwordStrengthScore("123456")).toBeLessThan(50);
  });

  it("password forte tem score alto", () => {
    expect(passwordStrengthScore("Xk9#mP2$vL8@nQ4!wR7")).toBeGreaterThan(70);
  });
});

describe("analyzeVaultHygiene (VAULT-008)", () => {
  it("deteta password reutilizada entre logins", async () => {
    const report = await analyzeVaultHygiene([
      { itemId: "a", title: "GitHub", password: "MesmaPass123!" },
      { itemId: "b", title: "Gmail", password: "MesmaPass123!" },
      { itemId: "c", title: "AWS", password: "UnicaForte#2024$xyz" },
    ]);

    expect(report.reusedCount).toBe(2);
    const github = report.items.find((i) => i.itemId === "a")!;
    expect(github.issues).toContain("reused");
    expect(github.reusedWith).toContain("b");
  });

  it("conta passwords fracas", async () => {
    const report = await analyzeVaultHygiene([
      { itemId: "1", title: "A", password: "123456" },
      { itemId: "2", title: "B", password: "Tr0ub4dor&3ExtraChars!" },
    ]);
    expect(report.weakCount).toBeGreaterThanOrEqual(1);
  });

  it("toServerSummary não inclui items (sem passwords)", () => {
    const summary = toServerSummary({
      overallScore: 75,
      totalLogins: 2,
      weakCount: 1,
      reusedCount: 0,
      items: [],
    });
    expect(summary).not.toHaveProperty("items");
    expect(summary.overallScore).toBe(75);
  });

  it("cofre vazio tem score 100", async () => {
    const report = await analyzeVaultHygiene([]);
    expect(report.overallScore).toBe(100);
    expect(report.totalLogins).toBe(0);
  });
});
