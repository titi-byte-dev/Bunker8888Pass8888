<script lang="ts">
  import { page } from "$app/state";
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { loginWithPasskeyOnly, loginWithPassword, navigateAfterAuth } from "$lib/auth/flow";
  import { resolveAuthRedirect } from "$lib/auth/guard";
  import { passkeysSupported } from "$lib/passkey";

  let email = $state("");
  let masterPassword = $state("");
  let busy = $state(false);
  let error = $state("");
  let deriving = $state(false);

  const webAuthnOk = passkeysSupported();
  const redirectTo = $derived(resolveAuthRedirect(page.url.searchParams));

  async function handlePasswordLogin() {
    busy = true;
    deriving = true;
    error = "";
    try {
      await loginWithPassword(email, masterPassword);
      await navigateAfterAuth(redirectTo);
    } catch (e) {
      error = e instanceof Error ? e.message : "Login falhou";
    } finally {
      busy = false;
      deriving = false;
    }
  }

  async function handlePasskeyLogin() {
    if (!email.trim()) {
      error = "Indica o email antes de usar passkey.";
      return;
    }
    busy = true;
    error = "";
    try {
      await loginWithPasskeyOnly(email);
      const unlockRedirect = encodeURIComponent(redirectTo);
      await navigateAfterAuth(`/auth/unlock?redirect=${unlockRedirect}`);
    } catch (e) {
      error = e instanceof Error ? e.message : "Passkey falhou";
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Iniciar sessão — AegisPass</title>
</svelte:head>

<AuthShell
  title="Iniciar sessão"
  subtitle="A Master Password desbloqueia o cofre localmente — o servidor nunca a vê."
>
  <form
    class="auth-form"
    onsubmit={(e) => {
      e.preventDefault();
      handlePasswordLogin();
    }}
  >
    <label>
      Email
      <input type="email" bind:value={email} autocomplete="username" required disabled={busy} />
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

    <ArgonProgress active={deriving} />

    {#if error}
      <p class="error" role="alert">{error}</p>
    {/if}

    <div class="actions">
      <button type="submit" disabled={busy || !email || !masterPassword}>Entrar</button>
      {#if webAuthnOk}
        <button type="button" class="secondary" disabled={busy || !email} onclick={handlePasskeyLogin}>
          Entrar com passkey
        </button>
      {/if}
    </div>

    {#if webAuthnOk}
      <p class="hint">
        Passkey autentica o servidor; depois precisas de desbloquear o cofre com a Master Password (Zero-Knowledge).
      </p>
    {/if}

    <div class="links">
      <a href="/auth/register">Criar conta</a>
      <a href="/auth/recovery">Recuperar acesso</a>
    </div>
  </form>
</AuthShell>
