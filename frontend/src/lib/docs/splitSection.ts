import type { DocFlow, DocSectionPart } from "./types";

const FLOW_RE = /<!--DOC_FLOW:([\w-]+)-->/g;

/**
 * Divide o HTML de uma secção em blocos de prosa e fluxos Mermaid.
 * Didático: o Markdown gera comentários `<!--DOC_FLOW:id-->`; o Svelte
 * monta componentes interactivos nos sítios certos — sem duplicar conteúdo.
 */
export function splitSectionHtml(html: string, flows: DocFlow[] = []): DocSectionPart[] {
  const byId = new Map(flows.map((f) => [f.id, f]));
  const parts: DocSectionPart[] = [];
  let last = 0;
  let match: RegExpExecArray | null;

  const re = new RegExp(FLOW_RE.source, "g");
  while ((match = re.exec(html)) !== null) {
    if (match.index > last) {
      const chunk = html.slice(last, match.index).trim();
      if (chunk) parts.push({ kind: "html", content: chunk });
    }
    const flow = byId.get(match[1]);
    if (flow) parts.push({ kind: "flow", flow });
    last = re.lastIndex;
  }

  const tail = html.slice(last).trim();
  if (tail) parts.push({ kind: "html", content: tail });

  // Só fallback se não havia placeholders Mermaid (evita mostrar comentários crus).
  const hasFlowMarkers = /<!--DOC_FLOW:[\w-]+-->/.test(html);
  if (parts.length === 0 && html.trim() && !hasFlowMarkers) {
    parts.push({ kind: "html", content: html });
  }

  return parts;
}
