<script lang="ts">
  interface Props {
    /** Texto honesto — Argon2id não expõe progresso real ao JS. */
    label?: string;
    active?: boolean;
  }

  let { label = "A derivar chaves (Argon2id)…", active = false }: Props = $props();
</script>

{#if active}
  <div class="argon" role="status" aria-live="polite">
    <div class="track" aria-hidden="true">
      <div class="bar"></div>
    </div>
    <p>{label}</p>
  </div>
{/if}

<style>
  .argon {
    margin: var(--space-4) 0;
  }

  .track {
    height: 4px;
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    overflow: hidden;
  }

  .bar {
    height: 100%;
    width: 40%;
    background: var(--color-accent);
    border-radius: var(--radius-sm);
    animation: argon-slide 1.2s var(--ease-out) infinite alternate;
  }

  p {
    margin: var(--space-2) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  @keyframes argon-slide {
    from {
      transform: translateX(-10%);
    }
    to {
      transform: translateX(250%);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .bar {
      animation: none;
      width: 100%;
      opacity: 0.6;
    }
  }
</style>
