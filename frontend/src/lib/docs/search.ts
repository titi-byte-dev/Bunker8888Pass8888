/**
 * Pesquisa full-text na documentação in-app (DOC-010).
 * Didático: indexamos em memória a partir dos JSON gerados em build — sem fetch extra.
 */
import { inAppDocs } from "./loader";
import { stripHtml } from "./stripHtml";
import type { DocPage } from "./types";
import type { CommandEntry } from "$lib/shell/commands";

const pageModules = import.meta.glob<DocPage>("./generated/*.json", {
  eager: true,
  import: "default",
}) as Record<string, DocPage>;

function allPages(): DocPage[] {
  return Object.values(pageModules).filter((p): p is DocPage => Boolean(p?.slug && p.in_app));
}

function pageSearchText(page: DocPage): string {
  const parts = [
    page.title,
    page.summary,
    page.actor ?? "",
    page.categoryLabel,
    page.feature ?? "",
    ...page.concepts.map((c) => `${c.title} ${stripHtml(c.html)}`),
    ...page.sections.map((s) => `${s.title} ${stripHtml(s.html)}`),
  ];
  return parts.join(" ").toLowerCase();
}

export type DocSearchHit = {
  slug: string;
  title: string;
  categoryLabel: string;
  summary: string;
  snippet: string;
  score: number;
};

function makeSnippet(text: string, query: string, fallback: string): string {
  const idx = text.indexOf(query);
  if (idx < 0) return fallback.slice(0, 140);
  const start = Math.max(0, idx - 48);
  const end = Math.min(text.length, idx + query.length + 72);
  const slice = text.slice(start, end).trim();
  return (start > 0 ? "…" : "") + slice + (end < text.length ? "…" : "");
}

/** Pesquisa páginas por título, resumo e corpo (HTML ignorado). */
export function searchDocs(query: string, limit = 12): DocSearchHit[] {
  const q = query.trim().toLowerCase();
  if (!q) return [];

  const words = q.split(/\s+/).filter(Boolean);
  const hits: DocSearchHit[] = [];

  for (const page of allPages()) {
    const text = pageSearchText(page);
    const title = page.title.toLowerCase();
    let score = 0;

    if (title === q) score += 20;
    else if (title.includes(q)) score += 10;

    if (page.summary.toLowerCase().includes(q)) score += 5;

    for (const w of words) {
      if (title.includes(w)) score += 3;
      if (page.summary.toLowerCase().includes(w)) score += 2;
      if (text.includes(w)) score += 1;
    }

    if (score === 0) continue;

    hits.push({
      slug: page.slug,
      title: page.title,
      categoryLabel: page.categoryLabel,
      summary: page.summary,
      snippet: makeSnippet(text, q, page.summary || page.title),
      score,
    });
  }

  return hits.sort((a, b) => b.score - a.score).slice(0, limit);
}

/** Comandos estáticos (títulos) para a palette quando há correspondência simples. */
export function buildDocCommands(): CommandEntry[] {
  return inAppDocs().map((doc) => ({
    id: `doc-${doc.slug}`,
    label: doc.title,
    group: "docs" as const,
    keywords: `${doc.summary} ${doc.categoryLabel} ${doc.feature ?? ""} documentação ajuda`,
    href: `/settings/docs/${doc.slug}`,
  }));
}

/** Resultados de pesquisa profunda para a command palette. */
export function buildDocSearchCommands(query: string): CommandEntry[] {
  return searchDocs(query, 10).map((hit) => ({
    id: `doc-search-${hit.slug}`,
    label: hit.title,
    group: "docs" as const,
    keywords: `${hit.snippet} ${hit.categoryLabel} ${hit.summary}`,
    href: `/settings/docs/${hit.slug}`,
  }));
}
