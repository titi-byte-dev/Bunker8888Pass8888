<script lang="ts">
  import { docsByCategory } from "./loader";
  import { LEVEL_LABELS, type DocComplexityLevel } from "./types";

  interface Props {
    pathname: string;
  }

  let { pathname }: Props = $props();

  const groups = docsByCategory();

  function isActive(slug: string): boolean {
    return pathname === `/settings/docs/${slug}`;
  }

  function shortTitle(title: string): string {
    return title.replace(/^Journey:\s*/i, "").replace(/^Jornada:\s*/i, "");
  }

  /** Indentação visual por nível de complexidade (L1/L2/L3). */
  function levelIndent(level: number): string {
    const base = 0.375;
    return `${base + (Math.max(1, level) - 1) * 0.625}rem`;
  }

  function levelLabel(level: number): string {
    return LEVEL_LABELS[Math.min(3, level) as DocComplexityLevel] ?? `L${level}`;
  }
</script>

<nav class="doc-nav" aria-label="Índice da documentação">
  <a href="/settings/docs" class="nav-home" class:active={pathname === "/settings/docs"}>
    Índice geral
  </a>

  {#each groups as group (group.id)}
    <section class="nav-group" aria-labelledby="nav-group-{group.id}">
      <h3 id="nav-group-{group.id}" class="nav-group-title">{group.label}</h3>
      <ul class="nav-list">
        {#each group.docs as doc (doc.slug)}
          <li class="nav-item" style:--nav-indent={levelIndent(doc.level)}>
            <a
              href="/settings/docs/{doc.slug}"
              class:active={isActive(doc.slug)}
              class:level-2={doc.level === 2}
              class:level-3={doc.level >= 3}
              title="{doc.title} · {levelLabel(doc.level)}"
            >
              <span class="doc-tree" aria-hidden="true"></span>
              <span class="doc-title">{shortTitle(doc.title)}</span>
              <span class="doc-level" title={levelLabel(doc.level)}>L{doc.level}</span>
            </a>
          </li>
        {/each}
      </ul>
    </section>
  {/each}
</nav>

<style>
  .doc-nav {
    font-size: var(--text-xs);
  }

  .nav-home {
    display: block;
    padding: var(--nav-item-py) var(--space-2);
    margin-bottom: var(--space-2);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    text-decoration: none;
    font-weight: 600;
    font-size: var(--text-sm);
    border-left: 3px solid var(--color-accent);
    padding-left: calc(var(--space-2) + 2px);
  }

  .nav-home:hover,
  .nav-home.active {
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .nav-group {
    margin-bottom: var(--space-3);
  }

  .nav-group-title {
    margin: 0 0 var(--space-1);
    padding: var(--space-1) var(--space-2);
    font-size: 0.65rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-muted);
    border-bottom: 1px solid color-mix(in srgb, var(--color-border) 80%, transparent);
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .nav-item {
    position: relative;
  }

  .nav-list a {
    display: flex;
    align-items: flex-start;
    gap: var(--space-1);
    padding: var(--nav-item-py) var(--space-2);
    padding-left: var(--nav-indent, var(--space-2));
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    text-decoration: none;
    line-height: var(--nav-item-leading);
    min-height: var(--nav-item-min-height);
    border-left: 2px solid transparent;
  }

  .nav-list a.level-2 {
    font-size: 0.8rem;
  }

  .nav-list a.level-3 {
    font-size: 0.75rem;
    color: color-mix(in srgb, var(--color-text-muted) 92%, var(--color-text));
  }

  .doc-tree {
    flex-shrink: 0;
    width: 0.5rem;
    height: 0.5rem;
    margin-top: 0.35em;
    border-left: 1px solid var(--color-border-strong);
    border-bottom: 1px solid var(--color-border-strong);
    opacity: 0.7;
  }

  .nav-list a:hover {
    background: var(--color-bg-surface);
    color: var(--color-text);
    border-left-color: color-mix(in srgb, var(--color-accent) 35%, transparent);
  }

  .nav-list a.active {
    background: var(--color-accent-muted);
    color: var(--color-accent);
    font-weight: 600;
    border-left-color: var(--color-accent);
  }

  .doc-title {
    flex: 1;
    min-width: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    word-break: break-word;
  }

  .doc-level {
    flex-shrink: 0;
    font-size: 0.55rem;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    margin-top: 0.15em;
    padding: 0 3px;
    border: 1px solid var(--color-border);
    border-radius: 3px;
    line-height: 1.2;
  }

  .nav-list a.active .doc-level {
    border-color: color-mix(in srgb, var(--color-accent) 40%, var(--color-border));
    color: var(--color-accent);
  }
</style>
