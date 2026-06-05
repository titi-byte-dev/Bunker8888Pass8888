<script lang="ts">
  import { onMount } from "svelte";
  import { filterLogins, loadDecodedLogins, type DecodedLogin } from "$lib/vault/ui";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { animateListStagger } from "$lib/motion/presets";
  import { runMotionScope } from "$lib/motion/scope";
  import {
    Button,
    EmptyState,
    ListRow,
    PageShell,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";

  let items = $state<DecodedLogin[]>([]);
  let query = $state("");
  let busy = $state(true);
  let error = $state("");

  const filtered = $derived(filterLogins(items, query));

  let listRoot = $state<HTMLElement | undefined>(undefined);
  let lastBusy = true;

  $effect(() => {
    if (lastBusy && !busy && filtered.length > 0 && listRoot) {
      runMotionScope(listRoot, () => animateListStagger(listRoot!, "li"));
    }
    lastBusy = busy;
  });

  async function refresh(showToast = false) {
    busy = true;
    error = "";
    try {
      items = await loadDecodedLogins();
      if (showToast) toast.success("Cofre actualizado.");
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar cofre";
    } finally {
      busy = false;
    }
  }

  onMount(() => {
    refresh();
  });
</script>

<svelte:head>
  <title>Cofre — AegisPass</title>
</svelte:head>

<PageShell
  title="Cofre"
  taskId="VAULT-001"
  description={`${items.length} login${items.length === 1 ? "" : "s"} cifrados — desbloqueados localmente com a Master Key.`}
 
  breadcrumb={false}
>
  {#snippet actions()}
    <DocHelpLink />
    <Button variant="secondary" size="sm" onclick={() => refresh(true)} loading={busy} disabled={busy}>
      Actualizar
    </Button>
    <Button href="/vault/new">Novo login</Button>
  {/snippet}

  <label class="search">
    <span class="sr-only">Pesquisar</span>
    <input
      type="search"
      bind:value={query}
      placeholder="Pesquisar título, utilizador ou URL…"
      disabled={busy}
    />
  </label>

  {#if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {:else if busy}
    <Skeleton variant="row" />
    <Skeleton variant="row" />
    <Skeleton variant="row" />
  {:else if filtered.length === 0}
    <EmptyState
      title={query ? "Nenhum resultado" : "Cofre vazio"}
      description={query ? "Tenta outro termo de pesquisa." : "Adiciona o primeiro login para começar."}
    >
      {#snippet action()}
        <Button href="/vault/new">Criar login</Button>
      {/snippet}
    </EmptyState>
  {:else}
    <ul class="list" bind:this={listRoot}>
      {#each filtered as { meta, login } (meta.id)}
        <li>
          <ListRow
            href="/vault/{meta.id}"
            title={login.title}
            meta={login.username || login.url || "—"}
          />
        </li>
      {/each}
    </ul>
  {/if}
</PageShell>

<style>
  .search {
    display: block;
  }

  .search input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: inherit;
    box-sizing: border-box;
    font-size: var(--text-sm);
  }

  .list {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    border: 0;
  }
</style>
