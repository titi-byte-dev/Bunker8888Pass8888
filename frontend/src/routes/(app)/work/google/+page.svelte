<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    getGoogleWorkspaceStatus,
    type GoogleWorkspaceStatus,
  } from "$lib/work/googleWorkspaceService";

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

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">GOOGLE-001</p>
      <h1>Google Workspace</h1>
      <DocHelpLink slug="journey-google-dev-stub" label="Stub de desenvolvimento" />
    </div>
    <a class="back" href="/work">← Trabalho</a>
  </header>

  {#if loading}
    <p class="muted">A verificar provider…</p>
  {:else if error}
    <p class="error" role="alert">{error}</p>
  {:else if status}
    <section class="panel">
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
    </section>

    {#if status.provider === "mock"}
      <section class="panel cta">
        <p>
          Em desenvolvimento local usa o <strong>stub</strong> — Drive cifrado e Sheets com
          masking sem ligar à Google.
        </p>
        <a class="btn primary" href="/work/google-dev">Abrir simulação dev</a>
      </section>
    {:else}
      <section class="panel cta">
        <p>Service Account configurada. GOOGLE-002/003 (Drive ZK + masking) em breve.</p>
      </section>
    {/if}
  {/if}
</section>

<style>
  .page { max-width: 42rem; }
  .page-head {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: var(--space-5);
  }
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  h1 { margin: 0; font-family: var(--font-display); }
  .back { font-size: var(--text-sm); color: var(--color-text-muted); }
  .panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--space-4);
    margin-bottom: var(--space-4);
  }
  .status-grid {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: var(--space-2) var(--space-4);
    margin: 0 0 var(--space-4);
  }
  dt { color: var(--color-text-muted); font-size: var(--text-sm); }
  dd { margin: 0; }
  .hint { font-size: var(--text-sm); color: var(--color-text-muted); }
  ul { margin: 0; padding-left: 1.25rem; font-size: var(--text-sm); }
  .cta p { margin: 0 0 var(--space-3); }
  .btn.primary {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    background: var(--color-accent);
    color: var(--color-on-accent);
    border-radius: var(--radius-sm);
    text-decoration: none;
    font-weight: 600;
  }
  .muted, .error { font-size: var(--text-sm); }
  .error { color: var(--color-danger); }
</style>
