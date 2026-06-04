<script lang="ts">
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { recoverMasterKeyFromEmail } from "$lib/vault/recovery";
  import { setMasterKey, purgeMasterKey } from "$lib/vault/masterKeyStore";

  let email = $state("");
  let recoveryCode = $state("");
  let busy = $state(false);
  let deriving = $state(false);
  let error = $state("");
  let success = $state(false);

  async function handleRecover() {
    busy = true;
    deriving = true;
    error = "";
    success = false;
    try {
      // Valida o backup cifrado — Master Key fica volátil até login completo.
      const masterKey = await recoverMasterKeyFromEmail("", email, recoveryCode);
      setMasterKey(masterKey);
      purgeMasterKey();
      success = true;
    } catch (e) {
      error = e instanceof Error ? e.message : "Recuperação falhou";
    } finally {
      busy = false;
      deriving = false;
    }
  }
</script>

<svelte:head>
  <title>Recuperar acesso — AegisPass</title>
</svelte:head>

<AuthShell
  title="Recuperar acesso"
  subtitle="Usa a chave de recuperação guardada offline. Não substitui login — valida o backup cifrado."
>
  {#if success}
    <div class="success" role="status">
      <p>Backup de recuperação válido. Inicia sessão para continuar.</p>
      <a class="cta" href="/auth/login">Ir para login</a>
    </div>
  {:else}
    <form
      class="auth-form"
      onsubmit={(e) => {
        e.preventDefault();
        handleRecover();
      }}
    >
      <label>
        Email
        <input type="email" bind:value={email} autocomplete="username" required disabled={busy} />
      </label>
      <label>
        Código de recuperação
        <input
          type="text"
          bind:value={recoveryCode}
          placeholder="XXXXX-XXXXX-XXXXX-XXXXX"
          required
          disabled={busy}
          autocomplete="off"
        />
      </label>

      <ArgonProgress active={deriving} label="A validar backup cifrado…" />

      {#if error}
        <p class="error" role="alert">{error}</p>
      {/if}

      <div class="actions">
        <button type="submit" disabled={busy || !email || !recoveryCode}>Validar backup</button>
      </div>

      <div class="links">
        <a href="/auth/login">Voltar ao login</a>
      </div>
    </form>
  {/if}
</AuthShell>

<style>
  .success {
    text-align: center;
  }

  .success p {
    margin: 0 0 var(--space-4);
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
  }

  .cta {
    display: inline-block;
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    text-decoration: none;
    font-size: var(--text-sm);
    font-weight: 600;
  }
</style>
