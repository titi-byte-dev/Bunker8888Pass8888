/**
 * Ícones da navegação — SVG inline (sem dependência extra).
 * Fonte única para sidebar desktop e tab bar mobile.
 */

export type NavIconName =
  | "vault"
  | "security"
  | "work"
  | "team"
  | "hr"
  | "mail"
  | "fin"
  | "crm"
  | "admin"
  | "settings"
  | "child";

/** Mapeia href de módulo de topo → ícone. */
const MODULE_ICONS: Record<string, NavIconName> = {
  "/vault": "vault",
  "/security": "security",
  "/work": "work",
  "/team": "team",
  "/hr": "hr",
  "/mail": "mail",
  "/fin": "fin",
  "/crm": "crm",
  "/admin": "admin",
  "/settings": "settings",
};

export function iconForHref(href: string, isChild = false): NavIconName {
  if (isChild) return "child";
  for (const [prefix, icon] of Object.entries(MODULE_ICONS)) {
    if (href === prefix || href.startsWith(`${prefix}/`)) {
      return MODULE_ICONS[prefix] ?? "child";
    }
  }
  return "child";
}

/** Paths SVG (viewBox 0 0 24 24, stroke, sem fill). */
export const NAV_ICON_PATHS: Record<NavIconName, string> = {
  vault:
    "M12 2a5 5 0 0 0-5 5v2H6a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-9a2 2 0 0 0-2-2h-1V7a5 5 0 0 0-5-5zm-3 7V7a3 3 0 1 1 6 0v2H9z",
  security:
    "M12 2l8 4v6c0 5.25-3.5 9.74-8 11-4.5-1.26-8-5.75-8-11V6l8-4zm0 3.2L6 8.1v3.9c0 3.86 2.55 7.2 6 8.35 3.45-1.15 6-4.49 6-8.35V8.1l-6-2.9z",
  work: "M10 2h4a2 2 0 0 1 2 2v2h4v14H4V6h4V4a2 2 0 0 1 2-2zm2 4V4h-2v2h2zm-2 6h8v2H10v-2z",
  team: "M16 11c1.66 0 3-1.34 3-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 3-1.34 3-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5C15 14.17 10.33 13 8 13zm8 0c-.29 0-.62.02-.97.06 1.16.84 1.97 1.97 1.97 3.44V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z",
  hr: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zm-1 2l5 5h-5V4zM8 12h8v2H8v-2zm0 4h6v2H8v-2z",
  mail: "M20 4H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2zm0 4-8 5L4 8V6l8 5 8-5v2z",
  fin: "M4 19V5h2v14H4zm14 0V9h-4v2h2v8h2zM11 19V11H9v8h2zm5-10h2v10h-2V9z",
  crm: "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 15v-4H8l5-9v4h3l-5 9z",
  admin:
    "M12 1l3 6 6 .9-4.5 4.4 1 6.2L12 16.8 6.5 18.5l1-6.2L3 7.9 9 7l3-6zm0 4.2L10.2 9H7.6l2.1 2-0.8 2.9L12 12.8l2.1 1.1-0.8-2.9 2.1-2h-2.6L12 5.2z",
  settings:
    "M12 8a4 4 0 1 1 0 8 4 4 0 0 1 0-8zm8.94 3a7.96 7.96 0 0 0 .05-.94 7.96 7.96 0 0 0-.05-.94l2.03-1.58a.5.5 0 0 0 .12-.64l-1.92-3.32a.5.5 0 0 0-.6-.22l-2.39.96a7.28 7.28 0 0 0-1.62-.94l-.36-2.54A.5.5 0 0 0 14 2h-4a.5.5 0 0 0-.49.42l-.36 2.54c-.58.22-1.12.52-1.62.94l-2.39-.96a.5.5 0 0 0-.6.22L2.62 8.9a.5.5 0 0 0 .12.64L4.77 11.1c-.03.31-.05.63-.05.94s.02.63.05.94l-2.03 1.58a.5.5 0 0 0-.12.64l1.92 3.32a.5.5 0 0 0 .6.22l2.39-.96c.5.42 1.04.77 1.62.99l.36 2.54A.5.5 0 0 0 10 22h4a.5.5 0 0 0 .49-.42l.36-2.54a7.28 7.28 0 0 0 1.62-.99l2.39.96a.5.5 0 0 0 .6-.22l1.92-3.32a.5.5 0 0 0-.12-.64l-2.03-1.58z",
  child: "M8 12h8v2H8v-2z",
};
