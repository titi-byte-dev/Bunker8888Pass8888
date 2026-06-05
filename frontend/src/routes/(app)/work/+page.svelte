<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import ShiftStatusCard from "$lib/work/ShiftStatusCard.svelte";
  import { HubLinks, PageShell } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";
  import { routeChildren } from "$lib/shell/routes";

  const descriptions: Record<string, string> = {
    "/work/shifts": "Estado horário, countdown, desvio NTP e zona GPS (VAULT-010 / VAULT-011).",
    "/work/sandbox": "Injeção de credenciais sem revelar a password no painel (VAULT-013).",
    "/work/cli": "Registo de dispositivo, listagem e injecção em scripts (VAULT-017).",
    "/work/inventory": "Stock operacional, alertas e ordens de compra sugeridas (AGENT-008).",
    "/work/google": "Estado do provider OAuth / Service Account (GOOGLE-001).",
    "/work/google-dev": "Drive cifrado + mascaramento Sheets sem OAuth real (DoD Fase 2).",
  };

  const items: HubLinkItem[] = [
    ...routeChildren("/work").map((c) => ({
      href: c.href,
      title: c.label,
      description: descriptions[c.href],
      taskId: c.taskId,
      comingSoon: c.comingSoon,
    })),
    {
      href: "/security/devices",
      title: "Dispositivos CLI activos",
      description: "Revogar certificados mTLS registados na conta.",
      taskId: "SEC-003",
    },
  ];
</script>

<svelte:head>
  <title>Trabalho — AegisPass</title>
</svelte:head>

<PageShell
  title="Trabalho"
  description="Turnos, sandbox browser, CLI e ferramentas operacionais BYOD."
 
>
  {#snippet actions()}
    <DocHelpLink />
  {/snippet}

  <ShiftStatusCard />
  <HubLinks {items} />
</PageShell>
