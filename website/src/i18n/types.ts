import type { Locale } from "../config";

export interface SiteCopy {
  meta: {
    title: string;
    description: string;
    ogTitle: string;
  };
  skip: string;
  nav: {
    zeroKnowledge: string;
    ecosystem: string;
    pillars: string;
    product: string;
    partners: string;
    enter: string;
  };
  hero: {
    eyebrow: string;
    headline: string;
    subline: string;
    ctaPrimary: string;
    ctaSecondary: string;
  };
  zk: {
    id: string;
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
  ecosystem: {
    id: string;
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
  pillars: {
    id: string;
    title: string;
    items: { fig: string; title: string; body: string }[];
  };
  product: {
    id: string;
    eyebrow: string;
    title: string;
    lead: string;
    mockVault: string;
    mockItems: string[];
    mockNavExtra: string[];
    mockStatus: string;
  };
  partners: {
    id: string;
    eyebrow: string;
    title: string;
    body: string;
    hint: string;
  };
  footer: {
    tagline: string;
    about: string;
    contact: string;
    legal: string;
    privacy: string;
    terms: string;
    app: string;
    rights: string;
  };
  lang: Record<Locale, string>;
}
