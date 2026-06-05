/**
 * Paletas de identidade visual (UI-013) — base do white-label por empresa.
 *
 * Didatico: separamos DOIS eixos ortogonais de aparencia.
 *
 *   modo    (theme.ts)   -> claro | escuro | sistema   ... attr data-theme
 *   paleta  (palette.ts) -> aegis | aurora | midnight | paper ... attr data-palette
 *
 * O CSS combina os dois: `[data-palette="aurora"][data-theme="dark"] { ... }`.
 * Assim cada empresa escolhe a SUA identidade (cor de marca) sem mexer no
 * layout nem perder o modo claro/escuro do colaborador.
 *
 * Seguranca/contraste: as paletas so re-tintam fundos + accent. Os tokens de
 * texto e de PERIGO (--color-danger) herdam do modo -> contraste WCAG AA fica
 * garantido em todas as combinacoes (ver palette.test.ts).
 */

export type PaletteId = "aegis" | "aurora" | "midnight" | "paper";

export const PALETTE_STORAGE_KEY = "aegis-palette";

export interface PaletteMeta {
  id: PaletteId;
  /** Nome de produto (PT-PT). */
  label: string;
  /** Uma linha de personalidade para o cartao em /settings. */
  personality: string;
  /** Cor de accent representativa (so para o swatch do cartao). */
  swatch: string;
}

/** Catalogo — a ORDEM aqui e a ordem dos cartoes em /settings. */
export const PALETTES: PaletteMeta[] = [
  {
    id: "aegis",
    label: "Aegis",
    personality: "Azul contido, confianca. A identidade padrao.",
    swatch: "#4da3ff",
  },
  {
    id: "aurora",
    label: "Aurora",
    personality: "Futurista suave — teal e violeta desaturado.",
    swatch: "#36c6d3",
  },
  {
    id: "midnight",
    label: "Midnight",
    personality: "Quase monocromatico, accent frio. Sessoes longas.",
    swatch: "#7c8db0",
  },
  {
    id: "paper",
    label: "Paper",
    personality: "Claro quente, sombras minimas. RH e documentos.",
    swatch: "#b5742a",
  },
];

const VALID = new Set<PaletteId>(PALETTES.map((p) => p.id));

/** Type guard — aceita so IDs conhecidos (defende contra storage corrompido). */
export function isPaletteId(value: unknown): value is PaletteId {
  return typeof value === "string" && VALID.has(value as PaletteId);
}

/**
 * Resolve a paleta efetiva.
 * Precedencia (white-label): escolha do colaborador > default da empresa > aegis.
 * `tenantDefault` chega do backend no futuro (config do tenant); por agora o
 * chamador passa undefined e usamos so localStorage.
 */
export function resolvePalette(stored: string | null, tenantDefault?: string): PaletteId {
  if (isPaletteId(stored)) return stored;
  if (isPaletteId(tenantDefault)) return tenantDefault;
  return "aegis";
}

/** Le a preferencia guardada (ou aegis). */
export function loadPalettePreference(tenantDefault?: string): PaletteId {
  const stored = typeof localStorage === "undefined" ? null : localStorage.getItem(PALETTE_STORAGE_KEY);
  return resolvePalette(stored, tenantDefault);
}

/** Aplica data-palette ao <html>. */
export function applyPalette(palette: PaletteId): void {
  if (typeof document === "undefined") return;
  document.documentElement.setAttribute("data-palette", palette);
}

/** Persiste a escolha do colaborador e aplica de imediato. */
export function setPalettePreference(palette: PaletteId): void {
  if (typeof localStorage !== "undefined") {
    localStorage.setItem(PALETTE_STORAGE_KEY, palette);
  }
  applyPalette(palette);
}

/** Arranque: le storage (com default do tenant) e aplica. */
export function initPalette(tenantDefault?: string): PaletteId {
  const palette = loadPalettePreference(tenantDefault);
  applyPalette(palette);
  return palette;
}

export function paletteLabel(palette: PaletteId): string {
  return PALETTES.find((p) => p.id === palette)?.label ?? "Aegis";
}
