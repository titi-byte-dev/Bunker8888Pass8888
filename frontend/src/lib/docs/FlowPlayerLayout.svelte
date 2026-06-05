<script lang="ts">
  import type { Snippet } from "svelte";
  import "./flowPlayerLayout.css";

  interface Props {
    visual: Snippet;
    steps?: Snippet;
  }

  let { visual, steps }: Props = $props();
</script>

<div class="flow-player-layout">
  <div class="flow-player-visual">
    {@render visual()}
  </div>
  {#if steps}
    <div class="flow-player-steps">
      {@render steps()}
    </div>
  {/if}
</div>
