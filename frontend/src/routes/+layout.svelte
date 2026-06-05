<script lang="ts">
  import { onMount } from "svelte";
  import { initTheme, initPalette } from "$lib/design";
  import "$lib/design/tokens.css";
  import "$lib/motion/motion.css";

  let { children } = $props();

  onMount(() => {
    initTheme();
    // tenantDefault chegara da config da empresa (white-label); por agora undefined.
    initPalette();
  });
</script>

<svelte:head>
  <title>AegisPass</title>
</svelte:head>

{@render children()}
