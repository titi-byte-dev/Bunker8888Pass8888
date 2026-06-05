<script lang="ts">
  /**
   * Demonstração dos presets GSAP (UI-005) — só em /dev.
   */
  import { onMount } from "svelte";
  import {
    animateFadeIn,
    animateListStagger,
    animatePanelEnter,
    animateSavedHighlight,
    prefersReducedMotion,
  } from "$lib/motion";
  import { runMotionScope } from "$lib/motion/scope";

  let demoRoot: HTMLElement | undefined;
  let reduced = $state(false);

  onMount(() => {
    reduced = prefersReducedMotion();
  });

  function replayPanel() {
    if (!demoRoot) return;
    const panel = demoRoot.querySelector("[data-demo-panel]");
    if (panel) animatePanelEnter(panel);
  }

  function replayList() {
    if (!demoRoot) return;
    const list = demoRoot.querySelector("[data-demo-list]");
    if (list) runMotionScope(list as HTMLElement, () => animateListStagger(list as HTMLElement, "li"));
  }

  function replaySaved() {
    if (!demoRoot) return;
    const row = demoRoot.querySelector("[data-demo-saved]");
    if (row) animateSavedHighlight(row);
  }

  function replayFade() {
    if (!demoRoot) return;
    const el = demoRoot.querySelector("[data-demo-fade]");
    if (el) animateFadeIn(el);
  }
</script>

<section class="motion-demo" bind:this={demoRoot}>
  <h2>Motion system (UI-005)</h2>
  <p class="hint">
    GSAP + <code>prefers-reduced-motion</code>
    {#if reduced}
      — <strong>activo</strong> no teu sistema (animações desactivadas).
    {:else}
      — inactivo (animações visíveis).
    {/if}
  </p>

  <div class="grid">
    <article class="card" data-demo-panel>
      <h3>Panel enter</h3>
      <button type="button" onclick={replayPanel}>Repetir</button>
    </article>

    <article class="card">
      <h3>List stagger</h3>
      <ul data-demo-list>
        <li>Item A</li>
        <li>Item B</li>
        <li>Item C</li>
      </ul>
      <button type="button" onclick={replayList}>Repetir</button>
    </article>

    <article class="card">
      <h3>Saved highlight</h3>
      <div class="saved-row" data-demo-saved>Login guardado</div>
      <button type="button" onclick={replaySaved}>Repetir</button>
    </article>

    <article class="card">
      <h3>Fade in</h3>
      <p data-demo-fade>Micro-transição</p>
      <button type="button" onclick={replayFade}>Repetir</button>
    </article>
  </div>
</section>

<style>
  .motion-demo {
    margin-top: var(--space-8);
    padding-top: var(--space-6);
    border-top: 1px solid var(--color-border);
  }

  h2 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-xl);
  }

  .hint {
    margin: 0 0 var(--space-4);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .grid {
    display: grid;
    gap: var(--space-4);
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
  }

  .card {
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
  }

  .card h3 {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    font-weight: 600;
  }

  ul {
    list-style: none;
    margin: 0 0 var(--space-3);
    padding: 0;
  }

  li {
    padding: var(--space-2);
    border-bottom: 1px solid var(--color-border);
    font-size: var(--text-sm);
  }

  .saved-row {
    padding: var(--space-2) var(--space-3);
    margin-bottom: var(--space-3);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    font-size: var(--text-sm);
  }

  button {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    color: var(--color-text);
    font-size: var(--text-sm);
    cursor: pointer;
  }
</style>
