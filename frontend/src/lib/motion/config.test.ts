import { describe, expect, it } from "vitest";
import { MOTION } from "./config";

describe("motion config (UI-005)", () => {
  it("presets respeitam limite de 400ms para acções frequentes", () => {
    expect(MOTION.micro.duration).toBeLessThanOrEqual(0.4);
    expect(MOTION.panel.duration).toBeLessThanOrEqual(0.4);
    expect(MOTION.list.duration).toBeLessThanOrEqual(0.4);
  });

  it("list stagger dentro de 30–50ms por item", () => {
    expect(MOTION.list.stagger).toBeGreaterThanOrEqual(0.03);
    expect(MOTION.list.stagger).toBeLessThanOrEqual(0.05);
  });
});
