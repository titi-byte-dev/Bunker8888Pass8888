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

export type DocSection = {
  level: number;
  title: string;
  html: string;
  collapsed: boolean;
};

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
