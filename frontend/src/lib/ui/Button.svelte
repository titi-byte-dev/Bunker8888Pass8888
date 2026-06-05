<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * Button (UI-012) — substitui .btn/.primary/.secondary/.danger inline.
   * Renderiza <a> quando `href` esta presente (mantem aparencia de botao).
   */
  type Variant = "primary" | "secondary" | "ghost" | "danger";
  type Size = "sm" | "md";

  interface Props {
    variant?: Variant;
    size?: Size;
    loading?: boolean;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    href?: string;
    onclick?: (e: MouseEvent) => void;
    children: Snippet;
  }

  let {
    variant = "primary",
    size = "md",
    loading = false,
    disabled = false,
    type = "button",
    href,
    onclick,
    children,
  }: Props = $props();

  const isDisabled = $derived(disabled || loading);
</script>

{#if href && !isDisabled}
  <a class="btn {variant} {size}" {href} {onclick} data-loading={loading}>
    {@render children()}
  </a>
{:else}
  <button
    class="btn {variant} {size}"
    {type}
    disabled={isDisabled}
    aria-busy={loading}
    {onclick}
  >
    {#if loading}<span class="spinner" aria-hidden="true"></span>{/if}
    {@render children()}
  </button>
{/if}

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    font-family: var(--font-ui);
    font-weight: 600;
    text-decoration: none;
    cursor: pointer;
    white-space: nowrap;
    transition:
      background-color var(--duration-fast) var(--ease-out),
      border-color var(--duration-fast) var(--ease-out),
      opacity var(--duration-fast) var(--ease-out);
  }
  .sm { padding: var(--space-1) var(--space-3); font-size: var(--text-xs); }
  .md { padding: var(--space-2) var(--space-4); font-size: var(--text-sm); }

  .primary { background: var(--color-accent); color: var(--color-accent-fg); }
  .primary:hover:not(:disabled) { filter: brightness(1.08); }

  .secondary {
    background: var(--color-bg-surface);
    color: var(--color-text);
    border-color: var(--color-border);
  }
  .secondary:hover:not(:disabled) { border-color: var(--color-border-strong); }

  .ghost { background: transparent; color: var(--color-text-muted); }
  .ghost:hover:not(:disabled) { background: var(--color-bg-surface); color: var(--color-text); }

  .danger { background: var(--color-danger); color: #fff; }
  .danger:hover:not(:disabled) { filter: brightness(1.08); }

  .btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .spinner {
    width: 0.85em;
    height: 0.85em;
    border: 2px solid currentColor;
    border-right-color: transparent;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }
  @media (prefers-reduced-motion: reduce) {
    .btn { transition: none; }
    .spinner { animation-duration: 1.5s; }
  }
</style>
