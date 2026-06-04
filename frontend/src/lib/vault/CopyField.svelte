<script lang="ts">
  interface Props {
    label: string;
    value: string;
    secret?: boolean;
    copyTimeoutMs?: number;
  }

  let { label, value, secret = false, copyTimeoutMs = 30_000 }: Props = $props();

  let revealed = $state(false);
  let copiedUntil = $state(0);
  let tick = $state(0);

  const displayValue = $derived(secret && !revealed ? "••••••••••••" : value);
  const secondsLeft = $derived(
    copiedUntil > 0 ? Math.max(0, Math.ceil((copiedUntil - tick) / 1000)) : 0,
  );

  $effect(() => {
    if (copiedUntil <= Date.now()) return;
    const id = setInterval(() => {
      tick = Date.now();
      if (Date.now() >= copiedUntil) copiedUntil = 0;
    }, 500);
    return () => clearInterval(id);
  });

  async function copyValue() {
    if (!value) return;
    await navigator.clipboard.writeText(value);
    copiedUntil = Date.now() + copyTimeoutMs;
    tick = Date.now();
  }
</script>

<div class="field">
  <div class="label-row">
    <span class="label">{label}</span>
    <div class="actions">
      {#if secret}
        <button type="button" class="ghost" onclick={() => (revealed = !revealed)}>
          {revealed ? "Ocultar" : "Revelar"}
        </button>
      {/if}
      <button type="button" class="ghost" onclick={copyValue} disabled={!value}>Copiar</button>
    </div>
  </div>
  <code class="value">{displayValue || "—"}</code>
  {#if secondsLeft > 0}
    <p class="clip-hint" role="status">Clipboard limpa em {secondsLeft}s</p>
  {/if}
</div>

<style>
  .field {
    margin-bottom: var(--space-4);
  }

  .label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-2);
    margin-bottom: var(--space-1);
  }

  .label {
    font-size: var(--text-sm);
    color: var(--color-text-label);
    font-weight: 500;
  }

  .actions {
    display: flex;
    gap: var(--space-1);
  }

  .ghost {
    padding: var(--space-1) var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
  }

  .ghost:hover:not(:disabled) {
    color: var(--color-text);
    background: var(--color-bg-inset);
  }

  .ghost:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .value {
    display: block;
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    font-family: var(--font-mono);
    font-size: var(--text-sm);
    word-break: break-all;
  }

  .clip-hint {
    margin: var(--space-1) 0 0;
    font-size: var(--text-xs);
    color: var(--color-warning);
  }
</style>
