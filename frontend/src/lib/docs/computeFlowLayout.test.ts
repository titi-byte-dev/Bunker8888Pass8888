import { describe, expect, it } from "vitest";
import { computeFlowLayout, rowHeightForLabel } from "./computeFlowLayout";
import type { DocFlowGraph } from "./types";

const graph: DocFlowGraph = {
  nodes: [
    { id: "A", label: "E-mail", x: 0 },
    { id: "B", label: "API", x: 200 },
    { id: "C", label: "Orquestrador", x: 400 },
  ],
  edges: [
    { id: "step-0", source: "A", target: "B", label: "ingest alias" },
    { id: "step-1", source: "B", target: "C", label: "mail.inbox.received" },
    { id: "step-2", source: "C", target: "C", label: "orchestrator.action.suggested" },
  ],
};

describe("computeFlowLayout (DOC-011)", () => {
  it("coloca mensagens em filas verticais distintas", () => {
    const layout = computeFlowLayout(graph);
    const ys = layout.edges.map((e) => e.rowY);
    expect(new Set(ys).size).toBe(ys.length);
    expect(layout.canvasHeight).toBeGreaterThan(240);
    expect(layout.canvasWidth).toBeGreaterThan(300);
  });

  it("cresce a altura com rótulos longos", () => {
    const short = rowHeightForLabel("ok");
    const long = rowHeightForLabel("POST /api/agent/prospection/run (utilizador)");
    expect(long).toBeGreaterThan(short);
  });

  it("distribui actores pela largura alvo do canvas", () => {
    const layout = computeFlowLayout(graph, 800);
    expect(layout.canvasWidth).toBe(800);
    const xs = layout.nodes.map((n) => n.columnCenterX);
    expect(xs[0]).toBeLessThan(xs[xs.length - 1]);
    expect(xs[0]).toBeGreaterThanOrEqual(32);
    expect(xs[xs.length - 1]).toBeLessThanOrEqual(768);
  });

  it("assigna lifeline a cada actor", () => {
    const layout = computeFlowLayout(graph);
    for (const node of layout.nodes) {
      expect(node.lifelineHeight).toBeGreaterThan(0);
      expect(node.columnCenterX).toBeGreaterThan(0);
    }
  });
});
