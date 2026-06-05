import { describe, expect, it } from "vitest";
import { buildDocCommands, searchDocs } from "./search";

describe("doc search (DOC-010)", () => {
  it("searchDocs encontra glossário por termo técnico", () => {
    const hits = searchDocs("zero-knowledge");
    expect(hits.some((h) => h.slug === "glossary")).toBe(true);
  });

  it("searchDocs encontra percursos SHARE", () => {
    const hits = searchDocs("secret link");
    expect(hits.some((h) => h.slug === "journey-secret-link")).toBe(true);
  });

  it("searchDocs devolve vazio sem query", () => {
    expect(searchDocs("")).toHaveLength(0);
    expect(searchDocs("   ")).toHaveLength(0);
  });

  it("buildDocCommands inclui páginas in-app", () => {
    const cmds = buildDocCommands();
    expect(cmds.length).toBeGreaterThan(10);
    expect(cmds.every((c) => c.group === "docs")).toBe(true);
    expect(cmds.some((c) => c.id === "doc-glossary")).toBe(true);
  });
});
