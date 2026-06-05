<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import AdminGate from "$lib/admin/AdminGate.svelte";
  import { hasAdminKey } from "$lib/admin/adminKey";
  import { HubLinks, PageShell } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";
  import { routeChildren } from "$lib/shell/routes";

  let unlocked = $state(hasAdminKey());

  function refreshGate() {
    unlocked = hasAdminKey();
  }

  const descriptions: Record<string, string> = {
    "/admin/users": "Lista, turnos, geofence e remote wipe por colaborador.",
    "/admin/audit": "Eventos de remote wipe (append-only).",
  };

  const items: HubLinkItem[] = routeChildren("/admin").map((c) => ({
    href: c.href,
    title: c.label,
    description: descriptions[c.href],
    taskId: c.taskId,
    comingSoon: c.comingSoon,
  }));
</script>

<svelte:head>
  <title>Administração — AegisPass</title>
</svelte:head>

<PageShell
  title="Administração"
  taskId="UI-008"
  description="Utilizadores, políticas de acesso, auditoria e remote wipe."
>
  {#snippet actions()}
    <DocHelpLink />
  {/snippet}

  <AdminGate onUnlocked={refreshGate} />

  {#if unlocked}
    <HubLinks {items} />
  {/if}
</PageShell>
