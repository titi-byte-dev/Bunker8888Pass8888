<script lang="ts">
  /**
   * Hub Financas (UI-014) — pagina-indice do modulo.
   * Didatico: o /fin deixou de ser a pagina de custos (movida para /fin/costs);
   * passa a ser um HUB que deriva os cartoes da ROUTE_TREE. Assim a navegacao
   * fica intuitiva (uma porta por sub-area) e a sidebar pode ficar limpa.
   */
  import { PageShell, HubLinks } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";
  import { routeChildren } from "$lib/shell/routes";

  const descriptions: Record<string, string> = {
    "/fin/costs": "Subscricoes SaaS cifradas, alertas de licencas esquecidas e poupanca potencial.",
    "/fin/fiscal": "Classificacao IRC calculada no cliente — dedutiveis estimados por subscricao.",
    "/fin/invoices": "Pro-forma, faturas e recibos com numeracao legal.",
    "/fin/commissions": "Comissoes de equipa e parceiros, ligadas ao CRM.",
    "/fin/banking": "Reconciliacao Open Banking contra subscricoes e faturas.",
  };

  const items: HubLinkItem[] = routeChildren("/fin").map((c) => ({
    href: c.href,
    title: c.label,
    description: descriptions[c.href],
    taskId: c.taskId,
    comingSoon: c.comingSoon,
  }));
</script>

<svelte:head><title>Financas — AegisPass</title></svelte:head>

<PageShell
  title="Financas"
  description="Custos, fiscalidade e faturacao — tudo cifrado com a tua Master Key. Escolhe uma area para comecar."
>
  <HubLinks {items} />
</PageShell>
