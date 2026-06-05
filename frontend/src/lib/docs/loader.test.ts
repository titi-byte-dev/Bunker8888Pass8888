import { describe, expect, it } from "vitest";
import { DOC_MANIFEST, inAppDocs, loadDocPage, docsByCategory } from "./loader";

describe("docs loader (DOC-002)", () => {
  it("manifest tem páginas in-app", () => {
    expect(inAppDocs().length).toBeGreaterThan(0);
    expect(DOC_MANIFEST.levelLabels[1]).toBe("Essencial");
  });

  it("carrega glossary com conceitos", () => {
    const page = loadDocPage("glossary");
    expect(page?.title).toContain("Glossário");
    expect(page?.concepts.length).toBeGreaterThan(0);
  });

  it("agrupa por categoria", () => {
    const groups = docsByCategory();
    expect(groups.some((g) => g.id === "concepts")).toBe(true);
    expect(groups.some((g) => g.id === "journeys")).toBe(true);
  });
});
