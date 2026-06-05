<script lang="ts">
  /**
   * Página-alvo demo do sandbox (VAULT-013) — same-origin, sem app shell.
   */
  import { onMount } from "svelte";
  import { isSandboxFillPayload, SANDBOX_READY_MESSAGE } from "$lib/sandbox";
  import "$lib/design/tokens.css";

  let username = $state("");
  let password = $state("");
  let filled = $state(false);

  onMount(() => {
    const origin = window.location.origin;

    window.parent.postMessage({ type: SANDBOX_READY_MESSAGE }, origin);

    function onMessage(e: MessageEvent) {
      if (e.origin !== origin) return;
      if (!isSandboxFillPayload(e.data)) return;

      username = e.data.username;
      password = e.data.password;
      filled = true;

      queueMicrotask(() => {
        const userInput = document.getElementById("sandbox-user") as HTMLInputElement | null;
        const passInput = document.getElementById("sandbox-pass") as HTMLInputElement | null;
        if (userInput) userInput.value = e.data.username;
        if (passInput) passInput.value = e.data.password;
      });
    }

    window.addEventListener("message", onMessage);
    return () => window.removeEventListener("message", onMessage);
  });

  function handleSubmit(e: Event) {
    e.preventDefault();
    password = "";
    filled = false;
  }
</script>

<svelte:head>
  <title>Alvo demo — Sandbox</title>
</svelte:head>

<div class="target">
  <p class="badge">Alvo demo · same-origin</p>
  <h1>Iniciar sessão</h1>
  <p class="hint">Formulário simulado — credenciais chegam por postMessage do AegisPass.</p>

  <form onsubmit={handleSubmit}>
    <label>
      Email / utilizador
      <input id="sandbox-user" type="text" autocomplete="username" bind:value={username} />
    </label>
    <label>
      Password
      <input
        id="sandbox-pass"
        type="password"
        autocomplete="current-password"
        bind:value={password}
        readonly={filled}
      />
    </label>
    <button type="submit">Entrar</button>
  </form>

  {#if filled}
    <p class="ok">✓ Campos preenchidos pelo sandbox (password oculta no painel pai).</p>
  {/if}
</div>

<style>
  :global(body) {
    margin: 0;
    background: var(--color-bg-base);
    color: var(--color-text);
    font-family: var(--font-ui);
  }

  .target {
    max-width: 20rem;
    margin: var(--space-8) auto;
    padding: var(--space-6);
  }

  .badge {
    font-size: var(--text-xs);
    color: var(--color-link);
    margin: 0 0 var(--space-2);
  }

  h1 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-xl);
    font-family: var(--font-display);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0 0 var(--space-6);
  }

  form {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-sm);
  }

  input {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: inherit;
  }

  button {
    padding: var(--space-2) var(--space-4);
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    font-weight: 600;
    cursor: pointer;
  }

  .ok {
    margin-top: var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-success-fg);
  }
</style>
