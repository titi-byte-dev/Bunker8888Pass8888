<script lang="ts">
  import type { DocFlow } from "./types";
  import { annotateGlossaryHtml } from "./glossary";
  import { splitSectionHtml } from "./splitSection";
  import FlowPlayer from "./FlowPlayer.svelte";

  interface Props {
    html: string;
    flows?: DocFlow[];
  }

  let { html, flows = [] }: Props = $props();

  const parts = $derived(splitSectionHtml(annotateGlossaryHtml(html), flows));
</script>

{#each parts as part, i (part.kind === "flow" ? part.flow.id : `html-${i}`)}
  {#if part.kind === "html"}
    <div class="prose section-html">
      {@html part.content}
    </div>
  {:else}
    <FlowPlayer flow={part.flow} />
  {/if}
{/each}
