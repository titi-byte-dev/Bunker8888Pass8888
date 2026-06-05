<script lang="ts">
  import PlaygroundDev from "$lib/playground/PlaygroundDev.svelte";
  import MotionDemo from "$lib/motion/MotionDemo.svelte";
</script>

<svelte:head>
  <title>Cofre playground — AegisPass</title>
</svelte:head>

<a href="/dev" class="back">← Desenvolvimento</a>
<PlaygroundDev />
<MotionDemo />

<style>
  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }
</style>
