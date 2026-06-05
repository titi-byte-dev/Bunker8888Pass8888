import type { DocFlowGraph } from "./types";

/** Constantes do layout — estilo sequence diagram (actores no topo, mensagens em filas). */
export const FLOW_LAYOUT = {
  padX: 32,
  padBottom: 40,
  headerY: 20,
  nodeHeight: 48,
  lifelineGap: 10,
  rowGapMin: 52,
  rowGapPerChar: 0.42,
  rowGapLong: 72,
  columnMin: 96,
  columnCharPx: 7,
  columnPad: 32,
  selfLoopWidth: 56,
  minCanvasWidth: 360,
} as const;

export type FlowLayoutNode = {
  id: string;
  label: string;
  x: number;
  y: number;
  columnCenterX: number;
  columnWidth: number;
  lifelineHeight: number;
};

export type FlowLayoutEdge = {
  id: string;
  rowIndex: number;
  rowY: number;
  rowHeight: number;
  sourceCenterX: number;
  targetCenterX: number;
  selfLoop: boolean;
};

export type FlowLayout = {
  nodes: FlowLayoutNode[];
  edges: FlowLayoutEdge[];
  canvasWidth: number;
  canvasHeight: number;
  lifelineTop: number;
};

/** Altura de cada fila conforme o comprimento do rótulo da mensagem. */
export function rowHeightForLabel(label: string): number {
  const len = label.length;
  if (len > 48) return FLOW_LAYOUT.rowGapLong + 16;
  if (len > 32) return FLOW_LAYOUT.rowGapLong;
  return Math.max(
    FLOW_LAYOUT.rowGapMin,
    Math.round(FLOW_LAYOUT.rowGapMin + len * FLOW_LAYOUT.rowGapPerChar),
  );
}

/** Largura mínima quando o contentor ainda não foi medido. */
function minContentWidth(graph: DocFlowGraph): number {
  const columnWidths = graph.nodes.map((n) =>
    Math.min(
      220,
      Math.max(
        FLOW_LAYOUT.columnMin,
        Math.ceil(n.label.length * FLOW_LAYOUT.columnCharPx) + FLOW_LAYOUT.columnPad,
      ),
    ),
  );
  let x = FLOW_LAYOUT.padX;
  for (const w of columnWidths) {
    x += w;
  }
  return Math.max(FLOW_LAYOUT.minCanvasWidth, x + FLOW_LAYOUT.padX);
}

/**
 * Distribui actores uniformemente pela largura do canvas (targetWidth).
 * Didático: o contentor do player mede-se com ResizeObserver; o grafo
 * estica-se para usar toda a área interactiva, não só o canto superior-esquerdo.
 */
export function computeFlowLayout(graph: DocFlowGraph, targetWidth?: number): FlowLayout {
  const lifelineTop = FLOW_LAYOUT.headerY + FLOW_LAYOUT.nodeHeight + FLOW_LAYOUT.lifelineGap;
  const canvasWidth = Math.max(minContentWidth(graph), targetWidth ?? minContentWidth(graph));

  const innerLeft = FLOW_LAYOUT.padX;
  const innerRight = canvasWidth - FLOW_LAYOUT.padX;
  const span = Math.max(0, innerRight - innerLeft);
  const nodeCount = graph.nodes.length;

  const columnCenters: number[] = [];
  if (nodeCount === 1) {
    columnCenters.push(canvasWidth / 2);
  } else {
    for (let i = 0; i < nodeCount; i++) {
      columnCenters.push(innerLeft + (span * i) / (nodeCount - 1));
    }
  }

  const maxNodeWidth =
    nodeCount > 1
      ? Math.min(220, Math.max(FLOW_LAYOUT.columnMin, span / nodeCount - 20))
      : Math.min(280, span * 0.55);

  const centerById = new Map<string, number>();
  const layoutNodes: FlowLayoutNode[] = graph.nodes.map((n, i) => {
    const columnCenterX = columnCenters[i];
    centerById.set(n.id, columnCenterX);
    const nodeWidth = Math.min(
      maxNodeWidth,
      Math.max(88, Math.min(maxNodeWidth, n.label.length * 6.5 + 28)),
    );
    return {
      id: n.id,
      label: n.label,
      x: columnCenterX - nodeWidth / 2,
      y: FLOW_LAYOUT.headerY,
      columnCenterX,
      columnWidth: maxNodeWidth,
      lifelineHeight: 0,
    };
  });

  let rowY = lifelineTop;
  const layoutEdges: FlowLayoutEdge[] = graph.edges.map((e, rowIndex) => {
    const rowHeight = rowHeightForLabel(e.label);
    const edgeLayout: FlowLayoutEdge = {
      id: e.id,
      rowIndex,
      rowY,
      rowHeight,
      sourceCenterX: centerById.get(e.source) ?? 0,
      targetCenterX: centerById.get(e.target) ?? 0,
      selfLoop: e.source === e.target,
    };
    rowY += rowHeight;
    return edgeLayout;
  });

  const timelineHeight = rowY - lifelineTop;
  const lifelineHeight = timelineHeight + FLOW_LAYOUT.padBottom;

  for (const node of layoutNodes) {
    node.lifelineHeight = lifelineHeight;
  }

  const canvasHeight = lifelineTop + timelineHeight + FLOW_LAYOUT.padBottom;

  return {
    nodes: layoutNodes,
    edges: layoutEdges,
    canvasWidth,
    canvasHeight,
    lifelineTop,
  };
}
