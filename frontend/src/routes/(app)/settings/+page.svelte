<script lang="ts">
  import { onMount } from "svelte";
  import {
    loadThemePreference,
    setThemePreference,
    themeModeLabel,
    type ThemeMode,
    PALETTES,
    loadPalettePreference,
    setPalettePreference,
    paletteLabel,
    type PaletteId,
  } from "$lib/design";
  import { fetchSessionProfile } from "$lib/auth/session-api";
  import { loadUserEmail } from "$lib/session";
  import {
    listRegisteredPasskeys,
  } from "$lib/security/api";
  import { passkeysSupported, registerPasskey, type PasskeyMeta } from "$lib/passkey";
  import { loadSessionToken } from "$lib/session";
  import { PageShell, StatusBanner } from "$lib/ui";

  let themeMode = $state<ThemeMode>(loadThemePreference());
  let palette = $state<PaletteId>(loadPalettePreference());
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

  function setPalette(id: PaletteId) {
    palette = id;
    setPalettePreference(id);
    status = `Paleta: ${paletteLabel(id)}.`;
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

<PageShell
  title="Definições"
  taskId="UI-001"
  description="Conta, aparência e autenticação sem password."
>
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

      <p class="hint">
        Paleta de identidade — muda só a cor de marca, não o layout. A empresa pode
        definir uma paleta por omissão (white-label); o modo claro/escuro mantém-se.
      </p>
      <div class="palette-grid" role="radiogroup" aria-label="Paleta">
        {#each PALETTES as p (p.id)}
          <button
            type="button"
            class="palette-card"
            class:active={palette === p.id}
            role="radio"
            aria-checked={palette === p.id}
            onclick={() => setPalette(p.id)}
          >
            <span class="swatch" style="background:{p.swatch}" aria-hidden="true"></span>
            <span class="p-name">
              {p.label}
              {#if palette === p.id}<span class="p-on" aria-hidden="true">✓</span>{/if}
            </span>
            <span class="p-desc">{p.personality}</span>
          </button>
        {/each}
      </div>
    </section>

    <section class="block">
      <h2>Documentação</h2>
      <p class="hint">
        Aprende a usar o AegisPass por níveis de complexidade — conceitos-chave em
        dropdowns, secções técnicas colapsáveis. Fonte única em <code>docs/</code>.
      </p>
      <a href="/settings/docs" class="docs-link">Abrir documentação</a>
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
</PageShell>

<style>
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

  .palette-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(13rem, 1fr));
    gap: var(--space-3);
    margin-top: var(--space-3);
  }

  .palette-card {
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: auto auto;
    column-gap: var(--space-3);
    row-gap: var(--space-1);
    align-items: center;
    text-align: left;
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    color: var(--color-text);
    cursor: pointer;
    transition: border-color var(--duration-fast) var(--ease-out);
  }

  .palette-card:hover {
    border-color: var(--color-border-strong);
  }

  .palette-card.active {
    border-color: var(--color-accent);
    box-shadow: 0 0 0 1px var(--color-accent);
  }

  .swatch {
    grid-row: 1 / span 2;
    width: 2rem;
    height: 2rem;
    border-radius: var(--radius-sm);
    border: 1px solid rgba(0, 0, 0, 0.2);
  }

  .p-name {
    font-size: var(--text-sm);
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .p-on {
    color: var(--color-accent);
  }

  .p-desc {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    line-height: 1.4;
  }

  @media (prefers-reduced-motion: reduce) {
    .palette-card {
      transition: none;
    }
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

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .docs-link {
    display: inline-flex;
    margin-top: var(--space-3);
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-accent);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    color: var(--color-text);
    font-size: var(--text-sm);
    font-weight: 500;
    text-decoration: none;
  }

  .docs-link:hover {
    background: color-mix(in srgb, var(--color-accent) 20%, var(--color-bg-surface));
  }
</style>
