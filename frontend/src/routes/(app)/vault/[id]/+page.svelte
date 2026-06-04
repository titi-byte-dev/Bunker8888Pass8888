<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import CopyField from "$lib/vault/CopyField.svelte";
  import { loadDecodedLogin, requireVaultAccess, type DecodedLogin } from "$lib/vault/ui";

  let item = $state<DecodedLogin | null>(null);
  let busy = $state(true);
  let error = $state("");
  let deleting = $state(false);

  const id = $derived(page.params.id ?? "");

  async function load() {
    busy = true;
    error = "";
    try {
      item = await loadDecodedLogin(id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Item não encontrado";
      item = null;
    } finally {
      busy = false;
    }
  }

  async function handleDelete() {
    if (!item || !confirm(`Apagar «${item.login.title}»? Esta acção é irreversível.`)) return;
    deleting = true;
    try {
      const { api } = requireVaultAccess();
      await api.delete(item.meta.id);
      await goto("/vault");
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao apagar";
      deleting = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>{item?.login.title ?? "Login"} — AegisPass</title>
</svelte:head>

<section class="vault-page">
  <a href="/vault" class="back">← Cofre</a>

  {#if busy}
    <p class="muted">A descifrar…</p>
  {:else if error}
    <p class="error" role="alert">{error}</p>
  {:else if item}
    <header class="page-header">
      <h1>{item.login.title}</h1>
      <a class="edit" href="/vault/{item.meta.id}/edit">Editar</a>
    </header>

    <CopyField label="Utilizador" value={item.login.username} />
    <CopyField label="Password" value={item.login.password} secret />

    {#if item.login.url}
      <CopyField label="URL" value={item.login.url} />
    {/if}

    {#if item.login.notes}
      <CopyField label="Notas" value={item.login.notes} />
    {/if}

    <p class="meta muted">
      Actualizado {new Date(item.meta.updated_at).toLocaleString("pt-PT")}
    </p>

    <button type="button" class="danger" onclick={handleDelete} disabled={deleting}>
      {deleting ? "A apagar…" : "Apagar login"}
    </button>
  {/if}
</section>

<style>
  .vault-page {
    max-width: 32rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    margin-bottom: var(--space-6);
  }

  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .edit {
    font-size: var(--text-sm);
    color: var(--color-link);
    text-decoration: none;
  }

  .meta {
    margin: var(--space-6) 0;
    font-size: var(--text-xs);
  }

  .danger {
    width: 100%;
    padding: var(--space-3);
    border: 1px solid color-mix(in srgb, var(--color-danger) 40%, transparent);
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-danger);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  .danger:disabled {
    opacity: 0.5;
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
  }
</style>
