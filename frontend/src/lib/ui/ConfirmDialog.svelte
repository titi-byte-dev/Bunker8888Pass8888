<script lang="ts">
  import { closeConfirm, confirmStore } from "./confirm";
  import Button from "./Button.svelte";

  const state = $derived($confirmStore);
  const opts = $derived(state.options);

  function onKey(e: KeyboardEvent) {
    if (!state.open) return;
    if (e.key === "Escape") closeConfirm(false);
  }

  function onBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) closeConfirm(false);
  }
</script>

<svelte:window onkeydown={onKey} />

{#if state.open && opts}
  <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
  <div class="backdrop" role="presentation" onclick={onBackdrop}>
    <div
      class="dialog"
      role="alertdialog"
      aria-modal="true"
      aria-labelledby="confirm-title"
      aria-describedby="confirm-msg"
    >
      <h2 id="confirm-title">{opts.title}</h2>
      <p id="confirm-msg">{opts.message}</p>
      <div class="actions">
        <Button variant="secondary" onclick={() => closeConfirm(false)}>
          {opts.cancelLabel ?? "Cancelar"}
        </Button>
        <Button
          variant={opts.variant === "danger" ? "danger" : "primary"}
          onclick={() => closeConfirm(true)}
        >
          {opts.confirmLabel ?? "Confirmar"}
        </Button>
      </div>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    z-index: 300;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--space-4);
    background: rgba(0, 0, 0, 0.45);
    backdrop-filter: blur(4px);
  }

  .dialog {
    width: 100%;
    max-width: 24rem;
    padding: var(--space-6);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    background: var(--color-bg-elevated);
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.35);
    animation: dialog-in var(--duration-normal) var(--ease-out);
  }

  h2 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-lg);
    line-height: var(--leading-tight);
  }

  p {
    margin: 0 0 var(--space-6);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    line-height: var(--leading-body);
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--space-2);
    flex-wrap: wrap;
  }

  @keyframes dialog-in {
    from {
      opacity: 0;
      transform: scale(0.96);
    }
    to {
      opacity: 1;
      transform: scale(1);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .dialog {
      animation: none;
    }
  }
</style>
