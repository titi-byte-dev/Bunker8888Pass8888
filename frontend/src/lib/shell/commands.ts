/**
 * Registo de comandos da command palette (UI-006).
 * Sem side-effects — navegação via `href` tratada no componente.
 */
import { ROUTE_TREE, type RouteNode } from "./routes";
import type { DecodedLogin } from "$lib/vault/ui";

export type CommandGroup = "navigation" | "vault" | "action" | "docs";

export type CommandEntry = {
  id: string;
  label: string;
  group: CommandGroup;
  keywords?: string;
  href: string;
};

const GROUP_LABELS: Record<CommandGroup, string> = {
  navigation: "Navegação",
  vault: "Cofre",
  action: "Acções",
  docs: "Documentação",
};

export function groupLabel(group: CommandGroup): string {
  return GROUP_LABELS[group];
}

export function buildNavigationCommands(): CommandEntry[] {
  // Achata o ROUTE_TREE em comandos hierarquicos: "Financas > Fiscal".
  const out: CommandEntry[] = [];
  const walk = (nodes: RouteNode[], prefix: string) => {
    for (const node of nodes) {
      const label = prefix ? `${prefix} › ${node.label}` : node.label;
      // Evita duplicado quando o hub e um filho partilham href (ex.: /fin, /hr).
      const childShadows = (node.children ?? []).some((c) => c.href === node.href);
      if (!node.comingSoon && !childShadows) {
        out.push({
          id: `nav-${node.href}`,
          label,
          group: "navigation",
          keywords: `${node.href} ${node.taskId ?? ""}`,
          href: node.href,
        });
      }
      if (node.children) walk(node.children, node.label);
    }
  };
  walk(ROUTE_TREE, "");
  return out;
}

export function buildVaultCommands(logins: DecodedLogin[]): CommandEntry[] {
  return logins.map(({ meta, login }) => ({
    id: `vault-${meta.id}`,
    label: login.title,
    group: "vault" as const,
    keywords: `${login.username} ${login.url ?? ""}`,
    href: `/vault/${meta.id}`,
  }));
}

export function buildActionCommands(): CommandEntry[] {
  return [
    {
      id: "action-new-login",
      label: "Novo login",
      group: "action",
      keywords: "criar adicionar vault cofre",
      href: "/vault/new",
    },
    {
      id: "action-sandbox",
      label: "Abrir browser sandbox",
      group: "action",
      keywords: "trabalho inject injeção",
      href: "/work/sandbox",
    },
    {
      id: "action-hygiene",
      label: "Saúde de segurança",
      group: "action",
      keywords: "segurança score breach dark web higiene",
      href: "/security/hygiene",
    },
    {
      id: "action-devices",
      label: "Dispositivos e sessões",
      group: "action",
      keywords: "segurança passkey cli sessão",
      href: "/security/devices",
    },
    {
      id: "action-shifts",
      label: "Turnos e geofence",
      group: "action",
      keywords: "trabalho horário turno geofence ntp",
      href: "/work/shifts",
    },
    {
      id: "action-settings",
      label: "Definições",
      group: "action",
      keywords: "tema aparência passkey conta settings",
      href: "/settings",
    },
    {
      id: "action-docs",
      label: "Documentação",
      group: "action",
      keywords: "ajuda help guia conceitos glossary aprender",
      href: "/settings/docs",
    },
    {
      id: "action-docs-glossary",
      label: "Glossário de conceitos",
      group: "action",
      keywords: "zero-knowledge nonce tenant totp conceito",
      href: "/settings/docs/glossary",
    },
    {
      id: "action-sentinel",
      label: "Sentinel Mode",
      group: "action",
      keywords: "segurança login impossível viagem step-up",
      href: "/security/sentinel",
    },
    {
      id: "action-admin",
      label: "Administração",
      group: "action",
      keywords: "admin utilizadores turno wipe auditoria",
      href: "/admin",
    },
    {
      id: "action-unlock",
      label: "Desbloquear cofre",
      group: "action",
      keywords: "master password unlock",
      href: "/auth/unlock",
    },
  ];
}

export function filterCommands(commands: CommandEntry[], query: string): CommandEntry[] {
  const q = query.trim().toLowerCase();
  if (!q) return commands;
  return commands.filter((cmd) => {
    const hay = `${cmd.label} ${cmd.keywords ?? ""} ${cmd.href}`.toLowerCase();
    return hay.includes(q);
  });
}

export function groupCommands(commands: CommandEntry[]): Map<CommandGroup, CommandEntry[]> {
  const order: CommandGroup[] = ["docs", "navigation", "vault", "action"];
  const map = new Map<CommandGroup, CommandEntry[]>();
  for (const g of order) map.set(g, []);
  for (const cmd of commands) {
    map.get(cmd.group)?.push(cmd);
  }
  for (const g of order) {
    if (map.get(g)?.length === 0) map.delete(g);
  }
  return map;
}
