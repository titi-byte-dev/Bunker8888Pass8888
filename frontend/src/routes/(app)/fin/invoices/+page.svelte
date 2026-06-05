<script lang="ts">
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

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let docs = $state<InvoiceDocument[]>([]);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  const DEFAULT_RATE = 10;
  const DEFAULT_BENEFICIARY = "Vendedor";

  // Formulario de emissao.
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
      error = (e as Error).message;
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
    } catch (err) {
      error = (err as Error).message;
    } finally {
      busy = false;
    }
  }

  // CRM-004: converte uma pro-forma emitida na fatura definitiva (FT).
  async function convertToInvoice(d: InvoiceDocument) {
    busy = true;
    error = "";
    try {
      await issueInvoice("invoice", invoiceFromProforma(d));
      await refresh();
    } catch (err) {
      error = (err as Error).message;
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
      }
      await refresh();
    } catch (err) {
      error = (err as Error).message;
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
    } catch (err) {
      error = (err as Error).message;
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
          window.location.href = "/fin/commissions";
        }
      }
    } catch (err) {
      error = (err as Error).message;
    } finally {
      decidingId = null;
    }
  }

  async function handleReject(ev: AgentEvent) {
    decidingId = ev.id;
    try {
      await rejectSuggestion(ev.id);
      await refreshEvents();
    } catch (err) {
      error = (err as Error).message;
    } finally {
      decidingId = null;
    }
  }
</script>

<svelte:head><title>Faturacao — AegisPass</title></svelte:head>

<section class="page">
  <header class="head">
    <div>
      <h1>Faturacao</h1>
      <p class="muted">Pro-forma, faturas e recibos com numeracao legal (FIN-006).</p>
    </div>
    <a class="link" href="/fin">&larr; Custos</a>
  </header>

  {#if locked}
    <p class="lock">Cofre bloqueado — desbloqueia para gerir faturas.</p>
  {:else}
    {#if error}<p class="err">{error}</p>{/if}

    <form class="card" onsubmit={handleIssue}>
      <h2>Emitir documento</h2>
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
          <tr><th>Descricao</th><th>Qtd</th><th>Preco</th><th>IVA %</th><th></th></tr>
        </thead>
        <tbody>
          {#each fLines as line, i}
            <tr>
              <td><input bind:value={line.description} placeholder="Artigo / servico" /></td>
              <td><input type="number" min="0" step="0.01" bind:value={line.quantity} /></td>
              <td><input type="number" min="0" step="0.01" bind:value={line.unitPrice} /></td>
              <td><input type="number" min="0" step="1" bind:value={line.vatRate} /></td>
              <td>
                <button type="button" class="mini" onclick={() => removeLine(i)} disabled={fLines.length === 1}>x</button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
      <button type="button" class="mini" onclick={addLine}>+ linha</button>

      <div class="totals">
        <span>Liquido <strong>{money(draftTotals.net)}</strong></span>
        <span>IVA <strong>{money(draftTotals.vat)}</strong></span>
        <span>Total <strong>{money(draftTotals.gross)}</strong></span>
      </div>

      <button class="primary" type="submit" disabled={busy}>
        {busy ? "A emitir..." : "Emitir"}
      </button>
    </form>

    <div class="card">
      <h2>Actividade dos agentes</h2>
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos.</p>
      {:else}
        <ul class="ev-list">
          {#each agentEvents.slice(0, 4) as ev (ev.id)}
            <li class:suggested={isPendingSuggestion(ev)}>
              <span>{ev.label}</span>
              {#if isPendingSuggestion(ev) && ev.payload.action === "calculate_commission"}
                <span class="ev-actions">
                  <button class="mini" disabled={decidingId !== null} onclick={() => handleApprove(ev)}>Aprovar</button>
                  <button class="mini" disabled={decidingId !== null} onclick={() => handleReject(ev)}>Rejeitar</button>
                </span>
              {/if}
            </li>
          {/each}
        </ul>
      {/if}
    </div>

    <div class="card">
      <h2>Documentos</h2>
      {#if loading}
        <p class="muted">A carregar...</p>
      {:else if docs.length === 0}
        <p class="muted">Sem documentos emitidos.</p>
      {:else}
        <table class="list">
          <thead>
            <tr><th>Numero</th><th>Tipo</th><th>Cliente</th><th>Total</th><th>Estado</th><th></th></tr>
          </thead>
          <tbody>
            {#each docs as d}
              {@const t = invoiceTotals(d)}
              <tr>
                <td class="mono">{d.number}</td>
                <td>{docTypeLabel(d.docType)}</td>
                <td>{d.client.name}</td>
                <td>{money(t.gross, d.currency)}</td>
                <td><span class="badge {d.status}">{d.status}</span></td>
                <td class="actions">
                  {#if canConvertToInvoice(d)}
                    <button class="mini" onclick={() => convertToInvoice(d)} disabled={busy}>Converter em fatura</button>
                  {/if}
                  {#if d.status === "issued"}
                    <button class="mini" onclick={() => setStatus(d.id, "paid")} disabled={busy}>Marcar pago</button>
                    <button class="mini" onclick={() => setStatus(d.id, "void")} disabled={busy}>Anular</button>
                  {/if}
                  {#if canIssueReceipt(d)}
                    <button class="mini" onclick={() => issueReceipt(d)} disabled={busy}>Emitir recibo</button>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>
  {/if}
</section>

<style>
  .page { max-width: 960px; margin: 0 auto; padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
  .head { display: flex; justify-content: space-between; align-items: flex-start; }
  h1 { font-size: 1.25rem; margin: 0; }
  h2 { font-size: 0.95rem; margin: 0 0 0.75rem; }
  .muted { color: var(--text-muted, #888); font-size: 0.85rem; }
  .link { font-size: 0.85rem; color: var(--accent, #4f7cff); text-decoration: none; }
  .lock, .err { padding: 0.75rem; border-radius: 6px; font-size: 0.85rem; }
  .lock { background: #2a2a33; color: #ccc; }
  .err { background: #3a1f24; color: #ffb4bd; }
  .card { background: var(--surface, #1b1b22); border: 1px solid var(--border, #2c2c36); border-radius: 8px; padding: 1rem; }
  .row { display: flex; flex-wrap: wrap; gap: 0.75rem; margin-bottom: 0.75rem; }
  label { display: flex; flex-direction: column; gap: 0.25rem; font-size: 0.75rem; color: #aaa; }
  input, select { background: #111; border: 1px solid #333; border-radius: 5px; padding: 0.35rem 0.5rem; color: #eee; font-size: 0.85rem; }
  table { width: 100%; border-collapse: collapse; font-size: 0.82rem; }
  .lines th, .lines td { padding: 0.25rem; text-align: left; }
  .lines input { width: 100%; box-sizing: border-box; }
  .list th, .list td { padding: 0.4rem 0.5rem; border-bottom: 1px solid #262630; text-align: left; }
  .mono { font-family: ui-monospace, monospace; }
  .totals { display: flex; gap: 1.25rem; margin: 0.75rem 0; font-size: 0.85rem; color: #bbb; }
  .actions { display: flex; gap: 0.4rem; }
  .mini { background: #23232c; border: 1px solid #34343f; color: #ddd; border-radius: 5px; padding: 0.2rem 0.5rem; font-size: 0.75rem; cursor: pointer; }
  .primary { background: var(--accent, #4f7cff); color: #fff; border: none; border-radius: 6px; padding: 0.45rem 1rem; font-size: 0.85rem; cursor: pointer; margin-top: 0.5rem; }
  .primary:disabled { opacity: 0.6; cursor: default; }
  .badge { padding: 0.1rem 0.45rem; border-radius: 999px; font-size: 0.7rem; text-transform: uppercase; }
  .badge.issued { background: #2b3147; color: #9db4ff; }
  .badge.paid { background: #1f3a2a; color: #8ff0b4; }
  .badge.void { background: #3a2a2a; color: #d99; }
  .ev-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.4rem; }
  .ev-list li { display: flex; justify-content: space-between; align-items: center; padding: 0.4rem; border: 1px solid #2c2c36; border-radius: 5px; font-size: 0.82rem; }
  .ev-list li.suggested { border-color: var(--accent, #4f7cff); }
  .ev-actions { display: flex; gap: 0.3rem; }
</style>
