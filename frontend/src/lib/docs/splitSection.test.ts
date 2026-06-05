import { describe, expect, it } from "vitest";
import { splitSectionHtml } from "./splitSection";
import type { DocFlow } from "./types";

describe("splitSectionHtml (DOC-008)", () => {
  const flow: DocFlow = {
    id: "flow-1",
    title: "Teste",
    type: "sequence",
    source: "sequenceDiagram\n  A->>B: hi",
    steps: [],
  };

  it("divide HTML e placeholders de fluxo", () => {
    const parts = splitSectionHtml(
      "<p>Intro</p><!--DOC_FLOW:flow-1--><p>Fim</p>",
      [flow],
    );
    expect(parts).toHaveLength(3);
    expect(parts[0]).toEqual({ kind: "html", content: "<p>Intro</p>" });
    expect(parts[1]).toEqual({ kind: "flow", flow });
    expect(parts[2]).toEqual({ kind: "html", content: "<p>Fim</p>" });
  });

  it("ignora placeholder sem fluxo correspondente", () => {
    const parts = splitSectionHtml("<!--DOC_FLOW:missing-->", []);
    expect(parts).toHaveLength(0);
  });
});
