<script lang="ts">
  import { onMount } from "svelte";
  import { filterLogins, loadDecodedLogins, type DecodedLogin } from "$lib/vault/ui";
  import { animateListStagger } from "$lib/motion/presets";
  import { runMotionScope } from "$lib/motion/scope";

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

  async function refresh() {
    busy = true;
    error = "";
    try {
      items = await loadDecodedLogins();
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

<section class="vault-page">
  <header class="page-header">
    <div>
      <h1>Cofre</h1>
      <p class="subtitle">{items.length} login{items.length === 1 ? "" : "s"} cifrados</p>
    </div>
    <div class="header-actions">
      <button type="button" class="secondary" onclick={refresh} disabled={busy}>Actualizar</button>
      <a class="primary" href="/vault/new">Novo login</a>
    </div>
  </header>

  <label class="search">
    <span class="sr-only">Pesquisar</span>
    <input type="search" bind:value={query} placeholder="Pesquisar título, utilizador ou URL…" disabled={busy} />
  </label>

  {#if error}
    <p class="error" role="alert">{error}</p>
  {:else if busy}
    <p class="muted">A descifrar itens localmente…</p>
  {:else if filtered.length === 0}
    <div class="empty">
      <p>{query ? "Nenhum resultado." : "Cofre vazio — adiciona o primeiro login."}</p>
      <a href="/vault/new">Criar login</a>
    </div>
  {:else}
    <ul class="list" bind:this={listRoot}>
      {#each filtered as { meta, login } (meta.id)}
        <li>
          <a href="/vault/{meta.id}" class="row-link">
            <span class="title">{login.title}</span>
            <span class="meta">{login.username || "—"}</span>
          </a>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .vault-page {
    max-width: 40rem;
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-4);
    margin-bottom: var(--space-6);
  }

  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .subtitle {
    margin: var(--space-1) 0 0;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .header-actions {
    display: flex;
    gap: var(--space-2);
    flex-shrink: 0;
  }

  .primary,
  .secondary {
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    font-weight: 600;
    text-decoration: none;
    border: none;
    cursor: pointer;
  }

  .primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
  }

  .secondary {
    background: var(--color-accent-muted);
    color: var(--color-text);
    border: 1px solid var(--color-border);
  }

  .search {
    display: block;
    margin-bottom: var(--space-4);
  }

  .search input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: inherit;
    box-sizing: border-box;
  }

  .list {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .row-link {
    display: flex;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-3) var(--space-4);
    text-decoration: none;
    color: inherit;
    border-bottom: 1px solid var(--color-border);
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .row-link:last-child {
    border-bottom: none;
  }

  .row-link:hover {
    background: var(--color-bg-surface);
  }

  .title {
    font-weight: 500;
  }

  .meta {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .empty {
    padding: var(--space-8);
    text-align: center;
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-md);
    color: var(--color-text-muted);
  }

  .empty a {
    color: var(--color-link);
  }

  .error {
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
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
