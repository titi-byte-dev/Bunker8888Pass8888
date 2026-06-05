<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * Panel (UI-012) — superficie elevada com cabecalho opcional.
   * Substitui ~85 blocos `.panel`/`.panel-head` re-definidos por pagina.
   */
  interface Props {
    title?: string;
    /** Padding interno. `none` para tabelas que gerem o seu proprio espaco. */
    padding?: "none" | "md";
    /** `inset` para fundo mais fundo (formularios dentro de painel). */
    variant?: "surface" | "inset";
    /** Accoes no cabecalho do painel. */
    actions?: Snippet;
    children: Snippet;
  }
  let { title, padding = "md", variant = "surface", actions, children }: Props = $props();
</script>

<section class="panel {variant}">
  {#if title || actions}
    <header class="panel-head">
      {#if title}<h2>{title}</h2>{/if}
      {#if actions}<div class="panel-actions">{@render actions()}</div>{/if}
    </header>
  {/if}
  <div class="panel-body" class:pad={padding === "md"}>
    {@render children()}
  </div>
</section>

<style>
  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-inset);
    overflow: hidden;
  }
  .surface { background: var(--color-bg-surface); }
  .inset { background: var(--color-bg-inset); }

  .panel-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
  }
  h2 {
    margin: 0;
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }
  .panel-actions { display: flex; gap: var(--space-2); }
  .panel-body.pad { padding: var(--space-3); }
</style>
