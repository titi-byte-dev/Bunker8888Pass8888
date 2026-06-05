<script lang="ts">
  import { generatePassword } from "$lib/vault/password";
  import type { LoginItem } from "$lib/vault/types";

  interface Props {
    initial?: Partial<LoginItem>;
    submitLabel?: string;
    busy?: boolean;
    focusPassword?: boolean;
    remediationNote?: string;
    onsubmit: (payload: LoginItem) => void | Promise<void>;
  }

  let {
    initial = {},
    submitLabel = "Guardar",
    busy = false,
    focusPassword = false,
    remediationNote = "",
    onsubmit,
  }: Props = $props();

  let passwordInput = $state<HTMLInputElement | undefined>(undefined);

  $effect(() => {
    if (focusPassword && passwordInput) {
      passwordInput.focus();
    }
  });

  let title = $state("");
  let username = $state("");
  let password = $state("");
  let url = $state("");
  let notes = $state("");

  $effect.pre(() => {
    title = initial.title ?? "";
    username = initial.username ?? "";
    password = initial.password ?? "";
    url = initial.url ?? "";
    notes = initial.notes ?? "";
  });

  function handleGenerate() {
    password = generatePassword({ length: 20 });
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    onsubmit({
      kind: "login",
      title: title.trim(),
      username: username.trim(),
      password,
      url: url.trim() || undefined,
      notes: notes.trim() || undefined,
    });
  }
</script>

<form class="vault-form" onsubmit={handleSubmit}>
  {#if remediationNote}
    <p class="remediation" role="alert">{remediationNote}</p>
  {/if}
  <label>
    Título
    <input type="text" bind:value={title} required disabled={busy} />
  </label>
  <label>
    Utilizador
    <input type="text" bind:value={username} autocomplete="username" disabled={busy} />
  </label>
  <label>
    Password
    <div class="row">
      <input type="text" bind:value={password} bind:this={passwordInput} required disabled={busy} />
      <button type="button" class="secondary" onclick={handleGenerate} disabled={busy}>Gerar</button>
    </div>
  </label>
  <label>
    URL
    <input type="text" bind:value={url} inputmode="url" disabled={busy} />
  </label>
  <label>
    Notas
    <textarea bind:value={notes} rows="3" disabled={busy}></textarea>
  </label>
  <button type="submit" disabled={busy || !title.trim() || !password}>{submitLabel}</button>
</form>

<style>
  .vault-form label {
    display: block;
    margin-bottom: var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text-label);
  }

  .vault-form input,
  .vault-form textarea {
    display: block;
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: inherit;
    box-sizing: border-box;
    font-family: var(--font-ui);
    font-size: var(--text-base);
  }

  .vault-form input:focus-visible,
  .vault-form textarea:focus-visible {
    outline: 2px solid var(--color-accent);
    outline-offset: 1px;
  }

  .row {
    display: flex;
    gap: var(--space-2);
  }

  .row input {
    flex: 1;
  }

  .vault-form button[type="submit"] {
    width: 100%;
    padding: var(--space-3);
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    font-weight: 600;
    cursor: pointer;
  }

  .vault-form button.secondary {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    color: var(--color-text);
    white-space: nowrap;
  }

  .vault-form button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .remediation {
    margin: 0 0 var(--space-4);
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-warning);
    background: color-mix(in srgb, var(--color-warning) 10%, transparent);
    font-size: var(--text-sm);
    color: var(--color-text);
  }
</style>
