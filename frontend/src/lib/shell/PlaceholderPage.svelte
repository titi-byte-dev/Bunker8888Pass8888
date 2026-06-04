<script lang="ts">
  interface Props {
    title: string;
    description?: string;
    taskId?: string;
  }

  let { title, description, taskId }: Props = $props();
</script>

<section class="placeholder">
  <p class="eyebrow">{taskId ?? "UI-002"}</p>
  <h1>{title}</h1>
  {#if description}
    <p class="lead">{description}</p>
  {/if}
  <p class="hint">
    Esta página será implementada nas tasks seguintes. Usa
    <a href="/dev">/dev</a> para testar o cofre end-to-end.
  </p>
</section>

<style>
  .placeholder {
    max-width: 36rem;
  }

  .eyebrow {
    margin: 0 0 var(--space-2);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--color-accent);
  }

  h1 {
    margin: 0 0 var(--space-3);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    line-height: var(--leading-tight);
  }

  .lead {
    margin: 0 0 var(--space-4);
    color: var(--color-text-muted);
    font-size: var(--text-base);
  }

  .hint {
    margin: 0;
    padding: var(--space-4);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  a {
    color: var(--color-link);
  }
</style>
