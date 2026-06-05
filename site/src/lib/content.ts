/**
 * content.ts — fonte unica do copy do site institucional (SITE-001).
 *
 * Didatico: tal como a app deriva a navegacao de ROUTE_TREE, o site deriva o
 * texto daqui. Centralizar o copy: (1) facilita i18n futuro (P2), (2) permite
 * testar a narrativa, (3) evita strings soltas espalhadas pelos componentes.
 *
 * Decisoes da spec (docs/roadmap/09-design/site-institucional-spec.md):
 *   - Frase-ancora: visao de plataforma ("Tudo assenta no Cofre").
 *   - v1 esconde a camada Ops (RH/Fin/CRM) -> so Core + Workspace.
 */

export const APP_URL = "https://app.aegispass.com";

export interface Hero {
  eyebrow: string;
  title: string;
  highlight: string;
  subtitle: string;
  primaryCta: { label: string; href: string };
  secondaryCta: { label: string; href: string };
}

export const HERO: Hero = {
  eyebrow: "Zero-knowledge · BYOD · White-label",
  title: "A camada de identidade e segredos da tua empresa.",
  highlight: "Tudo assenta no Cofre.",
  subtitle:
    "Os segredos da tua equipa cifrados de ponta a ponta. O servidor nunca os ve — e o resto da operacao cresce a partir daqui.",
  primaryCta: { label: "Entrar", href: APP_URL },
  secondaryCta: { label: "Ver como funciona", href: "#como-funciona" },
};

/** Provas/numeros curtos sob o hero (ajustar com dados reais). */
export interface Proof {
  label: string;
  value: string;
}

export const PROOF_POINTS: Proof[] = [
  { value: "0", label: "segredos visiveis no servidor" },
  { value: "E2E", label: "cifragem de ponta a ponta" },
  { value: "BYOD", label: "dispositivos da equipa, sob controlo" },
];
