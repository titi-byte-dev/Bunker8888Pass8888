<script lang="ts">
  import type { Snippet } from "svelte";
  import { page } from "$app/stores";
  import Breadcrumbs from "./Breadcrumbs.svelte";
  import Eyebrow from "./Eyebrow.svelte";

  /**
   * PageShell (UI-012) — cabecalho de pagina unico.
   * Layout fluido: ocupa toda a largura disponivel do shell (sem variantes narrow/wide).
   */
  interface Props {
    title: string;
    description?: string;
    /** ID de task — so em DEV. */
    taskId?: string;
    /** Esconde o trilho de breadcrumbs (hubs de topo). */
    breadcrumb?: boolean;
    /** Segmento final dinamico para o trilho (ex.: nome do item). */
    leaf?: string;
    /** Accoes alinhadas a direita do titulo. */
    actions?: Snippet;
    children: Snippet;
  }

  let {
    title,
    description,
    taskId,
    breadcrumb = true,
    leaf,
    actions,
    children,
  }: Props = $props();

  const pathname = $derived($page.url.pathname);
</script>

<section class="page">
  <header class="page-head">
    <div class="head-text">
      {#if breadcrumb}<Breadcrumbs {pathname} {leaf} />{/if}
      {#if taskId}<Eyebrow text={taskId} devOnly />{/if}
      <h1>{title}</h1>
      {#if description}<p class="lead">{description}</p>{/if}
    </div>
    {#if actions}
      <div class="head-actions">{@render actions()}</div>
    {/if}
  </header>

  {@render children()}
</section>

<style>
  .page {
    width: 100%;
    max-width: none;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--page-gap);
  }

  .page-head {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-2);
    flex-wrap: wrap;
  }
  .head-text {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }
  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    line-height: var(--leading-tight);
  }
  .lead {
    margin: 0;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    line-height: var(--leading-snug);
    max-width: var(--prose-max);
  }
  .head-actions {
    display: flex;
    gap: var(--space-2);
    flex-shrink: 0;
  }
</style>
