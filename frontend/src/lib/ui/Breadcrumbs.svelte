<script lang="ts">
  /**
   * Breadcrumbs (UI-011) — trilho "Financas > Fiscal" derivado do ROUTE_TREE.
   * Ultimo segmento e a pagina actual (sem link); anteriores sao <a>.
   */
  import { routeTrail, type RouteNode } from "$lib/shell/routes";

  interface Props {
    pathname: string;
    /** Acrescenta um segmento final dinamico (ex.: nome do item do cofre). */
    leaf?: string;
  }
  let { pathname, leaf }: Props = $props();

  const trail = $derived<RouteNode[]>(routeTrail(pathname));
</script>

{#if trail.length > 0}
  <nav class="crumbs" aria-label="Localizacao na app">
    <ol>
      {#each trail as node, i}
        {@const isLast = i === trail.length - 1 && !leaf}
        <li>
          {#if isLast}
            <span aria-current="page">{node.label}</span>
          {:else}
            <a href={node.href}>{node.label}</a>
            <span class="sep" aria-hidden="true">&rsaquo;</span>
          {/if}
        </li>
      {/each}
      {#if leaf}
        <li><span aria-current="page">{leaf}</span></li>
      {/if}
    </ol>
  </nav>
{/if}

<style>
  .crumbs ol {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-1);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
  .crumbs li { display: inline-flex; align-items: center; gap: var(--space-1); }
  .crumbs a { color: var(--color-text-muted); text-decoration: none; }
  .crumbs a:hover { color: var(--color-text); }
  .crumbs [aria-current="page"] { color: var(--color-text); font-weight: 600; }
  .sep { color: var(--color-border-strong); }
</style>
