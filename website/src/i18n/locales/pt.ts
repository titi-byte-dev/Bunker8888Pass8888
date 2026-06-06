import type { SiteCopy } from "./types";

const pt: SiteCopy = {
  meta: {
    title: "AegisPass — Identidade e segredos para a tua PME",
    description:
      "A camada de identidade e segredos da tua empresa. Cifragem zero-knowledge, BYOD e ecossistema de trabalho — tudo assenta no Cofre.",
    ogTitle: "AegisPass — Tudo assenta no Cofre",
  },
  skip: "Saltar para o conteúdo",
  nav: {
    zeroKnowledge: "Zero-Knowledge",
    ecosystem: "Ecossistema",
    pillars: "Princípios",
    product: "Produto",
    partners: "Parceiros",
    enter: "Entrar",
  },
  hero: {
    eyebrow: "Plataforma de identidade para PME",
    headline:
      "A camada de identidade e segredos da tua empresa. Tudo assenta no Cofre.",
    subline:
      "Palavras-passe, equipa e trabalho seguro no mesmo dispositivo — o servidor nunca vê os teus segredos em claro.",
    ctaPrimary: "Entrar na app",
    ctaSecondary: "Ver como funciona",
  },
  zk: {
    id: "zero-knowledge",
    title: "Zero-Knowledge, explicado sem jargão",
    lead: "Cifras no teu browser. O servidor guarda apenas blobs opacos — mesmo com acesso à base de dados, ninguém lê os teus dados.",
    diagram: {
      client: "Browser",
      clientAction: "AES-GCM",
      server: "Servidor",
      serverAction: "armazena opaco",
      payload: "blob cifrado",
    },
    steps: [
      {
        title: "1. Desbloqueias localmente",
        body: "A Master Password deriva a chave só no teu dispositivo. Nunca sai para a rede.",
      },
      {
        title: "2. O servidor guarda blobs",
        body: "Metadados visíveis (título, tipo) sim; conteúdo sensível — nunca em claro.",
      },
      {
        title: "3. Decifras quando precisas",
        body: "Sincronização, partilha e agentes respeitam a mesma regra: mínimo privilégio.",
      },
    ],
    note: "Mesmo nós, enquanto fornecedores, não conseguimos ler o teu cofre.",
  },
  ecosystem: {
    id: "ecosystem",
    title: "Um ecossistema em camadas",
    lead: "Não são ferramentas soltas. O Cofre é a base; o Workspace de trabalho assenta nele.",
    coreLabel: "Core",
    coreHint: "Identidade, segredos e confiança",
    coreItems: ["Cofre", "Equipa", "Segurança"],
    workspaceLabel: "Workspace",
    workspaceHint: "Trabalho diário no teu dispositivo",
    workspaceItems: [
      "Turnos e geofence",
      "Browser sandbox",
      "CLI mTLS",
      "Inventário",
      "Google Workspace",
    ],
    footnote:
      "RH, finanças e CRM existem como camada operacional — integrada, não exposta na v1 pública.",
  },
  pillars: {
    id: "pillars",
    title: "Feito de propósito para PME modernas",
    items: [
      {
        fig: "0.1",
        title: "Confiança visível",
        body: "Estados de sessão, turno e cifragem legíveis — calma operacional, zero alarmismo falso.",
      },
      {
        fig: "0.2",
        title: "Agentes com menor privilégio",
        body: "IA que só acede ao estritamente necessário, com aprovação humana em acções sensíveis.",
      },
      {
        fig: "0.3",
        title: "Velocidade silenciosa",
        body: "Interface fluida, teclado-first, densidade inteligente — admin vê mais, funcionário vê o essencial.",
      },
    ],
  },
  product: {
    id: "product",
    eyebrow: "Produto a mexer",
    title: "Calma, densa e pronta para o dia-a-dia",
    lead: "Pré-visualização da app real — cofre desbloqueado, lista legível, estados de segurança visíveis.",
    mockVault: "Cofre",
    mockItems: ["GitHub — engenharia", "Faturação SaaS", "Alias mail suporte"],
    mockNavExtra: ["Segurança", "Equipa"],
    mockStatus: "Cifrado · Sessão activa",
  },
  partners: {
    id: "partners",
    eyebrow: "Parceiros",
    title: "White-label para MSPs e integradores",
    body: "Revende a camada Core (Cofre + identidade) com a tua marca. Workspace e agentes ligam-se ao mesmo núcleo zero-knowledge.",
    hint: "Canal de parceiros — contacto early access.",
  },
  footer: {
    tagline: "AegisPass — identidade e segredos para equipas BYOD.",
    about: "Quem somos",
    contact: "Contacto",
    legal: "Legal",
    privacy: "Privacidade (em breve)",
    terms: "Termos (em breve)",
    app: "App",
    rights: "Todos os direitos reservados.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default pt;
