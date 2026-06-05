<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    adjustInventory,
    createInventoryItem,
    deleteInventoryItem,
    listInventory,
    type InventoryItem,
  } from "$lib/ops/inventory";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";
  import { buildPurchaseOrderDraft, type PurchaseOrderDraft } from "$lib/ops/purchaseOrder";
  import {
    Button,
    confirmDialog,
    PageShell,
    Panel,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";

  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let items = $state<InventoryItem[]>([]);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  let drafts = $state<PurchaseOrderDraft[]>([]);

  let fName = $state("");
  let fSku = $state("");
  let fQty = $state(10);
  let fReorder = $state(5);
  let fUnit = $state("un");

  onMount(() => {
    void refresh();
  });

  async function refresh() {
    loading = true;
    error = "";
    try {
      items = await listInventory();
      agentEvents = await listAgentEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao carregar inventário";
    } finally {
      loading = false;
    }
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
  }

  async function handleApproveSuggestion(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
    try {
      const result = await approveSuggestion(ev.id);
      await refresh();
      if (result.action === "create_purchase_order") {
        drafts = [buildPurchaseOrderDraft(ev.payload), ...drafts];
      }
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao aprovar";
    } finally {
      decidingId = null;
    }
  }

  async function handleRejectSuggestion(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
    try {
      await rejectSuggestion(ev.id);
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao rejeitar";
    } finally {
      decidingId = null;
    }
  }

  async function handleCreate(e: SubmitEvent) {
    e.preventDefault();
    if (!fName.trim()) return;
    busy = true;
    error = "";
    try {
      await createInventoryItem({
        name: fName.trim(),
        sku: fSku.trim(),
        quantity: fQty,
        reorderLevel: fReorder,
        unit: fUnit.trim() || "un",
      });
      fName = "";
      fSku = "";
      fQty = 10;
      fReorder = 5;
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao criar artigo";
    } finally {
      busy = false;
    }
  }

  async function handleAdjust(item: InventoryItem, delta: number) {
    busy = true;
    error = "";
    try {
      await adjustInventory(item.id, delta);
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao ajustar stock";
    } finally {
      busy = false;
    }
  }

  async function handleDelete(item: InventoryItem) {
    const ok = await confirmDialog({
      title: "Apagar artigo?",
      message: `Remove «${item.name}» do inventário.`,
      confirmLabel: "Apagar",
      variant: "danger",
    });
    if (!ok) return;
    busy = true;
    try {
      await deleteInventoryItem(item.id);
      await refresh();
      toast.success("Artigo apagado.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao apagar";
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Inventário — AegisPass</title>
</svelte:head>

<PageShell
  title="Inventário"
  taskId="AGENT-008"
  description="Gere stock operacional. Quando a quantidade desce ao nível de reordenação, o orquestrador sugere uma ordem de compra — aprova antes de gerar o rascunho."
 
>
  {#snippet actions()}
    <DocHelpLink slug="journey-ops-agent-inventory" label="Como funciona o agente de operações?" />
    <Button variant="ghost" size="sm" href="/work">← Trabalho</Button>
  {/snippet}

  {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

  <Panel title="Actividade dos agentes">
    {#if agentEvents.length === 0}
      <p class="muted">Sem eventos recentes.</p>
    {:else}
      <ul class="event-list">
        {#each agentEvents.slice(0, 6) as ev (ev.id)}
          <li class:suggested={isPendingSuggestion(ev)}>
            <div class="ev-body">
              <span class="ev-label">{ev.label}</span>
              {#if isPendingSuggestion(ev)}
                <div class="ev-actions">
                  <button
                    type="button"
                    class="btn approve"
                    disabled={decidingId !== null || busy}
                    onclick={() => handleApproveSuggestion(ev)}
                  >
                    {decidingId === ev.id ? "…" : "Aprovar"}
                  </button>
                  <button
                    type="button"
                    class="btn reject"
                    disabled={decidingId !== null}
                    onclick={() => handleRejectSuggestion(ev)}
                  >
                    Rejeitar
                  </button>
                </div>
              {/if}
            </div>
            <span class="ev-meta">{new Date(ev.createdAt).toLocaleString("pt-PT")}</span>
          </li>
        {/each}
      </ul>
    {/if}
  </Panel>

  {#if drafts.length > 0}
    <section class="panel drafts">
      <h2>Rascunhos de ordem de compra</h2>
      <ul>
        {#each drafts as d (d.createdAt + d.itemId)}
          <li>
            <strong>{d.itemName}</strong>
            {#if d.sku}<span class="muted"> ({d.sku})</span>{/if}
            — pedir <strong>{d.orderQty} {d.unit}</strong>
            <p class="muted sm">{d.note}</p>
          </li>
        {/each}
      </ul>
    </section>
  {/if}

  <section class="panel">
    <h2>Novo artigo</h2>
    <form class="grid" onsubmit={handleCreate}>
      <label>
        <span>Nome</span>
        <input bind:value={fName} placeholder="Toner HP" disabled={busy} />
      </label>
      <label>
        <span>SKU</span>
        <input bind:value={fSku} placeholder="TN-01" disabled={busy} />
      </label>
      <label>
        <span>Quantidade</span>
        <input type="number" min="0" bind:value={fQty} disabled={busy} />
      </label>
      <label>
        <span>Reordenar em</span>
        <input type="number" min="0" bind:value={fReorder} disabled={busy} />
      </label>
      <label>
        <span>Unidade</span>
        <input bind:value={fUnit} placeholder="un" disabled={busy} />
      </label>
      <button type="submit" class="btn primary" disabled={busy || !fName.trim()}>Adicionar</button>
    </form>
  </section>

  <section class="panel">
    <h2>Stock</h2>
    {#if loading}
      <Skeleton variant="row" />
      <Skeleton variant="row" />
    {:else if items.length === 0}
      <p class="muted">Sem artigos — adiciona o primeiro acima.</p>
    {:else}
      <ul class="stock-list">
        {#each items as item (item.id)}
          <li class:low={item.lowStock}>
            <div class="stock-main">
              <strong>{item.name}</strong>
              {#if item.sku}<span class="muted sm"> {item.sku}</span>{/if}
              <span class="qty">{item.quantity} {item.unit}</span>
              {#if item.lowStock}<span class="badge">Stock baixo</span>{/if}
            </div>
            <div class="stock-actions">
              <button type="button" class="btn sm" disabled={busy} onclick={() => handleAdjust(item, -1)}>−1</button>
              <button type="button" class="btn sm" disabled={busy} onclick={() => handleAdjust(item, 1)}>+1</button>
              <button type="button" class="btn sm danger" disabled={busy} onclick={() => handleDelete(item)}>Apagar</button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</PageShell>

<style>
  h2 { margin: 0 0 var(--space-3); font-size: var(--text-base); }
  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-5) var(--space-6);
    margin-bottom: var(--space-4);
  }
  .muted { color: var(--color-text-muted); font-size: var(--text-sm); }
  .sm { font-size: var(--text-xs); }
  .grid {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));
  }
  label > span {
    display: block;
    margin-bottom: var(--space-1);
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--color-text-label);
  }
  input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    box-sizing: border-box;
  }
  .btn {
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-size: var(--text-sm);
    cursor: pointer;
    background: var(--color-bg-elevated);
    color: var(--color-text);
  }
  .btn.primary { background: var(--color-accent); color: var(--color-accent-fg); border-color: transparent; }
  .btn.sm { padding: var(--space-1) var(--space-2); font-size: var(--text-xs); }
  .btn.danger { color: var(--color-danger); }
  .btn.approve { background: var(--color-success-bg); color: var(--color-success-fg); border-color: transparent; }
  .btn:disabled { opacity: 0.55; cursor: progress; }
  .event-list, .stock-list, .drafts ul { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: var(--space-2); }
  .event-list li, .stock-list li {
    display: flex;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    font-size: var(--text-sm);
  }
  .event-list li.suggested { border-color: var(--color-accent); }
  .stock-list li.low { border-color: var(--color-warning, var(--color-accent)); }
  .stock-main { display: flex; flex-wrap: wrap; align-items: center; gap: var(--space-2); }
  .qty { font-family: var(--font-mono); }
  .badge {
    font-size: var(--text-xs);
    padding: 0 var(--space-2);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }
  .stock-actions { display: flex; gap: var(--space-1); flex-shrink: 0; }
  .ev-body { display: flex; flex-direction: column; gap: var(--space-2); }
  .ev-actions { display: flex; gap: var(--space-2); }
  .ev-meta { font-size: var(--text-xs); color: var(--color-text-muted); white-space: nowrap; }
</style>
