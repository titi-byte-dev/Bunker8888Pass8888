import { describe, expect, it } from "vitest";
import { loadDocPage } from "./loader";
import type { DocFlow } from "./types";

/** Extrai todos os fluxos de uma página de documentação. */
function collectFlows(slug: string): DocFlow[] {
  const page = loadDocPage(slug);
  if (!page) return [];
  return page.sections.flatMap((s) => s.flows ?? []);
}

describe("Svelte Flow nos journeys (DOC-011/012)", () => {
  it("journeys do orquestrador usam renderer svelteflow com grafo", () => {
    for (const slug of [
      "journey-orchestrator",
      "journey-human-in-the-loop",
      "journey-ops-agent-inventory",
      "journey-hr-agent-recruitment",
      "journey-finance-agent-saas",
      "journey-finance-agent-reconcile",
    ]) {
      const flows = collectFlows(slug);
      const sequence = flows.find((f) => f.type === "sequence");
      expect(sequence, slug).toBeDefined();
      expect(sequence?.renderer).toBe("svelteflow");
      expect(sequence?.graph?.nodes.length).toBeGreaterThan(0);
      expect(sequence?.graph?.edges.length).toBe(sequence?.steps.length);
    }
  });

  it("payload DocFlowNodeData alinha com nós do grafo", () => {
    const flows = collectFlows("journey-orchestrator");
    const flow = flows[0];
    const node = flow.graph!.nodes[0];
    // O SvelteFlowPlayer mapeia graph.nodes → data.label; validamos o contrato.
    expect(node).toMatchObject({
      id: expect.any(String),
      label: expect.any(String),
      x: expect.any(Number),
    });
  });
});
