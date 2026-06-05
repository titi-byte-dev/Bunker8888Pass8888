<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { decryptSecret, keyFromFragment } from "$lib/share/secretLink";
  import { consumeSecretLink, SecretLinkGoneError } from "$lib/share/secretLinkApi";
  import type { Bytes } from "$lib/crypto";

  type State = "idle" | "revealing" | "revealed" | "gone" | "invalid" | "error";

  let phase = $state<State>("idle");
  let errorMsg = $state("");
  let secret = $state("");
  let copied = $state(false);

  let key: Bytes | null = null;
  const id = $derived(page.params.id ?? "");

  onMount(() => {
    key = keyFromFragment(window.location.hash);
    if (!key) phase = "invalid";
  });

  async function reveal() {
    if (!key) {
      phase = "invalid";
      return;
    }
    phase = "revealing";
    errorMsg = "";
    try {
      const ciphertext = await consumeSecretLink(id);
      secret = await decryptSecret(key, ciphertext);
      phase = "revealed";
    } catch (e) {
      if (e instanceof SecretLinkGoneError) {
        phase = "gone";
      } else {
        errorMsg = e instanceof Error ? e.message : "Falha ao abrir o segredo";
        phase = "error";
      }
    }
  }

  async function copy() {
    try {
      await navigator.clipboard.writeText(secret);
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      copied = false;
    }
  }
</script>

<svelte:head>
  <title>Segredo partilhado — AegisPass</title>
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<main class="wrap">
  <div class="card">
    <p class="brand">AegisPass · Secret Link</p>

    {#if phase === "invalid"}
      <h1>Link inválido</h1>
      <p class="muted">
        Falta a chave de cifra no link (a parte depois do #). Confirma que copiaste
        o link completo.
      </p>
    {:else if phase === "gone"}
      <h1>Já não está disponível</h1>
      <p class="muted">
        Este link não existe, expirou ou já foi utilizado. Por segurança, os
        segredos são de uso único / temporários e não deixam rasto.
      </p>
    {:else if phase === "revealed"}
      <h1>Segredo</h1>
      <p class="muted">Revelado uma vez. Já foi apagado do servidor.</p>
      <pre class="secret">{secret}</pre>
      <button type="button" class="btn primary" onclick={copy}>
        {copied ? "Copiado!" : "Copiar segredo"}
      </button>
    {:else if phase === "error"}
      <h1>Erro</h1>
      <p class="muted">{errorMsg}</p>
      <button type="button" class="btn secondary" onclick={reveal}>Tentar de novo</button>
    {:else}
      <h1>Tens um segredo à tua espera</h1>
      <p class="muted">
        Alguém partilhou um segredo cifrado contigo. Ao revelar, ele é mostrado
        <strong>uma única vez</strong> e apagado do servidor. A cifra é aberta no
        teu dispositivo.
      </p>
      <button type="button" class="btn primary" onclick={reveal} disabled={phase === "revealing"}>
        {phase === "revealing" ? "A abrir…" : "Revelar segredo"}
      </button>
    {/if}
  </div>
</main>

<style>
  .wrap {
    min-height: 100vh;
    display: grid;
    place-items: center;
    padding: var(--space-6);
    background: var(--color-bg);
  }
  .card {
    width: 100%;
    max-width: 30rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-8) var(--space-6);
  }
  .brand {
    margin: 0 0 var(--space-4);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  h1 {
    margin: 0 0 var(--space-3);
    font-family: var(--font-display);
    font-size: var(--text-xl);
    line-height: var(--leading-tight);
  }
  .muted {
    margin: 0 0 var(--space-5);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    line-height: var(--leading-body);
  }
  .secret {
    margin: 0 0 var(--space-4);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    white-space: pre-wrap;
    word-break: break-word;
  }
  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
  }
  .btn.primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
    border-color: transparent;
  }
  .btn.primary:hover:not(:disabled) {
    filter: brightness(1.08);
  }
  .btn.secondary {
    background: var(--color-bg-elevated);
    color: var(--color-text);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
</style>
