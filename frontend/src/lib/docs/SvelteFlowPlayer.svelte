<script lang="ts">
  import { browser } from "$app/environment";
  import { onMount } from "svelte";
  import { SvelteFlow, Background, Controls, type Node, type Edge } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import "./svelteFlowTheme.css";
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

  let scrollEl: HTMLDivElement | undefined = $state();
  let viewportWidth = $state(640);

  const canvas = $derived(flowCanvasSize(graph, viewportWidth));

  let stepIndex = $state(0);
  let playing = $state(false);

  const messageSteps = $derived(
    flow.steps.filter((s): s is Extract<DocFlowStep, { kind: "message" }> => s.kind === "message"),
  );

  let nodes = $state.raw<Node<DocFlowNodeData, "docFlow">[]>([]);
  let edges = $state.raw<Edge[]>([]);

  $effect(() => {
    const idx = stepIndex;
    const w = viewportWidth;
    const animate = !prefersReducedMotion();
    nodes = buildFlowNodes(graph, messageSteps, idx, w);
    edges = buildFlowEdges(graph, idx, animate, w);
  });

  onMount(() => {
    if (!scrollEl) return;
    const measure = () => {
      if (!scrollEl) return;
      viewportWidth = Math.max(360, Math.floor(scrollEl.clientWidth));
    };
    measure();
    const ro = new ResizeObserver(measure);
    ro.observe(scrollEl);
    return () => ro.disconnect();
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
    <div class="sf-scroll sf-theme" bind:this={scrollEl}>
      <div
        class="sf-canvas"
        style:width="{canvas.width}px"
        style:height="{canvas.height}px"
        style:min-height="{Math.max(canvas.height, 280)}px"
      >
        {#if browser}
          <SvelteFlow
            bind:nodes
            bind:edges
            {nodeTypes}
            {edgeTypes}
            width={canvas.width}
            height={canvas.height}
            defaultViewport={{ x: 0, y: 0, zoom: 1 }}
            nodesDraggable={false}
            nodesConnectable={false}
            elementsSelectable={false}
            panOnDrag
            panOnScroll
            zoomOnScroll
            minZoom={0.75}
            maxZoom={1.35}
            proOptions={{ hideAttribution: true }}
          >
            <Background gap={16} size={1} />
            <Controls showLock={false} showInteractive={false} />
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
    min-height: 300px;
    max-height: min(62vh, 520px);
    overflow: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-inset);
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
  }

  .sf-canvas {
    position: relative;
    box-sizing: border-box;
  }

  .sf-scroll :global(.svelte-flow) {
    width: 100% !important;
    height: 100% !important;
  }

  .sf-loading {
    margin: 0;
    padding: var(--space-6);
    text-align: center;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
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
</style>
