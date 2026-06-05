/**
 * ROUTE_TREE — fonte unica de verdade da arvore de navegacao (UI-011).
 *
 * Didatico: ate aqui a navegacao era PLANA (lib/shell/nav.ts) — sub-paginas como
 * "Fiscal" ou "Faturas" apareciam ao MESMO nivel que "Cofre". Isto contradiz
 * docs/roadmap/09-design/design-system.md §3. Aqui modelamos a hierarquia real
 * (modulo -> filhos) UMA vez e derivamos dela:
 *
 *   ROUTE_TREE ──► AppSidebar (seccoes colapsaveis, filho activo expande pai)
 *              ──► Breadcrumbs (trilho "Financas > Fiscal" pelo pathname)
 *              ──► CommandPalette (comandos hierarquicos)
 *
 * Regra: NUNCA duplicar labels/hrefs noutro sitio; importar daqui.
 */

export type RouteNode = {
  /** Nome de produto em PT-PT (o que o utilizador ve). */
  label: string;
  /** Caminho absoluto. Hubs apontam para a sua pagina indice. */
  href: string;
  /** ID de task (FIN-005...) — so renderizado em import.meta.env.DEV. */
  taskId?: string;
  /** Mostrar na tab bar mobile (max. 5). So nodes de topo. */
  tabBar?: boolean;
  /** Modulo futuro — link desactivado com badge "Em breve". */
  comingSoon?: boolean;
  /** Sub-rotas do modulo. */
  children?: RouteNode[];
};

/**
 * Arvore alvo — espelha docs/roadmap/09-design/design-system.md §3.2.
 * A ORDEM aqui e a ordem visual na sidebar.
 */
export const ROUTE_TREE: RouteNode[] = [
  { label: "Cofre", href: "/vault", taskId: "VAULT-001", tabBar: true },
  {
    label: "Seguranca",
    href: "/security",
    tabBar: true,
    children: [
      { label: "Saude de seguranca", href: "/security/hygiene", taskId: "SEC-002" },
      { label: "Dispositivos e sessoes", href: "/security/devices", taskId: "SEC-003" },
      { label: "Sentinel Mode", href: "/security/sentinel", taskId: "VAULT-014" },
      { label: "Acesso de emergencia", href: "/security/emergency", taskId: "VAULT-016" },
      { label: "Auditoria Guardiao", href: "/security/guardian", taskId: "AGENT-010" },
    ],
  },
  {
    label: "Trabalho",
    href: "/work",
    tabBar: true,
    children: [
      { label: "Turnos e geofence", href: "/work/shifts", taskId: "WORK-001" },
      { label: "Browser sandbox", href: "/work/sandbox", taskId: "VAULT-013" },
      { label: "CLI mTLS", href: "/work/cli", taskId: "CLI-001" },
      { label: "Inventario", href: "/work/inventory", taskId: "OPS-001" },
      { label: "Google Workspace", href: "/work/google", taskId: "GOOGLE-001" },
    ],
  },
  {
    label: "Equipa",
    href: "/team",
    children: [
      { label: "Shared vaults", href: "/team/vaults", taskId: "SHARE-002" },
      { label: "Notas partilhadas", href: "/team/notes", taskId: "SHARE-005" },
      { label: "Secret links", href: "/team/links", taskId: "SHARE-003" },
    ],
  },
  {
    label: "RH",
    href: "/hr",
    children: [
      { label: "Fichas e contratos", href: "/hr", taskId: "HR-001" },
      { label: "Onboarding", href: "/hr/onboarding", taskId: "HR-007" },
      { label: "Recrutamento", href: "/hr/recruitment", taskId: "HR-009" },
      { label: "Conformidade RGPD", href: "/hr/compliance", taskId: "HR-008" },
    ],
  },
  { label: "Mail", href: "/mail", taskId: "MAIL-001" },
  {
    label: "Financas",
    href: "/fin",
    children: [
      { label: "Custos SaaS", href: "/fin/costs", taskId: "FIN-001" },
      { label: "Fiscal", href: "/fin/fiscal", taskId: "FIN-005" },
      { label: "Faturas", href: "/fin/invoices", taskId: "FIN-006" },
      { label: "Comissoes", href: "/fin/commissions", taskId: "FIN-007" },
      { label: "Reconciliacao bancaria", href: "/fin/banking", taskId: "FIN-003" },
    ],
  },
  { label: "CRM", href: "/crm", taskId: "CRM-001" },
  {
    label: "Admin",
    href: "/admin",
    children: [
      { label: "Utilizadores", href: "/admin/users", taskId: "ADMIN-001" },
      { label: "Audit log", href: "/admin/audit", taskId: "HR-002" },
    ],
  },
  { label: "Definicoes", href: "/settings", tabBar: true, taskId: "UI-001" },
];

/** Item da sidebar/tab bar derivado da arvore (so node de topo). */
export type NavModule = RouteNode;

/** Modulos de topo (sidebar desktop). */
export function navModules(): NavModule[] {
  return ROUTE_TREE;
}

/** Subconjunto para tab bar mobile (max. 5). */
export function tabBarModules(): NavModule[] {
  return ROUTE_TREE.filter((n) => n.tabBar);
}

/** true se o pathname pertence a este node (exacto ou sub-rota). */
export function isRouteActive(pathname: string, href: string): boolean {
  if (href === "/") return pathname === "/";
  return pathname === href || pathname.startsWith(`${href}/`);
}

/**
 * Trilho de breadcrumbs para um pathname.
 * Devolve a cadeia raiz->folha do node que melhor corresponde:
 *   match exacto > prefixo mais longo (suporta rotas dinamicas /vault/[id]).
 * O ultimo segmento e a pagina actual (sem link no componente).
 */
export function routeTrail(pathname: string): RouteNode[] {
  let best: RouteNode[] | null = null;
  let bestLen = -1;

  const walk = (nodes: RouteNode[], acc: RouteNode[]) => {
    for (const node of nodes) {
      const trail = [...acc, node];
      if (isRouteActive(pathname, node.href) && node.href.length > bestLen) {
        best = trail;
        bestLen = node.href.length;
      }
      if (node.children) walk(node.children, trail);
    }
  };
  walk(ROUTE_TREE, []);

  return best ?? [];
}

/**
 * Filhos de um modulo (para paginas-hub via HubLinks).
 * Didatico: o hub NAO redeclara as suas sub-paginas — deriva-as da arvore,
 * por isso adicionar uma rota a ROUTE_TREE chega para aparecer no hub.
 * Exclui o filho-sombra que partilha href com o modulo (ex.: "Custos SaaS").
 */
export function routeChildren(href: string): RouteNode[] {
  const node = flattenRoutes().find((n) => n.href === href);
  return (node?.children ?? []).filter((c) => c.href !== href);
}

/** Lista plana de todos os nodes (para CommandPalette / pesquisa). */
export function flattenRoutes(nodes: RouteNode[] = ROUTE_TREE): RouteNode[] {
  return nodes.flatMap((n) => [n, ...(n.children ? flattenRoutes(n.children) : [])]);
}
