<script lang="ts">
  import { page } from "$app/state";
  import DocNav from "$lib/docs/DocNav.svelte";
  import { DOC_MANIFEST } from "$lib/docs/loader";

  let { children } = $props();

  const docTitle = $derived(
    page.params.slug
      ? DOC_MANIFEST.docs.find((d) => d.slug === page.params.slug)?.title
      : null,
  );
</script>
<div class="docs-layout">
  <aside class="docs-sidebar">
    <nav class="breadcrumbs" aria-label="Navegação">
      <a href="/settings">Definições</a>
      <span class="sep">/</span>
      <a href="/settings/docs" class:active={page.url.pathname === "/settings/docs"}>Documentação</a>
      {#if docTitle}
        <span class="sep">/</span>
        <span class="current">{docTitle}</span>
      {/if}
    </nav>
    <DocNav pathname={page.url.pathname} />
  </aside>
  <div class="docs-main">
    {@render children()}
  </div>
</div>

<style>
  .docs-layout {
    display: grid;
    gap: var(--space-6);
    max-width: none;
  }

  @media (min-width: 768px) {
    .docs-layout {
      grid-template-columns: 14rem 1fr;
      align-items: start;
    }
  }

  .docs-sidebar {
    position: sticky;
    top: var(--space-4);
  }

  .breadcrumbs {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-1);
    margin-bottom: var(--space-4);
    font-size: var(--text-xs);
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
  }

  .docs-main {
    min-width: 0;
    max-width: 42rem;
  }
</style>
