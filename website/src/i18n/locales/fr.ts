import type { SiteCopy } from "./types";

const fr: SiteCopy = {
  meta: {
    title: "AegisPass — Identité et secrets pour votre PME",
    description:
      "La couche d'identité et de secrets de votre entreprise. Chiffrement zero-knowledge, BYOD et écosystème de travail — tout repose sur le Coffre.",
    ogTitle: "AegisPass — Tout repose sur le Coffre",
  },
  skip: "Aller au contenu",
  nav: {
    zeroKnowledge: "Zero-Knowledge",
    ecosystem: "Écosystème",
    pillars: "Principes",
    product: "Produit",
    partners: "Partenaires",
    enter: "Connexion",
  },
  hero: {
    eyebrow: "Plateforme d'identité pour PME",
    headline:
      "La couche d'identité et de secrets de votre entreprise. Tout repose sur le Coffre.",
    subline:
      "Mots de passe, équipe et travail sécurisé sur le même appareil — le serveur ne voit jamais vos secrets en clair.",
    ctaPrimary: "Ouvrir l'app",
    ctaSecondary: "Voir comment ça marche",
  },
  zk: {
    id: "zero-knowledge",
    title: "Zero-Knowledge, sans jargon",
    lead: "Chiffrement dans votre navigateur. Le serveur ne stocke que des blobs opaques — même avec accès à la base, personne ne lit vos données.",
    diagram: {
      client: "Navigateur",
      clientAction: "AES-GCM",
      server: "Serveur",
      serverAction: "stockage opaque",
      payload: "blob chiffré",
    },
    steps: [
      {
        title: "1. Déverrouillage local",
        body: "Le mot de passe maître dérive la clé uniquement sur votre appareil. Elle ne quitte jamais le réseau.",
      },
      {
        title: "2. Le serveur garde des blobs",
        body: "Métadonnées visibles (titre, type) oui ; contenu sensible — jamais en clair.",
      },
      {
        title: "3. Déchiffrement à la demande",
        body: "Sync, partage et agents respectent la même règle : moindre privilège.",
      },
    ],
    note: "Même nous, en tant qu'éditeurs, ne pouvons pas lire votre coffre.",
  },
  ecosystem: {
    id: "ecosystem",
    title: "Un écosystème en couches",
    lead: "Pas des outils isolés. Le Coffre est la base ; l'espace de travail s'appuie dessus.",
    coreLabel: "Core",
    coreHint: "Identité, secrets et confiance",
    coreItems: ["Coffre", "Équipe", "Sécurité"],
    workspaceLabel: "Workspace",
    workspaceHint: "Travail quotidien sur votre appareil",
    workspaceItems: [
      "Plages horaires et géofence",
      "Navigateur sandbox",
      "CLI mTLS",
      "Inventaire",
      "Google Workspace",
    ],
    footnote:
      "RH, finance et CRM existent en couche opérationnelle — intégrée, non exposée en v1 publique.",
  },
  pillars: {
    id: "pillars",
    title: "Conçu pour les PME modernes",
    items: [
      {
        fig: "0.1",
        title: "Confiance visible",
        body: "États de session, plage horaire et chiffrement lisibles — calme opérationnelle, sans fausses alertes.",
      },
      {
        fig: "0.2",
        title: "Agents à moindre privilège",
        body: "IA qui n'accède qu'au strict nécessaire, avec approbation humaine pour les actions sensibles.",
      },
      {
        fig: "0.3",
        title: "Vitesse silencieuse",
        body: "Interface fluide, clavier d'abord, densité intelligente — l'admin voit plus, l'employé l'essentiel.",
      },
    ],
  },
  product: {
    id: "product",
    eyebrow: "Produit en action",
    title: "Calme, dense et prête au quotidien",
    lead: "Aperçu de l'app réelle — coffre déverrouillé, liste lisible, états de sécurité visibles.",
    mockVault: "Coffre",
    mockItems: ["GitHub — ingénierie", "Abonnement SaaS", "Alias mail support"],
    mockNavExtra: ["Sécurité", "Équipe"],
    mockStatus: "Chiffré · Session active",
  },
  partners: {
    id: "partners",
    eyebrow: "Partenaires",
    title: "White-label pour MSP et intégrateurs",
    body: "Revendez la couche Core (Coffre + identité) sous votre marque. Workspace et agents s'appuient sur le même noyau zero-knowledge.",
    hint: "Canal partenaires — contact early access.",
  },
  footer: {
    tagline: "AegisPass — identité et secrets pour équipes BYOD.",
    about: "Qui sommes-nous",
    contact: "Contact",
    legal: "Mentions légales",
    privacy: "Confidentialité (bientôt)",
    terms: "Conditions (bientôt)",
    app: "Application",
    rights: "Tous droits réservés.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default fr;
