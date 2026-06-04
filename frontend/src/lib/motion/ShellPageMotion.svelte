<script lang="ts">
  /**
   * Wrapper de conteúdo com cross-fade suave entre rotas (UI-005).
   */
  import { afterNavigate } from "$app/navigation";
  import type { Snippet } from "svelte";
  import { animatePageEnter } from "$lib/motion/presets";
  import { runMotionScope } from "$lib/motion/scope";

  interface Props {
    children: Snippet;
  }

  let { children }: Props = $props();
  let root: HTMLElement | undefined;

  let cleanup: (() => void) | undefined;

  function playEnter() {
    if (!root) return;
    cleanup?.();
    cleanup = runMotionScope(root, () => animatePageEnter(root!));
  }

  afterNavigate(() => {
    playEnter();
  });
</script>

<div class="shell-page-motion" bind:this={root}>
  {@render children()}
</div>

<style>
  .shell-page-motion {
    width: 100%;
  }
</style>
