import { describe, expect, it } from "vitest";
import { annotateGlossaryHtml, buildGlossaryIndex, glossaryTerms } from "./glossary";

describe("glossary inline (DOC-010)", () => {
  it("buildGlossaryIndex inclui conceitos do glossário", () => {
    const index = buildGlossaryIndex();
    expect(index.has("zero-knowledge")).toBe(true);
    expect(index.has("tenant")).toBe(true);
  });

  it("glossaryTerms ordena aliases longos primeiro", () => {
    const terms = glossaryTerms();
    expect(terms.length).toBeGreaterThan(5);
    for (let i = 1; i < terms.length; i++) {
      expect(terms[i - 1]!.alias.length).toBeGreaterThanOrEqual(terms[i]!.alias.length);
    }
  });

  it("annotateGlossaryHtml envolve termos fora de code", () => {
    const html = "<p>O modelo <strong>Zero-Knowledge</strong> protege o cofre.</p>";
    const out = annotateGlossaryHtml(html);
    expect(out).toContain("glossary-term");
    expect(out).toContain("role=\"tooltip\"");
  });

  it("annotateGlossaryHtml ignora conteúdo em code", () => {
    const html = "<p>Ver <code>tenant_id</code> na query.</p>";
    const out = annotateGlossaryHtml(html);
    expect(out).not.toContain('<code class="glossary-term"');
    expect(out).toContain("<code>tenant_id</code>");
  });
});
