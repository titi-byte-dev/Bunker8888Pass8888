<script lang="ts">
  import { docsByCategory } from "./loader";

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
</script>

<nav class="doc-nav" aria-label="Índice da documentação">
  <a href="/settings/docs" class="nav-home" class:active={pathname === "/settings/docs"}>
    Índice
  </a>
  {#each groups as group (group.id)}
    <div class="nav-group">
      <h3 class="nav-group-title">{group.label}</h3>
      <ul class="nav-list">
        {#each group.docs as doc (doc.slug)}
          <li>
            <a
              href="/settings/docs/{doc.slug}"
              class:active={isActive(doc.slug)}
              title={doc.title}
            >
              <span class="doc-title">{shortTitle(doc.title)}</span>
              {#if doc.level > 1}
                <span class="doc-level">L{doc.maxLevel}</span>
              {/if}
            </a>
          </li>
        {/each}
      </ul>
    </div>
  {/each}
</nav>

<style>
  .doc-nav {
    font-size: var(--text-xs);
  }

  .nav-home {
    display: block;
    padding: var(--space-1) var(--space-2);
    margin-bottom: var(--space-3);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    text-decoration: none;
    font-weight: 600;
    font-size: var(--text-sm);
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
    padding: 0 var(--space-2);
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .nav-list a {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-1);
    padding: 4px var(--space-2);
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    text-decoration: none;
    line-height: 1.35;
  }

  .nav-list a:hover {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }

  .nav-list a.active {
    background: var(--color-accent-muted);
    color: var(--color-accent);
    font-weight: 600;
  }

  .doc-title {
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
    font-size: 0.6rem;
    color: var(--color-text-muted);
    font-family: var(--font-mono);
    margin-top: 2px;
  }
</style>
