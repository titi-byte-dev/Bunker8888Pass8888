import { describe, expect, it } from "vitest";
import { LOCALES } from "../config";
import { getCopy } from "../i18n";
import pt from "./locales/pt";

function leafKeys(obj: unknown, prefix = ""): string[] {
  if (obj === null || typeof obj !== "object" || Array.isArray(obj)) {
    return prefix ? [prefix] : [];
  }
  return Object.entries(obj as Record<string, unknown>).flatMap(([k, v]) => {
    const path = prefix ? `${prefix}.${k}` : k;
    if (v !== null && typeof v === "object" && !Array.isArray(v)) {
      return leafKeys(v, path);
    }
    return [path];
  });
}

describe("i18n locales", () => {
  const ptKeys = new Set(leafKeys(pt));

  it.each(LOCALES)("locale %s tem as mesmas chaves que pt", (locale) => {
    expect(new Set(leafKeys(getCopy(locale)))).toEqual(ptKeys);
  });

  it("frase-âncora PT contém Cofre", () => {
    expect(pt.home.campaign.highlight).toContain("Cofre");
  });

  it("quatro produtos com meta e hero", () => {
    expect(Object.keys(pt.products)).toHaveLength(4);
    expect(pt.products.vault.meta.title.length).toBeGreaterThan(10);
  });

  it("três serviços com meta e hero", () => {
    expect(Object.keys(pt.services)).toHaveLength(3);
    expect(pt.services.agents.meta.title.length).toBeGreaterThan(10);
  });
});
