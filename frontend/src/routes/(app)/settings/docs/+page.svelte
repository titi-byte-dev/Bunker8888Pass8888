<script lang="ts">
  import { inAppDocs, docsByCategory } from "$lib/docs/loader";
  import { LEVEL_LABELS } from "$lib/docs/types";

  const groups = docsByCategory();
  const total = inAppDocs().length;
</script>

<svelte:head>
  <title>Documentação — AegisPass</title>
</svelte:head>

<section class="docs-index">
  <h1>Documentação</h1>
  <p class="lead">
    Conhecimento agregado por níveis — começa pelo essencial e expande quando precisares.
    Projecto didático: conceitos embutidos, sem poluir o ecrã.
  </p>

  <details class="how-it-works" open>
    <summary>Como usar esta documentação</summary>
    <ul>
      <li>
        <strong>Nível 1 — Essencial:</strong> o mínimo para usar a app com confiança.
      </li>
      <li>
        <strong>Nível 2 — Intermédio:</strong> aprofunda fluxos e políticas de segurança.
      </li>
      <li>
        <strong>Nível 3 — Técnico:</strong> código, API e arquitectura (programadores).
      </li>
    </ul>
    <p class="hint">
      Cada página tem cartões <em>Conceitos-chave</em> (dropdown) e secções colapsáveis.
      Usa o filtro «Mostrar até» para controlar a profundidade.
    </p>
  </details>

  <p class="stats">{total} páginas · fonte única em <code>docs/</code></p>

  {#each groups as group (group.id)}
    <section class="category-block">
      <h2>{group.label}</h2>
      <ul class="doc-cards">
        {#each group.docs as doc (doc.slug)}
          <li>
            <a href="/settings/docs/{doc.slug}" class="doc-card">
              <span class="card-title">{doc.title}</span>
              <span class="card-summary">{doc.summary}</span>
              <span class="card-meta">
                {LEVEL_LABELS[Math.min(doc.maxLevel, 3) as 1 | 2 | 3]} ·
                nv. 1–{doc.maxLevel}
              </span>
            </a>
          </li>
        {/each}
      </ul>
    </section>
  {/each}
</section>

<style>
  h1 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .lead {
    margin: 0 0 var(--space-6);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    line-height: 1.6;
  }

  .how-it-works {
    margin-bottom: var(--space-6);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    font-size: var(--text-sm);
  }

  .how-it-works summary {
    cursor: pointer;
    font-weight: 600;
    margin-bottom: var(--space-2);
  }

  .how-it-works ul {
    margin: var(--space-3) 0;
    padding-left: var(--space-6);
    line-height: 1.6;
  }

  .hint {
    margin: 0;
    color: var(--color-text-muted);
  }

  .stats {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    margin-bottom: var(--space-8);
  }

  .category-block {
    margin-bottom: var(--space-8);
  }

  h2 {
    margin: 0 0 var(--space-3);
    font-size: var(--text-lg);
    font-family: var(--font-display);
  }

  .doc-cards {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .doc-card {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    text-decoration: none;
    color: inherit;
    transition: border-color var(--duration-fast) var(--ease-out);
  }

  .doc-card:hover {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .card-title {
    font-weight: 600;
    font-size: var(--text-sm);
  }

  .card-summary {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    line-height: 1.5;
  }

  .card-meta {
    font-size: var(--text-xs);
    color: var(--color-accent);
    font-family: var(--font-mono);
    margin-top: var(--space-1);
  }

  @media (prefers-reduced-motion: reduce) {
    .doc-card {
      transition: none;
    }
  }
</style>
