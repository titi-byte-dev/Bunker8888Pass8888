import { describe, expect, it } from "vitest";
import { blindScreenCV } from "./blindScreen";

describe("blindScreenCV (AGENT-007)", () => {
  it("oculta género e etnia mas mantém skills", () => {
    const out = blindScreenCV("Nome: Ana\nGénero: F\nEtnia: Y\nSkills: TypeScript");
    expect(out).not.toContain("Género: F");
    expect(out).toContain("Ana");
    expect(out).toContain("TypeScript");
  });
});
