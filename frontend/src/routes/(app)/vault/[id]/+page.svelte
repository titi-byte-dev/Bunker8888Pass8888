<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import CopyField from "$lib/vault/CopyField.svelte";
  import { loadDecodedLogin, requireVaultAccess, type DecodedLogin } from "$lib/vault/ui";
  import {
    Button,
    confirmDialog,
    PageShell,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";

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
    if (!item) return;
    const ok = await confirmDialog({
      title: "Apagar login?",
      message: `Remove «${item.login.title}». Esta acção é irreversível.`,
      confirmLabel: "Apagar",
      variant: "danger",
    });
    if (!ok) return;
    deleting = true;
    try {
      const { api } = requireVaultAccess();
      await api.delete(item.meta.id);
      toast.success("Login apagado.");
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

<PageShell
  title={item?.login.title ?? "Login"}
  taskId="VAULT-001"
  leaf={item?.login.title}
  description={busy ? "A descifrar localmente…" : "Credenciais desbloqueadas com a Master Key — o servidor nunca vê estes valores."}
  breadcrumb={false}
>
  {#snippet actions()}
    <Button variant="ghost" size="sm" href="/vault">← Cofre</Button>
    {#if item}
      <Button variant="secondary" size="sm" href="/vault/{item.meta.id}/edit">Editar</Button>
    {/if}
  {/snippet}

  {#if busy}
    <Skeleton variant="block" height="2rem" />
    <Skeleton variant="block" height="4rem" />
    <Skeleton variant="block" height="4rem" />
  {:else if error}
    <StatusBanner variant="error">{error}</StatusBanner>
    <Button variant="ghost" href="/vault">Voltar ao cofre</Button>
  {:else if item}
    <CopyField label="Utilizador" value={item.login.username} />
    <CopyField label="Password" value={item.login.password} secret />

    <p class="sandbox-link">
      <a href="/work/sandbox?item={item.meta.id}">Abrir no browser sandbox</a>
    </p>

    {#if item.login.url}
      <CopyField label="URL" value={item.login.url} />
    {/if}

    {#if item.login.notes}
      <CopyField label="Notas" value={item.login.notes} />
    {/if}

    <p class="meta muted">
      Actualizado {new Date(item.meta.updated_at).toLocaleString("pt-PT")}
    </p>

    <Button variant="ghost" onclick={handleDelete} disabled={deleting} loading={deleting}>
      Apagar login
    </Button>
  {/if}
</PageShell>

<style>
  .sandbox-link {
    margin: var(--space-4) 0 0;
    font-size: var(--text-sm);
  }
  .sandbox-link a {
    color: var(--color-link);
    text-decoration: none;
  }
  .meta {
    margin: var(--space-6) 0 var(--space-4);
    font-size: var(--text-xs);
  }
  .muted {
    color: var(--color-text-muted);
  }
</style>
