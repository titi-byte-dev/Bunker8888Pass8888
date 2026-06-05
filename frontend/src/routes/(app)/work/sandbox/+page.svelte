<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import {
    SANDBOX_DEMO_TARGET_PATH,
    isSandboxReadyPayload,
    postSandboxFill,
  } from "$lib/sandbox";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { loadDecodedLogins, type DecodedLogin } from "$lib/vault/ui";

  type SandboxMode = "demo" | "external";

  let logins = $state<DecodedLogin[]>([]);
  let selectedId = $state("");
  let mode = $state<SandboxMode>("demo");
  let externalUrl = $state("https://example.com");
  let iframeEl = $state<HTMLIFrameElement | undefined>(undefined);
  let iframeReady = $state(false);
  let busy = $state(true);
  let injecting = $state(false);
  let error = $state("");
  let status = $state("");

  const selected = $derived(logins.find((l) => l.meta.id === selectedId) ?? null);

  const iframeSrc = $derived(
    mode === "demo" ? SANDBOX_DEMO_TARGET_PATH : externalUrl.trim() || "about:blank",
  );

  const targetOrigin = $derived(
    typeof window !== "undefined" ? window.location.origin : "",
  );

  onMount(() => {
    void (async () => {
      try {
        logins = await loadDecodedLogins();
        if (logins[0]) selectedId = logins[0].meta.id;
      } catch (e) {
        error = e instanceof Error ? e.message : "Cofre indisponível";
      } finally {
        busy = false;
      }
    })();

    function onMessage(e: MessageEvent) {
      if (e.origin !== targetOrigin) return;
      if (isSandboxReadyPayload(e.data)) {
        iframeReady = true;
      }
    }
    window.addEventListener("message", onMessage);
    return () => window.removeEventListener("message", onMessage);
  });

  function onIframeLoad() {
    iframeReady = false;
  }

  async function injectCredentials() {
    if (!selected || !iframeEl) return;
    if (mode !== "demo") {
      status = "Injeção automática só disponível no alvo demo (same-origin). Sites externos exigem extensão desktop.";
      return;
    }
    injecting = true;
    error = "";
    status = "";
    try {
      postSandboxFill(
        iframeEl,
        { username: selected.login.username, password: selected.login.password },
        targetOrigin,
      );
      status = "Credenciais injectadas — a password nunca foi mostrada neste painel.";
    } catch (e) {
      error = e instanceof Error ? e.message : "Injeção falhou";
    } finally {
      injecting = false;
    }
  }

  /** Pre-fill from vault detail link ?item=id */
  $effect(() => {
    const q = page.url.searchParams.get("item");
    if (q && logins.some((l) => l.meta.id === q)) selectedId = q;
  });
</script>

<svelte:head>
  <title>Sandbox — AegisPass</title>
</svelte:head>

<section class="page">
  <a href="/work" class="back">← Trabalho</a>
  <h1>Browser sandbox</h1>
  <DocHelpLink />
  <p class="lead">
    Contexto isolado para login. A Master Password <strong>nunca</strong> aparece
    nem é copiável neste painel — só é enviada ao iframe via postMessage.
  </p>

  {#if busy}
    <p class="muted">A carregar logins…</p>
  {:else if error && logins.length === 0}
    <p class="error" role="alert">{error}</p>
  {:else}
    <div class="controls">
      <label>
        Login do cofre
        <select bind:value={selectedId}>
          {#each logins as { meta, login } (meta.id)}
            <option value={meta.id}>{login.title}</option>
          {/each}
        </select>
      </label>

      <label>
        Modo
        <select bind:value={mode}>
          <option value="demo">Demo (injecção activa)</option>
          <option value="external">URL externa (só visualização)</option>
        </select>
      </label>

      {#if mode === "external"}
        <label>
          URL
          <input type="url" bind:value={externalUrl} placeholder="https://…" />
        </label>
        <p class="warn">
          Sites de terceiros bloqueiam injecção cross-origin. Usa o modo demo ou a
          extensão/desktop agent (VAULT-017) para produção.
        </p>
      {/if}

      {#if selected}
        <div class="cred-preview">
          <span>Utilizador: <strong>{selected.login.username || "—"}</strong></span>
          <span>Password: <strong class="masked" aria-label="Oculta">••••••••</strong></span>
        </div>
      {/if}

      <div class="actions">
        <button type="button" onclick={injectCredentials} disabled={injecting || !selected || mode !== "demo"}>
          {injecting ? "A injectar…" : "Injectar credenciais"}
        </button>
        {#if !iframeReady && mode === "demo"}
          <span class="muted">Aguarda iframe…</span>
        {/if}
      </div>

      {#if status}
        <p class="status">{status}</p>
      {/if}
      {#if error}
        <p class="error" role="alert">{error}</p>
      {/if}
    </div>

    <div class="frame-wrap">
      <iframe
        bind:this={iframeEl}
        title="Sandbox browser"
        src={iframeSrc}
        sandbox="allow-scripts allow-forms allow-same-origin allow-popups"
        referrerpolicy="no-referrer"
        onload={onIframeLoad}
      ></iframe>
    </div>
  {/if}
</section>

<style>
  .page {
    max-width: 52rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
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
    line-height: var(--leading-body);
  }

  .controls {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    margin-bottom: var(--space-4);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-sm);
    color: var(--color-text-label);
  }

  select,
  input {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: inherit;
    font-family: var(--font-ui);
  }

  .cred-preview {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-4);
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    border: 1px solid var(--color-border);
    font-size: var(--text-sm);
  }

  .masked {
    user-select: none;
    letter-spacing: 2px;
  }

  .warn {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-warning);
  }

  .actions {
    display: flex;
    align-items: center;
    gap: var(--space-3);
  }

  .actions button {
    padding: var(--space-2) var(--space-4);
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    font-weight: 600;
    cursor: pointer;
  }

  .actions button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .status {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-success-fg);
  }

  .frame-wrap {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
    min-height: 22rem;
    background: var(--color-bg-inset);
  }

  iframe {
    width: 100%;
    height: 22rem;
    border: none;
    display: block;
  }

  .error {
    padding: var(--space-3);
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
