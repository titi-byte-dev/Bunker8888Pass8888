<script lang="ts">
  import type { DocConcept } from "./types";
  import { annotateGlossaryHtml } from "./glossary";
  import { LEVEL_LABELS, type DocComplexityLevel } from "./types";

  interface Props {
    concept: DocConcept;
    /** Nível máximo visível — conceitos acima ficam ocultos */
    maxVisibleLevel?: DocComplexityLevel;
  }

  let { concept, maxVisibleLevel = 3 }: Props = $props();

  const visible = $derived(concept.level <= maxVisibleLevel);
  const levelLabel = $derived(LEVEL_LABELS[concept.level as DocComplexityLevel] ?? "");
  const bodyHtml = $derived(annotateGlossaryHtml(concept.html));
</script>

{#if visible}
  <details class="concept-card" data-level={concept.level} id="concept-{concept.id}">
    <summary class="concept-summary">
      <span class="concept-icon" aria-hidden="true">💡</span>
      <span class="concept-title">{concept.title}</span>
      {#if concept.level > 1}
        <span class="level-badge">{levelLabel}</span>
      {/if}
      <span class="chevron" aria-hidden="true">▾</span>
    </summary>
    <div class="concept-body prose">
      {@html bodyHtml}
    </div>
  </details>
{/if}

<style>
  .concept-card {
    border: 1px solid color-mix(in srgb, var(--color-accent) 25%, var(--color-border));
    border-radius: var(--radius-md);
    background: color-mix(in srgb, var(--color-accent) 6%, var(--color-bg-surface));
    margin-bottom: var(--space-3);
    overflow: hidden;
  }

  .concept-summary {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-4);
    cursor: pointer;
    list-style: none;
    font-size: var(--text-sm);
    font-weight: 500;
    user-select: none;
  }

  .concept-summary::-webkit-details-marker {
    display: none;
  }

  .concept-icon {
    flex-shrink: 0;
  }

  .concept-title {
    flex: 1;
    min-width: 0;
  }

  .level-badge {
    font-size: var(--text-xs);
    font-weight: 400;
    color: var(--color-text-muted);
    padding: 2px var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
  }

  .chevron {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    transition: transform var(--duration-fast) var(--ease-out);
  }

  .concept-card[open] .chevron {
    transform: rotate(180deg);
  }

  .concept-body {
    padding: 0 var(--space-4) var(--space-4);
    font-size: var(--text-sm);
    line-height: 1.6;
    color: var(--color-text-muted);
    border-top: 1px solid var(--color-border);
    padding-top: var(--space-3);
    margin: 0 var(--space-4) var(--space-4);
  }

  @media (prefers-reduced-motion: reduce) {
    .chevron {
      transition: none;
    }
  }
</style>
