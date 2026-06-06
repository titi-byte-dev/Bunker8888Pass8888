import type { SiteCopy } from "./types";

const es: SiteCopy = {
  meta: {
    title: "AegisPass — Identidad y secretos para tu PYME",
    description:
      "La capa de identidad y secretos de tu empresa. Cifrado zero-knowledge, BYOD y ecosistema de trabajo — todo se apoya en la Bóveda.",
    ogTitle: "AegisPass — Todo se apoya en la Bóveda",
  },
  skip: "Saltar al contenido",
  nav: {
    zeroKnowledge: "Zero-Knowledge",
    ecosystem: "Ecosistema",
    pillars: "Principios",
    product: "Producto",
    partners: "Socios",
    enter: "Entrar",
  },
  hero: {
    eyebrow: "Plataforma de identidad para PYME",
    headline:
      "La capa de identidad y secretos de tu empresa. Todo se apoya en la Bóveda.",
    subline:
      "Contraseñas, equipo y trabajo seguro en el mismo dispositivo — el servidor nunca ve tus secretos en claro.",
    ctaPrimary: "Abrir la app",
    ctaSecondary: "Ver cómo funciona",
  },
  zk: {
    id: "zero-knowledge",
    title: "Zero-Knowledge, sin jerga",
    lead: "Cifras en tu navegador. El servidor guarda solo blobs opacos — aunque accedan a la base de datos, nadie lee tus datos.",
    diagram: {
      client: "Navegador",
      clientAction: "AES-GCM",
      server: "Servidor",
      serverAction: "almacena opaco",
      payload: "blob cifrado",
    },
    steps: [
      {
        title: "1. Desbloqueas localmente",
        body: "La contraseña maestra deriva la clave solo en tu dispositivo. Nunca sale a la red.",
      },
      {
        title: "2. El servidor guarda blobs",
        body: "Metadatos visibles (título, tipo) sí; contenido sensible — nunca en claro.",
      },
      {
        title: "3. Descifras cuando lo necesitas",
        body: "Sincronización, compartición y agentes respetan la misma regla: mínimo privilegio.",
      },
    ],
    note: "Ni nosotros, como proveedores, podemos leer tu bóveda.",
  },
  ecosystem: {
    id: "ecosystem",
    title: "Un ecosistema en capas",
    lead: "No son herramientas sueltas. La Bóveda es la base; el espacio de trabajo se apoya en ella.",
    coreLabel: "Core",
    coreHint: "Identidad, secretos y confianza",
    coreItems: ["Bóveda", "Equipo", "Seguridad"],
    workspaceLabel: "Workspace",
    workspaceHint: "Trabajo diario en tu dispositivo",
    workspaceItems: [
      "Turnos y geofence",
      "Navegador sandbox",
      "CLI mTLS",
      "Inventario",
      "Google Workspace",
    ],
    footnote:
      "RRHH, finanzas y CRM existen como capa operativa — integrada, no expuesta en la v1 pública.",
  },
  pillars: {
    id: "pillars",
    title: "Hecho a propósito para PYME modernas",
    items: [
      {
        fig: "0.1",
        title: "Confianza visible",
        body: "Estados de sesión, turno y cifrado legibles — calma operativa, sin alarmismo falso.",
      },
      {
        fig: "0.2",
        title: "Agentes con mínimo privilegio",
        body: "IA que solo accede a lo estrictamente necesario, con aprobación humana en acciones sensibles.",
      },
      {
        fig: "0.3",
        title: "Velocidad silenciosa",
        body: "Interfaz fluida, teclado primero, densidad inteligente — el admin ve más, el empleado lo esencial.",
      },
    ],
  },
  product: {
    id: "product",
    eyebrow: "Producto en acción",
    title: "Calma, densa y lista para el día a día",
    lead: "Vista previa de la app real — bóveda desbloqueada, lista legible, estados de seguridad visibles.",
    mockVault: "Bóveda",
    mockItems: ["GitHub — ingeniería", "Suscripción SaaS", "Alias mail soporte"],
    mockNavExtra: ["Seguridad", "Equipo"],
    mockStatus: "Cifrado · Sesión activa",
  },
  partners: {
    id: "partners",
    eyebrow: "Socios",
    title: "White-label para MSP e integradores",
    body: "Revende la capa Core (Bóveda + identidad) con tu marca. Workspace y agentes se apoyan en el mismo núcleo zero-knowledge.",
    hint: "Canal de socios — contacto early access.",
  },
  footer: {
    tagline: "AegisPass — identidad y secretos para equipos BYOD.",
    about: "Quiénes somos",
    contact: "Contacto",
    legal: "Legal",
    privacy: "Privacidad (próximamente)",
    terms: "Términos (próximamente)",
    app: "App",
    rights: "Todos los derechos reservados.",
  },
  lang: { pt: "Português", fr: "Français", es: "Español", de: "Deutsch" },
};

export default es;
