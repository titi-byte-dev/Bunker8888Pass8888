import { describe, expect, it } from "vitest";
import { normalizeEmail, wrapNetworkError } from "./http";

describe("auth http helpers", () => {
  it("normalizeEmail trim e lowercase", () => {
    expect(normalizeEmail("  Dev@Test.COM ")).toBe("dev@test.com");
  });

  it("wrapNetworkError traduz falha de fetch", () => {
    const err = wrapNetworkError(new TypeError("Failed to fetch"), "Login");
    expect(err.message).toContain("servidor inacessível");
  });
});
