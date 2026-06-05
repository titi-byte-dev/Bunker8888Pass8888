<script lang="ts">
  import { BaseEdge, type Edge, type EdgeProps } from "@xyflow/svelte";
  import type { DocSequenceEdgeData } from "./types";

  type SeqEdge = Edge<DocSequenceEdgeData, "docSequence">;

  let {
    id,
    label,
    data,
    style,
    markerEnd,
    markerStart,
    interactionWidth = 20,
    class: className,
  }: EdgeProps<SeqEdge> = $props();

  const d = $derived(data!);

  const path = $derived.by(() => {
    const rowY = d.rowY;
    const sx = d.sourceCenterX;
    const tx = d.targetCenterX;

    if (d.selfLoop) {
      const w = d.selfLoopWidth ?? 52;
      const drop = Math.min(32, d.rowHeight * 0.5);
      const right = sx + w;
      return `M ${sx} ${rowY} H ${right} V ${rowY + drop} H ${sx}`;
    }

    const left = Math.min(sx, tx);
    const right = Math.max(sx, tx);
    return `M ${left} ${rowY} H ${right}`;
  });

  const labelX = $derived(
    d.selfLoop
      ? d.sourceCenterX + (d.selfLoopWidth ?? 52) / 2
      : (d.sourceCenterX + d.targetCenterX) / 2,
  );
  const labelY = $derived(d.selfLoop ? d.rowY - 10 : d.rowY - 12);

  const bandY = $derived(d.rowY - d.rowHeight / 2 + 6);
  const bandH = $derived(d.rowHeight - 12);
</script>

<!-- Faixa de realce na fila do passo actual -->
{#if d.state === "current"}
  <rect
    class="row-band row-band-current"
    x="0"
    y={bandY}
    width={d.canvasWidth}
    height={bandH}
    rx="4"
  />
{:else if d.state === "done"}
  <rect
    class="row-band row-band-done"
    x="0"
    y={bandY}
    width={d.canvasWidth}
    height={bandH}
    rx="4"
  />
{/if}

<BaseEdge
  {id}
  {path}
  {labelX}
  {labelY}
  {label}
  {style}
  {markerStart}
  {markerEnd}
  {interactionWidth}
  class={className}
/>

<style>
  .row-band {
    pointer-events: none;
  }

  .row-band-current {
    fill: color-mix(in srgb, var(--color-accent) 14%, transparent);
    stroke: color-mix(in srgb, var(--color-accent) 35%, transparent);
    stroke-width: 1;
  }

  .row-band-done {
    fill: color-mix(in srgb, var(--color-success-fg) 8%, transparent);
  }
</style>
