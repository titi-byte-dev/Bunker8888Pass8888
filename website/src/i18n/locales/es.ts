import type { SiteCopy } from "../types";

const es: SiteCopy = {
  site: {
    name: "AegisPass",
    defaultMeta: {
      title: "AegisPass — Identidad y secretos para PYME",
      description:
        "La capa de identidad y secretos de tu empresa. Zero-knowledge, BYOD y ecosistema de trabajo — todo se apoya en la Bóveda.",
    },
  },
  skip: "Saltar al contenido",
  nav: {
    products: "Productos",
    services: "Servicios",
    platform: "Plataforma",
    partners: "Socios",
    enter: "Entrar",
    productLabels: {
      vault: "Bóveda",
      security: "Seguridad",
      team: "Equipo",
      workspace: "Workspace",
    },
    serviceLabels: {
      agents: "Agentes IA",
      compliance: "Cumplimiento",
      deployment: "Implementación",
    },
  },
  home: {
    meta: {
      title: "AegisPass — El ecosistema de identidad para PYME",
      description:
        "Una bóveda zero-knowledge, seguridad visible y BYOD — una plataforma calmada, densa y hecha para crecer contigo.",
    },
    campaign: {
      eyebrow: "Campaña 2026 · Plataforma en construcción",
      headline: "La identidad de tu PYME.",
      highlight: "Todo se apoya en la Bóveda.",
      subline:
        "Secretos, equipo y trabajo en el mismo dispositivo — cifrado que el servidor nunca ve. Construimos contigo; valida cada módulo antes del lanzamiento.",
      ctaPrimary: "Abrir la app",
      ctaSecondary: "Explorar productos",
    },
    proof: [
      { value: "0", label: "secretos legibles en el servidor" },
      { value: "E2E", label: "cifrado zero-knowledge" },
      { value: "BYOD", label: "dispositivo personal, datos profesionales" },
    ],
    products: {
      title: "Productos anclados en la Bóveda",
      lead: "Cada módulo tiene su página — explora, da feedback, prioriza con nosotros.",
      explore: "Explorar",
      cards: {
        vault: {
          tagline: "Bóveda",
          description: "Logins, notas y tarjetas cifrados en el navegador. La Master Key nunca sale del dispositivo.",
          status: "live",
        },
        security: {
          tagline: "Seguridad",
          description: "Higiene, dispositivos, Sentinel y acceso de emergencia — confianza visible, sin alarmismo.",
          status: "preview",
        },
        team: {
          tagline: "Equipo",
          description: "Bóvedas compartidas, enlaces secretos efímeros y notas que desaparecen tras la lectura.",
          status: "preview",
        },
        workspace: {
          tagline: "Workspace",
          description: "Turnos, geofence, sandbox y CLI mTLS — fronteras claras entre lo personal y lo profesional.",
          status: "building",
        },
      },
    },
    services: {
      title: "Servicios profesionales",
      lead: "Implementación, cumplimiento RGPD y agentes con human-in-the-loop — para adoptar el ecosistema sin reinventar procesos.",
      explore: "Explorar",
      allLink: "Ver todos los servicios",
    },
    platformTeaser: {
      title: "Plataforma zero-knowledge",
      body: "Cómo funciona el cifrado, las capas Core + Workspace y por qué zero-knowledge es arquitectura — no eslogan.",
      link: "Ver plataforma",
    },
    ctaBand: {
      title: "Construimos en público",
      body: "La app ya corre en sandbox. Entra, prueba la bóveda y cuéntanos qué falta para tu PYME.",
      primary: "Abrir app",
    },
  },
  platform: {
    meta: {
      title: "Plataforma — Zero-knowledge explicado | AegisPass",
      description: "Cifrado en el cliente, blobs opacos en el servidor y ecosistema en capas — la base técnica de AegisPass.",
    },
    pageStatus: "live",
    hero: {
      eyebrow: "Plataforma",
      headline: "Zero-knowledge no es un eslogan. Es un contrato.",
      subline:
        "El servidor almacena metadatos y blobs ilegibles. Solo tu dispositivo tiene la clave — ni nosotros podemos leer la bóveda.",
    },
    zk: {
      title: "Cómo funciona, en tres pasos",
      lead: "Lenguaje de negocio — sin whitepaper. La misma lógica que los bancos, adaptada a tu PYME.",
      diagram: {
        client: "Navegador",
        clientAction: "AES-GCM",
        server: "Servidor",
        serverAction: "almacena opaco",
        payload: "blob cifrado",
      },
      steps: [
        { title: "Desbloqueo local", body: "La contraseña maestra deriva la clave con Argon2id — solo en memoria del navegador." },
        { title: "El servidor guarda blobs", body: "Títulos y tipos sí; contraseñas, IBAN y notas — nunca en claro." },
        { title: "Descifras cuando lo necesitas", body: "Sync, compartición y agentes respetan el mínimo privilegio." },
      ],
      note: "Incluso con acceso a la base de datos, un atacante solo ve ruido criptográfico.",
    },
    layers: {
      title: "Ecosistema en capas",
      lead: "No vendemos herramientas sueltas. La Bóveda es el núcleo; el Workspace se apoya en ella.",
      coreLabel: "Core",
      coreHint: "Identidad, secretos, confianza",
      coreItems: ["Bóveda", "Equipo", "Seguridad"],
      workspaceLabel: "Workspace",
      workspaceHint: "Trabajo diario en el dispositivo",
      workspaceItems: ["Turnos", "Sandbox", "CLI mTLS", "Inventario", "Google"],
      footnote: "RRHH, finanzas y CRM se integran como capa operativa — visible en la app, no en el sitio público v1.",
    },
  },
  products: {
    vault: {
      meta: {
        title: "Bóveda — Vault zero-knowledge | AegisPass",
        description: "Contraseñas, notas y tarjetas cifradas en el cliente. Sync en tiempo real sin contenido visible en servidor.",
      },
      pageStatus: "live",
      hero: {
        eyebrow: "Producto · Bóveda",
        headline: "El corazón del ecosistema.",
        subline: "Guarda, sincroniza y comparte secretos — el servidor nunca ve la Master Key ni el contenido en claro.",
      },
      features: {
        title: "Lo que ya puedes validar",
        items: [
          { title: "Cifrado AES-GCM en el navegador", body: "Cada ítem es un blob con nonce único. Argon2id para Master Key y Auth Hash distinto.", status: "live" },
          { title: "Sync WebSocket en tiempo real", body: "Los cambios se propagan; el backend reenvía eventos, no descifra.", status: "live" },
          { title: "TOTP inline y generador", body: "2FA y contraseñas fuertes sin salir de la bóveda.", status: "live" },
          { title: "Score de higiene", body: "Contraseñas débiles y reutilizadas — calculado solo en el cliente.", status: "preview" },
        ],
      },
      cta: { primary: "Abrir bóveda", secondary: "Ver plataforma" },
    },
    security: {
      meta: {
        title: "Seguridad — Postura visible | AegisPass",
        description: "Higiene, dispositivos, Sentinel Mode y acceso de emergencia para PYME exigentes.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Producto · Seguridad",
        headline: "Confianza que se ve.",
        subline: "Postura de seguridad legible — sesiones, alertas y step-up auth sin teatro de alarmas.",
      },
      features: {
        title: "Capacidades",
        items: [
          { title: "Salud de seguridad", body: "Score de higiene y acciones recomendadas — calculado en el navegador.", status: "preview" },
          { title: "Sentinel Mode", body: "Detecta inicios de sesión geográficamente imposibles.", status: "preview" },
          { title: "Dispositivos y passkeys", body: "Revoca sesiones, gestiona passkeys, remote wipe.", status: "building" },
          { title: "Acceso de emergencia", body: "Heredero digital con periodo de espera.", status: "building" },
        ],
      },
      cta: { primary: "Explorar en la app", secondary: "Ver plataforma" },
    },
    team: {
      meta: {
        title: "Equipo — Compartición efímera | AegisPass",
        description: "Bóvedas compartidas, enlaces secretos de un uso y notas efímeras — zero-knowledge de extremo a extremo.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Producto · Equipo",
        headline: "Comparte sin dejar rastro.",
        subline: "Credenciales para freelancers, enlaces secretos en RAM y bóvedas compartidas con permisos granulares.",
      },
      features: {
        title: "Capacidades",
        items: [
          { title: "Enlaces secretos efímeros", body: "Blob cifrado en RAM; clave en el fragment URL — leído una vez, enlace muerto.", status: "preview" },
          { title: "Shared Vaults", body: "Bóvedas de equipo con claves asimétricas por usuario.", status: "preview" },
          { title: "Notas autodestructivas", body: "Cuenta atrás visible — «quema» criptográfica tras TTL.", status: "building" },
          { title: "Adjuntos cifrados", body: "Un blob por archivo — clave nunca en el servidor.", status: "building" },
        ],
      },
      cta: { primary: "Probar en la app", secondary: "Ver plataforma" },
    },
    workspace: {
      meta: {
        title: "Workspace — BYOD profesional | AegisPass",
        description: "Turnos, geofence, navegador sandbox y CLI mTLS — trabajo seguro en dispositivo personal.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Producto · Workspace",
        headline: "BYOD con fronteras nítidas.",
        subline: "Fuera del turno, la bóveda se expulsa de memoria. Dentro de las reglas, herramientas de trabajo aisladas.",
      },
      features: {
        title: "Roadmap visible",
        items: [
          { title: "Turnos y geofence", body: "Master Key solo dentro del horario y lugar autorizados.", status: "preview" },
          { title: "Navegador sandbox", body: "Inyección de credenciales sin revelar la contraseña al SO.", status: "building" },
          { title: "CLI mTLS", body: "Secretos para pipelines y dev — autenticación mutua.", status: "building" },
          { title: "Google Workspace ZK", body: "Proxy con blobs cifrados para Drive y Docs.", status: "building" },
        ],
      },
      cta: { primary: "Ver progreso en la app", secondary: "Ver plataforma" },
    },
  },
  productsIndex: {
    meta: {
      title: "Productos — Ecosistema AegisPass",
      description:
        "Cofre, Seguridad, Equipo y Workspace — cuatro módulos zero-knowledge sobre el mismo núcleo de identidad.",
    },
    hero: {
      eyebrow: "Productos",
      headline: "Cuatro módulos. Un núcleo.",
      subline:
        "Cada producto tiene su página — explora capacidades, estado de lanzamiento y da feedback antes del go-live.",
    },
  },
  servicesIndex: {
    meta: {
      title: "Servicios — Implementación y cumplimiento | AegisPass",
      description:
        "Agentes IA, auditoría RGPD e implementación BYOD — servicios profesionales sobre la plataforma zero-knowledge.",
    },
    hero: {
      eyebrow: "Servicios",
      headline: "De la auditoría a la operación.",
      subline:
        "Ayudamos a pymes y socios a adoptar el ecosistema — entregables claros, human-in-the-loop y sin exponer secretos.",
    },
    explore: "Explorar",
    cards: {
      agents: {
        tagline: "Agentes IA",
        description:
          "Automatización de onboarding, triaje de solicitudes y asistentes internos — siempre con aprobación humana y datos cifrados.",
        status: "preview",
      },
      compliance: {
        tagline: "Cumplimiento",
        description:
          "Gap analysis RGPD, registros de tratamiento, crypto-shredding y certificados de eliminación — alineado al módulo RRHH.",
        status: "preview",
      },
      deployment: {
        tagline: "Implementación",
        description:
          "Rollout BYOD, WireGuard, multi-tenant e integración Google Workspace — del piloto a producción.",
        status: "building",
      },
    },
  },
  services: {
    agents: {
      meta: {
        title: "Agentes IA — Human-in-the-loop | AegisPass",
        description:
          "Asistentes y automatizaciones sobre datos cifrados — aprobaciones, audit trail y mínimo privilegio por defecto.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Servicio · Agentes IA",
        headline: "Automatización sin perder control.",
        subline:
          "Los agentes leen metadatos y blobs solo en el cliente — las acciones sensibles pasan por cola de aprobación.",
      },
      features: {
        title: "Entregables",
        items: [
          { title: "Onboarding asistido", body: "Checklists, alias y cofres iniciales — el agente propone, el admin aprueba.", status: "preview" },
          { title: "Triaje de solicitudes", body: "Clasifica tickets internos sin exportar secretos a LLM externos.", status: "preview" },
          { title: "Audit trail inmutable", body: "Cada acción del agente queda encadenada en log — listo para auditoría.", status: "building" },
          { title: "Políticas por tenant", body: "Límites de scope, horarios y geofence aplicados antes de ejecutar.", status: "building" },
        ],
      },
      cta: { primary: "Solicitar early access", secondary: "Ver plataforma" },
    },
    compliance: {
      meta: {
        title: "Cumplimiento RGPD — Auditoría | AegisPass",
        description:
          "Gap analysis, registros de tratamiento y pruebas criptográficas de eliminación — RGPD by design.",
      },
      pageStatus: "preview",
      hero: {
        eyebrow: "Servicio · Cumplimiento",
        headline: "RGPD que se demuestra.",
        subline:
          "Mapeamos tratamientos, minimizamos datos y documentamos crypto-shredding — sin diapositivas genéricas.",
      },
      features: {
        title: "Entregables",
        items: [
          { title: "Gap analysis inicial", body: "Inventario de flujos, bases legales y riesgos — priorizado para pymes.", status: "preview" },
          { title: "Registro de tratamientos", body: "Art. 30 alineado a módulos RRHH y Cofre — campos sensibles cifrados.", status: "preview" },
          { title: "Certificado de eliminación", body: "Prueba criptográfica de erasure tras crypto-shredding.", status: "building" },
          { title: "Formación DPO light", body: "Talleres para equipos pequeños — lenguaje de negocio, no juridiqués.", status: "building" },
        ],
      },
      cta: { primary: "Agendar conversación", secondary: "Ver plataforma" },
    },
    deployment: {
      meta: {
        title: "Implementación — BYOD y tenant | AegisPass",
        description:
          "Piloto, WireGuard, multi-tenant e integraciones — rollout seguro del ecosistema AegisPass.",
      },
      pageStatus: "building",
      hero: {
        eyebrow: "Servicio · Implementación",
        headline: "Del piloto a producción.",
        subline:
          "Configuramos tenant, políticas BYOD y red — con runbooks y traspaso a tu equipo o MSP socio.",
      },
      features: {
        title: "Fases típicas",
        items: [
          { title: "Piloto sandbox", body: "10–30 usuarios, cofre live, métricas de adopción.", status: "preview" },
          { title: "Red WireGuard", body: "Túnel zero-trust entre dispositivos y API — sin VPN legacy.", status: "building" },
          { title: "Multi-tenant + RLS", body: "Aislamiento PostgreSQL validado con tests negativos.", status: "building" },
          { title: "Integración Workspace", body: "Google, alias y CLI mTLS — roadmap alineado al producto.", status: "building" },
        ],
      },
      cta: { primary: "Contactar equipo", secondary: "Ver socios" },
    },
  },
  partners: {
    meta: {
      title: "Socios — White-label | AegisPass",
      description: "MSP e integradores: revende la capa Core con tu marca sobre zero-knowledge real.",
    },
    pageStatus: "preview",
    hero: {
      eyebrow: "Socios",
      headline: "Tu marca. Nuestro núcleo.",
      subline: "White-label de la Bóveda e identidad — Workspace y agentes en el mismo ecosistema. Canal early access abierto.",
    },
    benefits: {
      title: "Por qué ser socio",
      items: [
        "Core zero-knowledge listo para revender — no reinventes el cifrado.",
        "Paletas y tenant white-label en la app (roadmap v2 sitio).",
        "Documentación in-app didáctica — menos soporte L1.",
        "Agentes con human-in-the-loop — automatización con control.",
      ],
    },
    contact: { label: "Early access socios", email: "partners@aegispass.com" },
  },
  construction: {
    badge: "En construcción",
    title: "Estamos puliendo esta página",
    body: "El núcleo ya funciona en sandbox app — esta página pública acompaña el lanzamiento. Tu feedback orienta nuestro roadmap.",
  },
  statusLabels: { live: "Disponible", preview: "Vista previa", building: "En construcción" },
  footer: {
    tagline: "Identidad y secretos para equipos BYOD.",
    columns: { products: "Productos", services: "Servicios", company: "Empresa", legal: "Legal" },
    platform: "Plataforma",
    partners: "Socios",
    contact: "Contacto",
    privacy: "Privacidad (próximamente)",
    terms: "Términos (próximamente)",
    app: "App",
    rights: "Todos los derechos reservados.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default es;
