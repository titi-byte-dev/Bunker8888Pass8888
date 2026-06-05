import { afterEach, describe, expect, it, vi } from "vitest";
import { motionDuration, motionStagger, prefersReducedMotion } from "./reduced";

describe("reduced motion helpers", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("prefersReducedMotion lê matchMedia", () => {
    vi.stubGlobal("window", {
      matchMedia: vi.fn(() => ({ matches: true })),
    });
    expect(prefersReducedMotion()).toBe(true);
  });

  it("motionDuration devolve 0 quando reduced motion activo", () => {
    vi.stubGlobal("window", {
      matchMedia: vi.fn(() => ({ matches: true })),
    });
    expect(motionDuration(0.32)).toBe(0);
    expect(motionStagger(0.04)).toBe(0);
  });

  it("motionDuration mantém valor quando animações permitidas", () => {
    vi.stubGlobal("window", {
      matchMedia: vi.fn(() => ({ matches: false })),
    });
    expect(motionDuration(0.32)).toBe(0.32);
    expect(motionStagger(0.04)).toBe(0.04);
  });
});
