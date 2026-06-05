/**
 * content.test.ts — guarda a narrativa do site (SITE-001).
 *
 * A frase-ancora e os CTAs sao decisoes de negocio (spec). Testa-los evita que
 * uma edicao acidental quebre o posicionamento ou aponte o "Entrar" para o sitio
 * errado. Tambem garante que a camada Ops continua ESCONDIDA na v1.
 */
import { describe, it, expect } from "vitest";
import { HERO, PROOF_POINTS, APP_URL } from "./content";

describe("HERO", () => {
  it("usa a frase-ancora de plataforma", () => {
    expect(HERO.highlight).toBe("Tudo assenta no Cofre.");
    expect(HERO.title.toLowerCase()).toContain("identidade e segredos");
  });

  it("CTA primario aponta para a app (nao para uma rota interna)", () => {
    expect(HERO.primaryCta.href).toBe(APP_URL);
    expect(APP_URL.startsWith("https://app.")).toBe(true);
  });

  it("CTA secundario e uma ancora na propria pagina", () => {
    expect(HERO.secondaryCta.href.startsWith("#")).toBe(true);
  });

  it("nao expoe a camada Ops (RH/Financas/CRM/Mail) na v1", () => {
    const blob = JSON.stringify(HERO).toLowerCase();
    for (const banned of ["rh", "financas", "crm", "mail", "faturas"]) {
      expect(blob).not.toContain(banned);
    }
  });
});

describe("PROOF_POINTS", () => {
  it("tem 3 provas com valor e label", () => {
    expect(PROOF_POINTS).toHaveLength(3);
    for (const p of PROOF_POINTS) {
      expect(p.value.length).toBeGreaterThan(0);
      expect(p.label.length).toBeGreaterThan(0);
    }
  });
});
