import type { SiteCopy } from "../types";

const fr: SiteCopy = {
  site: {
    name: "AegisPass",
    defaultMeta: {
      title: "AegisPass — Identité et secrets pour PME",
      description:
        "La couche d'identité et de secrets de votre entreprise. Zero-knowledge, BYOD et écosystème de travail — tout repose sur le Coffre.",
    },
  },
  skip: "Aller au contenu",
  nav: {
    products: "Produits",
    platform: "Plateforme",
    partners: "Partenaires",
    enter: "Connexion",
    productLabels: {
      vault: "Coffre",
      security: "Sécurité",
      team: "Équipe",
      workspace: "Workspace",
    },
  },
  home: {
    meta: {
      title: "AegisPass — L'écosystème d'identité pour PME",
      description:
        "Un coffre zero-knowledge, une sécurité lisible et le BYOD — une plateforme calme, dense et faite pour grandir avec vous.",
    },
    campaign: {
      eyebrow: "Campagne 2026 · Plateforme en construction",
      headline: "L'identité de votre PME.",
      highlight: "Tout repose sur le Coffre.",
      subline:
        "Secrets, équipe et travail sur le même appareil — chiffrement que le serveur ne voit jamais. Nous construisons avec vous ; validez chaque module avant le lancement.",
      ctaPrimary: "Ouvrir l'app",
      ctaSecondary: "Explorer les produits",
    },
    proof: [
      { value: "0", label: "secret lisible côté serveur" },
      { value: "E2E", label: "chiffrement zero-knowledge" },
      { value: "BYOD", label: "appareil perso, données pro" },
    ],
    products: {
      title: "Des produits ancrés sur le Coffre",
      lead: "Chaque module a sa page — explorez, donnez votre avis, priorisez avec nous.",
      explore: "Explorer",
      cards: {
        vault: {
          tagline: "Coffre",
          description: "Identifiants, notes et cartes chiffrés dans le navigateur. La Master Key ne quitte jamais l'appareil.",
          status: "live",
        },
        security: {
          tagline: "Sécurité",
          description: "Hygiène, appareils, Sentinel et accès d'urgence — confiance visible, sans alarmisme.",
          status: "preview",
        },
        team: {
          tagline: "Équipe",
          description: "Coffres partagés, liens secrets éphémères et notes qui disparaissent après lecture.",
          status: "preview",
        },
        workspace: {
          tagline: "Workspace",
          description: "Plages horaires, géofence, sandbox et CLI mTLS — frontières nettes entre vie perso et pro.",
          status: "building",
        },
      },
    },
    platformTeaser: {
      title: "Plateforme zero-knowledge",
      body: "Comment fonctionne le chiffrement, les couches Core + Workspace, et pourquoi le zero-knowledge est une architecture — pas un slogan.",
      link: "Voir la plateforme",
    },
    ctaBand: {
      title: "Nous construisons en public",
      body: "L'app tourne déjà en sandbox. Connectez-vous, testez le coffre et dites-nous ce qui manque à votre PME.",
      primary: "Ouvrir l'app",
    },
  },
  platform: {
    meta: {
      title: "Plateforme — Zero-knowledge expliqué | AegisPass",
      description:
        "Chiffrement côté client, blobs opaques côté serveur et écosystème en couches — la base technique d'AegisPass.",
    },
    pageStatus: "live",
    hero: {
      eyebrow: "Plateforme",
      headline: "Zero-knowledge n'est pas un slogan. C'est un contrat.",
      subline:
        "Le serveur stocke métadonnées et blobs illisibles. Seul votre appareil détient la clé — même nous ne pouvons pas lire le coffre.",
    },
    zk: {
      title: "Comment ça marche, en trois étapes",
      lead: "Langage métier — sans whitepaper. La même logique que les banques, adaptée à votre PME.",
      diagram: {
        client: "Navigateur",
        clientAction: "AES-GCM",
        server: "Serveur",
        serverAction: "stockage opaque",
        payload: "blob chiffré",
      },
      steps: [
        {
          title: "Déverrouillage local",
          body: "Le mot de passe maître dérive la clé avec Argon2id — uniquement en mémoire navigateur.",
        },
        {
          title: "Le serveur garde des blobs",
          body: "Titres et types oui ; mots de passe, IBAN et notes — jamais en clair.",
        },
        {
          title: "Déchiffrement à la demande",
          body: "Sync, partage et agents respectent le moindre privilège.",
        },
      ],
      note: "Même avec accès à la base, un attaquant ne voit que du bruit cryptographique.",
    },
    layers: {
      title: "Écosystème en couches",
      lead: "Pas d'outils isolés. Le Coffre est le noyau ; le Workspace s'appuie dessus.",
      coreLabel: "Core",
      coreHint: "Identité, secrets, confiance",
      coreItems: ["Coffre", "Équipe", "Sécurité"],
      workspaceLabel: "Workspace",
      workspaceHint: "Travail quotidien sur l'appareil",
      workspaceItems: ["Plages horaires", "Sandbox", "CLI mTLS", "Inventaire", "Google"],
      footnote: "RH, finance et CRM s'intègrent en couche opérationnelle — visible dans l'app, pas sur le site public v1.",
    },
  },
  products: {
    vault: {
      meta: {
        title: "Coffre — Vault zero-knowledge | AegisPass",
        description: "Mots de passe, notes et cartes chiffrés côté client. Sync temps réel sans contenu visible serveur.",
      },
      pageStatus: "live",
      hero: {
        eyebrow: "Produit · Coffre",
        headline: "Le cœur de l'écosystème.",
        subline: "Stockez, synchronisez et partagez des secrets — le serveur ne voit jamais la Master Key ni le contenu en clair.",
      },
      features: {
        title: "Ce que vous pouvez déjà valider",
        items: [
          { title: "Chiffrement AES-GCM dans le navigateur", body: "Chaque item est un blob à nonce unique. Argon2id pour Master Key et Auth Hash distinct.", status: "live" },
          { title: "Sync WebSocket temps réel", body: "Les changements se propagent ; le backend relaie, ne déchiffre pas.", status: "live" },
          { title: "TOTP inline et générateur", body: "2FA et mots de passe forts sans quitter le coffre.", status: "live" },
          { title: "Score d'hygiène", body: "Mots de passe faibles et réutilisés — calculés uniquement côté client.", status: "preview" },
        ],
      },
      cta: { primary: "Ouvrir le coffre", secondary: "Voir la plateforme" },
    },
    security: {
      meta: {
        title: "Sécurité — Posture visible | AegisPass",
        description: "Hygiène, appareils, Sentinel Mode et accès d'urgence pour PME exigeantes.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produit · Sécurité",
        headline: "Une confiance qui se voit.",
        subline: "Posture de sécurité lisible — sessions, alertes et step-up auth sans théâtre d'alarmes.",
      },
      features: {
        title: "Capacités",
        items: [
          { title: "Santé de sécurité", body: "Score d'hygiène et actions recommandées — calculés dans le navigateur.", status: "preview" },
          { title: "Sentinel Mode", body: "Détecte les connexions géographiquement impossibles.", status: "preview" },
          { title: "Appareils et passkeys", body: "Révoquez sessions, gérez passkeys, remote wipe.", status: "building" },
          { title: "Accès d'urgence", body: "Héritier digital avec délai d'attente.", status: "building" },
        ],
      },
      cta: { primary: "Explorer dans l'app", secondary: "Voir la plateforme" },
    },
    team: {
      meta: {
        title: "Équipe — Partage éphémère | AegisPass",
        description: "Coffres partagés, liens secrets à usage unique et notes éphémères — zero-knowledge bout en bout.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produit · Équipe",
        headline: "Partager sans laisser de trace.",
        subline: "Identifiants pour freelances, liens secrets en RAM et coffres partagés à permissions fines.",
      },
      features: {
        title: "Capacités",
        items: [
          { title: "Liens secrets éphémères", body: "Blob chiffré en RAM ; clé dans le fragment URL — lu une fois, lien mort.", status: "preview" },
          { title: "Shared Vaults", body: "Coffres d'équipe avec clés asymétriques par utilisateur.", status: "preview" },
          { title: "Notes auto-destructives", body: "Compte à rebours visible — « destruction » cryptographique après TTL.", status: "building" },
          { title: "Pièces jointes chiffrées", body: "Un blob par fichier — clé jamais côté serveur.", status: "building" },
        ],
      },
      cta: { primary: "Tester dans l'app", secondary: "Voir la plateforme" },
    },
    workspace: {
      meta: {
        title: "Workspace — BYOD professionnel | AegisPass",
        description: "Plages horaires, géofence, sandbox navigateur et CLI mTLS — travail sécurisé sur appareil personnel.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Produit · Workspace",
        headline: "BYOD aux frontières nettes.",
        subline: "Hors plage horaire, le coffre est expurgé de la mémoire. Dans les règles, outils de travail isolés.",
      },
      features: {
        title: "Roadmap visible",
        items: [
          { title: "Plages et géofence", body: "Master Key seulement dans horaire et lieu autorisés.", status: "preview" },
          { title: "Sandbox navigateur", body: "Injection d'identifiants sans révéler le mot de passe à l'OS.", status: "building" },
          { title: "CLI mTLS", body: "Secrets pour pipelines et dev — authentification mutuelle.", status: "building" },
          { title: "Google Workspace ZK", body: "Proxy avec blobs chiffrés pour Drive et Docs.", status: "building" },
        ],
      },
      cta: { primary: "Voir l'avancement", secondary: "Voir la plateforme" },
    },
  },
  partners: {
    meta: {
      title: "Partenaires — White-label | AegisPass",
      description: "MSP et intégrateurs : revendez la couche Core sous votre marque, sur du zero-knowledge réel.",
    },
    pageStatus: "preview",
    hero: {
      eyebrow: "Partenaires",
      headline: "Votre marque. Notre noyau.",
      subline: "White-label du Coffre et de l'identité — Workspace et agents sur le même écosystème. Canal early access ouvert.",
    },
    benefits: {
      title: "Pourquoi devenir partenaire",
      items: [
        "Core zero-knowledge prêt à revendre — ne réinventez pas le chiffrement.",
        "Palettes et tenant white-label dans l'app (roadmap v2 site).",
        "Documentation in-app didactique — moins de support L1.",
        "Agents avec human-in-the-loop — automatisation maîtrisée.",
      ],
    },
    contact: { label: "Early access partenaires", email: "partners@aegispass.com" },
  },
  construction: {
    badge: "En construction",
    title: "Nous peaufinons cette page",
    body: "Le cœur fonctionne déjà en sandbox app — cette page publique suit le lancement. Vos retours orientent notre roadmap.",
  },
  statusLabels: { live: "Disponible", preview: "Aperçu", building: "En construction" },
  footer: {
    tagline: "Identité et secrets pour équipes BYOD.",
    columns: { products: "Produits", company: "Entreprise", legal: "Mentions" },
    platform: "Plateforme",
    partners: "Partenaires",
    contact: "Contact",
    privacy: "Confidentialité (bientôt)",
    terms: "Conditions (bientôt)",
    app: "Application",
    rights: "Tous droits réservés.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default fr;
