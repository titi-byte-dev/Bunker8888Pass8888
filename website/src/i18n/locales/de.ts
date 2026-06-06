import type { SiteCopy } from "../types";

const de: SiteCopy = {
  site: {
    name: "AegisPass",
    defaultMeta: {
      title: "AegisPass — Identität und Geheimnisse für KMU",
      description:
        "Die Identitäts- und Geheimnisschicht für Ihr Unternehmen. Zero-Knowledge, BYOD und Arbeits-Ökosystem — alles basiert auf dem Tresor.",
    },
  },
  skip: "Zum Inhalt springen",
  nav: {
    products: "Produkte",
    services: "Services",
    platform: "Plattform",
    partners: "Partner",
    enter: "Anmelden",
    productLabels: {
      vault: "Tresor",
      security: "Sicherheit",
      team: "Team",
      workspace: "Workspace",
    },
    serviceLabels: {
      agents: "KI-Agenten",
      compliance: "Compliance",
      deployment: "Implementierung",
    },
  },
  home: {
    meta: {
      title: "AegisPass — Das Identitäts-Ökosystem für KMU",
      description:
        "Ein Zero-Knowledge-Tresor, sichtbare Sicherheit und BYOD — eine ruhige, dichte Plattform, die mit Ihnen wächst.",
    },
    campaign: {
      eyebrow: "Kampagne 2026 · Plattform im Aufbau",
      headline: "Die Identität Ihres KMU.",
      highlight: "Alles basiert auf dem Tresor.",
      subline:
        "Geheimnisse, Team und Arbeit auf demselben Gerät — Verschlüsselung, die der Server nie sieht. Wir bauen mit Ihnen; validieren Sie jedes Modul vor dem Launch.",
      ctaPrimary: "App öffnen",
      ctaSecondary: "Produkte entdecken",
    },
    proof: [
      { value: "0", label: "Geheimnisse lesbar auf dem Server" },
      { value: "E2E", label: "Zero-Knowledge-Verschlüsselung" },
      { value: "BYOD", label: "Privatgerät, Firmendaten" },
    ],
    products: {
      title: "Produkte, die auf dem Tresor basieren",
      lead: "Jedes Modul hat eine eigene Seite — erkunden, Feedback geben, Prioritäten setzen.",
      explore: "Entdecken",
      cards: {
        vault: {
          tagline: "Tresor",
          description: "Logins, Notizen und Karten im Browser verschlüsselt. Der Master Key verlässt nie das Gerät.",
          status: "live",
        },
        security: {
          tagline: "Sicherheit",
          description: "Hygiene, Geräte, Sentinel und Notfallzugriff — sichtbares Vertrauen ohne Alarmismus.",
          status: "preview",
        },
        team: {
          tagline: "Team",
          description: "Geteilte Tresore, ephemere Secret Links und Notizen, die nach dem Lesen verschwinden.",
          status: "preview",
        },
        workspace: {
          tagline: "Workspace",
          description: "Schichten, Geofence, Sandbox und CLI mTLS — klare Grenzen zwischen privat und beruflich.",
          status: "building",
        },
      },
    },
    services: {
      title: "Professionelle Services",
      lead: "Implementierung, DSGVO-Compliance und Agenten mit Human-in-the-loop — das Ökosystem ohne Prozess-Neuerfindung adoptieren.",
      explore: "Entdecken",
      allLink: "Alle Services ansehen",
    },
    platformTeaser: {
      title: "Zero-Knowledge-Plattform",
      body: "Wie Verschlüsselung funktioniert, Core + Workspace Schichten — und warum Zero-Knowledge Architektur ist, kein Marketing.",
      link: "Plattform ansehen",
    },
    ctaBand: {
      title: "Wir bauen in der Öffentlichkeit",
      body: "Die App läuft bereits in der Sandbox. Einloggen, Tresor testen — sagen Sie uns, was Ihrem KMU fehlt.",
      primary: "App öffnen",
    },
  },
  platform: {
    meta: {
      title: "Plattform — Zero-Knowledge erklärt | AegisPass",
      description: "Client-seitige Verschlüsselung, opake Blobs auf dem Server und Schicht-Ökosystem — die technische Basis von AegisPass.",
    },
    pageStatus: "live",
    hero: {
      eyebrow: "Plattform",
      headline: "Zero-Knowledge ist kein Slogan. Es ist ein Vertrag.",
      subline:
        "Der Server speichert Metadaten und unlesbare Blobs. Nur Ihr Gerät hält den Schlüssel — selbst wir können den Tresor nicht lesen.",
    },
    zk: {
      title: "So funktioniert es — in drei Schritten",
      lead: "Business-Sprache — kein Whitepaper. Dieselbe Logik wie bei Banken, für Ihr KMU.",
      diagram: {
        client: "Browser",
        clientAction: "AES-GCM",
        server: "Server",
        serverAction: "speichert opak",
        payload: "verschlüsselter Blob",
      },
      steps: [
        { title: "Lokal entsperren", body: "Master-Passwort leitet den Schlüssel mit Argon2id ab — nur im Browser-Speicher." },
        { title: "Server speichert Blobs", body: "Titel und Typ ja; Passwörter, IBAN und Notizen — nie im Klartext." },
        { title: "Entschlüsseln bei Bedarf", body: "Sync, Freigabe und Agenten folgen dem geringsten Privileg." },
      ],
      note: "Selbst mit Datenbankzugriff sieht ein Angreifer nur kryptographisches Rauschen.",
    },
    layers: {
      title: "Ökosystem in Schichten",
      lead: "Keine losen Werkzeuge. Der Tresor ist der Kern; der Workspace baut darauf auf.",
      coreLabel: "Core",
      coreHint: "Identität, Geheimnisse, Vertrauen",
      coreItems: ["Tresor", "Team", "Sicherheit"],
      workspaceLabel: "Workspace",
      workspaceHint: "Tägliche Arbeit auf dem Gerät",
      workspaceItems: ["Schichten", "Sandbox", "CLI mTLS", "Inventar", "Google"],
      footnote: "HR, Finanzen und CRM integrieren sich als operative Schicht — in der App sichtbar, nicht auf der öffentlichen v1-Website.",
    },
  },
  products: {
    vault: {
      meta: {
        title: "Tresor — Zero-Knowledge Vault | AegisPass",
        description: "Passwörter, Notizen und Karten client-seitig verschlüsselt. Echtzeit-Sync ohne Server-Inhalt.",
      },
      pageStatus: "live",
      hero: {
        eyebrow: "Produkt · Tresor",
        headline: "Das Herz des Ökosystems.",
        subline: "Speichern, synchronisieren und teilen — der Server sieht nie Master Key noch Klartext.",
      },
      features: {
        title: "Was Sie bereits validieren können",
        items: [
          { title: "AES-GCM im Browser", body: "Jedes Item ist ein Blob mit einmaligem Nonce. Argon2id für Master Key und separaten Auth Hash.", status: "live" },
          { title: "WebSocket-Echtzeit-Sync", body: "Änderungen propagieren; Backend leitet weiter, entschlüsselt nicht.", status: "live" },
          { title: "TOTP inline & Generator", body: "2FA und starke Passwörter ohne den Tresor zu verlassen.", status: "live" },
          { title: "Hygiene-Score", body: "Schwache und wiederverwendete Passwörter — nur client-seitig berechnet.", status: "preview" },
        ],
      },
      cta: { primary: "Tresor in der App öffnen", secondary: "Plattform ansehen" },
    },
    security: {
      meta: {
        title: "Sicherheit — Sichtbare Postur | AegisPass",
        description: "Hygiene, Geräte, Sentinel Mode und Notfallzugriff für anspruchsvolle KMU.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produkt · Sicherheit",
        headline: "Vertrauen, das man sieht.",
        subline: "Lesbare Sicherheitspostur — Sessions, Alerts und Step-up-Auth ohne Alarmschau.",
      },
      features: {
        title: "Fähigkeiten",
        items: [
          { title: "Sicherheitsgesundheit", body: "Hygiene-Score und empfohlene Aktionen — im Browser berechnet.", status: "preview" },
          { title: "Sentinel Mode", body: "Erkennt geografisch unmögliche Logins.", status: "preview" },
          { title: "Geräte & Passkeys", body: "Sessions widerrufen, Passkeys verwalten, Remote Wipe.", status: "building" },
          { title: "Notfallzugriff", body: "Digitaler Erbe mit Wartezeit.", status: "building" },
        ],
      },
      cta: { primary: "In der App erkunden", secondary: "Plattform ansehen" },
    },
    team: {
      meta: {
        title: "Team — Ephemere Freigabe | AegisPass",
        description: "Geteilte Tresore, Einmal-Secret-Links und vergängliche Notizen — Zero-Knowledge End-to-End.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produkt · Team",
        headline: "Teilen ohne Spuren.",
        subline: "Credentials für Freelancer, Secret Links im RAM und geteilte Tresore mit feinen Berechtigungen.",
      },
      features: {
        title: "Fähigkeiten",
        items: [
          { title: "Ephemere Secret Links", body: "Verschlüsselter Blob im RAM; Schlüssel im URL-Fragment — einmal lesen, Link tot.", status: "preview" },
          { title: "Shared Vaults", body: "Team-Tresore mit asymmetrischen Schlüsseln pro Nutzer.", status: "preview" },
          { title: "Selbstzerstörende Notizen", body: "Sichtbarer Countdown — kryptographisches «Verbrennen» nach TTL.", status: "building" },
          { title: "Verschlüsselte Anhänge", body: "Ein Blob pro Datei — Schlüssel nie auf dem Server.", status: "building" },
        ],
      },
      cta: { primary: "In der App testen", secondary: "Plattform ansehen" },
    },
    workspace: {
      meta: {
        title: "Workspace — Professionelles BYOD | AegisPass",
        description: "Schichten, Geofence, Sandbox-Browser und CLI mTLS — sicheres Arbeiten auf dem Privatgerät.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Produkt · Workspace",
        headline: "BYOD mit klaren Grenzen.",
        subline: "Außerhalb der Schicht wird der Tresor aus dem Speicher entfernt. Innerhalb der Regeln isolierte Arbeitswerkzeuge.",
      },
      features: {
        title: "Sichtbare Roadmap",
        items: [
          { title: "Schichten & Geofence", body: "Master Key nur innerhalb von Zeit und Ort.", status: "preview" },
          { title: "Sandbox-Browser", body: "Credential-Injection ohne Passwort an das OS.", status: "building" },
          { title: "CLI mTLS", body: "Geheimnisse für Pipelines und Dev — gegenseitige Auth.", status: "building" },
          { title: "Google Workspace ZK", body: "Proxy mit verschlüsselten Blobs für Drive und Docs.", status: "building" },
        ],
      },
      cta: { primary: "Fortschritt in der App", secondary: "Plattform ansehen" },
    },
  },
  productsIndex: {
    meta: {
      title: "Produkte — AegisPass-Ökosystem",
      description:
        "Tresor, Sicherheit, Team und Workspace — vier Zero-Knowledge-Module auf demselben Identitätskern.",
    },
    hero: {
      eyebrow: "Produkte",
      headline: "Vier Module. Ein Kern.",
      subline:
        "Jedes Produkt hat eine eigene Seite — Funktionen, Launch-Status erkunden und Feedback vor Go-live geben.",
    },
  },
  servicesIndex: {
    meta: {
      title: "Services — Implementierung und Compliance | AegisPass",
      description:
        "KI-Agenten, DSGVO-Audit und BYOD-Rollout — professionelle Services auf der Zero-Knowledge-Plattform.",
    },
    hero: {
      eyebrow: "Services",
      headline: "Vom Audit zum Betrieb.",
      subline:
        "Wir helfen KMU und Partnern beim Ökosystem-Adopt — klare Deliverables, Human-in-the-loop, ohne Secrets preiszugeben.",
    },
    explore: "Entdecken",
    cards: {
      agents: {
        tagline: "KI-Agenten",
        description:
          "Onboarding-Automatisierung, Ticket-Triage und interne Assistenten — immer mit menschlicher Freigabe und verschlüsselten Daten.",
        status: "preview",
      },
      compliance: {
        tagline: "Compliance",
        description:
          "DSGVO-Gap-Analyse, Verarbeitungsverzeichnis, Crypto-Shredding und Löschzertifikate — abgestimmt auf HR-Modul.",
        status: "preview",
      },
      deployment: {
        tagline: "Implementierung",
        description:
          "BYOD-Rollout, WireGuard, Multi-Tenant und Google-Workspace-Integration — vom Pilot bis Produktion.",
        status: "building",
      },
    },
  },
  services: {
    agents: {
      meta: {
        title: "KI-Agenten — Human-in-the-loop | AegisPass",
        description:
          "Assistenten und Automatisierung auf verschlüsselten Daten — Freigaben, Audit Trail und Least Privilege standardmäßig.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Service · KI-Agenten",
        headline: "Automatisierung ohne Kontrollverlust.",
        subline:
          "Agenten lesen Metadaten und Blobs nur clientseitig — sensible Aktionen durchlaufen eine Freigabe-Warteschlange.",
      },
      features: {
        title: "Deliverables",
        items: [
          { title: "Assistiertes Onboarding", body: "Checklisten, Aliase und Start-Tresore — Agent schlägt vor, Admin genehmigt.", status: "preview" },
          { title: "Anfragen-Triage", body: "Interne Tickets klassifizieren ohne Secrets an externe LLMs zu exportieren.", status: "preview" },
          { title: "Unveränderlicher Audit Trail", body: "Jede Agentenaktion verkettet im Log — auditbereit.", status: "building" },
          { title: "Richtlinien pro Tenant", body: "Scope-Grenzen, Zeiten und Geofence vor Ausführung.", status: "building" },
        ],
      },
      cta: { primary: "Early Access anfragen", secondary: "Plattform ansehen" },
    },
    compliance: {
      meta: {
        title: "DSGVO-Compliance — Audit | AegisPass",
        description:
          "Gap-Analyse, Verarbeitungsverzeichnis und kryptographische Löschnachweise — DSGVO by design.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Service · Compliance",
        headline: "DSGVO, die man beweisen kann.",
        subline:
          "Wir kartieren Verarbeitungen, minimieren Daten und dokumentieren Crypto-Shredding — ohne generische Slides.",
      },
      features: {
        title: "Deliverables",
        items: [
          { title: "Initiale Gap-Analyse", body: "Fluss-Inventar, Rechtsgrundlagen und Risiken — für KMU priorisiert.", status: "preview" },
          { title: "Verarbeitungsverzeichnis", body: "Art. 30 abgestimmt auf HR- und Tresor-Module — sensible Felder verschlüsselt.", status: "preview" },
          { title: "Löschzertifikat", body: "Kryptographischer Erasure-Nachweis nach Crypto-Shredding.", status: "building" },
          { title: "DPO-light-Schulung", body: "Workshops für kleine Teams — Business-Sprache, kein Juristenjargon.", status: "building" },
        ],
      },
      cta: { primary: "Gespräch vereinbaren", secondary: "Plattform ansehen" },
    },
    deployment: {
      meta: {
        title: "Implementierung — BYOD und Tenant | AegisPass",
        description:
          "Pilot, WireGuard, Multi-Tenant und Integrationen — sicherer Rollout des AegisPass-Ökosystems.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Service · Implementierung",
        headline: "Vom Pilot zur Produktion.",
        subline:
          "Wir konfigurieren Tenant, BYOD-Richtlinien und Netz — mit Runbooks und Übergabe an Ihr Team oder MSP-Partner.",
      },
      features: {
        title: "Typische Phasen",
        items: [
          { title: "Sandbox-Pilot", body: "10–30 Nutzer, live Tresor, Adoptionsmetriken.", status: "preview" },
          { title: "WireGuard-Netz", body: "Zero-Trust-Tunnel zwischen Geräten und API — ohne Legacy-VPN.", status: "building" },
          { title: "Multi-Tenant + RLS", body: "PostgreSQL-Isolation mit negativen Tests validiert.", status: "building" },
          { title: "Workspace-Integration", body: "Google, Aliase und mTLS-CLI — Roadmap am Produkt ausgerichtet.", status: "building" },
        ],
      },
      cta: { primary: "Team kontaktieren", secondary: "Partner ansehen" },
    },
  },
  partners: {
    meta: {
      title: "Partner — White-label | AegisPass",
      description: "MSPs und Integratoren: Core-Schicht unter Ihrer Marke auf echtem Zero-Knowledge verkaufen.",
    },
    pageStatus: "preview",
    hero: {
      eyebrow: "Partner",
      headline: "Ihre Marke. Unser Kern.",
      subline: "White-label Tresor und Identität — Workspace und Agenten auf demselben Ökosystem. Early-Access-Kanal offen.",
    },
    benefits: {
      title: "Warum Partner werden",
      items: [
        "Zero-Knowledge-Core verkaufsfertig — keine eigene Kryptographie.",
        "Paletten und White-label-Tenant in der App (Site-Roadmap v2).",
        "Didaktische In-App-Doku — weniger L1-Support.",
        "Agenten mit Human-in-the-loop — Automatisierung mit Kontrolle.",
      ],
    },
    contact: { label: "Partner Early Access", email: "partners@aegispass.com" },
  },
  construction: {
    badge: "Im Aufbau",
    title: "Wir verfeinern diese Seite",
    body: "Der Kern läuft bereits in der App-Sandbox — diese öffentliche Seite folgt dem Launch. Ihr Feedback steuert unsere Roadmap.",
  },
  statusLabels: { live: "Verfügbar", preview: "Vorschau", building: "Im Aufbau" },
  footer: {
    tagline: "Identität und Geheimnisse für BYOD-Teams.",
    columns: { products: "Produkte", services: "Services", company: "Unternehmen", legal: "Rechtliches" },
    platform: "Plattform",
    partners: "Partner",
    contact: "Kontakt",
    privacy: "Datenschutz (demnächst)",
    terms: "AGB (demnächst)",
    app: "App",
    rights: "Alle Rechte vorbehalten.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default de;
