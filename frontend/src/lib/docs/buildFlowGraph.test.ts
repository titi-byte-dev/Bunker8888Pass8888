import { describe, expect, it } from "vitest";
import { buildFlowEdges, buildFlowNodes } from "./buildFlowGraph";
import type { DocFlowGraph } from "./types";

const sampleGraph: DocFlowGraph = {
  nodes: [
    { id: "SMTP", label: "E-mail", x: 0 },
    { id: "API", label: "API", x: 200 },
    { id: "Orq", label: "Orquestrador", x: 400 },
  ],
  edges: [
    { id: "step-0", source: "SMTP", target: "API", label: "ingest" },
    { id: "step-1", source: "API", target: "Orq", label: "event" },
    { id: "step-2", source: "Orq", target: "Orq", label: "self" },
  ],
};

const messageSteps = [
  { kind: "message" as const, from: "SMTP", to: "API", arrow: "->>", label: "ingest" },
  { kind: "message" as const, from: "API", to: "Orq", arrow: "->>", label: "event" },
  { kind: "message" as const, from: "Orq", to: "Orq", arrow: "->>", label: "self" },
];

describe("buildFlowGraph (DOC-011)", () => {
  it("destaca nós do passo actual com papéis from/to", () => {
    const n0 = buildFlowNodes(sampleGraph, messageSteps, 0);
    expect(n0.find((n) => n.id === "SMTP")?.data.role).toBe("from");
    expect(n0.find((n) => n.id === "API")?.data.role).toBe("to");
    expect(n0.find((n) => n.id === "Orq")?.data.role).toBe(null);

    const n1 = buildFlowNodes(sampleGraph, messageSteps, 1);
    expect(n1.find((n) => n.id === "SMTP")?.data.active).toBe(false);
    expect(n1.find((n) => n.id === "API")?.data.active).toBe(true);
    expect(n1.find((n) => n.id === "Orq")?.data.active).toBe(true);
  });

  it("anima só a aresta do passo actual", () => {
    const e0 = buildFlowEdges(sampleGraph, 0, true);
    expect(e0[0].animated).toBe(true);
    expect(e0[1].animated).toBe(false);
    expect(e0[0].class).toBe("flow-edge-current");
    expect(e0[1].class).toBe("flow-edge-pending");
  });

  it("self-loop usa aresta docSequence com flag", () => {
    const edges = buildFlowEdges(sampleGraph, 2, true);
    expect(edges[2].type).toBe("docSequence");
    expect(edges[2].data?.selfLoop).toBe(true);
    expect(edges[2].data?.rowY).toBeGreaterThan(0);
  });

  it("sem animação quando reduced motion", () => {
    const edges = buildFlowEdges(sampleGraph, 1, false);
    expect(edges[1].animated).toBe(false);
  });
});
