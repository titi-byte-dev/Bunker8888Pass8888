<script lang="ts">
  import type { DocFlow, DocFlowStep } from "./types";
  import MermaidBlock from "./MermaidBlock.svelte";
  import SvelteFlowPlayer from "./SvelteFlowPlayer.svelte";
  import FlowStepControls from "./FlowStepControls.svelte";
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
  const current = $derived(messageSteps[stepIndex]);

  function goTo(idx: number) {
    stepIndex = Math.max(0, Math.min(idx, messageSteps.length - 1));
  }

  function togglePlay() {
    if (prefersReducedMotion()) return;
    playing = !playing;
  }

  $effect(() => {
    if (!playing || messageSteps.length === 0 || useSvelteFlow) return;
    const id = setInterval(() => {
      if (stepIndex >= messageSteps.length - 1) {
        playing = false;
        return;
      }
      goTo(stepIndex + 1);
    }, 2800);
    return () => clearInterval(id);
  });
</script>

{#if useSvelteFlow}
  <SvelteFlowPlayer {flow} />
{:else}
  <div class="flow-player">
    <MermaidBlock source={flow.source} id={flow.id} title={flow.title || undefined} />

    {#if hasSteps}
      <FlowStepControls
        {messageSteps}
        {stepIndex}
        {playing}
        onStep={goTo}
        onTogglePlay={togglePlay}
      />
    {:else if flow.type === "flowchart"}
      <p class="flow-hint">
        Diagrama estático — segue as setas do fluxograma na ordem lógica.
      </p>
    {/if}
  </div>
{/if}

<style>
  .flow-player {
    margin: var(--space-4) 0;
  }

  .flow-hint {
    margin: var(--space-2) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
