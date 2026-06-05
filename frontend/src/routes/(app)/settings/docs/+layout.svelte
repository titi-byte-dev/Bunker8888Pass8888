<script lang="ts">
  import { page } from "$app/state";
  import DocNav from "$lib/docs/DocNav.svelte";

  let { children } = $props();
</script>

<div class="docs-layout">
  <aside class="docs-sidebar">
    <a href="/settings" class="back-link">← Definições</a>
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

  .back-link {
    display: inline-block;
    margin-bottom: var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-link);
    text-decoration: none;
  }

  .back-link:hover {
    text-decoration: underline;
  }

  .docs-main {
    min-width: 0;
    max-width: 42rem;
  }
</style>
