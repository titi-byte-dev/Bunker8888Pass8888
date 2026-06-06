import type { SiteCopy } from "../types";

const pt: SiteCopy = {
  site: {
    name: "AegisPass",
    defaultMeta: {
      title: "AegisPass — Identidade e segredos para PME",
      description:
        "A camada de identidade e segredos da tua empresa. Zero-knowledge, BYOD e ecossistema de trabalho — tudo assenta no Cofre.",
    },
  },
  skip: "Saltar para o conteúdo",
  nav: {
    products: "Produtos",
    services: "Serviços",
    platform: "Plataforma",
    partners: "Parceiros",
    enter: "Entrar",
    productLabels: {
      vault: "Cofre",
      security: "Segurança",
      team: "Equipa",
      workspace: "Workspace",
    },
    serviceLabels: {
      agents: "Agentes IA",
      compliance: "Conformidade",
      deployment: "Implementação",
    },
  },
  home: {
    meta: {
      title: "AegisPass — O ecossistema de identidade para PME",
      description:
        "Um cofre zero-knowledge, segurança visível e trabalho BYOD — numa plataforma calma, densa e pronta para crescer contigo.",
    },
    campaign: {
      eyebrow: "Campanha 2026 · Plataforma em construção",
      headline: "A identidade da tua PME.",
      highlight: "Tudo assenta no Cofre.",
      subline:
        "Segredos, equipa e trabalho no mesmo dispositivo — com cifragem que o servidor nunca vê. Construímos contigo; valida cada módulo antes do lançamento.",
      ctaPrimary: "Entrar na app",
      ctaSecondary: "Explorar produtos",
    },
    proof: [
      { value: "0", label: "segredos legíveis no servidor" },
      { value: "E2E", label: "cifragem zero-knowledge" },
      { value: "BYOD", label: "dispositivo pessoal, dados profissionais" },
    ],
    products: {
      title: "Produtos que assentam no Cofre",
      lead: "Cada módulo é uma página — explora, dá feedback, ajuda-nos a priorizar.",
      explore: "Explorar",
      cards: {
        vault: {
          tagline: "Cofre",
          description: "Logins, notas e cartões cifrados no browser. A Master Key nunca sai do dispositivo.",
          status: "live",
        },
        security: {
          tagline: "Segurança",
          description: "Higiene, dispositivos, Sentinel e acesso de emergência — confiança visível, sem alarmismo.",
          status: "preview",
        },
        team: {
          tagline: "Equipa",
          description: "Cofres partilhados, secret links efémeros e notas que desaparecem após leitura.",
          status: "preview",
        },
        workspace: {
          tagline: "Workspace",
          description: "Turnos, geofence, sandbox e CLI mTLS — fronteiras claras entre vida pessoal e profissional.",
          status: "building",
        },
      },
    },
    services: {
      title: "Serviços profissionais",
      lead: "Implementação, conformidade RGPD e agentes com human-in-the-loop — para adoptar o ecossistema sem reinventar processos.",
      explore: "Explorar",
      allLink: "Ver todos os serviços",
    },
    platformTeaser: {
      title: "Plataforma zero-knowledge",
      body: "Como funciona a cifragem, as camadas Core + Workspace e porque «zero-knowledge» não é marketing — é arquitectura.",
      link: "Ver plataforma",
    },
    ctaBand: {
      title: "Construímos em público",
      body: "A app já corre em sandbox. Entra, testa o cofre e diz-nos o que falta para a tua PME.",
      primary: "Abrir app",
    },
  },
  platform: {
    meta: {
      title: "Plataforma — Zero-knowledge explicado | AegisPass",
      description:
        "Cifragem no cliente, blobs opacos no servidor e ecossistema em camadas — a base técnica do AegisPass.",
    },
    pageStatus: "live",
    hero: {
      eyebrow: "Plataforma",
      headline: "Zero-knowledge não é slogan. É contrato.",
      subline:
        "O servidor armazena metadados e blobs ilegíveis. Só o teu dispositivo detém a chave — mesmo nós não conseguimos ler o cofre.",
    },
    zk: {
      title: "Como funciona, em três passos",
      lead: "Linguagem de negócio — sem whitepaper. A mesma lógica que protege bancos, aplicada à tua PME.",
      diagram: {
        client: "Browser",
        clientAction: "AES-GCM",
        server: "Servidor",
        serverAction: "armazena opaco",
        payload: "blob cifrado",
      },
      steps: [
        {
          title: "Desbloqueias localmente",
          body: "A Master Password deriva a chave com Argon2id — só na memória do browser.",
        },
        {
          title: "O servidor guarda blobs",
          body: "Títulos e tipos sim; passwords, IBANs e notas — nunca em claro.",
        },
        {
          title: "Decifras quando precisas",
          body: "Sync, partilha e agentes obedecem ao menor privilégio.",
        },
      ],
      note: "Mesmo com acesso à base de dados, um atacante só vê ruído criptográfico.",
    },
    layers: {
      title: "Ecossistema em camadas",
      lead: "Não vendemos ferramentas soltas. O Cofre é o núcleo; o Workspace assenta nele.",
      coreLabel: "Core",
      coreHint: "Identidade, segredos, confiança",
      coreItems: ["Cofre", "Equipa", "Segurança"],
      workspaceLabel: "Workspace",
      workspaceHint: "Trabalho diário no dispositivo",
      workspaceItems: ["Turnos", "Sandbox", "CLI mTLS", "Inventário", "Google"],
      footnote: "RH, finanças e CRM integram-se como camada operacional — visível na app, não no site público v1.",
    },
  },
  products: {
    vault: {
      meta: {
        title: "Cofre — Zero-knowledge vault | AegisPass",
        description: "Palavras-passe, notas e cartões cifrados no cliente. Sync em tempo real sem o servidor ver conteúdo.",
      },
      pageStatus: "live",
      hero: {
        eyebrow: "Produto · Cofre",
        headline: "O coração do ecossistema.",
        subline: "Guarda, sincroniza e partilha segredos — o servidor nunca vê a Master Key nem o conteúdo em claro.",
      },
      features: {
        title: "O que já podes validar",
        items: [
          {
            title: "Cifragem AES-GCM no browser",
            body: "Cada item é um blob com nonce único. Derivação Argon2id para Master Key e Auth Hash distinto.",
            status: "live",
          },
          {
            title: "Sync WebSocket em tempo real",
            body: "Alterações propagam entre dispositivos; o backend reencaminha eventos, não decifra payloads.",
            status: "live",
          },
          {
            title: "TOTP inline e gerador",
            body: "2FA e passwords fortes sem sair do cofre.",
            status: "live",
          },
          {
            title: "Score de higiene",
            body: "Passwords fracas e reutilizadas — calculado só no cliente.",
            status: "preview",
          },
        ],
      },
      cta: { primary: "Abrir cofre na app", secondary: "Ver plataforma" },
    },
    security: {
      meta: {
        title: "Segurança — Postura visível | AegisPass",
        description: "Higiene, dispositivos, Sentinel Mode e acesso de emergência para PME exigentes.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produto · Segurança",
        headline: "Confiança que se vê.",
        subline: "Postura de segurança legível — sessões, alertas e step-up auth sem teatro de alarmes.",
      },
      features: {
        title: "Capacidades",
        items: [
          {
            title: "Saúde de segurança",
            body: "Score de higiene e acções recomendadas — tudo calculado no browser.",
            status: "preview",
          },
          {
            title: "Sentinel Mode",
            body: "Deteta logins geograficamente impossíveis e exige confirmação extra.",
            status: "preview",
          },
          {
            title: "Dispositivos e passkeys",
            body: "Revoga sessões, gere passkeys, remote wipe de emergência.",
            status: "building",
          },
          {
            title: "Acesso de emergência",
            body: "Herdeiro digital com período de espera — confiança empresarial.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Explorar na app", secondary: "Ver plataforma" },
    },
    team: {
      meta: {
        title: "Equipa — Partilha efémera | AegisPass",
        description: "Shared vaults, secret links de um clique e notas que desaparecem — zero-knowledge end-to-end.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Produto · Equipa",
        headline: "Partilha sem deixar rasto.",
        subline: "Credenciais para freelancers, links secretos em RAM e cofres partilhados com permissões granulares.",
      },
      features: {
        title: "Capacidades",
        items: [
          {
            title: "Secret links efémeros",
            body: "Blob cifrado em RAM; chave no fragment URL — lido uma vez, link morto.",
            status: "preview",
          },
          {
            title: "Shared Vaults",
            body: "Cofres de equipa com chaves assimétricas por utilizador.",
            status: "preview",
          },
          {
            title: "Notas auto-destrutivas",
            body: "Countdown visível — «queima» criptográfica após TTL.",
            status: "building",
          },
          {
            title: "Anexos cifrados",
            body: "Um blob por ficheiro — chave nunca no servidor.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Testar na app", secondary: "Ver plataforma" },
    },
    workspace: {
      meta: {
        title: "Workspace — BYOD profissional | AegisPass",
        description: "Turnos, geofence, browser sandbox e CLI mTLS — trabalho seguro no dispositivo pessoal.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Produto · Workspace",
        headline: "BYOD com fronteiras nítidas.",
        subline: "Fora do turno, o cofre expurga-se da memória. Dentro das regras, ferramentas de trabalho isoladas.",
      },
      features: {
        title: "Roadmap visível",
        items: [
          {
            title: "Turnos e geofence",
            body: "Master Key só dentro do horário e local autorizados.",
            status: "preview",
          },
          {
            title: "Browser sandbox",
            body: "Injeção de credenciais sem revelar password ao OS.",
            status: "building",
          },
          {
            title: "CLI mTLS",
            body: "Segredos para pipelines e dev — autenticação mútua.",
            status: "building",
          },
          {
            title: "Google Workspace ZK",
            body: "Proxy com blobs cifrados para Drive e Docs.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Ver progresso na app", secondary: "Ver plataforma" },
    },
  },
  productsIndex: {
    meta: {
      title: "Produtos — Ecossistema AegisPass",
      description:
        "Cofre, Segurança, Equipa e Workspace — quatro módulos zero-knowledge que assentam no mesmo núcleo de identidade.",
    },
    hero: {
      eyebrow: "Produtos",
      headline: "Quatro módulos. Um núcleo.",
      subline:
        "Cada produto tem página própria — explora capacidades, estado de lançamento e dá feedback antes do go-live.",
    },
  },
  servicesIndex: {
    meta: {
      title: "Serviços — Implementação e conformidade | AegisPass",
      description:
        "Agentes IA, auditoria RGPD e implementação BYOD — serviços profissionais sobre a plataforma zero-knowledge.",
    },
    hero: {
      eyebrow: "Serviços",
      headline: "Da auditoria à operação.",
      subline:
        "Ajudamos PME e parceiros a adoptar o ecossistema — com entregáveis claros, human-in-the-loop e sem expor segredos.",
    },
    explore: "Explorar",
    cards: {
      agents: {
        tagline: "Agentes IA",
        description:
          "Automação de onboarding, triagem de pedidos e assistentes internos — sempre com aprovação humana e dados cifrados.",
        status: "preview",
      },
      compliance: {
        tagline: "Conformidade",
        description:
          "Gap analysis RGPD, registos de tratamento, crypto-shredding e certificados de eliminação — alinhado ao produto HR.",
        status: "preview",
      },
      deployment: {
        tagline: "Implementação",
        description:
          "Rollout BYOD, WireGuard, tenant multi-empresa e integração com Google Workspace — do piloto à produção.",
        status: "building",
      },
    },
  },
  services: {
    agents: {
      meta: {
        title: "Agentes IA — Human-in-the-loop | AegisPass",
        description:
          "Assistentes e automações sobre dados cifrados — aprovações, audit trail e menor privilégio por defeito.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Serviço · Agentes IA",
        headline: "Automação sem perder controlo.",
        subline:
          "Agentes leem metadados e blobs só no cliente — pedidos sensíveis passam por fila de aprovação antes de agir.",
      },
      features: {
        title: "Entregáveis",
        items: [
          {
            title: "Onboarding assistido",
            body: "Checklists, aliases e cofres iniciais — o agente propõe, o admin aprova.",
            status: "preview",
          },
          {
            title: "Triagem de pedidos",
            body: "Classifica tickets internos sem exportar segredos para LLMs externos.",
            status: "preview",
          },
          {
            title: "Audit trail imutável",
            body: "Cada acção do agente fica encadeada em log — pronto para auditoria.",
            status: "building",
          },
          {
            title: "Políticas por tenant",
            body: "Limites de scope, horários e geofence aplicados antes de executar.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Pedir early access", secondary: "Ver plataforma" },
    },
    compliance: {
      meta: {
        title: "Conformidade RGPD — Auditoria | AegisPass",
        description:
          "Gap analysis, registos de tratamento e evidências criptográficas de eliminação — RGPD by design.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Serviço · Conformidade",
        headline: "RGPD que se prova.",
        subline:
          "Mapeamos tratamentos, minimizamos dados e documentamos crypto-shredding — sem slides genéricos.",
      },
      features: {
        title: "Entregáveis",
        items: [
          {
            title: "Gap analysis inicial",
            body: "Inventário de fluxos, bases legais e riscos — priorizado para PME.",
            status: "preview",
          },
          {
            title: "Registo de tratamentos",
            body: "Art. 30 alinhado aos módulos HR e Cofre — campos sensíveis cifrados.",
            status: "preview",
          },
          {
            title: "Certificado de eliminação",
            body: "Prova criptográfica de erasure após crypto-shredding.",
            status: "building",
          },
          {
            title: "Formação DPO light",
            body: "Workshops para equipas pequenas — linguagem de negócio, não juridiquês.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Agendar conversa", secondary: "Ver plataforma" },
    },
    deployment: {
      meta: {
        title: "Implementação — BYOD e tenant | AegisPass",
        description:
          "Piloto, WireGuard, multi-tenant e integrações — rollout seguro do ecossistema AegisPass.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Serviço · Implementação",
        headline: "Do piloto à produção.",
        subline:
          "Configuramos tenant, políticas BYOD e rede — com runbooks e handover para a tua equipa ou MSP parceiro.",
      },
      features: {
        title: "Fases típicas",
        items: [
          {
            title: "Piloto sandbox",
            body: "10–30 utilizadores, cofre live, métricas de adopção.",
            status: "preview",
          },
          {
            title: "Rede WireGuard",
            body: "Túnel zero-trust entre dispositivos e API — sem VPN legacy.",
            status: "building",
          },
          {
            title: "Multi-tenant + RLS",
            body: "Isolamento PostgreSQL validado com testes negativos.",
            status: "building",
          },
          {
            title: "Integração Workspace",
            body: "Google, aliases e CLI mTLS — roadmap alinhado ao produto.",
            status: "building",
          },
        ],
      },
      cta: { primary: "Contactar equipa", secondary: "Ver parceiros" },
    },
  },
  partners: {
    meta: {
      title: "Parceiros — White-label | AegisPass",
      description: "MSPs e integradores: revende a camada Core com a tua marca sobre zero-knowledge real.",
    },
    pageStatus: "preview",
    hero: {
      eyebrow: "Parceiros",
      headline: "A tua marca. O nosso núcleo.",
      subline:
        "White-label do Cofre e identidade — Workspace e agentes ligam-se ao mesmo ecossistema. Canal early access aberto.",
    },
    benefits: {
      title: "Porquê parceiro",
      items: [
        "Core zero-knowledge pronto a revender — não reinventas cifragem.",
        "Paletas e tenant white-label na app (roadmap v2 no site).",
        "Documentação in-app didática — menos suporte L1 teu.",
        "Agentes com human-in-the-loop — automação sem perder controlo.",
      ],
    },
    contact: {
      label: "Early access parceiros",
      email: "partners@aegispass.com",
    },
  },
  construction: {
    badge: "Em construção",
    title: "Estamos a polir esta página",
    body: "A funcionalidade core já existe na app sandbox — esta página pública acompanha o lançamento. O teu feedback orienta a ordem do roadmap.",
  },
  statusLabels: {
    live: "Disponível",
    preview: "Pré-visualização",
    building: "Em construção",
  },
  footer: {
    tagline: "Identidade e segredos para equipas BYOD.",
    columns: { products: "Produtos", services: "Serviços", company: "Empresa", legal: "Legal" },
    platform: "Plataforma",
    partners: "Parceiros",
    contact: "Contacto",
    privacy: "Privacidade (em breve)",
    terms: "Termos (em breve)",
    app: "App",
    rights: "Todos os direitos reservados.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default pt;
