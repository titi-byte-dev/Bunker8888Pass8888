<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * ListRow (UI-015) — linha de lista densa (cofre, admin).
   * Pode ser link (`href`) ou botão (`onclick`).
   */
  interface Props {
    title: string;
    meta?: string;
    href?: string;
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    trailing?: Snippet;
  }

  let { title, meta, href, disabled = false, onclick, trailing }: Props = $props();
</script>

{#if href && !disabled}
  <a class="row" {href} {onclick}>
    <span class="main">
      <span class="title">{title}</span>
      {#if meta}<span class="meta">{meta}</span>{/if}
    </span>
    {#if trailing}<span class="trail">{@render trailing()}</span>{/if}
    <span class="chevron" aria-hidden="true">›</span>
  </a>
{:else}
  <button type="button" class="row" {disabled} {onclick}>
    <span class="main">
      <span class="title">{title}</span>
      {#if meta}<span class="meta">{meta}</span>{/if}
    </span>
    {#if trailing}<span class="trail">{@render trailing()}</span>{/if}
  </button>
{/if}

<style>
  .row {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    width: 100%;
    padding: var(--nav-item-py) var(--space-4);
    min-height: var(--nav-item-min-height);
    border: none;
    border-bottom: 1px solid var(--color-border);
    background: transparent;
    text-decoration: none;
    color: inherit;
    font: inherit;
    text-align: left;
    cursor: pointer;
    box-sizing: border-box;
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .row:last-child {
    border-bottom: none;
  }

  .row:hover:not(:disabled) {
    background: var(--color-bg-surface);
  }

  .row:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .main {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .title {
    font-weight: 600;
    font-size: var(--text-sm);
    line-height: var(--nav-item-leading);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .meta {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .trail {
    flex-shrink: 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .chevron {
    flex-shrink: 0;
    color: var(--color-text-muted);
    font-size: var(--text-lg);
    line-height: 1;
  }

  @media (prefers-reduced-motion: reduce) {
    .row {
      transition: none;
    }
  }
</style>
