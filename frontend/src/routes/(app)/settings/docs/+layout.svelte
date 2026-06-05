<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import DocNav from "$lib/docs/DocNav.svelte";
  import DocSearch from "$lib/docs/DocSearch.svelte";
  import { DOC_MANIFEST } from "$lib/docs/loader";
  import {
    clampDocNavWidth,
    loadDocNavCollapsed,
    loadDocNavWidth,
    saveDocNavCollapsed,
    saveDocNavWidth,
  } from "$lib/docs/docNavState";

  let { children } = $props();

  let navCollapsed = $state(loadDocNavCollapsed());
  let navWidth = $state(loadDocNavWidth());
  let resizing = $state(false);
  let layoutEl: HTMLDivElement | undefined = $state();

  const docTitle = $derived(
    page.params.slug
      ? DOC_MANIFEST.docs.find((d) => d.slug === page.params.slug)?.title
      : null,
  );

  function toggleDocNav() {
    navCollapsed = !navCollapsed;
    saveDocNavCollapsed(navCollapsed);
  }

  function startResize(e: MouseEvent) {
    if (navCollapsed) return;
    resizing = true;
    e.preventDefault();
  }

  onMount(() => {
    function onMove(e: MouseEvent) {
      if (!resizing || !layoutEl) return;
      const left = layoutEl.getBoundingClientRect().left;
      navWidth = clampDocNavWidth(e.clientX - left);
    }

    function onUp() {
      if (!resizing) return;
      resizing = false;
      saveDocNavWidth(navWidth);
    }

    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => {
      window.removeEventListener("mousemove", onMove);
      window.removeEventListener("mouseup", onUp);
    };
  });
</script>

<div
  class="docs-layout"
  class:nav-collapsed={navCollapsed}
  class:resizing
  bind:this={layoutEl}
  style:--doc-nav-width="{navCollapsed ? '2.75rem' : `${navWidth}px`}"
>
  <aside class="docs-sidebar" aria-label="Índice da documentação">
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
      <div class="docs-sidebar-scroll">
        <DocSearch />
        <DocNav pathname={page.url.pathname} />
      </div>
    {/if}
  </aside>

  <!-- Separador arrastável — ajusta largura do índice (desktop). -->
  <button
    type="button"
    class="docs-resizer"
    class:hidden={navCollapsed}
    onmousedown={startResize}
    aria-label="Redimensionar índice"
    title="Arrastar para ajustar largura do índice"
  ></button>

  <div class="docs-main">
    <div class="docs-main-scroll">
      {@render children()}
    </div>
  </div>
</div>

<style>
  .docs-layout {
    display: flex;
    flex: 1;
    min-height: 0;
    min-width: 0;
    width: 100%;
    overflow: hidden;
  }

  .docs-sidebar {
    flex: 0 0 var(--doc-nav-width);
    width: var(--doc-nav-width);
    min-width: 0;
    display: flex;
    flex-direction: column;
    min-height: 0;
    border-right: 1px solid var(--color-border);
    background: var(--color-bg-elevated);
    box-sizing: border-box;
  }

  .docs-sidebar-head {
    flex-shrink: 0;
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--space-1);
    padding: var(--space-2) var(--space-2) var(--space-2) var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }

  .docs-sidebar-scroll {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
    overscroll-behavior: contain;
    padding: var(--space-2) var(--space-2) var(--space-3);
    scrollbar-gutter: stable;
  }

  .docs-resizer {
    flex: 0 0 6px;
    width: 6px;
    margin: 0;
    padding: 0;
    border: none;
    cursor: col-resize;
    background: transparent;
    position: relative;
    z-index: 2;
  }

  .docs-resizer.hidden {
    display: none;
  }

  .docs-resizer::after {
    content: "";
    position: absolute;
    inset: 0;
    width: 2px;
    margin: 0 auto;
    border-radius: 1px;
    background: var(--color-border);
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .docs-resizer:hover::after,
  .docs-layout.resizing .docs-resizer::after {
    background: var(--color-accent);
  }

  .docs-main {
    flex: 1;
    min-width: 0;
    min-height: 0;
    display: flex;
    flex-direction: column;
  }

  .docs-main-scroll {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    overflow-x: hidden;
    overscroll-behavior: contain;
    padding: var(--space-3) 10px var(--space-6);
    scrollbar-gutter: stable;
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
    max-width: 10rem;
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

  @media (max-width: 767px) {
    .docs-layout {
      flex-direction: column;
      overflow: visible;
    }

    .docs-sidebar {
      flex: none;
      width: 100%;
      max-height: 42vh;
      border-right: none;
      border-bottom: 1px solid var(--color-border);
    }

    .docs-resizer {
      display: none;
    }

    .docs-main-scroll {
      overflow: visible;
      padding-inline: 10px;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .docs-resizer::after {
      transition: none;
    }
  }
</style>
