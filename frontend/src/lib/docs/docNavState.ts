const WIDTH_KEY = "aegispass-doc-nav-width";
const COLLAPSED_KEY = "aegispass-doc-nav-collapsed";

export const DOC_NAV_WIDTH_MIN = 168;
export const DOC_NAV_WIDTH_MAX = 420;
export const DOC_NAV_WIDTH_DEFAULT = 240;

/** Largura do índice de documentação (px) — persistida no browser. */
export function loadDocNavWidth(): number {
  if (typeof localStorage === "undefined") return DOC_NAV_WIDTH_DEFAULT;
  const raw = localStorage.getItem(WIDTH_KEY);
  if (!raw) return DOC_NAV_WIDTH_DEFAULT;
  const n = Number.parseInt(raw, 10);
  if (Number.isNaN(n)) return DOC_NAV_WIDTH_DEFAULT;
  return Math.min(DOC_NAV_WIDTH_MAX, Math.max(DOC_NAV_WIDTH_MIN, n));
}

export function saveDocNavWidth(width: number): void {
  if (typeof localStorage === "undefined") return;
  localStorage.setItem(WIDTH_KEY, String(Math.round(width)));
}

export function loadDocNavCollapsed(): boolean {
  if (typeof localStorage === "undefined") return false;
  return localStorage.getItem(COLLAPSED_KEY) === "1";
}

export function saveDocNavCollapsed(collapsed: boolean): void {
  if (typeof localStorage === "undefined") return;
  localStorage.setItem(COLLAPSED_KEY, collapsed ? "1" : "0");
}

export function clampDocNavWidth(width: number): number {
  return Math.min(DOC_NAV_WIDTH_MAX, Math.max(DOC_NAV_WIDTH_MIN, width));
}
