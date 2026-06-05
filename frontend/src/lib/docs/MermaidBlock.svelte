<script lang="ts">
  import { onMount } from "svelte";
  import { renderMermaid } from "./mermaid";

  interface Props {
    source: string;
    id: string;
    title?: string;
  }

  let { source, id, title = "" }: Props = $props();

  let container = $state<HTMLDivElement | null>(null);
  let error = $state("");
  let busy = $state(true);

  onMount(() => {
    let cancelled = false;

    (async () => {
      if (!container) return;
      busy = true;
      error = "";
      try {
        await renderMermaid(container, source, id);
      } catch (e) {
        if (!cancelled) {
          error = e instanceof Error ? e.message : "Falha ao renderizar diagrama";
        }
      } finally {
        if (!cancelled) busy = false;
      }
    })();

    return () => {
      cancelled = true;
    };
  });
</script>

<figure class="mermaid-block" aria-label={title || "Diagrama"}>
  {#if title}
    <figcaption class="mermaid-caption">{title}</figcaption>
  {/if}
  <div class="mermaid-canvas" class:busy bind:this={container}>
    {#if busy && !error}
      <p class="mermaid-loading">A desenhar fluxo…</p>
    {/if}
  </div>
  {#if error}
    <p class="mermaid-error" role="alert">{error}</p>
    <details class="mermaid-fallback">
      <summary>Ver código Mermaid</summary>
      <pre><code>{source}</code></pre>
    </details>
  {/if}
</figure>

<style>
  .mermaid-block {
    margin: var(--space-4) 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-base);
    overflow: hidden;
  }

  .mermaid-caption {
    margin: 0;
    padding: var(--space-2) var(--space-4);
    font-size: var(--text-xs);
    font-weight: 600;
    color: var(--color-text-muted);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-bg-surface);
  }

  .mermaid-canvas {
    padding: var(--space-4);
    overflow-x: auto;
    min-height: 4rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .mermaid-canvas :global(svg) {
    display: block;
    margin: 0 auto;
  }

  .mermaid-loading {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .mermaid-error {
    margin: 0;
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
    border-top: 1px solid var(--color-border);
  }

  .mermaid-fallback {
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-xs);
    border-top: 1px solid var(--color-border);
  }

  .mermaid-fallback pre {
    margin: var(--space-2) 0 0;
    overflow-x: auto;
    font-family: var(--font-mono);
    font-size: var(--text-xs);
  }
</style>
