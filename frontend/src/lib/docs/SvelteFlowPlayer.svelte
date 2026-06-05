<script lang="ts">
  import { browser } from "$app/environment";
  import { SvelteFlow, Background, Controls, type Node, type Edge } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import type { DocFlow, DocFlowNodeData, DocFlowStep } from "./types";
  import DocFlowNode from "./DocFlowNode.svelte";
  import DocSequenceEdge from "./DocSequenceEdge.svelte";
  import FlowStepControls from "./FlowStepControls.svelte";
  import FlowPlayerLayout from "./FlowPlayerLayout.svelte";
  import { prefersReducedMotion } from "$lib/motion/reduced";
  import { buildFlowEdges, buildFlowNodes, flowCanvasSize } from "./buildFlowGraph";

  interface Props {
    flow: DocFlow;
  }

  let { flow }: Props = $props();

  const nodeTypes = { docFlow: DocFlowNode };
  const edgeTypes = { docSequence: DocSequenceEdge };
  const graph = $derived(flow.graph!);
  const canvas = $derived(flowCanvasSize(graph));

  let stepIndex = $state(0);
  let playing = $state(false);

  const messageSteps = $derived(
    flow.steps.filter((s): s is Extract<DocFlowStep, { kind: "message" }> => s.kind === "message"),
  );

  let nodes = $state.raw<Node<DocFlowNodeData, "docFlow">[]>([]);
  let edges = $state.raw<Edge[]>([]);

  $effect(() => {
    const idx = stepIndex;
    const animate = !prefersReducedMotion();
    nodes = buildFlowNodes(graph, messageSteps, idx);
    edges = buildFlowEdges(graph, idx, animate);
  });

  function goTo(idx: number) {
    stepIndex = Math.max(0, Math.min(idx, messageSteps.length - 1));
  }

  function togglePlay() {
    if (prefersReducedMotion()) return;
    playing = !playing;
  }

  $effect(() => {
    if (!playing || messageSteps.length === 0) return;
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

<FlowPlayerLayout>
  {#snippet visual()}
    <div class="sf-scroll">
      <div class="sf-canvas" style:width="{canvas.width}px" style:height="{canvas.height}px">
        {#if browser}
          <SvelteFlow
            bind:nodes
            bind:edges
            {nodeTypes}
            {edgeTypes}
            width={canvas.width}
            height={canvas.height}
            nodesDraggable={false}
            nodesConnectable={false}
            elementsSelectable={false}
            panOnDrag
            zoomOnScroll
            minZoom={0.65}
            maxZoom={1.25}
            proOptions={{ hideAttribution: true }}
          >
            <Background gap={20} size={1} />
            <Controls showLock={false} />
          </SvelteFlow>
        {:else}
          <p class="sf-loading" aria-hidden="true">A carregar diagrama…</p>
        {/if}
      </div>
    </div>
  {/snippet}

  {#snippet steps()}
    <FlowStepControls
      {messageSteps}
      {stepIndex}
      {playing}
      onStep={goTo}
      onTogglePlay={togglePlay}
    />
  {/snippet}
</FlowPlayerLayout>

<style>
  .sf-scroll {
    width: 100%;
    overflow-x: auto;
    overflow-y: hidden;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-base);
    -webkit-overflow-scrolling: touch;
  }

  .sf-canvas {
    min-width: 100%;
  }

  .sf-loading {
    margin: 0;
    padding: var(--space-6);
    text-align: center;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .sf-scroll :global(.svelte-flow) {
    --xy-edge-stroke-default: var(--color-border);
    --xy-node-background-color-default: var(--color-bg-surface);
    --xy-background-color-default: var(--color-bg-base);
    --xy-background-pattern-color-default: color-mix(
      in srgb,
      var(--color-border) 40%,
      transparent
    );
  }

  .sf-scroll :global(.svelte-flow__controls) {
    box-shadow: none;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .sf-scroll :global(.svelte-flow__controls-button) {
    background: var(--color-bg-surface);
    border-bottom: 1px solid var(--color-border);
    fill: var(--color-text-muted);
  }

  .sf-scroll :global(.svelte-flow__controls-button:hover) {
    background: var(--color-accent-muted);
  }

  .sf-scroll :global(.svelte-flow__handle) {
    width: 1px;
    height: 1px;
    min-width: 0;
    min-height: 0;
    opacity: 0;
    border: none;
    background: transparent;
    pointer-events: none;
  }

  .sf-scroll :global(.flow-edge-current path) {
    stroke: var(--color-accent);
    stroke-width: 2.5;
  }

  .sf-scroll :global(.flow-edge-done path) {
    stroke: var(--color-success-fg);
    stroke-width: 2;
    opacity: 0.85;
  }

  .sf-scroll :global(.flow-edge-pending path) {
    stroke: var(--color-border);
    stroke-width: 1.5;
    opacity: 0.35;
  }

  .sf-scroll :global(.flow-edge-current .svelte-flow__edge-text) {
    fill: var(--color-text);
    font-weight: 600;
    font-size: 11px;
  }

  .sf-scroll :global(.flow-edge-done .svelte-flow__edge-text) {
    fill: var(--color-success-fg);
    font-size: 10px;
    font-weight: 500;
  }

  .sf-scroll :global(.flow-edge-pending .svelte-flow__edge-text) {
    fill: var(--color-text-muted);
    font-size: 10px;
    opacity: 0.65;
  }

  .sf-scroll :global(.svelte-flow__edge-textbg) {
    fill: var(--color-bg-base);
    fill-opacity: 0.92;
  }
</style>
