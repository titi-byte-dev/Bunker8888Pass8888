<script lang="ts">
  import { LEVEL_LABELS, type DocComplexityLevel } from "./types";

  interface Props {
    value: DocComplexityLevel;
    maxLevel?: DocComplexityLevel;
    onchange: (level: DocComplexityLevel) => void;
  }

  let { value, maxLevel = 3, onchange }: Props = $props();

  const levels = $derived(
    ([1, 2, 3] as DocComplexityLevel[]).filter((l) => l <= maxLevel),
  );
</script>

<div class="level-filter" role="group" aria-label="Nível de complexidade">
  <span class="filter-label">Mostrar até:</span>
  {#each levels as level (level)}
    <button
      type="button"
      class="level-btn"
      class:active={value === level}
      aria-pressed={value === level}
      onclick={() => onchange(level)}
    >
      <span class="level-num">{level}</span>
      {LEVEL_LABELS[level]}
    </button>
  {/each}
</div>

<style>
  .level-filter {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-2);
    margin-bottom: var(--space-6);
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
  }

  .filter-label {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin-right: var(--space-1);
  }

  .level-btn {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    cursor: pointer;
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .level-btn:hover {
    background: var(--color-bg-base);
    color: var(--color-text);
  }

  .level-btn.active {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .level-num {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 1.25rem;
    height: 1.25rem;
    border-radius: 50%;
    background: var(--color-bg-base);
    font-size: var(--text-xs);
    font-weight: 600;
  }

  .level-btn.active .level-num {
    background: var(--color-accent);
    color: var(--color-bg-base);
  }

  @media (prefers-reduced-motion: reduce) {
    .level-btn {
      transition: none;
    }
  }
</style>
