<script lang="ts">
  /**
   * Eyebrow (UI-012) — pequena etiqueta acima do titulo.
   * Substitui ~20 copias de `.eyebrow` espalhadas pelas paginas.
   * `devOnly`: so aparece em import.meta.env.DEV (ex.: IDs de task FIN-005).
   */
  interface Props {
    text: string;
    devOnly?: boolean;
  }
  let { text, devOnly = false }: Props = $props();
  const show = $derived(!devOnly || import.meta.env.DEV);
</script>

{#if show}
  <p class="eyebrow" class:dev={devOnly}>{text}</p>
{/if}

<style>
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  .eyebrow.dev {
    font-family: var(--font-mono);
    letter-spacing: 0.04em;
  }
</style>
