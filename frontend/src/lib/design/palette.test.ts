import { describe, it, expect, beforeEach } from "vitest";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import {
  PALETTES,
  PALETTE_STORAGE_KEY,
  isPaletteId,
  resolvePalette,
  loadPalettePreference,
  setPalettePreference,
  initPalette,
  paletteLabel,
} from "./palette";

// Stubs minimos (env de teste e "node", sem jsdom) — chega para o round-trip.
function installDomStubs() {
  const store = new Map<string, string>();
  (globalThis as Record<string, unknown>).localStorage = {
    getItem: (k: string) => store.get(k) ?? null,
    setItem: (k: string, v: string) => void store.set(k, v),
    removeItem: (k: string) => void store.delete(k),
    clear: () => store.clear(),
  };
  const attrs = new Map<string, string>();
  (globalThis as Record<string, unknown>).document = {
    documentElement: {
      setAttribute: (k: string, v: string) => void attrs.set(k, v),
      getAttribute: (k: string) => attrs.get(k) ?? null,
      removeAttribute: (k: string) => void attrs.delete(k),
    },
  };
}

describe("palette — modelo", () => {
  beforeEach(() => {
    installDomStubs();
    localStorage.clear();
    document.documentElement.removeAttribute("data-palette");
  });

  it("isPaletteId aceita so IDs do catalogo", () => {
    expect(isPaletteId("aegis")).toBe(true);
    expect(isPaletteId("aurora")).toBe(true);
    expect(isPaletteId("ratoeira")).toBe(false);
    expect(isPaletteId(42)).toBe(false);
    expect(isPaletteId(null)).toBe(false);
  });

  it("resolvePalette respeita precedencia colaborador > tenant > aegis", () => {
    expect(resolvePalette("midnight", "paper")).toBe("midnight"); // escolha vence
    expect(resolvePalette(null, "paper")).toBe("paper"); // default da empresa
    expect(resolvePalette(null, undefined)).toBe("aegis"); // fallback
    expect(resolvePalette("lixo", "lixo-tambem")).toBe("aegis"); // storage corrompido
  });

  it("set/load fazem round-trip por localStorage", () => {
    setPalettePreference("aurora");
    expect(localStorage.getItem(PALETTE_STORAGE_KEY)).toBe("aurora");
    expect(loadPalettePreference()).toBe("aurora");
  });

  it("applyPalette/init escrevem data-palette no <html>", () => {
    setPalettePreference("paper");
    expect(document.documentElement.getAttribute("data-palette")).toBe("paper");
    localStorage.clear();
    initPalette("midnight"); // sem storage -> usa default do tenant
    expect(document.documentElement.getAttribute("data-palette")).toBe("midnight");
  });

  it("paletteLabel devolve nome de produto", () => {
    expect(paletteLabel("aurora")).toBe("Aurora");
    expect(paletteLabel("aegis")).toBe("Aegis");
  });
});

/* ---- Contraste WCAG AA (checkpoint design:accessibility-review) ----------
 * Validamos o par critico de cada bloco de paleta: accent-fg sobre accent
 * (texto de botoes primarios). Lemos tokens.css e exigimos racio >= 4.5.
 * Os tokens de texto/perigo herdam do modo (ja AA), por isso o risco real das
 * paletas esta no accent novo de cada uma.
 */
function srgbToLin(c: number): number {
  const s = c / 255;
  return s <= 0.03928 ? s / 12.92 : ((s + 0.055) / 1.055) ** 2.4;
}
function luminance(hex: string): number {
  const h = hex.replace("#", "");
  const r = parseInt(h.slice(0, 2), 16);
  const g = parseInt(h.slice(2, 4), 16);
  const b = parseInt(h.slice(4, 6), 16);
  return 0.2126 * srgbToLin(r) + 0.7152 * srgbToLin(g) + 0.0722 * srgbToLin(b);
}
function contrast(a: string, b: string): number {
  const la = luminance(a);
  const lb = luminance(b);
  const [hi, lo] = la > lb ? [la, lb] : [lb, la];
  return (hi + 0.05) / (lo + 0.05);
}

function readBlocks(): { selector: string; accent: string; fg: string }[] {
  const css = readFileSync(fileURLToPath(new URL("./tokens.css", import.meta.url)), "utf8");
  const blocks: { selector: string; accent: string; fg: string }[] = [];
  const re = /\[data-palette[^{]*\{([^}]*)\}/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(css))) {
    const body = m[1];
    const accent = /--color-accent:\s*(#[0-9a-fA-F]{6})/.exec(body)?.[1];
    const fg = /--color-accent-fg:\s*(#[0-9a-fA-F]{6})/.exec(body)?.[1];
    if (accent && fg) {
      blocks.push({ selector: m[0].slice(0, m[0].indexOf("{")).trim(), accent, fg });
    }
  }
  return blocks;
}

describe("palette — contraste AA do accent", () => {
  const blocks = readBlocks();

  it("encontra blocos de paleta em tokens.css", () => {
    // aurora/midnight/paper x (dark+light) = 6 blocos com accent definido.
    expect(blocks.length).toBeGreaterThanOrEqual(6);
  });

  it.each(readBlocks())("accent-fg sobre accent >= 4.5 em %s", ({ selector, accent, fg }) => {
    const ratio = contrast(accent, fg);
    expect(ratio, `${selector}: accent ${accent} / fg ${fg} = ${ratio.toFixed(2)}`).toBeGreaterThanOrEqual(4.5);
  });
});
