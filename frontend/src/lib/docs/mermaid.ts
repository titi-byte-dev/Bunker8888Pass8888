/**
 * Renderização lazy de diagramas Mermaid (DOC-008).
 *
 * Didático: import dinâmico evita ~200KB no bundle inicial; só carrega quando
 * abres uma página com fluxograma ou sequence diagram.
 */
let initTheme: string | null = null;

function currentTheme(): "dark" | "neutral" {
  if (typeof document === "undefined") return "neutral";
  const t = document.documentElement.getAttribute("data-theme");
  return t === "dark" ? "dark" : "neutral";
}

async function ensureMermaid() {
  const mermaid = (await import("mermaid")).default;
  const theme = currentTheme();
  if (initTheme !== theme) {
    mermaid.initialize({
      startOnLoad: false,
      theme,
      securityLevel: "strict",
      fontFamily: "var(--font-ui, system-ui, sans-serif)",
      sequence: {
        diagramMarginX: 12,
        diagramMarginY: 12,
        actorMargin: 48,
        messageMargin: 40,
      },
    });
    initTheme = theme;
  }
  return mermaid;
}

/** Renderiza código Mermaid dentro de um contentor e devolve o SVG injectado. */
export async function renderMermaid(
  container: HTMLElement,
  source: string,
  uniqueId: string,
): Promise<void> {
  const mermaid = await ensureMermaid();
  const renderId = `aegis-mmd-${uniqueId}-${crypto.randomUUID().slice(0, 8)}`;
  const { svg } = await mermaid.render(renderId, source);
  container.innerHTML = svg;
  const svgEl = container.querySelector("svg");
  if (svgEl) {
    svgEl.setAttribute("role", "img");
    svgEl.removeAttribute("height");
    svgEl.style.maxWidth = "100%";
    svgEl.style.height = "auto";
  }
}

/** Força re-inicialização quando o tema muda (light ↔ dark). */
export function resetMermaidTheme(): void {
  initTheme = null;
}
