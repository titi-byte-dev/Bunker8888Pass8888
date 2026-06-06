import type { Locale } from "../config";
import type { SiteCopy } from "./types";
import pt from "./locales/pt";
import fr from "./locales/fr";
import es from "./locales/es";
import de from "./locales/de";

const bundles: Record<Locale, SiteCopy> = { pt, fr, es, de };

export function getCopy(locale: Locale): SiteCopy {
  return bundles[locale] ?? bundles.pt;
}

export function htmlLang(locale: Locale): string {
  const map: Record<Locale, string> = {
    pt: "pt-PT",
    fr: "fr",
    es: "es",
    de: "de",
  };
  return map[locale];
}

export type { SiteCopy };
