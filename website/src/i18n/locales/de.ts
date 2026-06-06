import type { SiteCopy } from "./types";

const de: SiteCopy = {
  meta: {
    title: "AegisPass — Identität und Geheimnisse für Ihr KMU",
    description:
      "Die Identitäts- und Geheimnisschicht für Ihr Unternehmen. Zero-Knowledge-Verschlüsselung, BYOD und Arbeits-Ökosystem — alles basiert auf dem Tresor.",
    ogTitle: "AegisPass — Alles basiert auf dem Tresor",
  },
  skip: "Zum Inhalt springen",
  nav: {
    zeroKnowledge: "Zero-Knowledge",
    ecosystem: "Ökosystem",
    pillars: "Prinzipien",
    product: "Produkt",
    partners: "Partner",
    enter: "Anmelden",
  },
  hero: {
    eyebrow: "Identitätsplattform für KMU",
    headline:
      "Die Identitäts- und Geheimnisschicht für Ihr Unternehmen. Alles basiert auf dem Tresor.",
    subline:
      "Passwörter, Team und sicheres Arbeiten auf demselben Gerät — der Server sieht Ihre Geheimnisse nie im Klartext.",
    ctaPrimary: "App öffnen",
    ctaSecondary: "So funktioniert es",
  },
  zk: {
    id: "zero-knowledge",
    title: "Zero-Knowledge, ohne Fachjargon",
    lead: "Verschlüsselung im Browser. Der Server speichert nur undurchsichtige Blobs — selbst mit Datenbankzugriff liest niemand Ihre Daten.",
    diagram: {
      client: "Browser",
      clientAction: "AES-GCM",
      server: "Server",
      serverAction: "speichert opak",
      payload: "verschlüsselter Blob",
    },
    steps: [
      {
        title: "1. Lokal entsperren",
        body: "Das Master-Passwort leitet den Schlüssel nur auf Ihrem Gerät ab. Er verlässt nie das Netzwerk.",
      },
      {
        title: "2. Der Server speichert Blobs",
        body: "Sichtbare Metadaten (Titel, Typ) ja; sensibler Inhalt — nie im Klartext.",
      },
      {
        title: "3. Entschlüsseln bei Bedarf",
        body: "Sync, Freigabe und Agenten folgen derselben Regel: geringstes Privileg.",
      },
    ],
    note: "Selbst wir als Anbieter können Ihren Tresor nicht lesen.",
  },
  ecosystem: {
    id: "ecosystem",
    title: "Ein Ökosystem in Schichten",
    lead: "Keine losen Werkzeuge. Der Tresor ist die Basis; der Arbeitsbereich baut darauf auf.",
    coreLabel: "Core",
    coreHint: "Identität, Geheimnisse und Vertrauen",
    coreItems: ["Tresor", "Team", "Sicherheit"],
    workspaceLabel: "Workspace",
    workspaceHint: "Tägliche Arbeit auf Ihrem Gerät",
    workspaceItems: [
      "Schichten und Geofence",
      "Sandbox-Browser",
      "CLI mTLS",
      "Inventar",
      "Google Workspace",
    ],
    footnote:
      "HR, Finanzen und CRM existieren als operative Schicht — integriert, in v1 nicht öffentlich.",
  },
  pillars: {
    id: "pillars",
    title: "Für moderne KMU gebaut",
    items: [
      {
        fig: "0.1",
        title: "Sichtbares Vertrauen",
        body: "Session-, Schicht- und Verschlüsselungsstatus lesbar — operative Ruhe, kein falsches Alarmismus.",
      },
      {
        fig: "0.2",
        title: "Agenten mit geringstem Privileg",
        body: "KI greift nur auf das Nötigste zu — mit menschlicher Freigabe bei sensiblen Aktionen.",
      },
      {
        fig: "0.3",
        title: "Stille Geschwindigkeit",
        body: "Flüssige UI, Tastatur zuerst, intelligente Dichte — Admin sieht mehr, Mitarbeiter das Wesentliche.",
      },
    ],
  },
  product: {
    id: "product",
    eyebrow: "Produkt in Aktion",
    title: "Ruhig, dicht und alltagstauglich",
    lead: "Vorschau der echten App — Tresor entsperrt, lesbare Liste, sichtbare Sicherheitszustände.",
    mockVault: "Tresor",
    mockItems: ["GitHub — Engineering", "SaaS-Abo", "Mail-Alias Support"],
    mockNavExtra: ["Sicherheit", "Team"],
    mockStatus: "Verschlüsselt · Session aktiv",
  },
  partners: {
    id: "partners",
    eyebrow: "Partner",
    title: "White-label für MSPs und Integratoren",
    body: "Verkaufen Sie die Core-Schicht (Tresor + Identität) unter Ihrer Marke. Workspace und Agenten nutzen denselben Zero-Knowledge-Kern.",
    hint: "Partnerkanal — Early-Access-Kontakt.",
  },
  footer: {
    tagline: "AegisPass — Identität und Geheimnisse für BYOD-Teams.",
    about: "Über uns",
    contact: "Kontakt",
    legal: "Rechtliches",
    privacy: "Datenschutz (demnächst)",
    terms: "AGB (demnächst)",
    app: "App",
    rights: "Alle Rechte vorbehalten.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default de;
