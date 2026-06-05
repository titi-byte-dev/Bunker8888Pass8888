<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import LoginForm from "$lib/vault/LoginForm.svelte";
  import { loadDecodedLogin, requireVaultAccess, type DecodedLogin } from "$lib/vault/ui";
  import { blobToBase64, sealItem } from "$lib/vault/items";
  import type { LoginItem } from "$lib/vault/types";

  let item = $state<DecodedLogin | null>(null);
  let busy = $state(true);
  let saving = $state(false);
  let error = $state("");

  const id = $derived(page.params.id ?? "");
  const remediate = $derived(page.url.searchParams.get("remediate"));

  const remediationNote = $derived(
    remediate === "breach" || remediate === "weak_and_breach"
      ? "⚠️ Esta password apareceu em fugas de dados. Escolhe uma password nova e única antes de guardar."
      : "",
  );

  async function load() {
    busy = true;
    error = "";
    try {
      item = await loadDecodedLogin(id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Item não encontrado";
    } finally {
      busy = false;
    }
  }

  async function handleSave(payload: LoginItem) {
    if (!item) return;
    saving = true;
    error = "";
    try {
      const { api, key } = requireVaultAccess();
      const blob = blobToBase64(await sealItem(key, payload));
      await api.update(item.meta.id, { type: "login", blob });
      await goto(`/vault/${item.meta.id}`);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao guardar";
      saving = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Editar — AegisPass</title>
</svelte:head>

<section class="vault-page">
  <a href="/vault/{id}" class="back">← Voltar</a>

  {#if busy}
    <p class="muted">A carregar…</p>
  {:else if error && !item}
    <p class="error" role="alert">{error}</p>
  {:else if item}
    <h1>Editar login</h1>
    {#if error}
      <p class="error" role="alert">{error}</p>
    {/if}
    <LoginForm
      initial={item.login}
      submitLabel="Guardar alterações"
      busy={saving}
      focusPassword={!!remediationNote}
      remediationNote={remediationNote}
      onsubmit={handleSave}
    />
  {/if}
</section>

<style>
  .vault-page {
    max-width: 28rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  h1 {
    margin: 0 0 var(--space-6);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .error {
    margin-bottom: var(--space-4);
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
