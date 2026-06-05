<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    getGoogleWorkspaceStatus,
    type GoogleWorkspaceStatus,
  } from "$lib/work/googleWorkspaceService";
  import { Button, PageShell, Panel, Skeleton, StatusBanner } from "$lib/ui";

  let loading = $state(true);
  let error = $state("");
  let status = $state<GoogleWorkspaceStatus | null>(null);

  onMount(async () => {
    try {
      status = await getGoogleWorkspaceStatus();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao ler estado Google";
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head><title>Google Workspace — AegisPass</title></svelte:head>

<PageShell
  title="Google Workspace"
  taskId="GOOGLE-001"
  description="Estado do provider e scopes planeados para Drive ZK e masking de Sheets."
 
>
  {#snippet actions()}
    <DocHelpLink slug="journey-google-dev-stub" label="Stub de desenvolvimento" />
    <Button variant="ghost" size="sm" href="/work">← Trabalho</Button>
  {/snippet}

  {#if loading}
    <Skeleton variant="block" height="6rem" />
  {:else if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {:else if status}
    <Panel title="Estado do provider">
      <dl class="status-grid">
        <dt>Provider</dt>
        <dd><code>{status.provider}</code></dd>
        <dt>Pronto</dt>
        <dd>{status.ready ? "sim" : "não"}</dd>
        <dt>Activado</dt>
        <dd>{status.enabled ? "sim" : "não"}</dd>
        {#if status.delegated_user}
          <dt>Utilizador delegado</dt>
          <dd>{status.delegated_user}</dd>
        {/if}
      </dl>
      {#if status.message}
        <p class="hint">{status.message}</p>
      {/if}
      <h2>Scopes planeados</h2>
      <ul>
        {#each status.scopes as scope}
          <li><code>{scope}</code></li>
        {/each}
      </ul>
    </Panel>

    {#if status.provider === "mock"}
      <Panel title="Desenvolvimento local">
        <p>
          Em desenvolvimento local usa o <strong>stub</strong> — Drive cifrado e Sheets com
          masking sem ligar à Google.
        </p>
        <Button href="/work/google-dev">Abrir simulação dev</Button>
      </Panel>
    {:else}
      <Panel title="Produção">
        <p>Service Account configurada. GOOGLE-002/003 (Drive ZK + masking) em breve.</p>
      </Panel>
    {/if}
  {/if}
</PageShell>

<style>
  .status-grid {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: var(--space-2) var(--space-4);
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
  }
  dt {
    color: var(--color-text-label);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
  }
  dd {
    margin: 0;
  }
  h2 {
    margin: var(--space-4) 0 var(--space-2);
    font-size: var(--text-sm);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  ul {
    margin: 0;
    padding-left: var(--space-4);
    font-size: var(--text-sm);
  }
  .hint {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
