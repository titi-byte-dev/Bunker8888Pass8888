/** Tipos da documentação gerada por scripts/build-docs.mjs (SSOT: docs/). */

export type DocCategory =
  | "concepts"
  | "product"
  | "developer"
  | "competitive"
  | "journeys";

export type DocAudience = "user" | "developer" | "admin";

export type DocConcept = {
  id: string;
  title: string;
  level: number;
  html: string;
};

export type DocFlowStep =
  | { kind: "message"; from: string; to: string; arrow: string; label: string }
  | { kind: "branch"; label: string };

export type DocFlowGraphNode = {
  id: string;
  label: string;
  x: number;
};

/** Payload de cada nó no Svelte Flow (DOC-011). */
export type DocFlowNodeRole = "from" | "to" | "both" | null;

export type DocFlowNodeData = {
  label: string;
  active?: boolean;
  /** Origem / destino do passo actual — cores distintas. */
  role?: DocFlowNodeRole;
  /** Passo actual — força re-render quando o highlight muda. */
  stepIndex?: number;
  /** Linha vertical (lifeline) até ao fim da timeline. */
  lifelineHeight?: number;
};

export type DocFlowEdgeState = "current" | "done" | "pending";

/** Dados da aresta em fila (sequence diagram). */
export type DocSequenceEdgeData = {
  rowIndex: number;
  rowY: number;
  rowHeight: number;
  sourceCenterX: number;
  targetCenterX: number;
  selfLoop: boolean;
  selfLoopWidth?: number;
  state: DocFlowEdgeState;
  canvasWidth: number;
};

export type DocFlowGraphEdge = {
  id: string;
  source: string;
  target: string;
  label: string;
  dashed?: boolean;
};

export type DocFlowGraph = {
  nodes: DocFlowGraphNode[];
  edges: DocFlowGraphEdge[];
};

export type DocFlowRenderer = "mermaid" | "svelteflow";

export type DocFlow = {
  id: string;
  title: string;
  type: "sequence" | "flowchart" | "diagram";
  source: string;
  steps: DocFlowStep[];
  renderer?: DocFlowRenderer;
  graph?: DocFlowGraph;
};

export type DocSection = {
  level: number;
  title: string;
  html: string;
  flows?: DocFlow[];
  collapsed: boolean;
};

export type DocSectionPart =
  | { kind: "html"; content: string }
  | { kind: "flow"; flow: DocFlow };

export type DocMeta = {
  slug: string;
  title: string;
  category: DocCategory;
  categoryLabel: string;
  audience: DocAudience[];
  layer: string[];
  feature?: string;
  level: number;
  maxLevel: number;
  in_app: boolean;
  summary: string;
  related: string[];
  actor?: string;
  order: number;
};

export type DocPage = DocMeta & {
  concepts: DocConcept[];
  sections: DocSection[];
};

export type DocManifest = {
  generatedAt: string;
  categories: { id: DocCategory; label: string }[];
  levelLabels: Record<number, string>;
  docs: DocMeta[];
};

export type DocComplexityLevel = 1 | 2 | 3;

export const LEVEL_LABELS: Record<DocComplexityLevel, string> = {
  1: "Essencial",
  2: "Intermédio",
  3: "Técnico",
};
