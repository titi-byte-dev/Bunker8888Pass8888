<script lang="ts">
  import { page } from "$app/state";
  import AuthShell from "$lib/auth/AuthShell.svelte";
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import { loginGeoHeaders } from "$lib/auth/loginGeo";
  import { loginWithPasskeyOnly, loginWithPassword, navigateAfterAuth, persistSession } from "$lib/auth/flow";
  import { resolveAuthRedirect } from "$lib/auth/guard";
  import { normalizeEmail } from "$lib/auth/http";
  import { passkeysSupported } from "$lib/passkey";
  import {
    completeSentinelStepUp,
    isSentinelStepUp,
    type SentinelStepUp,
  } from "$lib/sentinel/api";
  import { setMasterKey } from "$lib/vault/masterKeyStore";
  import { deriveMasterKeyBytes, fetchKdfParams, importAesKeyFromBytes } from "$lib/auth";

  let email = $state("");
  let masterPassword = $state("");
  let busy = $state(false);
  let error = $state("");
  let deriving = $state(false);
  let sentinelStepUp = $state<SentinelStepUp | null>(null);
  let pendingMasterKey = $state<CryptoKey | null>(null);

  const webAuthnOk = passkeysSupported();
  const redirectTo = $derived(resolveAuthRedirect(page.url.searchParams));

  async function handlePasswordLogin() {
    busy = true;
    deriving = true;
    error = "";
    sentinelStepUp = null;
    try {
      await loginWithPassword(email, masterPassword);
      await navigateAfterAuth(redirectTo);
    } catch (e) {
      if (isSentinelStepUp(e)) {
        sentinelStepUp = e.stepUp;
        try {
          const kdf = await fetchKdfParams("", normalizeEmail(email));
          const mk = await deriveMasterKeyBytes(masterPassword, kdf.salt, kdf);
          pendingMasterKey = await importAesKeyFromBytes(mk);
        } catch {
          pendingMasterKey = null;
        }
        error = "";
      } else {
        error = e instanceof Error ? e.message : "Login falhou";
      }
    } finally {
      busy = false;
      deriving = false;
    }
  }

  async function handleSentinelVerify() {
    if (!sentinelStepUp) return;
    busy = true;
    error = "";
    try {
      const geo = await loginGeoHeaders();
      const token = await completeSentinelStepUp(sentinelStepUp.challengeId, geo);
      persistSession(normalizeEmail(email), token);
      if (pendingMasterKey) {
        setMasterKey(pendingMasterKey);
        await navigateAfterAuth(redirectTo);
      } else {
        const unlockRedirect = encodeURIComponent(redirectTo);
        await navigateAfterAuth(`/auth/unlock?redirect=${unlockRedirect}`);
      }
    } catch (e) {
      error = e instanceof Error ? e.message : "Verificação falhou";
    } finally {
      busy = false;
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
  {#if sentinelStepUp}
    <section class="sentinel-banner" role="alert">
      <h2>Sentinel Mode</h2>
      <p>
        Detetámos um padrão de login suspeito ({sentinelStepUp.detail || "viagem geograficamente implausível"}).
        Confirma a tua identidade com passkey para continuar.
      </p>
      {#if error}
        <p class="error">{error}</p>
      {/if}
      {#if webAuthnOk}
        <button type="button" disabled={busy} onclick={handleSentinelVerify}>
          {busy ? "A verificar…" : "Verificar com passkey"}
        </button>
      {:else}
        <p class="hint">WebAuthn não disponível — regista uma passkey em Definições noutro dispositivo.</p>
      {/if}
      <button type="button" class="linkish" disabled={busy} onclick={() => (sentinelStepUp = null)}>
        Voltar
      </button>
    </section>
  {:else}
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

      <p class="hint">
        Primeira vez? <a href="/auth/register">Criar conta</a>. O Sentinel usa GPS opcional para detetar logins impossíveis.
      </p>

      <div class="links">
        <a href="/auth/register">Criar conta</a>
        <a href="/auth/recovery">Recuperar acesso</a>
      </div>
    </form>
  {/if}
</AuthShell>

<style>
  .sentinel-banner {
    padding: var(--space-4);
    border: 1px solid color-mix(in srgb, var(--color-danger) 35%, var(--color-border));
    border-radius: var(--radius-md);
    background: color-mix(in srgb, var(--color-danger) 8%, var(--color-bg-surface));
  }

  .sentinel-banner h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-lg);
  }

  .sentinel-banner p {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    line-height: 1.5;
  }

  .sentinel-banner button {
    margin-right: var(--space-2);
    margin-bottom: var(--space-2);
  }

  .linkish {
    border: none;
    background: transparent;
    color: var(--color-link);
    cursor: pointer;
    font-size: var(--text-sm);
    text-decoration: underline;
  }

  .error {
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
