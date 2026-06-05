/**
 * Ligações «Como funciona?» — rota da app → slug de documentação.
 * Didático: prefixos mais longos primeiro para `/team/links` não cair em `/team`.
 */
export type DocHelpTarget = {
  slug: string;
  label: string;
};

type DocHelpRoute = DocHelpTarget & { prefix: string };

/** Ordem irrelevante — resolveDocHelp ordena por comprimento do prefixo. */
export const DOC_HELP_ROUTES: DocHelpRoute[] = [
  { prefix: "/team/links", slug: "journey-secret-link", label: "Como funcionam os secret links?" },
  { prefix: "/team/vaults", slug: "journey-shared-vault", label: "Como funcionam os cofres partilhados?" },
  { prefix: "/team/notes", slug: "team-sharing", label: "Como funcionam as notas temporárias?" },
  { prefix: "/team", slug: "team-sharing", label: "Como funciona a partilha?" },
  { prefix: "/work/google-dev", slug: "journey-google-dev-stub", label: "Como funciona o Google em dev?" },
  { prefix: "/work/google", slug: "journey-google-dev-stub", label: "Como ligar o Google Workspace?" },
  { prefix: "/fin/fiscal", slug: "journey-fiscal-categorization", label: "Como funciona a categorização fiscal?" },
  { prefix: "/security/guardian", slug: "journey-guardian-audit", label: "Como auditar o Guardião?" },
  { prefix: "/fin/invoices", slug: "journey-erp-flow-dev", label: "Como funciona o fluxo ERP?" },
  { prefix: "/fin/banking", slug: "journey-finance-agent-reconcile", label: "Como funciona a reconciliação?" },
  { prefix: "/fin/costs", slug: "journey-finance-agent-saas", label: "Como funciona o agente financeiro?" },
  { prefix: "/fin", slug: "journey-finance-agent-saas", label: "Como funciona o agente financeiro?" },
  { prefix: "/crm", slug: "journey-crm-prospection", label: "Como funciona a prospeção CRM?" },
  { prefix: "/mail", slug: "journey-mail-alias-relay", label: "Como funcionam aliases e relay?" },
  { prefix: "/security/sentinel", slug: "journey-sentinel", label: "Como funciona o Sentinel?" },
  { prefix: "/security/emergency", slug: "journey-emergency-access", label: "Como funciona o acesso de emergência?" },
  { prefix: "/security/hygiene", slug: "security", label: "Como funciona a higiene?" },
  { prefix: "/security/devices", slug: "journey-passkey", label: "Como funcionam passkeys e sessões?" },
  { prefix: "/security", slug: "security", label: "Como funciona a segurança?" },
  { prefix: "/hr/recruitment", slug: "journey-hr-agent-recruitment", label: "Como funciona a triagem às cegas?" },
  { prefix: "/hr/onboarding", slug: "journey-hr-agent-onboarding", label: "Como funciona o agente RH?" },
  { prefix: "/hr/compliance", slug: "journey-rgpd-erasure", label: "Como funciona o direito ao esquecimento?" },
  { prefix: "/hr", slug: "hr-rgpd", label: "Como funciona o RH cifrado?" },
  { prefix: "/work/inventory", slug: "journey-ops-agent-inventory", label: "Como funciona o inventário?" },
  { prefix: "/work/shifts", slug: "security", label: "Como funcionam turnos e geofence?" },
  { prefix: "/work/sandbox", slug: "vault", label: "Como funciona o browser sandbox?" },
  { prefix: "/work/cli", slug: "developer-crypto", label: "Como funciona a CLI mTLS?" },
  { prefix: "/work", slug: "security", label: "Como funciona o trabalho BYOD?" },
  { prefix: "/admin/audit", slug: "journey-remote-wipe", label: "Como funciona a auditoria?" },
  { prefix: "/admin", slug: "journey-admin-onboarding", label: "Como funciona a administração?" },
  { prefix: "/vault", slug: "vault", label: "Como funciona o cofre?" },
];

export function resolveDocHelp(pathname: string, slugOverride?: string): DocHelpTarget | null {
  if (slugOverride) {
    const hit = DOC_HELP_ROUTES.find((r) => r.slug === slugOverride);
    return hit ?? { slug: slugOverride, label: "Como funciona?" };
  }

  const sorted = [...DOC_HELP_ROUTES].sort((a, b) => b.prefix.length - a.prefix.length);
  for (const route of sorted) {
    if (pathname === route.prefix || pathname.startsWith(`${route.prefix}/`)) {
      return { slug: route.slug, label: route.label };
    }
  }
  return null;
}
