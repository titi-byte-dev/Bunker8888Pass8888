/** URL da app autenticada — override em build com PUBLIC_APP_URL. */
export const APP_URL =
  import.meta.env.PUBLIC_APP_URL ?? "https://app.aegispass.com";

export const SITE_ORIGIN = "https://aegispass.com";

export const LOCALES = ["pt", "fr", "es", "de"] as const;
export type Locale = (typeof LOCALES)[number];

export const DEFAULT_LOCALE: Locale = "pt";

/** Prefixo de path por locale (pt = raiz). */
export function localePath(locale: Locale): string {
  if (locale === DEFAULT_LOCALE) return "/";
  return `/${locale}/`;
}
