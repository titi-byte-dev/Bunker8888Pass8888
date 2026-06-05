<script lang="ts">
  import { clearAdminKey, hasAdminKey, loadAdminKey, saveAdminKey } from "$lib/admin/adminKey";
  import { verifyAdminKey } from "$lib/admin/api";

  interface Props {
    onUnlocked?: () => void;
  }

  let { onUnlocked }: Props = $props();

  let keyInput = $state(loadAdminKey() ?? "");
  let busy = $state(false);
  let error = $state("");

  async function handleUnlock() {
    busy = true;
    error = "";
    try {
      const ok = await verifyAdminKey(keyInput);
      if (!ok) {
        error = "Chave inválida ou admin desactivado no servidor.";
        return;
      }
      saveAdminKey(keyInput);
      onUnlocked?.();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha na verificação";
    } finally {
      busy = false;
    }
  }

  function handleLock() {
    clearAdminKey();
    keyInput = "";
    onUnlocked?.();
  }
</script>

<section class="gate">
  <h2>Acesso administrativo</h2>
  <p class="hint">
    Introduz a chave definida em <code>AEGIS_ADMIN_KEY</code>. Guardada só nesta sessão do browser
    (sessionStorage).
  </p>

  {#if hasAdminKey()}
    <p class="status">Chave admin activa nesta sessão.</p>
    <button type="button" class="secondary" onclick={handleLock}>Terminar sessão admin</button>
  {:else}
    <form
      class="form"
      onsubmit={(e) => {
        e.preventDefault();
        handleUnlock();
      }}
    >
      <label>
        Chave admin
        <input type="password" bind:value={keyInput} autocomplete="off" required />
      </label>
      {#if error}
        <p class="error" role="alert">{error}</p>
      {/if}
      <button type="submit" disabled={busy || !keyInput.trim()}>
        {busy ? "A verificar…" : "Desbloquear admin"}
      </button>
    </form>
  {/if}
</section>

<style>
  .gate {
    max-width: 28rem;
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    margin-bottom: var(--space-6);
  }

  h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-lg);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0 0 var(--space-4);
    line-height: 1.5;
  }

  .form label {
    display: block;
    font-size: var(--text-sm);
    margin-bottom: var(--space-3);
  }

  input {
    display: block;
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }

  button {
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  button.secondary {
    background: var(--color-bg-base);
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .status {
    font-size: var(--text-sm);
    color: var(--color-success-fg);
    margin: 0 0 var(--space-3);
  }

  .error {
    font-size: var(--text-sm);
    color: var(--color-danger);
    margin: 0 0 var(--space-3);
  }

  code {
    font-family: var(--font-mono);
    font-size: 0.9em;
  }
</style>
