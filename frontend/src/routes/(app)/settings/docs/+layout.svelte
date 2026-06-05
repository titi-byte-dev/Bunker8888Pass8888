<script lang="ts">
  import { page } from "$app/state";
  import DocNav from "$lib/docs/DocNav.svelte";
  import DocSearch from "$lib/docs/DocSearch.svelte";
  import { DOC_MANIFEST } from "$lib/docs/loader";
  let { children } = $props();

  let navCollapsed = $state(false);

  const docTitle = $derived(
    page.params.slug
      ? DOC_MANIFEST.docs.find((d) => d.slug === page.params.slug)?.title
      : null,
  );

  function toggleDocNav() {
    navCollapsed = !navCollapsed;
  }
</script>

<div class="docs-layout" class:nav-collapsed={navCollapsed}>
  <aside class="docs-sidebar">
    <div class="docs-sidebar-head">
      <nav class="breadcrumbs" aria-label="Navegação">
        {#if !navCollapsed}
          <a href="/settings">Definições</a>
          <span class="sep">/</span>
          <a href="/settings/docs" class:active={page.url.pathname === "/settings/docs"}>Documentação</a>
          {#if docTitle}
            <span class="sep">/</span>
            <span class="current">{docTitle}</span>
          {/if}
        {:else}
          <a href="/settings/docs" class="nav-icon-link" title="Documentação">Docs</a>
        {/if}
      </nav>
      <button
        type="button"
        class="nav-toggle"
        onclick={toggleDocNav}
        aria-expanded={!navCollapsed}
        aria-label={navCollapsed ? "Expandir índice" : "Recolher índice"}
        title={navCollapsed ? "Expandir índice" : "Recolher índice"}
      >
        {navCollapsed ? "»" : "«"}
      </button>
    </div>
    {#if !navCollapsed}
      <DocSearch />
      <DocNav pathname={page.url.pathname} />
    {/if}
  </aside>
  <div class="docs-main">
    {@render children()}
  </div>
</div>

<style>
  .docs-layout {
    display: grid;
    gap: var(--space-4);
    max-width: none;
    width: 100%;
  }

  @media (min-width: 768px) {
    .docs-layout {
      grid-template-columns: 9.5rem 1fr;
      align-items: start;
      gap: var(--space-3);
    }

    .docs-layout.nav-collapsed {
      grid-template-columns: 2.5rem 1fr;
    }
  }

  @media (min-width: 1200px) {
    .docs-layout {
      grid-template-columns: 10.5rem 1fr;
    }
  }

  .docs-sidebar {
    position: sticky;
    top: var(--space-4);
    min-width: 0;
  }

  .docs-sidebar-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-1);
    margin-bottom: var(--space-3);
  }

  .nav-toggle {
    flex-shrink: 0;
    width: 1.75rem;
    height: 1.75rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    line-height: 1;
  }

  .nav-toggle:hover {
    color: var(--color-text);
    background: var(--color-accent-muted);
  }

  .breadcrumbs {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-1);
    font-size: var(--text-xs);
    min-width: 0;
  }

  .breadcrumbs a {
    color: var(--color-link);
    text-decoration: none;
  }

  .breadcrumbs a:hover,
  .breadcrumbs a.active {
    text-decoration: underline;
  }

  .breadcrumbs .sep {
    color: var(--color-text-muted);
  }

  .breadcrumbs .current {
    color: var(--color-text-muted);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 12rem;
  }

  .nav-icon-link {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2rem;
    height: 2rem;
    text-decoration: none;
    border-radius: var(--radius-sm);
  }

  .nav-icon-link:hover {
    background: var(--color-accent-muted);
  }

  .docs-main {
    min-width: 0;
    max-width: none;
    width: 100%;
  }
</style>
