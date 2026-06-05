<script lang="ts">
  import { onMount } from "svelte";
  import type { Snippet } from "svelte";
  import { animatePanelEnter } from "$lib/motion/presets";
  import { runMotionScope } from "$lib/motion/scope";

  interface Props {
    title?: string;
    subtitle?: string;
    children: Snippet;
  }

  let { title, subtitle, children }: Props = $props();
  let card: HTMLElement | undefined;

  onMount(() => {
    if (!card) return;
    return runMotionScope(card, () => animatePanelEnter(card!));
  });
</script>

<div class="auth-page">
  <a href="/auth" class="brand">
    <span class="mark" aria-hidden="true">◆</span>
    <span class="name">AegisPass</span>
  </a>

  <div class="card" bind:this={card}>
    {#if title}
      <h1>{title}</h1>
    {/if}
    {#if subtitle}
      <p class="subtitle">{subtitle}</p>
    {/if}
    {@render children()}
  </div>
</div>

<style>
  .auth-page {
    min-height: 100dvh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-6) var(--space-4);
    box-sizing: border-box;
  }

  .brand {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    margin-bottom: var(--space-6);
    text-decoration: none;
    color: inherit;
  }

  .mark {
    color: var(--color-accent);
    font-size: var(--text-xl);
  }

  .name {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: var(--text-xl);
  }

  .card {
    width: 100%;
    max-width: 24rem;
    padding: var(--space-6);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-lg);
    background: var(--color-bg-surface);
    box-shadow: var(--shadow-inset);
    box-sizing: border-box;
  }

  h1 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    line-height: var(--leading-tight);
  }

  .subtitle {
    margin: 0 0 var(--space-6);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    line-height: var(--leading-body);
  }
</style>
