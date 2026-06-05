import type { Edge, Node } from "@xyflow/svelte";
import { computeFlowLayout } from "./computeFlowLayout";
import type {
  DocFlowGraph,
  DocFlowNodeData,
  DocFlowNodeRole,
  DocFlowStep,
  DocSequenceEdgeData,
} from "./types";

type MessageStep = Extract<DocFlowStep, { kind: "message" }>;

function nodeRole(step: MessageStep | undefined, nodeId: string): DocFlowNodeRole {
  if (!step) return null;
  const isFrom = nodeId === step.from;
  const isTo = nodeId === step.to;
  if (isFrom && isTo) return "both";
  if (isFrom) return "from";
  if (isTo) return "to";
  return null;
}

/** Constrói nós XYFlow com highlight do passo actual e layout em colunas. */
export function buildFlowNodes(
  graph: DocFlowGraph,
  messageSteps: MessageStep[],
  stepIndex: number,
): Node<DocFlowNodeData, "docFlow">[] {
  const step = messageSteps[stepIndex];
  const layout = computeFlowLayout(graph);

  return layout.nodes.map((n) => {
    const role = nodeRole(step, n.id);
    return {
      id: n.id,
      type: "docFlow" as const,
      position: { x: n.x, y: n.y },
      data: {
        label: n.label,
        active: role !== null,
        role,
        stepIndex,
        lifelineHeight: n.lifelineHeight,
      },
      draggable: false,
      selectable: false,
      focusable: false,
    };
  });
}

/** Constrói arestas em filas verticais — uma mensagem por linha (legível em canvas largo). */
export function buildFlowEdges(
  graph: DocFlowGraph,
  stepIndex: number,
  animate: boolean,
): Edge<DocSequenceEdgeData>[] {
  const layout = computeFlowLayout(graph);

  return graph.edges.map((e, i) => {
    const meta = layout.edges[i];
    const isCurrent = i === stepIndex;
    const isDone = i < stepIndex;
    const rtl = meta.targetCenterX < meta.sourceCenterX;
    const state = isCurrent ? "current" : isDone ? "done" : "pending";

    return {
      id: e.id,
      source: e.source,
      target: e.target,
      type: "docSequence",
      label: e.label,
      data: {
        rowIndex: meta.rowIndex,
        rowY: meta.rowY,
        rowHeight: meta.rowHeight,
        sourceCenterX: meta.sourceCenterX,
        targetCenterX: meta.targetCenterX,
        selfLoop: meta.selfLoop,
        selfLoopWidth: 52,
        state,
        canvasWidth: layout.canvasWidth,
      } satisfies DocSequenceEdgeData,
      animated: isCurrent && animate,
      class: `flow-edge-${state}`,
      markerEnd: !meta.selfLoop && !rtl ? "url(#svelte-flow__arrowclosed)" : undefined,
      markerStart: !meta.selfLoop && rtl ? "url(#svelte-flow__arrowclosed)" : undefined,
      selectable: false,
      focusable: false,
    };
  });
}

/** Dimensões do canvas para o grafo (partilhado pelo player). */
export function flowCanvasSize(graph: DocFlowGraph): { width: number; height: number } {
  const layout = computeFlowLayout(graph);
  return { width: layout.canvasWidth, height: layout.canvasHeight };
}
