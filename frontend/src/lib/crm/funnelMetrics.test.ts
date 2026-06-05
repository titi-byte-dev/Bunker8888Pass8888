import { describe, expect, it } from "vitest";
import { computeFunnelMetrics } from "./funnelMetrics";
import type { Lead } from "./leads";

const base = (stage: Lead["stage"]): Lead => ({
  id: "1",
  name: "A",
  email: "a@x.pt",
  stage,
  source: "manual",
});

describe("computeFunnelMetrics", () => {
  it("agrega por estágio e calcula conversão", () => {
    const m = computeFunnelMetrics([
      base("new"),
      base("won"),
      base("won"),
      base("lost"),
    ]);
    expect(m.total).toBe(4);
    expect(m.won).toBe(2);
    expect(m.lost).toBe(1);
    expect(m.open).toBe(1);
    expect(m.conversionPct).toBe(67);
  });
});
