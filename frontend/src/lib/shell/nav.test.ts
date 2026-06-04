import { describe, expect, it } from "vitest";
import { isNavActive, tabBarItems, visibleNavItems } from "./nav";

describe("shell nav (UI-002)", () => {
  it("visibleNavItems omite devOnly fora de dev", () => {
    const prod = visibleNavItems(false);
    expect(prod.some((i) => i.id === "dev")).toBe(false);
    expect(prod.some((i) => i.id === "vault")).toBe(true);
  });

  it("tabBarItems inclui cofre, segurança, trabalho e definições", () => {
    const tabs = tabBarItems(true).map((i) => i.id);
    expect(tabs).toEqual(["vault", "security", "work", "settings"]);
  });

  it("isNavActive distingue prefixos de rota", () => {
    expect(isNavActive("/vault", "/vault")).toBe(true);
    expect(isNavActive("/vault/item/1", "/vault")).toBe(true);
    expect(isNavActive("/security", "/vault")).toBe(false);
  });
});
