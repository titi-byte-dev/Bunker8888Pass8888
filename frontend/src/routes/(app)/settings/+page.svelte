<script lang="ts">
  import { onMount } from "svelte";
  import {
    loadThemePreference,
    setThemePreference,
    themeModeLabel,
    type ThemeMode,
  } from "$lib/design";
  import { fetchSessionProfile } from "$lib/auth/session-api";
  import { loadUserEmail } from "$lib/session";
  import {
    listRegisteredPasskeys,
  } from "$lib/security/api";
  import { passkeysSupported, registerPasskey, type PasskeyMeta } from "$lib/passkey";
  import { loadSessionToken } from "$lib/session";

  let themeMode = $state<ThemeMode>(loadThemePreference());
  let email = $state(loadUserEmail() ?? "");
  let passkeys = $state<PasskeyMeta[]>([]);
  let passkeyName = $state("");
  let busy = $state(true);
  let passkeyBusy = $state(false);
  let status = $state("");
  let error = $state("");

  const webauthnOk = passkeysSupported();

  function setTheme(mode: ThemeMode) {
    themeMode = mode;
    setThemePreference(mode);
    status = `Tema: ${themeModeLabel(mode)}.`;
  }

  async function refreshPasskeys() {
    passkeys = await listRegisteredPasskeys();
  }

  async function load() {
    busy = true;
    error = "";
    try {
      const profile = await fetchSessionProfile();
      email = profile.email;
      await refreshPasskeys();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar definições";
    } finally {
      busy = false;
    }
  }

  async function handleRegisterPasskey() {
    const name = passkeyName.trim() || "Passkey";
    const token = loadSessionToken();
    if (!token) {
      error = "Sessão inválida";
      return;
    }
    passkeyBusy = true;
    error = "";
    status = "";
    try {
      await registerPasskey("", token, name);
      passkeyName = "";
      status = "Passkey registada com sucesso.";
      await refreshPasskeys();
    } catch (e) {
      error = e instanceof Error ? e.message : "Registo falhou";
    } finally {
      passkeyBusy = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Definições — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Definições</h1>
  <p class="lead">Conta, aparência e autenticação sem password.</p>

  {#if status}
    <p class="status">{status}</p>
  {/if}
  {#if error}
    <p class="error" role="alert">{error}</p>
  {/if}

  {#if busy}
    <p class="muted">A carregar…</p>
  {:else}
    <section class="block">
      <h2>Conta</h2>
      <dl class="kv">
        <dt>Email</dt>
        <dd>{email || "—"}</dd>
      </dl>
      <p class="hint">
        A Master Password nunca é guardada no servidor (Zero-Knowledge). Para alterar credenciais
        expostas, usa <a href="/security/hygiene">Saúde de segurança</a>.
      </p>
    </section>

    <section class="block">
      <h2>Aparência</h2>
      <p class="hint">Preferência de tema — «Sistema» segue o modo claro/escuro do SO.</p>
      <div class="theme-row" role="radiogroup" aria-label="Tema">
        {#each (["light", "dark", "system"] as const) as mode (mode)}
          <button
            type="button"
            class="theme-opt"
            class:active={themeMode === mode}
            role="radio"
            aria-checked={themeMode === mode}
            onclick={() => setTheme(mode)}
          >
            {themeModeLabel(mode)}
          </button>
        {/each}
      </div>
    </section>

    <section class="block">
      <h2>Passkeys</h2>
      {#if !webauthnOk}
        <p class="muted">WebAuthn não suportado neste browser.</p>
      {:else}
        <p class="hint">
          Passkeys autenticam a sessão HTTP; o cofre continua a exigir Master Password para
          desbloqueio local.
        </p>
        {#if passkeys.length === 0}
          <p class="muted">Nenhuma passkey registada.</p>
        {:else}
          <ul class="list">
            {#each passkeys as pk (pk.id)}
              <li>
                <strong>{pk.name}</strong>
                <span class="meta">{new Date(pk.created_at).toLocaleDateString("pt-PT")}</span>
              </li>
            {/each}
          </ul>
        {/if}
        <div class="passkey-form">
          <input
            type="text"
            bind:value={passkeyName}
            placeholder="Nome da passkey (ex: MacBook)"
            aria-label="Nome da passkey"
          />
          <button type="button" disabled={passkeyBusy} onclick={handleRegisterPasskey}>
            {passkeyBusy ? "A registar…" : "Registar passkey"}
          </button>
        </div>
        <p class="hint">
          Gerir sessões e revogar dispositivos em
          <a href="/security/devices">Dispositivos e sessões</a>.
        </p>
      {/if}
    </section>
  {/if}
</section>

<style>
  .page {
    max-width: 40rem;
  }

  h1 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .lead {
    color: var(--color-text-muted);
    margin: 0 0 var(--space-6);
    font-size: var(--text-sm);
  }

  .block {
    margin-bottom: var(--space-8);
  }

  h2 {
    margin: 0 0 var(--space-3);
    font-size: var(--text-lg);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: var(--space-2) 0 0;
    line-height: 1.5;
  }

  .hint a {
    color: var(--color-link);
  }

  .kv {
    display: grid;
    grid-template-columns: 6rem 1fr;
    gap: var(--space-2);
    margin: 0;
    font-size: var(--text-sm);
  }

  dt {
    color: var(--color-text-muted);
  }

  dd {
    margin: 0;
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    word-break: break-all;
  }

  .theme-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    margin-top: var(--space-3);
  }

  .theme-opt {
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    color: var(--color-text);
    font-size: var(--text-sm);
    cursor: pointer;
  }

  .theme-opt.active {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .list {
    list-style: none;
    margin: var(--space-3) 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .list li {
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
    font-size: var(--text-sm);
  }

  .list li:last-child {
    border-bottom: none;
  }

  .meta {
    display: block;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    margin-top: var(--space-1);
  }

  .passkey-form {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    margin-top: var(--space-4);
  }

  .passkey-form input {
    flex: 1;
    min-width: 12rem;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    color: var(--color-text);
    font-size: var(--text-sm);
  }

  .passkey-form button {
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  .passkey-form button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .status {
    padding: var(--space-3);
    margin-bottom: var(--space-4);
    border-radius: var(--radius-sm);
    background: var(--color-success-bg);
    color: var(--color-success-fg);
    font-size: var(--text-sm);
  }

  .error {
    padding: var(--space-3);
    margin-bottom: var(--space-4);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
</style>
