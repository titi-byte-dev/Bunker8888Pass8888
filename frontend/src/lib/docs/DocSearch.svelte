<script lang="ts">
  /**
   * Pesquisa local na documentação (DOC-010).
   * Didático: filtramos o índice em memória — adequado ao tamanho actual (~20 páginas).
   */
  import { searchDocs, type DocSearchHit } from "./search";

  interface Props {
    /** Placeholder do campo de pesquisa */
    placeholder?: string;
    /** Máximo de resultados visíveis */
    limit?: number;
  }

  let { placeholder = "Pesquisar documentação…", limit = 8 }: Props = $props();

  let query = $state("");
  let activeIndex = $state(0);
  let inputEl = $state<HTMLInputElement | undefined>(undefined);

  const hits = $derived(searchDocs(query, limit));
  const open = $derived(query.trim().length > 0 && hits.length > 0);

  $effect(() => {
    query;
    activeIndex = 0;
  });

  function onKeydown(e: KeyboardEvent) {
    if (!open) return;
    if (e.key === "ArrowDown") {
      e.preventDefault();
      activeIndex = Math.min(activeIndex + 1, hits.length - 1);
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      activeIndex = Math.max(activeIndex - 1, 0);
    } else if (e.key === "Enter" && hits[activeIndex]) {
      e.preventDefault();
      window.location.href = `/settings/docs/${hits[activeIndex]!.slug}`;
    } else if (e.key === "Escape") {
      query = "";
      inputEl?.blur();
    }
  }
</script>

<div class="doc-search">
  <label class="sr-only" for="doc-search-input">Pesquisar documentação</label>
  <input
    id="doc-search-input"
    bind:this={inputEl}
    type="search"
    {placeholder}
    bind:value={query}
    onkeydown={onKeydown}
    autocomplete="off"
    spellcheck="false"
  />

  {#if open}
    <ul class="results" role="listbox" aria-label="Resultados da pesquisa">
      {#each hits as hit, i (hit.slug)}
        <li role="option" aria-selected={i === activeIndex}>
          <a
            href="/settings/docs/{hit.slug}"
            class:active={i === activeIndex}
            onmouseenter={() => (activeIndex = i)}
          >
            <span class="hit-title">{hit.title}</span>
            <span class="hit-meta">{hit.categoryLabel}</span>
            <span class="hit-snippet">{hit.snippet}</span>
          </a>
        </li>
      {/each}
    </ul>
  {:else if query.trim() && hits.length === 0}
    <p class="empty">Sem resultados para «{query.trim()}».</p>
  {/if}
</div>

<style>
  .doc-search {
    position: relative;
    margin-bottom: var(--space-4);
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    border: 0;
  }

  input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    color: inherit;
    font-size: var(--text-sm);
    font-family: var(--font-ui);
  }

  input:focus {
    outline: 2px solid var(--color-accent);
    outline-offset: 1px;
  }

  .results {
    position: absolute;
    z-index: 20;
    top: calc(100% + var(--space-1));
    left: 0;
    right: 0;
    list-style: none;
    margin: 0;
    padding: var(--space-1);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-elevated);
    box-shadow: 0 8px 24px color-mix(in srgb, #000 20%, transparent);
    max-height: 16rem;
    overflow-y: auto;
  }

  .results a {
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    text-decoration: none;
    color: inherit;
    font-size: var(--text-sm);
  }

  .results a:hover,
  .results a.active {
    background: var(--color-accent-muted);
  }

  .hit-title {
    font-weight: 600;
  }

  .hit-meta {
    font-size: var(--text-xs);
    color: var(--color-accent);
    font-family: var(--font-mono);
  }

  .hit-snippet {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    line-height: 1.4;
  }

  .empty {
    margin: var(--space-2) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
