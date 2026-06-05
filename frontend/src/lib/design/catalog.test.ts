import { describe, expect, it } from "vitest";
import { CATALOG_SECTIONS, MOCK_HEALTH_REPORT, TYPE_SCALE } from "./catalog";

describe("design catalog (UI-010)", () => {
  it("CATALOG_SECTIONS tem ids únicos", () => {
    const ids = CATALOG_SECTIONS.map((s) => s.id);
    expect(new Set(ids).size).toBe(ids.length);
    expect(ids.length).toBeGreaterThan(3);
  });

  it("MOCK_HEALTH_REPORT é score válido", () => {
    expect(MOCK_HEALTH_REPORT.compositeScore).toBeGreaterThan(0);
    expect(MOCK_HEALTH_REPORT.compositeScore).toBeLessThanOrEqual(100);
  });

  it("TYPE_SCALE inclui display sizes", () => {
    expect(TYPE_SCALE.some((t) => t.token === "--text-2xl")).toBe(true);
  });
});
