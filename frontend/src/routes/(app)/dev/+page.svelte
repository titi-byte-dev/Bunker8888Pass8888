<script lang="ts">
  import { HubLinks, PageShell } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";

  const items: HubLinkItem[] = [
    {
      href: "/dev/playground",
      title: "Cofre playground",
      description: "CRUD end-to-end legacy (VAULT-019).",
      taskId: "VAULT-019",
    },
    {
      href: "/dev/components",
      title: "Catálogo de componentes",
      description: "Tokens, botões, forms e previews (UI-010).",
      taskId: "UI-010",
    },
  ];
</script>

<svelte:head>
  <title>Desenvolvimento — AegisPass</title>
</svelte:head>

<PageShell
  title="Desenvolvimento"
  description="Ferramentas de teste e referência visual — só em build dev."
>
  <HubLinks {items} />
  <p class="hint">Motion demo e ArgonProgress estão no catálogo de componentes.</p>
</PageShell>

<style>
  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
