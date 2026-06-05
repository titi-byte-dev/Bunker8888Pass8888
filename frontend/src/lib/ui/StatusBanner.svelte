<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * StatusBanner (UI-012) — mensagem inline nao-bloqueante.
   * Substitui `.status`/`.error`/`.muted` dispersos. Calmo, nunca alarmista
   * (North Star: confianca visivel). Erros usam role="alert".
   */
  type Variant = "info" | "success" | "warning" | "error";

  interface Props {
    variant?: Variant;
    children: Snippet;
  }
  let { variant = "info", children }: Props = $props();

  const icon: Record<Variant, string> = {
    info: "i",
    success: "✓",
    warning: "!",
    error: "×",
  };
</script>

<div class="banner {variant}" role={variant === "error" ? "alert" : "status"}>
  <span class="dot" aria-hidden="true">{icon[variant]}</span>
  <div class="msg">{@render children()}</div>
</div>

<style>
  .banner {
    display: flex;
    align-items: flex-start;
    gap: var(--space-2);
    padding: var(--space-3);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    font-size: var(--text-sm);
  }
  .dot {
    flex-shrink: 0;
    width: 1.25rem;
    height: 1.25rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    font-size: var(--text-xs);
    font-weight: 700;
  }
  .msg { line-height: var(--leading-body); }

  .info { background: var(--color-accent-muted); color: var(--color-text); }
  .info .dot { background: var(--color-accent); color: var(--color-accent-fg); }

  .success { background: var(--color-success-bg); color: var(--color-success-fg); }
  .success .dot { background: var(--color-success-fg); color: var(--color-success-bg); }

  .warning { background: var(--color-bg-surface); color: var(--color-text); border-color: var(--color-warning); }
  .warning .dot { background: var(--color-warning); color: #1a1205; }

  .error { background: var(--color-bg-surface); color: var(--color-text); border-color: var(--color-danger); }
  .error .dot { background: var(--color-danger); color: #fff; }
</style>
