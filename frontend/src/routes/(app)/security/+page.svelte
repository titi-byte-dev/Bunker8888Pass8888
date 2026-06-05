<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { HubLinks, PageShell } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";
  import { routeChildren } from "$lib/shell/routes";

  const descriptions: Record<string, string> = {
    "/security/hygiene": "Score composto, fugas k-anonymity, alteração forçada (UI-007 / DW-002 / DW-003).",
    "/security/devices": "Sessões HTTP, passkeys e certificados CLI.",
    "/security/sentinel": "Logins impossíveis, step-up passkey e histórico (DW-004).",
    "/security/emergency": "Herdeiro digital, countdown, aprovar/rejeitar (VAULT-016).",
    "/security/guardian": "Tools de agentes sem exposição da Master Key (DoD Fase 3).",
  };

  const items: HubLinkItem[] = routeChildren("/security").map((c) => ({
    href: c.href,
    title: c.label,
    description: descriptions[c.href],
    taskId: c.taskId,
    comingSoon: c.comingSoon,
  }));
</script>

<svelte:head>
  <title>Segurança — AegisPass</title>
</svelte:head>

<PageShell
  title="Segurança"
  description="Higiene, dispositivos, sessões e confiança operacional."
 
>
  {#snippet actions()}
    <DocHelpLink />
  {/snippet}

  <HubLinks {items} />
</PageShell>
