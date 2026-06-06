import { describe, expect, it } from "vitest";
import { LOCALES } from "../config";
import { getCopy } from "../i18n";
import pt from "../i18n/locales/pt";

/** Garante paridade de chaves entre locales (i18n não fica a meio). */
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
    const keys = new Set(leafKeys(getCopy(locale)));
    expect(keys).toEqual(ptKeys);
  });

  it("frase-âncora PT contém Cofre", () => {
    expect(pt.hero.headline).toContain("Cofre");
  });
});
