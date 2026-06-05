<script lang="ts">
  /**
   * Skeleton (UI-017) — placeholder de carregamento.
   * Substitui texto «A carregar…» por forma que sugere o layout final.
   */
  type Variant = "text" | "block" | "circle" | "row" | "table";

  interface Props {
    variant?: Variant;
    /** Largura CSS (ex. 60%, 8rem). Omissão = 100% */
    width?: string;
    /** Altura CSS para block/row */
    height?: string;
    /** Repetições para table (linhas) */
    rows?: number;
    /** Colunas para table */
    cols?: number;
  }

  let {
    variant = "text",
    width,
    height,
    rows = 4,
    cols = 3,
  }: Props = $props();
</script>

{#if variant === "table"}
  <div class="sk-table" aria-hidden="true" role="presentation">
    <div class="sk-head">
      {#each Array(cols) as _, i (i)}
        <span class="sk cell"></span>
      {/each}
    </div>
    {#each Array(rows) as _, r (r)}
      <div class="sk-row">
        {#each Array(cols) as _, c (c)}
          <span class="sk cell" style:width={c === cols - 1 ? "40%" : undefined}></span>
        {/each}
      </div>
    {/each}
  </div>
{:else if variant === "row"}
  <div class="sk-row-layout" style:height={height} aria-hidden="true">
    <span class="sk title"></span>
    <span class="sk meta"></span>
  </div>
{:else}
  <span
    class="sk {variant}"
    style:width
    style:height
    aria-hidden="true"
    role="presentation"
  ></span>
{/if}

<style>
  .sk {
    display: block;
    border-radius: var(--radius-sm);
    background: linear-gradient(
      90deg,
      var(--color-bg-surface) 0%,
      var(--color-bg-elevated) 50%,
      var(--color-bg-surface) 100%
    );
    background-size: 200% 100%;
    animation: shimmer 1.4s ease-in-out infinite;
  }

  .text {
    height: 0.875rem;
    width: 100%;
    max-width: 12rem;
  }

  .block {
    height: 4rem;
    width: 100%;
  }

  .circle {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
  }

  .sk-row-layout {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
  }

  .title {
    height: 0.9rem;
    width: 40%;
    max-width: 10rem;
  }

  .meta {
    height: 0.75rem;
    width: 25%;
    max-width: 6rem;
  }

  .sk-table {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .sk-head,
  .sk-row {
    display: flex;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }

  .sk-head .cell {
    height: 0.65rem;
    flex: 1;
    opacity: 0.7;
  }

  .sk-row .cell {
    height: 0.8rem;
    flex: 1;
  }

  .sk-row:last-child {
    border-bottom: none;
  }

  @keyframes shimmer {
    0% { background-position: 100% 0; }
    100% { background-position: -100% 0; }
  }

  @media (prefers-reduced-motion: reduce) {
    .sk {
      animation: none;
      background: var(--color-bg-surface);
    }
  }
</style>
