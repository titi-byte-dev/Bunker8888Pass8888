/**
 * Carrega manifest e páginas geradas em build time.
 * Didático: import estático — o Vite inclui JSON no bundle; sem fetch extra em runtime.
 */
import manifest from "./generated/manifest.json";
import type { DocManifest, DocPage } from "./types";

export const DOC_MANIFEST = manifest as DocManifest;

/** Páginas geradas (exclui manifest.json) */
const pageModules = import.meta.glob<DocPage>("./generated/*.json", {
  eager: true,
  import: "default",
}) as Record<string, DocPage>;

/** Lista de metadados visíveis na app (Definições → Documentação) */
export function inAppDocs() {
  return DOC_MANIFEST.docs.filter((d) => d.in_app);
}

/** Agrupa documentos por categoria para a navegação lateral */
export function docsByCategory() {
  const groups = new Map<string, typeof DOC_MANIFEST.docs>();
  for (const doc of inAppDocs()) {
    const list = groups.get(doc.category) ?? [];
    list.push(doc);
    groups.set(doc.category, list);
  }
  for (const [, list] of groups) {
    list.sort((a, b) => a.order - b.order);
  }
  return DOC_MANIFEST.categories
    .map((cat) => ({
      ...cat,
      docs: groups.get(cat.id) ?? [],
    }))
    .filter((g) => g.docs.length > 0);
}

/** Carrega uma página pelo slug; undefined se não existir */
export function loadDocPage(slug: string): DocPage | undefined {
  const doc = pageModules[`./generated/${slug}.json`];
  return doc?.slug === slug ? doc : undefined;
}

/** Resolve slugs relacionados para links internos */
export function relatedDocs(slugs: string[]) {
  return slugs
    .map((s) => DOC_MANIFEST.docs.find((d) => d.slug === s))
    .filter((d): d is NonNullable<typeof d> => Boolean(d));
}
