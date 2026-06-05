<script lang="ts">
  import type { DocFlow, DocFlowStep } from "./types";
  import MermaidBlock from "./MermaidBlock.svelte";
  import SvelteFlowPlayer from "./SvelteFlowPlayer.svelte";
  import FlowStepControls from "./FlowStepControls.svelte";
  import FlowPlayerLayout from "./FlowPlayerLayout.svelte";
  import { prefersReducedMotion } from "$lib/motion/reduced";

  interface Props {
    flow: DocFlow;
  }

  let { flow }: Props = $props();

  const useSvelteFlow = $derived(
    flow.renderer === "svelteflow" && (flow.graph?.nodes?.length ?? 0) > 0,
  );

  let stepIndex = $state(0);
  let playing = $state(false);

  const messageSteps = $derived(
    flow.steps.filter((s): s is Extract<DocFlowStep, { kind: "message" }> => s.kind === "message"),
  );
  const hasSteps = $derived(messageSteps.length > 0);

  function goTo(idx: number) {
    stepIndex = Math.max(0, Math.min(idx, messageSteps.length - 1));
  }

  function togglePlay() {
    if (prefersReducedMotion()) return;
    playing = !playing;
  }

  $effect(() => {
    if (!playing || messageSteps.length === 0 || useSvelteFlow) return;
    if (stepIndex >= messageSteps.length - 1) {
      playing = false;
      return;
    }
    const timer = setTimeout(() => {
      if (playing) goTo(stepIndex + 1);
    }, 2800);
    return () => clearTimeout(timer);
  });
</script>

{#if useSvelteFlow}
  <SvelteFlowPlayer {flow} />
{:else if hasSteps || flow.type === "flowchart"}
  <FlowPlayerLayout>
    {#snippet visual()}
      <MermaidBlock source={flow.source} id={flow.id} title={flow.title || undefined} />
    {/snippet}
    {#snippet steps()}
      {#if hasSteps}
        <FlowStepControls
          {messageSteps}
          {stepIndex}
          {playing}
          onStep={goTo}
          onTogglePlay={togglePlay}
        />
      {:else}
        <p class="flow-hint">
          Diagrama estático — segue as setas do fluxograma na ordem lógica.
        </p>
      {/if}
    {/snippet}
  </FlowPlayerLayout>
{:else}
  <FlowPlayerLayout>
    {#snippet visual()}
      <MermaidBlock source={flow.source} id={flow.id} title={flow.title || undefined} />
    {/snippet}
  </FlowPlayerLayout>
{/if}

<style>
  .flow-hint {
    margin: 0;
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
