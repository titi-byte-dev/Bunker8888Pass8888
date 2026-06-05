<script lang="ts">
  import { dismissToast, toastStore, type ToastVariant } from "./toast";

  const labels: Record<ToastVariant, string> = {
    info: "Informação",
    success: "Sucesso",
    warning: "Atenção",
    error: "Erro",
  };
</script>

<div class="toast-host" aria-live="polite" aria-relevant="additions">
  {#each $toastStore as item (item.id)}
    <div class="toast {item.variant}" role="status">
      <span class="sr-only">{labels[item.variant]}:</span>
      <p class="msg">{item.message}</p>
      <button
        type="button"
        class="close"
        aria-label="Fechar notificação"
        onclick={() => dismissToast(item.id)}
      >
        ×
      </button>
    </div>
  {/each}
</div>

<style>
  .toast-host {
    position: fixed;
    bottom: calc(var(--shell-tab-height, 3.5rem) + var(--space-4));
    right: var(--space-4);
    z-index: 200;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    max-width: min(22rem, calc(100vw - var(--space-8)));
    pointer-events: none;
  }

  @media (min-width: 768px) {
    .toast-host {
      bottom: var(--space-4);
    }
  }

  .toast {
    pointer-events: auto;
    display: flex;
    align-items: flex-start;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-elevated);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
    animation: toast-in var(--duration-fast) var(--ease-out);
  }

  .info { border-left: 3px solid var(--color-accent); }
  .success { border-left: 3px solid var(--color-success-fg); }
  .warning { border-left: 3px solid var(--color-warning); }
  .error { border-left: 3px solid var(--color-danger); }

  .msg {
    margin: 0;
    flex: 1;
    font-size: var(--text-sm);
    line-height: var(--leading-body);
    color: var(--color-text);
  }

  .close {
    flex-shrink: 0;
    border: none;
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-lg);
    line-height: 1;
    cursor: pointer;
    padding: 0;
    margin: -2px 0 0;
  }

  .close:hover {
    color: var(--color-text);
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }

  @keyframes toast-in {
    from {
      opacity: 0;
      transform: translateY(0.5rem);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .toast {
      animation: none;
    }
  }
</style>
