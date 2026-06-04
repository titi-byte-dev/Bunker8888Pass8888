<script lang="ts">
  import { onMount } from "svelte";
  import { initTheme } from "$lib/design";
  import "$lib/design/tokens.css";

  let { children } = $props();

  onMount(() => {
    initTheme();
  });
</script>

<svelte:head>
  <title>AegisPass</title>
</svelte:head>

{@render children()}
