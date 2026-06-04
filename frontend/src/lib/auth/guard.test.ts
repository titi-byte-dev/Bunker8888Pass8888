import { describe, expect, it } from "vitest";
import { isGuestAuthPath, isPublicAppPath, sanitizeRedirect } from "./guard";

describe("auth guard (UI-003)", () => {
  it("sanitizeRedirect bloqueia URLs externas e /auth", () => {
    expect(sanitizeRedirect("https://evil.test")).toBe("/vault");
    expect(sanitizeRedirect("//evil.test")).toBe("/vault");
    expect(sanitizeRedirect("/auth/login")).toBe("/vault");
    expect(sanitizeRedirect("/vault/item/1")).toBe("/vault/item/1");
    expect(sanitizeRedirect(null)).toBe("/vault");
  });

  it("isPublicAppPath inclui /dev", () => {
    expect(isPublicAppPath("/dev")).toBe(true);
    expect(isPublicAppPath("/vault")).toBe(false);
  });

  it("isGuestAuthPath permite login e recovery", () => {
    expect(isGuestAuthPath("/auth/login")).toBe(true);
    expect(isGuestAuthPath("/auth/unlock")).toBe(false);
  });
});
