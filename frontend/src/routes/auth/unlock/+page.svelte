<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { navigateAfterAuth, unlockWithPassword } from "$lib/auth/flow";
  import { resolveAuthRedirect, resolveUnlockEmail } from "$lib/auth/guard";
  import { fetchSessionProfile } from "$lib/auth/session-api";
  import { loadSessionToken, saveUserEmail } from "$lib/session";

  let masterPassword = $state("");
  let email = $state("");
  /** true quando o email veio de storage/servidor — não editável. */
  let emailLocked = $state(false);
  let hydrating = $state(true);
  let busy = $state(false);
  let deriving = $state(false);
  let error = $state("");

  const redirectTo = $derived(resolveAuthRedirect(page.url.searchParams));
  const hasSession = $derived(!!loadSessionToken());

  onMount(async () => {
    const cached = resolveUnlockEmail(page.url.searchParams);
    if (cached) {
      email = cached;
      emailLocked = true;
      hydrating = false;
      return;
    }
    if (!loadSessionToken()) {
      hydrating = false;
      return;
    }
    try {
      const profile = await fetchSessionProfile();
      email = profile.email;
      saveUserEmail(profile.email);
      emailLocked = true;
    } catch {
      // Token sem email em storage — utilizador confirma manualmente.
      emailLocked = false;
    } finally {
      hydrating = false;
    }
  });

  async function handleUnlock() {
    const normalized = email.trim().toLowerCase();
    if (!normalized) {
      error = "Indica o email da conta ou inicia sessão novamente.";
      return;
    }
    busy = true;
    deriving = true;
    error = "";
    try {
      await unlockWithPassword(normalized, masterPassword);
      saveUserEmail(normalized);
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
        <input
          type="email"
          bind:value={email}
          autocomplete="username"
          required
          disabled={emailLocked || busy || hydrating}
          placeholder={hydrating ? "A carregar…" : "email@empresa.pt"}
        />
      </label>
      <label>
        Master Password
        <input
          type="password"
          bind:value={masterPassword}
          autocomplete="current-password"
          required
          disabled={busy || hydrating}
        />
      </label>

      <ArgonProgress active={deriving} label="A derivar Master Key (Argon2id)…" />

      {#if error}
        <p class="error" role="alert">{error}</p>
      {/if}

      <div class="actions">
        <button type="submit" disabled={busy || hydrating || !masterPassword || !email.trim()}>
          Desbloquear
        </button>
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
