import { describe, expect, it } from "vitest";
import {
  ROUTE_TREE,
  navModules,
  tabBarModules,
  isRouteActive,
  routeTrail,
  flattenRoutes,
  routeChildren,
} from "./routes";

describe("ROUTE_TREE (UI-011)", () => {
  it("navModules devolve os modulos de topo (sem os hideFromNav) na ordem definida", () => {
    const mods = navModules();
    expect(mods[0].label).toBe("Cofre");
    // Definicoes vive no menu de perfil (hideFromNav) -> fora da navegacao.
    expect(mods.every((m) => !m.hideFromNav)).toBe(true);
    expect(mods.some((m) => m.href === "/settings")).toBe(false);
    // Ordem preservada face a ROUTE_TREE visivel.
    expect(mods.map((m) => m.href)).toEqual(
      ROUTE_TREE.filter((m) => !m.hideFromNav).map((m) => m.href),
    );
  });

  it("tabBarModules: cofre, seguranca, trabalho (max 5, sem Definicoes)", () => {
    const tabs = tabBarModules().map((m) => m.href);
    expect(tabs).toEqual(["/vault", "/security", "/work"]);
    expect(tabs.length).toBeLessThanOrEqual(5);
  });

  it("isRouteActive distingue exacto e prefixo de seccao", () => {
    expect(isRouteActive("/fin", "/fin")).toBe(true);
    expect(isRouteActive("/fin/fiscal", "/fin")).toBe(true);
    expect(isRouteActive("/finance", "/fin")).toBe(false); // nao e sub-rota
    expect(isRouteActive("/crm", "/fin")).toBe(false);
  });

  it("routeTrail constroi o trilho raiz->folha (Financas > Fiscal)", () => {
    const trail = routeTrail("/fin/fiscal").map((n) => n.label);
    expect(trail).toEqual(["Financas", "Fiscal"]);
  });

  it("routeTrail usa prefixo mais longo para rotas dinamicas", () => {
    const trail = routeTrail("/security/devices/abc").map((n) => n.label);
    expect(trail).toEqual(["Seguranca", "Dispositivos e sessoes"]);
  });

  it("routeTrail devolve [] quando nada corresponde", () => {
    expect(routeTrail("/desconhecido")).toEqual([]);
  });

  it("flattenRoutes inclui modulos e filhos sem perder nenhum href", () => {
    const flat = flattenRoutes();
    expect(flat.some((n) => n.href === "/fin/commissions")).toBe(true);
    expect(flat.length).toBeGreaterThan(ROUTE_TREE.length);
  });

  it("routeChildren deriva os cartoes do hub (UI-014)", () => {
    const fin = routeChildren("/fin").map((c) => c.href);
    expect(fin).toContain("/fin/costs");
    expect(fin).toContain("/fin/fiscal");
    // exclui o filho-sombra que partilharia href com o modulo
    expect(fin).not.toContain("/fin");
    expect(routeChildren("/vault")).toEqual([]); // modulo-folha sem filhos
  });
});
