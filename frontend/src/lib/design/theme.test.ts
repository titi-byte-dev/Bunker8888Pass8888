import { describe, expect, it } from "vitest";
import { cycleThemePreference, resolveTheme, themeModeLabel } from "./theme";

describe("resolveTheme", () => {
  it("system segue preferência do SO", () => {
    expect(resolveTheme("system", true)).toBe("dark");
    expect(resolveTheme("system", false)).toBe("light");
  });

  it("modos explícitos", () => {
    expect(resolveTheme("dark", false)).toBe("dark");
    expect(resolveTheme("light", true)).toBe("light");
  });
});

describe("cycleThemePreference", () => {
  it("roda light → dark → system → light", () => {
    expect(cycleThemePreference("light")).toBe("dark");
    expect(cycleThemePreference("dark")).toBe("system");
    expect(cycleThemePreference("system")).toBe("light");
  });
});

describe("themeModeLabel", () => {
  it("rótulos PT", () => {
    expect(themeModeLabel("light")).toBe("Claro");
    expect(themeModeLabel("dark")).toBe("Escuro");
    expect(themeModeLabel("system")).toBe("Sistema");
  });
});
