<script lang="ts">
  import { SvelteFlow, Background, Controls, type Node, type Edge } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import type { DocFlow, DocFlowStep } from "./types";
  import DocFlowNode from "./DocFlowNode.svelte";
  import FlowStepControls from "./FlowStepControls.svelte";
  import { prefersReducedMotion } from "$lib/motion/reduced";

  interface Props {
    flow: DocFlow;
  }

  let { flow }: Props = $props();

  const nodeTypes = { docFlow: DocFlowNode };
  const graph = $derived(flow.graph!);

  let stepIndex = $state(0);
  let playing = $state(false);

  const messageSteps = $derived(
    flow.steps.filter((s): s is Extract<DocFlowStep, { kind: "message" }> => s.kind === "message"),
  );

  const nodes = $derived.by(() => {
    const step = messageSteps[stepIndex];
    return graph.nodes.map(
      (n): Node => ({
        id: n.id,
        type: "docFlow",
        position: { x: n.x, y: 24 },
        data: {
          label: n.label,
          active: Boolean(step && (n.id === step.from || n.id === step.to)),
        },
        draggable: false,
        selectable: false,
        focusable: false,
      }),
    );
  });

  const edges = $derived.by(() => {
    return graph.edges.map(
      (e, i): Edge => ({
        id: e.id,
        source: e.source,
        target: e.target,
        label: e.label,
        type: "smoothstep",
        animated: i === stepIndex && !prefersReducedMotion(),
        style:
          i === stepIndex
            ? "stroke: var(--color-accent); stroke-width: 2.5"
            : i < stepIndex
              ? "stroke: var(--color-text-muted); stroke-width: 1.5; opacity: 0.55"
              : "stroke: var(--color-border); stroke-width: 1.5; opacity: 0.35",
        labelStyle:
          i === stepIndex
            ? "fill: var(--color-text); font-weight: 600; font-size: 11px"
            : "fill: var(--color-text-muted); font-size: 10px",
        selectable: false,
        focusable: false,
      }),
    );
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

<div class="sf-player">
  <div class="sf-canvas" style:min-height={`${Math.max(200, graph.edges.length * 28 + 120)}px`}>
    <SvelteFlow
      {nodes}
      {edges}
      {nodeTypes}
      fitView
      fitViewOptions={{ padding: 0.2 }}
      nodesDraggable={false}
      nodesConnectable={false}
      elementsSelectable={false}
      panOnDrag
      zoomOnScroll
      minZoom={0.4}
      maxZoom={1.4}
      proOptions={{ hideAttribution: true }}
    >
      <Background gap={16} size={1} />
      <Controls showInteractive={false} />
    </SvelteFlow>
  </div>

  <FlowStepControls
    {messageSteps}
    {stepIndex}
    {playing}
    onStep={goTo}
    onTogglePlay={togglePlay}
  />
</div>

<style>
  .sf-player {
    margin: var(--space-4) 0;
  }

  .sf-canvas {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
    background: var(--color-bg-base);
  }

  /* Tokens AegisPass sobre o tema default do Svelte Flow */
  .sf-canvas :global(.svelte-flow) {
    --xy-edge-stroke-default: var(--color-border);
    --xy-node-background-color-default: var(--color-bg-surface);
    --xy-background-color-default: var(--color-bg-base);
    --xy-background-pattern-color-default: color-mix(
      in srgb,
      var(--color-border) 40%,
      transparent
    );
  }

  .sf-canvas :global(.svelte-flow__controls) {
    box-shadow: none;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .sf-canvas :global(.svelte-flow__controls-button) {
    background: var(--color-bg-surface);
    border-bottom: 1px solid var(--color-border);
    fill: var(--color-text-muted);
  }

  .sf-canvas :global(.svelte-flow__controls-button:hover) {
    background: var(--color-accent-muted);
  }
</style>
