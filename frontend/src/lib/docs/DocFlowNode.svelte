<script lang="ts">
  import type { NodeProps } from "@xyflow/svelte";

  type DocFlowNodeData = {
    label: string;
    active?: boolean;
  };

  let { data }: NodeProps<DocFlowNodeData> = $props();
</script>

<div class="doc-flow-node" class:active={data.active}>
  <span class="node-label">{data.label}</span>
</div>

<style>
  .doc-flow-node {
    padding: var(--space-2) var(--space-3);
    min-width: 5.5rem;
    max-width: 8rem;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    font-size: var(--text-xs);
    font-weight: 500;
    text-align: center;
    line-height: 1.3;
    transition:
      border-color var(--duration-fast) var(--ease-out),
      box-shadow var(--duration-fast) var(--ease-out);
  }

  .doc-flow-node.active {
    border-color: var(--color-accent);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--color-accent) 25%, transparent);
    background: var(--color-accent-muted);
  }

  .node-label {
    display: block;
    word-break: break-word;
  }

  @media (prefers-reduced-motion: reduce) {
    .doc-flow-node {
      transition: none;
    }
  }
</style>
