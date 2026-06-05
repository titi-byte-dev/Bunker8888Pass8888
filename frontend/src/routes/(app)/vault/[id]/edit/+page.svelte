<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import LoginForm from "$lib/vault/LoginForm.svelte";
  import { loadDecodedLogin, requireVaultAccess, type DecodedLogin } from "$lib/vault/ui";
  import { blobToBase64, sealItem } from "$lib/vault/items";
  import type { LoginItem } from "$lib/vault/types";
  import { Button, PageShell, Skeleton, StatusBanner } from "$lib/ui";

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

<PageShell
  title="Editar login"
  taskId="VAULT-001"
  leaf={item?.login.title}
  breadcrumb={false}
  description="Altera credenciais localmente — o servidor só recebe o blob cifrado."
>
  {#snippet actions()}
    <Button variant="ghost" size="sm" href="/vault/{id}">← Voltar</Button>
  {/snippet}

  {#if busy}
    <Skeleton variant="block" height="2rem" />
    <Skeleton variant="block" height="6rem" />
  {:else if error && !item}
    <StatusBanner variant="error">{error}</StatusBanner>
    <Button variant="ghost" href="/vault">Voltar ao cofre</Button>
  {:else if item}
    {#if error}
      <StatusBanner variant="error">{error}</StatusBanner>
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
</PageShell>
