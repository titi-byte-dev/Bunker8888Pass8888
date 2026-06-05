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
              title={doc.summary}
            >
              <span class="doc-title">{doc.title}</span>
              {#if doc.level > 1}
                <span class="doc-level">Nv.{doc.maxLevel}</span>
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
    font-size: var(--text-sm);
  }

  .nav-home {
    display: block;
    padding: var(--space-2) var(--space-3);
    margin-bottom: var(--space-4);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    text-decoration: none;
    font-weight: 600;
  }

  .nav-home:hover,
  .nav-home.active {
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .nav-group {
    margin-bottom: var(--space-4);
  }

  .nav-group-title {
    margin: 0 0 var(--space-2);
    padding: 0 var(--space-3);
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  .nav-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .nav-list a {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    color: var(--color-text-muted);
    text-decoration: none;
  }

  .nav-list a:hover {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }

  .nav-list a.active {
    background: var(--color-accent-muted);
    color: var(--color-text);
    font-weight: 500;
  }

  .doc-title {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .doc-level {
    flex-shrink: 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }
</style>
