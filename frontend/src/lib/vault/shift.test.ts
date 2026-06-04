import { describe, expect, it } from "vitest";
import { isClockSkewAcceptable, isWithinShift, msUntilShiftEnd } from "./shift";
import type { ShiftPolicy } from "./shift";

const policy: ShiftPolicy = {
  enabled: true,
  timezone: "UTC",
  max_clock_skew_seconds: 300,
  schedule: {
    wed: [{ start: "09:00", end: "17:00" }],
  },
};

describe("isWithinShift", () => {
  it("permite quando disabled", () => {
    expect(isWithinShift(new Date(), { ...policy, enabled: false })).toBe(true);
  });

  it("dentro da janela quarta 10:00 UTC", () => {
    const now = new Date("2026-06-03T10:00:00.000Z");
    expect(isWithinShift(now, policy)).toBe(true);
  });

  it("fora da janela quarta 18:00 UTC", () => {
    const now = new Date("2026-06-03T18:00:00.000Z");
    expect(isWithinShift(now, policy)).toBe(false);
  });
});

describe("msUntilShiftEnd", () => {
  it("devolve ms positivos dentro do turno", () => {
    const now = new Date("2026-06-03T10:00:00.000Z");
    const ms = msUntilShiftEnd(now, policy);
    expect(ms).not.toBeNull();
    expect(ms!).toBeGreaterThan(0);
  });
});

describe("isClockSkewAcceptable", () => {
  it("aceita desvio dentro do limite", () => {
    expect(isClockSkewAcceptable(1_000_000, 1_000_100, 300)).toBe(true);
  });

  it("rejeita desvio excessivo", () => {
    expect(isClockSkewAcceptable(0, 600_000, 300)).toBe(false);
  });
});
