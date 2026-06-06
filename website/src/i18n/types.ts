import type { Locale } from "../config";
import { DEFAULT_LOCALE } from "../config";

/** Produtos públicos (estilo Linear — cada um com subpágina). */
export const PRODUCT_SLUGS = ["vault", "security", "team", "workspace"] as const;
export type ProductSlug = (typeof PRODUCT_SLUGS)[number];

export type PageStatus = "live" | "preview" | "building";

export interface PageMeta {
  title: string;
  description: string;
}

export interface FeatureItem {
  title: string;
  body: string;
  status: PageStatus;
}

export interface ProductPageCopy {
  meta: PageMeta;
  pageStatus: PageStatus;
  hero: {
    eyebrow: string;
    headline: string;
    subline: string;
  };
  features: {
    title: string;
    items: FeatureItem[];
  };
  cta: {
    primary: string;
    secondary: string;
  };
}

export interface SiteCopy {
  site: {
    name: string;
    defaultMeta: PageMeta;
  };
  skip: string;
  nav: {
    products: string;
    platform: string;
    partners: string;
    enter: string;
    productLabels: Record<ProductSlug, string>;
  };
  home: {
    meta: PageMeta;
    campaign: {
      eyebrow: string;
      headline: string;
      highlight: string;
      subline: string;
      ctaPrimary: string;
      ctaSecondary: string;
    };
    proof: { value: string; label: string }[];
    products: {
      title: string;
      lead: string;
      explore: string;
      cards: Record<
        ProductSlug,
        { tagline: string; description: string; status: PageStatus }
      >;
    };
    platformTeaser: {
      title: string;
      body: string;
      link: string;
    };
    ctaBand: {
      title: string;
      body: string;
      primary: string;
    };
  };
  platform: {
    meta: PageMeta;
    pageStatus: PageStatus;
    hero: { eyebrow: string; headline: string; subline: string };
    zk: {
      title: string;
      lead: string;
      diagram: {
        client: string;
        clientAction: string;
        server: string;
        serverAction: string;
        payload: string;
      };
      steps: { title: string; body: string }[];
      note: string;
    };
    layers: {
      title: string;
      lead: string;
      coreLabel: string;
      coreHint: string;
      coreItems: string[];
      workspaceLabel: string;
      workspaceHint: string;
      workspaceItems: string[];
      footnote: string;
    };
  };
  products: Record<ProductSlug, ProductPageCopy>;
  partners: {
    meta: PageMeta;
    pageStatus: PageStatus;
    hero: { eyebrow: string; headline: string; subline: string };
    benefits: { title: string; items: string[] };
    contact: { label: string; email: string };
  };
  construction: {
    badge: string;
    title: string;
    body: string;
  };
  statusLabels: Record<PageStatus, string>;
  footer: {
    tagline: string;
    columns: {
      products: string;
      company: string;
      legal: string;
    };
    platform: string;
    partners: string;
    contact: string;
    privacy: string;
    terms: string;
    app: string;
    rights: string;
  };
  lang: Record<Locale, string>;
}

/** Path localizado: pagePath('fr','products','vault') → /fr/products/vault */
export function pagePath(locale: Locale, ...segments: string[]): string {
  const prefix = locale === DEFAULT_LOCALE ? "" : `/${locale}`;
  if (segments.length === 0) return prefix || "/";
  const tail = segments.join("/");
  return `${prefix}/${tail}`.replace(/\/+/g, "/");
}

export function productPath(locale: Locale, slug: ProductSlug): string {
  return pagePath(locale, "products", slug);
}
