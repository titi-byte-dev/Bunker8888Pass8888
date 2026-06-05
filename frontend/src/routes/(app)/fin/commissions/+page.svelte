<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { listInvoices } from "$lib/fin/invoicesService";
  import {
    createCommission,
    listCommissions,
    updateCommissionStatus,
  } from "$lib/fin/commissionsService";
  import {
    commissionAmount,
    commissionFromInvoice,
    commissionStatusLabel,
    type CommissionDocument,
    type CommissionStatus,
  } from "$lib/fin/commissions";
  import { invoiceTotals, type InvoiceDocument } from "$lib/fin/invoices";
  import { reportCommissionRecorded } from "$lib/fin/erpFlowService";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";
  import {
    Button,
    DataTable,
    EmptyState,
    MetricCard,
    PageShell,
    Panel,
    StatusBanner,
    toast,
    type DataColumn,
  } from "$lib/ui";

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let commissions = $state<CommissionDocument[]>([]);
  let paidInvoices = $state<InvoiceDocument[]>([]);

  let fInvoiceId = $state("");
  let fBeneficiary = $state("");
  let fRate = $state(10);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);

  const selectedInvoice = $derived(paidInvoices.find((i) => i.id === fInvoiceId));
  const preview = $derived(
    selectedInvoice
      ? commissionAmount(commissionFromInvoice(selectedInvoice, fRate, fBeneficiary || "—"))
      : 0,
  );

  const totalPending = $derived(commissions.filter((c) => c.status === "pending").length);
  const totalPaid = $derived(commissions.filter((c) => c.status === "paid").length);

  const commissionColumns = $derived<DataColumn<CommissionDocument>[]>([
    { id: "beneficiary", label: "Beneficiário", accessor: (c) => c.beneficiary },
    { id: "invoice", label: "Fatura", mono: true, accessor: (c) => c.invoiceNumber ?? "—" },
    { id: "rate", label: "Taxa", accessor: (c) => `${c.ratePct}%` },
    {
      id: "amount",
      label: "Valor",
      align: "right",
      mono: true,
      accessor: (c) => money(commissionAmount(c), c.currency),
    },
    { id: "status", label: "Estado", accessor: (c) => commissionStatusLabel(c.status) },
  ]);

  function money(n: number, currency = "EUR"): string {
    try {
      return new Intl.NumberFormat("pt-PT", { style: "currency", currency }).format(n);
    } catch {
      return `${n.toFixed(2)} ${currency}`;
    }
  }

  async function refresh() {
    loading = true;
    error = "";
    try {
      const [coms, invs] = await Promise.all([listCommissions(), listInvoices()]);
      commissions = coms;
      paidInvoices = invs.filter((i) => i.docType === "invoice" && i.status === "paid");
      try {
        agentEvents = await listAgentEvents();
      } catch {
        agentEvents = [];
      }
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    if (!getMasterKey()) {
      locked = true;
      loading = false;
      return;
    }
    void refresh();
  });

  async function handleCreate(e: Event) {
    e.preventDefault();
    if (!selectedInvoice) {
      error = "Escolhe uma fatura paga.";
      return;
    }
    if (!fBeneficiary.trim()) {
      error = "Indica o beneficiário.";
      return;
    }
    busy = true;
    error = "";
    try {
      const payload = commissionFromInvoice(selectedInvoice, fRate, fBeneficiary);
      await createCommission(selectedInvoice.id, payload);
      await reportCommissionRecorded(selectedInvoice.id);
      agentEvents = await listAgentEvents();
      fBeneficiary = "";
      fRate = 10;
      fInvoiceId = "";
      await refresh();
      toast.success("Comissão registada.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao registar";
    } finally {
      busy = false;
    }
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
  }

  async function handleApprove(ev: AgentEvent) {
    decidingId = ev.id;
    try {
      await approveSuggestion(ev.id);
      if (ev.payload.action === "generate_rgpd_report") {
        toast.success("Relatório RGPD — a abrir conformidade.");
        await goto("/hr/compliance");
      }
      agentEvents = await listAgentEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao aprovar";
    } finally {
      decidingId = null;
    }
  }

  async function handleReject(ev: AgentEvent) {
    decidingId = ev.id;
    try {
      await rejectSuggestion(ev.id);
      agentEvents = await listAgentEvents();
      toast.info("Sugestão rejeitada.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao rejeitar";
    } finally {
      decidingId = null;
    }
  }

  async function setStatus(id: string, status: CommissionStatus) {
    busy = true;
    error = "";
    try {
      await updateCommissionStatus(id, status);
      await refresh();
      toast.success(status === "paid" ? "Comissão liquidada." : "Comissão anulada.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao actualizar";
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head><title>Comissões — AegisPass</title></svelte:head>

<PageShell
  title="Comissões"
  taskId="FIN-007"
  description="Comissões de vendas sobre faturas pagas — calculadas e cifradas no cliente."
 
>
  {#snippet actions()}
    <Button variant="ghost" size="sm" href="/fin/invoices">← Faturas</Button>
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para gerir comissões.">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else}
    {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

    {#if !loading && commissions.length > 0}
      <section class="metrics">
        <MetricCard label="registadas" value={String(commissions.length)} />
        <MetricCard label="pendentes" value={String(totalPending)} />
        <MetricCard label="liquidadas" value={String(totalPaid)} variant="success" />
      </section>
    {/if}

    <Panel title="Actividade dos agentes">
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos.</p>
      {:else}
        <ul class="ev-list">
          {#each agentEvents.slice(0, 4) as ev (ev.id)}
            <li class:suggested={isPendingSuggestion(ev)}>
              <span>{ev.label}</span>
              {#if isPendingSuggestion(ev) && ev.payload.action === "generate_rgpd_report"}
                <span class="ev-actions">
                  <Button variant="primary" size="sm" disabled={decidingId !== null} onclick={() => handleApprove(ev)}>Aprovar</Button>
                  <Button variant="ghost" size="sm" disabled={decidingId !== null} onclick={() => handleReject(ev)}>Rejeitar</Button>
                </span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </Panel>

    <Panel title="Gerar comissão">
      {#if paidInvoices.length === 0}
        <p class="muted">Sem faturas pagas. Marca uma fatura como paga primeiro.</p>
      {:else}
        <form class="form" onsubmit={handleCreate}>
          <div class="row">
            <label>
              Fatura paga
              <select bind:value={fInvoiceId}>
                <option value="">— escolher —</option>
                {#each paidInvoices as inv}
                  <option value={inv.id}>{inv.number} · {money(invoiceTotals(inv).net, inv.currency)} líq.</option>
                {/each}
              </select>
            </label>
            <label>Beneficiário<input bind:value={fBeneficiary} placeholder="Nome do vendedor" /></label>
            <label>Taxa %<input type="number" min="0" step="0.5" bind:value={fRate} /></label>
          </div>
          <p class="preview">Comissão estimada <strong>{money(preview, selectedInvoice?.currency)}</strong></p>
          <Button type="submit" disabled={busy || !fInvoiceId} loading={busy}>Registar comissão</Button>
        </form>
      {/if}
    </Panel>

    <Panel title="Registo" padding="none">
      <DataTable
        columns={commissionColumns}
        rows={commissions}
        keyFn={(c) => c.id}
        loading={loading}
        dense
        emptyTitle="Sem comissões registadas"
      >
        {#snippet actions(c)}
          {#if c.status === "pending"}
            <Button variant="ghost" size="sm" onclick={() => setStatus(c.id, "paid")} disabled={busy}>Liquidar</Button>
            <Button variant="ghost" size="sm" onclick={() => setStatus(c.id, "void")} disabled={busy}>Anular</Button>
          {/if}
        {/snippet}
      </DataTable>
    </Panel>
  {/if}
</PageShell>

<style>
  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
    gap: var(--space-3);
  }

  .muted { margin: 0; color: var(--color-text-muted); font-size: var(--text-sm); }

  .row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-3);
    margin-bottom: var(--space-3);
  }

  label {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-xs);
    color: var(--color-text-label);
  }

  input, select {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-size: var(--text-sm);
    min-width: 10rem;
  }

  .preview {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .ev-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .ev-list li {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
  }

  .ev-list li.suggested { border-color: var(--color-accent); }
  .ev-actions { display: flex; gap: var(--space-2); }
</style>
