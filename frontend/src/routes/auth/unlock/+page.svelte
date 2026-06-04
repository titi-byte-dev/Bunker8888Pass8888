<script lang="ts">
  import { page } from "$app/state";
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { navigateAfterAuth, unlockWithPassword } from "$lib/auth/flow";
  import { resolveAuthRedirect, resolveUnlockEmail } from "$lib/auth/guard";
  import { loadSessionToken } from "$lib/session";

  let masterPassword = $state("");
  let busy = $state(false);
  let deriving = $state(false);
  let error = $state("");

  const redirectTo = $derived(resolveAuthRedirect(page.url.searchParams));
  const email = $derived(resolveUnlockEmail(page.url.searchParams));
  const hasSession = $derived(!!loadSessionToken());

  async function handleUnlock() {
    if (!email) {
      error = "Sessão inválida — inicia sessão novamente.";
      return;
    }
    busy = true;
    deriving = true;
    error = "";
    try {
      await unlockWithPassword(email, masterPassword);
      await navigateAfterAuth(redirectTo);
    } catch (e) {
      error = e instanceof Error ? e.message : "Desbloqueio falhou";
    } finally {
      busy = false;
      deriving = false;
    }
  }
</script>

<svelte:head>
  <title>Desbloquear cofre — AegisPass</title>
</svelte:head>

<AuthShell
  title="Desbloquear cofre"
  subtitle="Sessão activa — a Master Password deriva a chave de cifragem só no teu dispositivo."
>
  {#if !hasSession}
    <p class="error" role="alert">Sem sessão activa. <a href="/auth/login">Iniciar sessão</a></p>
  {:else}
    <form
      class="auth-form"
      onsubmit={(e) => {
        e.preventDefault();
        handleUnlock();
      }}
    >
      <label>
        Email
        <input type="email" value={email} disabled />
      </label>
      <label>
        Master Password
        <input
          type="password"
          bind:value={masterPassword}
          autocomplete="current-password"
          required
          disabled={busy}
        />
      </label>

      <ArgonProgress active={deriving} label="A derivar Master Key (Argon2id)…" />

      {#if error}
        <p class="error" role="alert">{error}</p>
      {/if}

      <div class="actions">
        <button type="submit" disabled={busy || !masterPassword || !email}>Desbloquear</button>
      </div>

      <p class="hint">
        Unlock ≠ login: o servidor já confia na tua identidade; aqui só descifras o cofre localmente.
      </p>
    </form>
  {/if}
</AuthShell>

<style>
  .error a {
    color: var(--color-link);
  }
</style>
