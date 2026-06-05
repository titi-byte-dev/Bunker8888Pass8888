<script lang="ts">
  import type { DocSection } from "./types";
  import { LEVEL_LABELS, type DocComplexityLevel } from "./types";

  interface Props {
    section: DocSection;
    maxVisibleLevel?: DocComplexityLevel;
    /** Secções nível 1 nunca colapsam por defeito */
    forceOpen?: boolean;
  }

  let { section, maxVisibleLevel = 3, forceOpen = false }: Props = $props();

  const visible = $derived(section.level <= maxVisibleLevel);
  const levelLabel = $derived(LEVEL_LABELS[section.level as DocComplexityLevel] ?? "");
  const useDetails = $derived(
    !forceOpen && section.level > 1 && (section.collapsed || Boolean(section.title)),
  );
</script>

{#if visible}
  {#if useDetails}
    <details class="doc-section" data-level={section.level} open={section.level === 1}>
      <summary class="section-summary">
        {#if section.title}
          <span class="section-title">{section.title}</span>
        {:else}
          <span class="section-title">{levelLabel}</span>
        {/if}
        <span class="level-badge">{levelLabel}</span>
        <span class="hint">clica para expandir</span>
        <span class="chevron" aria-hidden="true">▾</span>
      </summary>
      <div class="section-body prose">
        {@html section.html}
      </div>
    </details>
  {:else}
    <section class="doc-section flat" data-level={section.level}>
      {#if section.title}
        <h2 class="section-heading">{section.title}</h2>
      {/if}
      <div class="section-body prose">
        {@html section.html}
      </div>
    </section>
  {/if}
{/if}

<style>
  .doc-section {
    margin-bottom: var(--space-6);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    overflow: hidden;
  }

  .doc-section.flat {
    border: none;
    background: transparent;
    overflow: visible;
  }

  .section-summary {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: var(--space-2);
    padding: var(--space-4);
    cursor: pointer;
    list-style: none;
  }

  .section-summary::-webkit-details-marker {
    display: none;
  }

  .section-heading {
    margin: 0 0 var(--space-3);
    font-size: var(--text-lg);
    font-family: var(--font-display);
  }

  .section-title {
    font-size: var(--text-base);
    font-weight: 600;
    flex: 1;
    min-width: 8rem;
  }

  .level-badge {
    font-size: var(--text-xs);
    color: var(--color-accent);
    background: var(--color-accent-muted);
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
  }

  .hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .chevron {
    margin-left: auto;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    transition: transform var(--duration-fast) var(--ease-out);
  }

  .doc-section[open] .chevron {
    transform: rotate(180deg);
  }

  .section-body {
    padding: 0 var(--space-4) var(--space-4);
    font-size: var(--text-sm);
    line-height: 1.65;
  }

  .doc-section:not(.flat) .section-body {
    border-top: 1px solid var(--color-border);
    padding-top: var(--space-4);
    margin: 0 var(--space-4) var(--space-4);
  }

  @media (prefers-reduced-motion: reduce) {
    .chevron {
      transition: none;
    }
  }
</style>
