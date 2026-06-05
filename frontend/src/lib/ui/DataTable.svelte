<script lang="ts" generics="T">
  import type { Snippet } from "svelte";
  import EmptyState from "./EmptyState.svelte";
  import Skeleton from "./Skeleton.svelte";

  import type { DataColumn } from "./data-table";

  interface Props {
    columns: DataColumn<T>[];
    rows: T[];
    keyFn: (row: T) => string;
    loading?: boolean;
    dense?: boolean;
    emptyTitle?: string;
    emptyDescription?: string;
    /** Classe extra por linha (ex. inactiva) */
    rowClass?: (row: T) => string | undefined;
    actions?: Snippet<[T]>;
  }

  let {
    columns,
    rows,
    keyFn,
    loading = false,
    dense = false,
    emptyTitle = "Sem dados",
    emptyDescription,
    rowClass,
    actions,
  }: Props = $props();

  function cellText(row: T, col: DataColumn<T>): string {
    if (col.accessor) return col.accessor(row);
    const v = (row as Record<string, unknown>)[col.id];
    return v == null ? "—" : String(v);
  }
</script>

<div class="table-wrap" class:dense>
  {#if loading}
    <Skeleton variant="table" rows={dense ? 6 : 4} cols={columns.length} />
  {:else if rows.length === 0}
    <EmptyState title={emptyTitle} description={emptyDescription} />
  {:else}
    <table>
      <thead>
        <tr>
          {#each columns as col (col.id)}
            <th class:align-right={col.align === "right"}>{col.label}</th>
          {/each}
          {#if actions}<th class="actions-col"><span class="sr-only">Acções</span></th>{/if}
        </tr>
      </thead>
      <tbody>
        {#each rows as row (keyFn(row))}
          <tr class={rowClass?.(row) ?? undefined}>
            {#each columns as col (col.id)}
              <td
                class:align-right={col.align === "right"}
                class:mono={col.mono}
                class:muted={col.muted}
              >
                {cellText(row, col)}
              </td>
            {/each}
            {#if actions}
              <td class="actions">
                {@render actions(row)}
              </td>
            {/if}
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</div>

<style>
  .table-wrap {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
    background: var(--color-bg-surface);
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--text-sm);
  }

  th {
    text-align: left;
    padding: var(--space-2) var(--space-3);
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-label);
    background: var(--color-bg-elevated);
    border-bottom: 1px solid var(--color-border);
  }

  td {
    padding: var(--space-3);
    border-bottom: 1px solid var(--color-border);
    vertical-align: middle;
  }

  tr:last-child td {
    border-bottom: none;
  }

  .dense th,
  .dense td {
    padding: var(--space-2) var(--space-3);
  }

  .align-right {
    text-align: right;
  }

  .mono {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }

  .actions {
    text-align: right;
    white-space: nowrap;
  }

  .actions-col {
    width: 1%;
  }

  :global(tr.off) td {
    opacity: 0.55;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }
</style>
