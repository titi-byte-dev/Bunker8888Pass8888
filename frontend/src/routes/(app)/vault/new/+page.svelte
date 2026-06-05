<script lang="ts">
  import { goto } from "$app/navigation";
  import LoginForm from "$lib/vault/LoginForm.svelte";
  import { requireVaultAccess } from "$lib/vault/ui";
  import { blobToBase64, sealItem } from "$lib/vault/items";
  import type { LoginItem } from "$lib/vault/types";
  import { Button, PageShell, StatusBanner } from "$lib/ui";

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

<PageShell
  title="Novo login"
  taskId="VAULT-001"
  description="Cifragem AES-GCM no browser antes de enviar ao servidor — o plaintext nunca sai do dispositivo."
 
  breadcrumb={false}
>
  {#snippet actions()}
    <Button variant="ghost" size="sm" href="/vault">← Cofre</Button>
  {/snippet}

  {#if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {/if}

  <LoginForm submitLabel="Guardar cifrado" busy={busy} onsubmit={handleCreate} />
</PageShell>
