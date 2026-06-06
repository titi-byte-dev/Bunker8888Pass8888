import type { Locale } from "../config";
import type { ProductSlug, ServiceSlug, SiteCopy } from "./types";
import pt from "./locales/pt";
import fr from "./locales/fr";
import es from "./locales/es";
import de from "./locales/de";

const bundles: Record<Locale, SiteCopy> = { pt, fr, es, de };

export function getCopy(locale: Locale): SiteCopy {
  return bundles[locale] ?? bundles.pt;
}

export function getProduct(locale: Locale, slug: ProductSlug) {
  return getCopy(locale).products[slug];
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

export function getService(locale: Locale, slug: ServiceSlug) {
  return getCopy(locale).services[slug];
}

export type { SiteCopy, ProductSlug, ServiceSlug, PageStatus } from "./types";
export { PRODUCT_SLUGS, SERVICE_SLUGS, pagePath, productPath, servicePath } from "./types";
