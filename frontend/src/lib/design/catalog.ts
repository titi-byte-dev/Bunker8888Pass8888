/**
 * Metadados do catálogo de componentes (UI-010).
 * Usado pela página /dev/components e testes de regressão leve.
 */

export type CatalogSection = {
  id: string;
  title: string;
};

export const CATALOG_SECTIONS: CatalogSection[] = [
  { id: "typography", title: "Tipografia" },
  { id: "colors", title: "Cores semânticas" },
  { id: "spacing", title: "Spacing & radius" },
  { id: "buttons", title: "Botões" },
  { id: "forms", title: "Formulários" },
  { id: "components", title: "Componentes vivos" },
];

export const SEMANTIC_COLORS = [
  { name: "accent", var: "--color-accent", bg: "--color-accent-muted" },
  { name: "success", var: "--color-success-fg", bg: "--color-success-bg" },
  { name: "warning", var: "--color-warning", bg: "--color-bg-surface" },
  { name: "danger", var: "--color-danger", bg: "--color-bg-surface" },
] as const;

export const TYPE_SCALE = [
  { token: "--text-xs", sample: "Texto xs" },
  { token: "--text-sm", sample: "Texto sm" },
  { token: "--text-base", sample: "Texto base" },
  { token: "--text-lg", sample: "Texto lg" },
  { token: "--text-xl", sample: "Texto xl" },
  { token: "--text-2xl", sample: "Texto 2xl" },
  { token: "--text-3xl", sample: "Texto 3xl" },
] as const;

/** Mock para preview do SecurityHealthCard no catálogo. */
export const MOCK_HEALTH_REPORT = {
  at: new Date().toISOString(),
  hygieneScore: 82,
  exposedCount: 0,
  weakCount: 1,
  reusedCount: 0,
  compositeScore: 78,
  trend: "up" as const,
  trendDelta: 4,
  totalLogins: 12,
  scannedCount: 12,
};
