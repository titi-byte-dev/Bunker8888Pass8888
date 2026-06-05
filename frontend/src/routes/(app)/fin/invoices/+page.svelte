<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import {
    issueInvoice,
    listInvoices,
    updateInvoiceStatus,
  } from "$lib/fin/invoicesService";
  import {
    DOC_TYPES,
    docTypeLabel,
    invoiceTotals,
    type DocStatus,
    type DocType,
    type InvoiceDocument,
    type InvoiceLine,
  } from "$lib/fin/invoices";
  import { canConvertToInvoice, invoiceFromProforma } from "$lib/crm/proformaToInvoice";
  import { canIssueReceipt, receiptFromInvoice } from "$lib/fin/receiptFromInvoice";
  import { reportInvoicePaid } from "$lib/fin/erpFlowService";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";
  import { commissionFromInvoice } from "$lib/fin/commissions";
  import { createCommission } from "$lib/fin/commissionsService";
  import {
    Button,
    DataTable,
    EmptyState,
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

  let docs = $state<InvoiceDocument[]>([]);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  const DEFAULT_RATE = 10;
  const DEFAULT_BENEFICIARY = "Vendedor";

  let fDocType = $state<DocType>("invoice");
  let fClientName = $state("");
  let fClientTaxId = $state("");
  let fClientEmail = $state("");
  let fCurrency = $state("EUR");
  let fLines = $state<InvoiceLine[]>([
    { description: "", quantity: 1, unitPrice: 0, vatRate: 23 },
  ]);

  const draftTotals = $derived(
    invoiceTotals({
      client: { name: fClientName },
      lines: fLines,
      currency: fCurrency,
      issuedAt: new Date().toISOString(),
    }),
  );

  const docColumns = $derived<DataColumn<InvoiceDocument>[]>([
    { id: "number", label: "Número", mono: true, accessor: (d) => d.number },
    { id: "type", label: "Tipo", accessor: (d) => docTypeLabel(d.docType) },
    { id: "client", label: "Cliente", accessor: (d) => d.client.name },
    {
      id: "total",
      label: "Total",
      align: "right",
      mono: true,
      accessor: (d) => money(invoiceTotals(d).gross, d.currency),
    },
    { id: "status", label: "Estado", accessor: (d) => d.status },
  ]);

  function money(n: number, currency = fCurrency): string {
    try {
      return new Intl.NumberFormat("pt-PT", { style: "currency", currency }).format(n);
    } catch {
      return `${n.toFixed(2)} ${currency}`;
    }
  }

  function addLine() {
    fLines = [...fLines, { description: "", quantity: 1, unitPrice: 0, vatRate: 23 }];
  }

  function removeLine(i: number) {
    fLines = fLines.filter((_, idx) => idx !== i);
  }

  async function refreshEvents() {
    try {
      agentEvents = await listAgentEvents();
    } catch {
      agentEvents = [];
    }
  }

  async function refresh() {
    loading = true;
    error = "";
    try {
      docs = await listInvoices();
      await refreshEvents();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      loading = false;
    }
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
  }

  onMount(() => {
    if (!getMasterKey()) {
      locked = true;
      loading = false;
      return;
    }
    void refresh();
  });

  async function handleIssue(e: Event) {
    e.preventDefault();
    if (!fClientName.trim()) {
      error = "Indica o nome do cliente.";
      return;
    }
    busy = true;
    error = "";
    try {
      await issueInvoice(fDocType, {
        client: {
          name: fClientName,
          taxId: fClientTaxId || undefined,
          email: fClientEmail || undefined,
        },
        lines: fLines,
        currency: fCurrency,
        issuedAt: new Date().toISOString(),
      });
      fClientName = "";
      fClientTaxId = "";
      fClientEmail = "";
      fLines = [{ description: "", quantity: 1, unitPrice: 0, vatRate: 23 }];
      await refresh();
      toast.success("Documento emitido.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao emitir";
    } finally {
      busy = false;
    }
  }

  async function convertToInvoice(d: InvoiceDocument) {
    busy = true;
    error = "";
    try {
      await issueInvoice("invoice", invoiceFromProforma(d));
      await refresh();
      toast.success("Pro-forma convertida em fatura.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha na conversão";
    } finally {
      busy = false;
    }
  }

  async function setStatus(id: string, status: DocStatus) {
    busy = true;
    error = "";
    try {
      const doc = docs.find((d) => d.id === id);
      await updateInvoiceStatus(id, status);
      if (status === "paid" && doc?.docType === "invoice") {
        await reportInvoicePaid(doc.id, doc.number);
        await refreshEvents();
        toast.success("Fatura marcada como paga.");
      } else if (status === "void") {
        toast.info("Documento anulado.");
      }
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao actualizar";
    } finally {
      busy = false;
    }
  }

  async function issueReceipt(d: InvoiceDocument) {
    busy = true;
    error = "";
    try {
      await issueInvoice("receipt", receiptFromInvoice(d));
      await refresh();
      toast.success("Recibo emitido.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao emitir recibo";
    } finally {
      busy = false;
    }
  }

  async function handleApprove(ev: AgentEvent) {
    decidingId = ev.id;
    try {
      const result = await approveSuggestion(ev.id);
      await refreshEvents();
      if (result.action === "calculate_commission") {
        const invId = String(ev.payload.invoice_id ?? "");
        const doc = docs.find((d) => d.id === invId);
        if (doc) {
          await createCommission(invId, commissionFromInvoice(doc, DEFAULT_RATE, DEFAULT_BENEFICIARY));
          toast.success("Comissão calculada — a abrir registo.");
          await goto("/fin/commissions");
        }
      }
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
      await refreshEvents();
      toast.info("Sugestão rejeitada.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao rejeitar";
    } finally {
      decidingId = null;
    }
  }
</script>

<svelte:head><title>Faturação — AegisPass</title></svelte:head>

<PageShell
  title="Faturação"
  taskId="FIN-006"
  description="Pro-forma, faturas e recibos com numeração legal — emitidos e cifrados no cliente."
 
>
  {#snippet actions()}
    <Button variant="ghost" size="sm" href="/fin/costs">← Custos</Button>
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para gerir faturas.">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else}
    {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

    <Panel title="Emitir documento">
      <form class="issue-form" onsubmit={handleIssue}>
        <div class="row">
          <label>
            Tipo
            <select bind:value={fDocType}>
              {#each DOC_TYPES as t}<option value={t.id}>{t.label}</option>{/each}
            </select>
          </label>
          <label>Cliente<input bind:value={fClientName} placeholder="Nome / empresa" /></label>
          <label>NIF<input bind:value={fClientTaxId} placeholder="opcional" /></label>
          <label>Email<input bind:value={fClientEmail} placeholder="opcional" /></label>
          <label>Moeda<input bind:value={fCurrency} maxlength="3" /></label>
        </div>

        <table class="lines">
          <thead>
            <tr><th>Descrição</th><th>Qtd</th><th>Preço</th><th>IVA %</th><th></th></tr>
          </thead>
          <tbody>
            {#each fLines as line, i}
              <tr>
                <td><input bind:value={line.description} placeholder="Artigo / serviço" /></td>
                <td><input type="number" min="0" step="0.01" bind:value={line.quantity} /></td>
                <td><input type="number" min="0" step="0.01" bind:value={line.unitPrice} /></td>
                <td><input type="number" min="0" step="1" bind:value={line.vatRate} /></td>
                <td>
                  <Button variant="ghost" size="sm" onclick={() => removeLine(i)} disabled={fLines.length === 1}>×</Button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
        <Button variant="ghost" size="sm" type="button" onclick={addLine}>+ linha</Button>

        <div class="totals">
          <span>Líquido <strong>{money(draftTotals.net)}</strong></span>
          <span>IVA <strong>{money(draftTotals.vat)}</strong></span>
          <span>Total <strong>{money(draftTotals.gross)}</strong></span>
        </div>

        <Button type="submit" loading={busy} disabled={busy}>Emitir</Button>
      </form>
    </Panel>

    <Panel title="Actividade dos agentes">
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos.</p>
      {:else}
        <ul class="ev-list">
          {#each agentEvents.slice(0, 4) as ev (ev.id)}
            <li class:suggested={isPendingSuggestion(ev)}>
              <span>{ev.label}</span>
              {#if isPendingSuggestion(ev) && ev.payload.action === "calculate_commission"}
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

    <Panel title="Documentos" padding="none">
      <DataTable
        columns={docColumns}
        rows={docs}
        keyFn={(d) => d.id}
        loading={loading}
        dense
        emptyTitle="Sem documentos emitidos"
      >
        {#snippet actions(d)}
          <div class="doc-actions">
            {#if canConvertToInvoice(d)}
              <Button variant="ghost" size="sm" onclick={() => convertToInvoice(d)} disabled={busy}>Converter</Button>
            {/if}
            {#if d.status === "issued"}
              <Button variant="ghost" size="sm" onclick={() => setStatus(d.id, "paid")} disabled={busy}>Pago</Button>
              <Button variant="ghost" size="sm" onclick={() => setStatus(d.id, "void")} disabled={busy}>Anular</Button>
            {/if}
            {#if canIssueReceipt(d)}
              <Button variant="ghost" size="sm" onclick={() => issueReceipt(d)} disabled={busy}>Recibo</Button>
            {/if}
          </div>
        {/snippet}
      </DataTable>
    </Panel>
  {/if}
</PageShell>

<style>
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
  }

  .lines {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--text-sm);
    margin-bottom: var(--space-2);
  }

  .lines th, .lines td { padding: var(--space-1); text-align: left; }
  .lines input { width: 100%; box-sizing: border-box; }

  .totals {
    display: flex;
    gap: var(--space-6);
    margin: var(--space-3) 0;
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
  .ev-actions { display: flex; gap: var(--space-2); flex-shrink: 0; }

  .doc-actions {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-1);
    justify-content: flex-end;
  }
</style>
