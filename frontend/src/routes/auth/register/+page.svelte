<script lang="ts">
  import { page } from "$app/state";
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { navigateAfterAuth, registerAndLogin } from "$lib/auth/flow";
  import { resolveAuthRedirect } from "$lib/auth/guard";

  let email = $state("");
  let masterPassword = $state("");
  let confirmPassword = $state("");
  let busy = $state(false);
  let deriving = $state(false);
  let error = $state("");

  const redirectTo = $derived(resolveAuthRedirect(page.url.searchParams));

  async function handleRegister() {
    if (masterPassword !== confirmPassword) {
      error = "As passwords não coincidem.";
      return;
    }
    if (masterPassword.length < 12) {
      error = "A Master Password deve ter pelo menos 12 caracteres.";
      return;
    }

    busy = true;
    deriving = true;
    error = "";
    try {
      await registerAndLogin(email, masterPassword);
      await navigateAfterAuth(redirectTo);
    } catch (e) {
      error = e instanceof Error ? e.message : "Registo falhou";
    } finally {
      busy = false;
      deriving = false;
    }
  }
</script>

<svelte:head>
  <title>Criar conta — AegisPass</title>
</svelte:head>

<AuthShell
  title="Criar conta"
  subtitle="A tua Master Password gera chaves no browser. Guarda-a offline — não há reset pelo servidor."
>
  <form
    class="auth-form"
    onsubmit={(e) => {
      e.preventDefault();
      handleRegister();
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
        autocomplete="new-password"
        required
        disabled={busy}
        minlength="12"
      />
    </label>
    <label>
      Confirmar Master Password
      <input
        type="password"
        bind:value={confirmPassword}
        autocomplete="new-password"
        required
        disabled={busy}
        minlength="12"
      />
    </label>

    <ArgonProgress active={deriving} label="A derivar chaves (Argon2id) — registo…" />

    {#if error}
      <p class="error" role="alert">{error}</p>
    {/if}

    <div class="actions">
      <button type="submit" disabled={busy || !email || !masterPassword || !confirmPassword}>
        Registar
      </button>
    </div>

    <div class="links">
      <a href="/auth/login">Já tens conta? Entrar</a>
    </div>
  </form>
</AuthShell>
