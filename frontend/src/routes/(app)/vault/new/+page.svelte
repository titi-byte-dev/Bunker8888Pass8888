<script lang="ts">
  import { goto } from "$app/navigation";
  import LoginForm from "$lib/vault/LoginForm.svelte";
  import { requireVaultAccess } from "$lib/vault/ui";
  import { blobToBase64, sealItem } from "$lib/vault/items";
  import type { LoginItem } from "$lib/vault/types";

  let busy = $state(false);
  let error = $state("");

  async function handleCreate(payload: LoginItem) {
    busy = true;
    error = "";
    try {
      const { api, key } = requireVaultAccess();
      const blob = blobToBase64(await sealItem(key, payload));
      const meta = await api.create({ type: "login", blob });
      await goto(`/vault/${meta.id}`);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao guardar";
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Novo login — AegisPass</title>
</svelte:head>

<section class="vault-page">
  <header class="page-header">
    <a href="/vault" class="back">← Cofre</a>
    <h1>Novo login</h1>
    <p class="subtitle">Cifragem AES-GCM no browser antes de enviar ao servidor.</p>
  </header>

  {#if error}
    <p class="error" role="alert">{error}</p>
  {/if}

  <LoginForm submitLabel="Guardar cifrado" busy={busy} onsubmit={handleCreate} />
</section>

<style>
  .vault-page {
    max-width: 28rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-2);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .subtitle {
    margin: var(--space-2) 0 var(--space-6);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .error {
    margin-bottom: var(--space-4);
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }
</style>
