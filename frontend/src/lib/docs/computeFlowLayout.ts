import type { DocFlowGraph } from "./types";

/** Constantes do layout — estilo sequence diagram (actores no topo, mensagens em filas). */
export const FLOW_LAYOUT = {
  padX: 28,
  padBottom: 36,
  headerY: 16,
  nodeHeight: 48,
  lifelineGap: 8,
  rowGapMin: 56,
  rowGapPerChar: 0.45,
  rowGapLong: 76,
  columnMin: 100,
  columnCharPx: 7.2,
  columnPad: 36,
  selfLoopWidth: 52,
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

/**
 * Distribui actores em colunas e mensagens em filas verticais.
 * Evita sobrepor todas as arestas na mesma linha (problema do layout horizontal antigo).
 */
export function computeFlowLayout(graph: DocFlowGraph): FlowLayout {
  const lifelineTop = FLOW_LAYOUT.headerY + FLOW_LAYOUT.nodeHeight + FLOW_LAYOUT.lifelineGap;

  const columnWidths = graph.nodes.map((n) =>
    Math.min(
      220,
      Math.max(FLOW_LAYOUT.columnMin, Math.ceil(n.label.length * FLOW_LAYOUT.columnCharPx) + FLOW_LAYOUT.columnPad),
    ),
  );

  const columnStarts: number[] = [];
  let x = FLOW_LAYOUT.padX;
  for (const w of columnWidths) {
    columnStarts.push(x);
    x += w;
  }
  const canvasWidth = Math.max(320, x + FLOW_LAYOUT.padX);

  const centerById = new Map<string, number>();
  const layoutNodes: FlowLayoutNode[] = graph.nodes.map((n, i) => {
    const columnWidth = columnWidths[i];
    const columnCenterX = columnStarts[i] + columnWidth / 2;
    centerById.set(n.id, columnCenterX);
    const nodeWidth = Math.min(columnWidth - 12, Math.max(88, n.label.length * 6.5 + 28));
    return {
      id: n.id,
      label: n.label,
      x: columnCenterX - nodeWidth / 2,
      y: FLOW_LAYOUT.headerY,
      columnCenterX,
      columnWidth,
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
