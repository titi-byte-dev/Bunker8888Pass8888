<script lang="ts">
  import { onMount } from "svelte";
  import {
    listCliDevices,
    listHttpSessions,
    listRegisteredPasskeys,
    revokeCliDevice,
    revokeHttpSession,
    revokeOtherHttpSessions,
    type CliDevice,
    type HttpSession,
  } from "$lib/security/api";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import type { PasskeyMeta } from "$lib/passkey";
  import { PageShell, StatusBanner } from "$lib/ui";

  let sessions = $state<HttpSession[]>([]);
  let passkeys = $state<PasskeyMeta[]>([]);
  let cliDevices = $state<CliDevice[]>([]);
  let busy = $state(true);
  let error = $state("");
  let status = $state("");

  async function refresh() {
    busy = true;
    error = "";
    status = "";
    try {
      [sessions, passkeys, cliDevices] = await Promise.all([
        listHttpSessions(),
        listRegisteredPasskeys(),
        listCliDevices(),
      ]);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar dispositivos";
    } finally {
      busy = false;
    }
  }

  async function handleRevokeSession(id: string) {
    try {
      await revokeHttpSession(id);
      status = "Sessão revogada.";
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao revogar";
    }
  }

  async function handleRevokeOthers() {
    try {
      const n = await revokeOtherHttpSessions();
      status = `${n} sessão(ões) revogada(s).`;
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao revogar";
    }
  }

  async function handleRevokeCli(id: string) {
    if (!confirm("Revogar este dispositivo CLI? O certificado deixa de funcionar.")) return;
    try {
      await revokeCliDevice(id);
      status = "Dispositivo CLI revogado.";
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao revogar";
    }
  }

  onMount(refresh);
</script>

<svelte:head>
  <title>Dispositivos e sessões — AegisPass</title>
</svelte:head>

<PageShell
  title="Dispositivos e sessões"
  taskId="SEC-003"
  description="Sessões HTTP, passkeys WebAuthn e certificados CLI activos na tua conta."
>
  {#snippet actions()}
    <DocHelpLink />
  {/snippet}

  {#if status}
    <StatusBanner variant="success">{status}</StatusBanner>
  {/if}
  {#if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {/if}

  {#if busy}
    <p class="muted">A carregar…</p>
  {:else}
    <section class="block">
      <div class="block-header">
        <h2>Sessões HTTP</h2>
        {#if sessions.length > 1}
          <button type="button" class="linkish" onclick={handleRevokeOthers}>Terminar outras sessões</button>
        {/if}
      </div>
      {#if sessions.length === 0}
        <p class="muted">Nenhuma sessão activa.</p>
      {:else}
        <ul>
          {#each sessions as s (s.id)}
            <li>
              <div>
                <strong>{s.current ? "Este browser" : "Outra sessão"}</strong>
                <span class="meta">Criada {new Date(s.created_at).toLocaleString("pt-PT")}</span>
              </div>
              {#if !s.current}
                <button type="button" onclick={() => handleRevokeSession(s.id)}>Revogar</button>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section class="block">
      <h2>Passkeys</h2>
      {#if passkeys.length === 0}
        <p class="muted">
          Nenhuma passkey registada.
          <a href="/settings">Registar em Definições</a>.
        </p>
      {:else}
        <ul>
          {#each passkeys as pk (pk.id)}
            <li>
              <strong>{pk.name}</strong>
              <span class="meta">{new Date(pk.created_at).toLocaleDateString("pt-PT")}</span>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section class="block">
      <h2>Dispositivos CLI (mTLS)</h2>
      {#if cliDevices.length === 0}
        <p class="muted">Nenhum dispositivo CLI registado.</p>
      {:else}
        <ul>
          {#each cliDevices as d (d.id)}
            <li>
              <div>
                <strong>{d.name}</strong>
                <span class="meta">{new Date(d.created_at).toLocaleDateString("pt-PT")}</span>
              </div>
              <button type="button" onclick={() => handleRevokeCli(d.id)}>Revogar</button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {/if}
</PageShell>

<style>
  .block {
    margin-bottom: var(--space-8);
  }

  .block-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    margin-bottom: var(--space-3);
  }

  h2 {
    margin: 0;
    font-size: var(--text-lg);
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
    font-size: var(--text-sm);
  }

  li:last-child {
    border-bottom: none;
  }

  .meta {
    display: block;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    margin-top: var(--space-1);
  }

  button {
    padding: var(--space-1) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    color: var(--color-text);
    font-size: var(--text-xs);
    cursor: pointer;
  }

  .linkish {
    border: none;
    background: transparent;
    color: var(--color-link);
    font-size: var(--text-sm);
    cursor: pointer;
    text-decoration: underline;
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
</style>
