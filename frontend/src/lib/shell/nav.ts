/**
 * Navegação top-level da app (UI-002).
 * Alinhada com docs/roadmap/09-design/product-ui-vision.md §4.1
 */

export type NavItem = {
  id: string;
  label: string;
  href: string;
  /** Mostrar na tab bar mobile (máx. ~4 itens) */
  tabBar?: boolean;
  /** Só visível em `import.meta.env.DEV` */
  devOnly?: boolean;
  /** Módulo futuro — link desactivado com badge */
  comingSoon?: boolean;
};

/** Itens principais da sidebar (desktop) */
export const NAV_ITEMS: NavItem[] = [
  { id: "vault", label: "Cofre", href: "/vault", tabBar: true },
  { id: "security", label: "Segurança", href: "/security", tabBar: true },
  { id: "work", label: "Trabalho", href: "/work", tabBar: true },
  { id: "team", label: "Equipa", href: "/team" },
  { id: "hr", label: "RH", href: "/hr" },
  { id: "mail", label: "Aliases", href: "/mail" },
  { id: "fin", label: "Custos", href: "/fin" },
  { id: "fiscal", label: "Fiscal", href: "/fin/fiscal" },
  { id: "invoices", label: "Faturas", href: "/fin/invoices" },
  { id: "commissions", label: "Comissoes", href: "/fin/commissions" },
  { id: "crm", label: "CRM", href: "/crm" },
  { id: "admin", label: "Admin", href: "/admin" },
  { id: "settings", label: "Definições", href: "/settings", tabBar: true },
  { id: "dev", label: "Playground", href: "/dev", devOnly: true },
];

/** Filtra itens visíveis consoante ambiente de build */
export function visibleNavItems(dev = import.meta.env.DEV): NavItem[] {
  return NAV_ITEMS.filter((item) => !item.devOnly || dev);
}

/** Subconjunto para tab bar mobile */
export function tabBarItems(dev = import.meta.env.DEV): NavItem[] {
  return visibleNavItems(dev).filter((item) => item.tabBar);
}

/** Verifica se o pathname corresponde ao href (exact ou prefixo de secção) */
export function isNavActive(pathname: string, href: string): boolean {
  if (href === "/") return pathname === "/";
  return pathname === href || pathname.startsWith(`${href}/`);
}
