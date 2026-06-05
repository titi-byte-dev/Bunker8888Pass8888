<script lang="ts">
  import DocPageView from "$lib/docs/DocPageView.svelte";

  let { data } = $props();
</script>

<svelte:head>
  <title>{data.page.title} — Documentação</title>
</svelte:head>

<DocPageView page={data.page} />
