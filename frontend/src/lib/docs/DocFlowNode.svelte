<script lang="ts">
  import { Handle, Position, type Node, type NodeProps } from "@xyflow/svelte";
  import type { DocFlowNodeData } from "./types";

  type DocFlowNodeType = Node<DocFlowNodeData, "docFlow">;

  let { data }: NodeProps<DocFlowNodeType> = $props();
</script>

<Handle type="target" position={Position.Left} id="in" />
<Handle type="source" position={Position.Right} id="out" />

<div class="participant">
  <div
    class="doc-flow-node"
    class:role-from={data.role === "from"}
    class:role-to={data.role === "to"}
    class:role-both={data.role === "both"}
    class:inactive={data.role === null}
  >
    <span class="node-label">{data.label}</span>
  </div>
  {#if data.lifelineHeight}
    <div
      class="lifeline"
      class:lit={data.role !== null && data.role !== undefined}
      style:height="{data.lifelineHeight}px"
      aria-hidden="true"
    ></div>
  {/if}
</div>

<style>
  .participant {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .lifeline {
    width: 2px;
    flex-shrink: 0;
    margin-top: var(--space-1);
    border-radius: 1px;
    background: color-mix(in srgb, var(--color-border) 65%, transparent);
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .lifeline.lit {
    background: color-mix(in srgb, var(--color-accent) 55%, var(--color-border));
  }

  .doc-flow-node {
    padding: var(--space-2) var(--space-3);
    min-width: 5.5rem;
    max-width: 14rem;
    border: 2px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    font-size: var(--text-xs);
    font-weight: 500;
    text-align: center;
    line-height: 1.3;
    transition:
      border-color var(--duration-fast) var(--ease-out),
      box-shadow var(--duration-fast) var(--ease-out),
      background-color var(--duration-fast) var(--ease-out);
  }

  .doc-flow-node.inactive {
    opacity: 0.72;
  }

  /* Origem da mensagem — âmbar */
  .doc-flow-node.role-from {
    border-color: var(--color-warning);
    background: color-mix(in srgb, var(--color-warning) 12%, var(--color-bg-surface));
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-warning) 22%, transparent);
    opacity: 1;
  }

  /* Destino da mensagem — verde */
  .doc-flow-node.role-to {
    border-color: var(--color-success-fg);
    background: var(--color-success-bg);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-success-fg) 22%, transparent);
    opacity: 1;
  }

  /* Self-loop — accent */
  .doc-flow-node.role-both {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-accent) 28%, transparent);
    opacity: 1;
  }

  .node-label {
    display: block;
    word-break: break-word;
  }

  @media (prefers-reduced-motion: reduce) {
    .doc-flow-node,
    .lifeline {
      transition: none;
    }
  }
</style>
