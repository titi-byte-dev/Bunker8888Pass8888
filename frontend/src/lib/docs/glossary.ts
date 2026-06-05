/**
 * Glossário inline (DOC-010): termos das páginas com tooltip e link.
 * Didático: reunimos todos os blocos :::concept do build-docs num índice único.
 */
import { stripHtml } from "./stripHtml";
import type { DocPage } from "./types";

const pageModules = import.meta.glob<DocPage>("./generated/*.json", {
  eager: true,
  import: "default",
}) as Record<string, DocPage>;

/** Aliases manuais → id de conceito (termos usados no texto mas definidos noutro sítio). */
const ALIAS_TO_CONCEPT: Record<string, string> = {
  argon2id: "master-key",
  "auth hash": "master-key",
  "master password": "master-key",
  "master key": "master-key",
  rls: "tenant",
  "row-level security": "tenant",
  "shared vault": "shared-vault",
  "shared vaults": "shared-vault",
  "secret link": "secret-link",
  "secret links": "secret-link",
  sentinel: "sentinel",
  "sentinel mode": "sentinel",
  zk: "zero-knowledge",
  "zero knowledge": "zero-knowledge",
  "k-anonymity": "hygiene-score",
  "k anonymity": "hygiene-score",
};

export type GlossaryEntry = {
  id: string;
  title: string;
  plain: string;
  aliases: string[];
  href: string;
};

function conceptAliases(title: string, id: string): string[] {
  const aliases = new Set<string>();
  aliases.add(title);

  const paren = title.match(/^(.+?)\s*\(([^)]+)\)/);
  if (paren) {
    aliases.add(paren[1]!.trim());
    aliases.add(paren[2]!.trim());
  }

  const beforeDash = title.split("—")[0]?.trim();
  if (beforeDash && beforeDash !== title) aliases.add(beforeDash);

  const idWords = id.replace(/-/g, " ");
  if (idWords.length >= 3) aliases.add(idWords);

  return [...aliases].filter((a) => a.length >= 3);
}

/** Índice único de conceitos (glossary tem prioridade sobre páginas de produto). */
export function buildGlossaryIndex(): Map<string, GlossaryEntry> {
  const entries = new Map<string, GlossaryEntry>();

  for (const page of Object.values(pageModules)) {
    if (!page?.concepts?.length) continue;
    for (const concept of page.concepts) {
      const href =
        page.slug === "glossary"
          ? `/settings/docs/glossary#concept-${concept.id}`
          : `/settings/docs/${page.slug}#concept-${concept.id}`;

      const candidate: GlossaryEntry = {
        id: concept.id,
        title: concept.title,
        plain: stripHtml(concept.html).slice(0, 240),
        aliases: conceptAliases(concept.title, concept.id),
        href,
      };

      const prev = entries.get(concept.id);
      if (!prev || page.slug === "glossary") {
        entries.set(concept.id, candidate);
      }
    }
  }

  for (const [alias, conceptId] of Object.entries(ALIAS_TO_CONCEPT)) {
    const entry = entries.get(conceptId);
    if (entry && !entry.aliases.some((a) => a.toLowerCase() === alias)) {
      entry.aliases.push(alias);
    }
  }

  return entries;
}

let cachedTerms: GlossaryTerm[] | null = null;

export type GlossaryTerm = {
  id: string;
  alias: string;
  title: string;
  plain: string;
  href: string;
};

/** Lista plana de aliases ordenada por comprimento (evita substituições parciais). */
export function glossaryTerms(): GlossaryTerm[] {
  if (cachedTerms) return cachedTerms;

  const terms: GlossaryTerm[] = [];
  for (const entry of buildGlossaryIndex().values()) {
    for (const alias of entry.aliases) {
      terms.push({
        id: entry.id,
        alias,
        title: entry.title,
        plain: entry.plain,
        href: entry.href,
      });
    }
  }

  cachedTerms = terms.sort((a, b) => b.alias.length - a.alias.length);
  return cachedTerms;
}

function escapeRegex(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

function escapeAttr(s: string): string {
  return s
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;");
}

const SKIP_TAGS =
  /^(code|pre|a|script|style|kbd|glossary-term)$/i;

/** Injeta spans com tooltip nos segmentos de texto (fora de tags sensíveis). */
export function annotateGlossaryHtml(html: string): string {
  if (!html.trim()) return html;

  const terms = glossaryTerms();
  if (terms.length === 0) return html;

  const tagPattern = /(<\/?[a-z][^>]*>)/gi;
  const parts = html.split(tagPattern);
  const stack: string[] = [];

  return parts
    .map((part) => {
      if (part.startsWith("<")) {
        const open = part.match(/^<([a-z0-9-]+)/i);
        const close = part.match(/^<\/([a-z0-9-]+)/i);
        if (open) stack.push(open[1]!.toLowerCase());
        if (close) stack.pop();
        return part;
      }

      if (stack.some((t) => SKIP_TAGS.test(t))) return part;
      if (part.includes("glossary-term")) return part;

      let text = part;
      for (const term of terms) {
        const re = new RegExp(`(?<![\\w-])(${escapeRegex(term.alias)})(?![\\w-])`, "gi");
        text = text.replace(re, (match) => {
          const tip = escapeAttr(`${term.title}: ${term.plain}`);
          const href = escapeAttr(term.href);
          return `<span class="glossary-term" tabindex="0" role="term" data-href="${href}" aria-label="${tip}">${match}<span class="glossary-tip" role="tooltip">${tip}<a href="${href}">Saber mais</a></span></span>`;
        });
      }
      return text;
    })
    .join("");
}
