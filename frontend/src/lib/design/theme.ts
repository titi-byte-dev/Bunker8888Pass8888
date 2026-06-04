/**
 * Gestão de tema light / dark / system (UI-001).
 *
 * Didático: guardamos a *preferência* do utilizador em localStorage; o tema
 * *resolvido* (light ou dark) aplica-se via data-theme no <html>.
 */

export type ThemeMode = "light" | "dark" | "system";
export type ResolvedTheme = "light" | "dark";

export const THEME_STORAGE_KEY = "aegis-theme";

/** Resolve system → light ou dark consoante preferência do SO. */
export function resolveTheme(mode: ThemeMode, prefersDark = false): ResolvedTheme {
  if (mode === "system") {
    return prefersDark ? "dark" : "light";
  }
  return mode;
}

/** Lê preferência guardada ou devolve "system". */
export function loadThemePreference(): ThemeMode {
  if (typeof localStorage === "undefined") return "system";
  const v = localStorage.getItem(THEME_STORAGE_KEY);
  if (v === "light" || v === "dark" || v === "system") return v;
  return "system";
}

/** Persiste preferência e aplica ao documento. */
export function setThemePreference(mode: ThemeMode): ResolvedTheme {
  if (typeof localStorage !== "undefined") {
    localStorage.setItem(THEME_STORAGE_KEY, mode);
  }
  const prefersDark =
    typeof window !== "undefined" &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;
  const resolved = resolveTheme(mode, prefersDark);
  applyResolvedTheme(resolved);
  return resolved;
}

/** Aplica data-theme + color-scheme no elemento raiz. */
export function applyResolvedTheme(resolved: ResolvedTheme): void {
  if (typeof document === "undefined") return;
  document.documentElement.setAttribute("data-theme", resolved);
}

/** Arranque: lê storage, resolve system, regista listener de mudança do SO. */
export function initTheme(): () => void {
  const mode = loadThemePreference();
  setThemePreference(mode);

  if (typeof window === "undefined") return () => {};

  const mq = window.matchMedia("(prefers-color-scheme: dark)");
  const onChange = () => {
    if (loadThemePreference() === "system") {
      applyResolvedTheme(mq.matches ? "dark" : "light");
    }
  };
  mq.addEventListener("change", onChange);
  return () => mq.removeEventListener("change", onChange);
}

/** Cicla light → dark → system (útil no playground até haver settings). */
export function cycleThemePreference(current: ThemeMode): ThemeMode {
  if (current === "light") return "dark";
  if (current === "dark") return "system";
  return "light";
}

export function themeModeLabel(mode: ThemeMode): string {
  switch (mode) {
    case "light":
      return "Claro";
    case "dark":
      return "Escuro";
    default:
      return "Sistema";
  }
}
