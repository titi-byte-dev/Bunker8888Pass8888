<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import { decryptNoteContent, parseNoteFragment, type NoteFragment } from "$lib/share/burnNote";
  import { consumeBurnNote, BurnNoteGoneError } from "$lib/share/burnNoteApi";

  type Phase = "prompt" | "revealing" | "revealed" | "gone" | "invalid" | "error";

  let phase = $state<Phase>("prompt");
  let errorMsg = $state("");
  let note = $state("");
  let passphrase = $state("");
  let copied = $state(false);

  let frag = $state<NoteFragment | null>(null);
  const id = $derived(page.params.id ?? "");

  onMount(() => {
    frag = parseNoteFragment(window.location.hash);
    if (!frag) phase = "invalid";
  });

  async function reveal() {
    if (!frag) {
      phase = "invalid";
      return;
    }
    if (frag.requiresPassphrase && !passphrase) {
      errorMsg = "Esta nota precisa de passphrase.";
      return;
    }
    phase = "revealing";
    errorMsg = "";
    try {
      const ciphertext = await consumeBurnNote(id);
      note = await decryptNoteContent(
        frag.key,
        ciphertext,
        frag.requiresPassphrase ? passphrase : null,
        frag.salt,
      );
      passphrase = "";
      phase = "revealed";
    } catch (e) {
      if (e instanceof BurnNoteGoneError) {
        phase = "gone";
      } else {
        // A nota foi consumida (ardeu) mas nao decifrou — passphrase errada.
        errorMsg = frag.requiresPassphrase
          ? "Passphrase errada. A nota foi destruida ao abrir e nao pode ser recuperada."
          : e instanceof Error
            ? e.message
            : "Falha ao abrir a nota.";
        phase = "error";
      }
    }
  }

  async function burnWithoutReading() {
    // Sem o burn_token (que só o autor tem), a unica forma de "destruir" no lado
    // do leitor é consumir a nota. Lemos e descartamos para garantir que arde.
    phase = "revealing";
    try {
      await consumeBurnNote(id);
    } catch {
      // ja nao existia — tambem serve.
    }
    note = "";
    phase = "gone";
  }

  async function copy() {
    try {
      await navigator.clipboard.writeText(note);
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      copied = false;
    }
  }
</script>

<svelte:head>
  <title>Nota auto-destrutiva — AegisPass</title>
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<main class="wrap">
  <div class="card">
    <p class="brand">AegisPass · Nota auto-destrutiva</p>

    {#if phase === "invalid"}
      <h1>Link inválido</h1>
      <p class="muted">
        Falta a chave de cifra no link (a parte depois do #). Confirma que copiaste
        o link completo.
      </p>
    {:else if phase === "gone"}
      <h1>Já não está disponível</h1>
      <p class="muted">
        Esta nota não existe, expirou ou já foi lida. As notas são de leitura única
        e não deixam rasto.
      </p>
    {:else if phase === "revealed"}
      <h1>Nota</h1>
      <p class="muted">Revelada uma vez. Já foi apagada do servidor.</p>
      <pre class="secret">{note}</pre>
      <button type="button" class="btn primary" onclick={copy}>
        {copied ? "Copiado!" : "Copiar nota"}
      </button>
    {:else if phase === "error"}
      <h1>Não foi possível abrir</h1>
      <p class="muted">{errorMsg}</p>
    {:else}
      <h1>Tens uma nota à tua espera</h1>
      <p class="muted">
        Ao revelar, a nota é mostrada <strong>uma única vez</strong> e destruída no
        servidor. A cifra é aberta no teu dispositivo.
        {#if frag?.requiresPassphrase}
          Esta nota está protegida por uma <strong>passphrase</strong> — pede-a a
          quem ta enviou.
        {/if}
      </p>
      {#if frag?.requiresPassphrase}
        <input
          class="pass"
          type="text"
          bind:value={passphrase}
          placeholder="Passphrase"
          autocomplete="off"
          spellcheck="false"
          disabled={phase === "revealing"}
        />
      {/if}
      {#if errorMsg}<p class="inline-error" role="alert">{errorMsg}</p>{/if}
      <div class="actions">
        <button
          type="button"
          class="btn primary"
          onclick={reveal}
          disabled={phase === "revealing" || (frag?.requiresPassphrase && !passphrase)}
        >
          {phase === "revealing" ? "A abrir…" : "Revelar nota"}
        </button>
        <button
          type="button"
          class="btn ghost"
          onclick={burnWithoutReading}
          disabled={phase === "revealing"}
        >
          Destruir sem ler
        </button>
      </div>
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
  .pass {
    width: 100%;
    margin: 0 0 var(--space-4);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }
  .pass:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  .actions {
    display: flex;
    gap: var(--space-2);
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
  .btn.ghost {
    background: none;
    color: var(--color-text-muted);
  }
  .btn.ghost:hover:not(:disabled) {
    color: var(--color-text);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
</style>
