<script lang="ts">
  import type { DocPage } from "./types";
  import type { DocComplexityLevel } from "./types";
  import { relatedDocs } from "./loader";
  import ConceptCard from "./ConceptCard.svelte";
  import DocSection from "./DocSection.svelte";
  import DocLevelFilter from "./DocLevelFilter.svelte";
  import DocProse from "./DocProse.svelte";
  import { LEVEL_LABELS } from "./types";

  interface Props {
    page: DocPage;
  }

  let { page }: Props = $props();

  const pageCap = $derived(Math.min(3, page.maxLevel) as DocComplexityLevel);
  let maxLevel = $state<DocComplexityLevel>(1);

  $effect(() => {
    maxLevel = pageCap;
  });

  const related = $derived(relatedDocs(page.related));
  const visibleConcepts = $derived(
    page.concepts.filter((c) => c.level <= maxLevel),
  );
</script>

<DocProse />

<article class="doc-page">
  <header class="doc-header">
    <p class="doc-category">{page.categoryLabel}</p>
    <h1>{page.title}</h1>
    {#if page.actor}
      <p class="doc-actor">Actor: {page.actor}</p>
    {/if}
    {#if page.summary}
      <p class="doc-summary">{page.summary}</p>
    {/if}
    <div class="doc-meta">
      <span class="meta-chip">Até {LEVEL_LABELS[maxLevel]}</span>
      {#each page.audience as aud (aud)}
        <span class="meta-chip">{aud}</span>
      {/each}
    </div>
  </header>

  <DocLevelFilter value={maxLevel} maxLevel={pageCap} onchange={(l) => (maxLevel = l)} />

  {#if visibleConcepts.length > 0}
    <section class="concepts-block" aria-labelledby="concepts-heading">
      <h2 id="concepts-heading" class="block-title">Conceitos-chave</h2>
      <p class="block-hint">Expande cada cartão para aprender sem sobrecarregar o ecrã.</p>
      {#each page.concepts as concept (concept.id)}
        <ConceptCard {concept} maxVisibleLevel={maxLevel} />
      {/each}
    </section>
  {/if}

  {#each page.sections as section, i (i)}
    <DocSection {section} maxVisibleLevel={maxLevel} />
  {/each}

  {#if related.length > 0}
    <footer class="doc-related">
      <h2 class="block-title">Continuar a aprender</h2>
      <ul>
        {#each related as doc (doc.slug)}
          <li>
            <a href="/settings/docs/{doc.slug}">
              <span class="rel-title">{doc.title}</span>
              <span class="rel-cat">{doc.categoryLabel}</span>
            </a>
          </li>
        {/each}
      </ul>
    </footer>
  {/if}
</article>

<style>
  .doc-page {
    min-width: 0;
  }

  .doc-header {
    margin-bottom: var(--space-6);
  }

  .doc-category {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-accent);
    font-weight: 600;
  }

  h1 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    line-height: 1.2;
  }

  .doc-actor {
    margin: 0 0 var(--space-2);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .doc-summary {
    margin: 0 0 var(--space-3);
    font-size: var(--text-base);
    color: var(--color-text-muted);
    line-height: 1.6;
  }

  .doc-meta {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
  }

  .meta-chip {
    font-size: var(--text-xs);
    padding: 2px var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  .concepts-block {
    margin-bottom: var(--space-8);
  }

  .block-title {
    margin: 0 0 var(--space-1);
    font-size: var(--text-lg);
    font-family: var(--font-display);
  }

  .block-hint {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .doc-related {
    margin-top: var(--space-8);
    padding-top: var(--space-6);
    border-top: 1px solid var(--color-border);
  }

  .doc-related ul {
    list-style: none;
    margin: var(--space-3) 0 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .doc-related a {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    text-decoration: none;
    color: var(--color-text);
    background: var(--color-bg-surface);
    font-size: var(--text-sm);
  }

  .doc-related a:hover {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .rel-title {
    font-weight: 500;
  }

  .rel-cat {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
