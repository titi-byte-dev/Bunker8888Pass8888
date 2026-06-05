<script lang="ts">
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

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let commissions = $state<CommissionDocument[]>([]);
  let paidInvoices = $state<InvoiceDocument[]>([]);

  // Formulario: gerar comissao a partir de uma fatura paga.
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
      // So faturas (FT) pagas geram comissao.
      paidInvoices = invs.filter((i) => i.docType === "invoice" && i.status === "paid");
      try {
        agentEvents = await listAgentEvents();
      } catch {
        agentEvents = [];
      }
    } catch (e) {
      error = (e as Error).message;
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
      error = "Indica o beneficiario.";
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
    } catch (err) {
      error = (err as Error).message;
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
        window.location.href = "/hr/compliance";
      }
      agentEvents = await listAgentEvents();
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
      agentEvents = await listAgentEvents();
    } catch (err) {
      error = (err as Error).message;
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
    } catch (err) {
      error = (err as Error).message;
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head><title>Comissoes — AegisPass</title></svelte:head>

<section class="page">
  <header class="head">
    <div>
      <h1>Comissoes</h1>
      <p class="muted">Comissoes de vendas sobre faturas pagas (FIN-007).</p>
    </div>
    <a class="link" href="/fin/invoices">&larr; Faturas</a>
  </header>

  {#if locked}
    <p class="lock">Cofre bloqueado — desbloqueia para gerir comissoes.</p>
  {:else}
    {#if error}<p class="err">{error}</p>{/if}

    <div class="card">
      <h2>Actividade dos agentes</h2>
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos.</p>
      {:else}
        <ul class="ev-list">
          {#each agentEvents.slice(0, 4) as ev (ev.id)}
            <li class:suggested={isPendingSuggestion(ev)}>
              <span>{ev.label}</span>
              {#if isPendingSuggestion(ev) && ev.payload.action === "generate_rgpd_report"}
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

    <form class="card" onsubmit={handleCreate}>
      <h2>Gerar comissao</h2>
      {#if paidInvoices.length === 0}
        <p class="muted">Sem faturas pagas. Marca uma fatura como paga primeiro.</p>
      {:else}
        <div class="row">
          <label>
            Fatura paga
            <select bind:value={fInvoiceId}>
              <option value="">— escolher —</option>
              {#each paidInvoices as inv}
                <option value={inv.id}>{inv.number} · {money(invoiceTotals(inv).net, inv.currency)} liq.</option>
              {/each}
            </select>
          </label>
          <label>Beneficiario<input bind:value={fBeneficiary} placeholder="Nome do vendedor" /></label>
          <label>Taxa %<input type="number" min="0" step="0.5" bind:value={fRate} /></label>
        </div>
        <div class="totals">
          <span>Comissao estimada <strong>{money(preview, selectedInvoice?.currency)}</strong></span>
        </div>
        <button class="primary" type="submit" disabled={busy || !fInvoiceId}>
          {busy ? "A registar..." : "Registar comissao"}
        </button>
      {/if}
    </form>

    <div class="card">
      <h2>Registo</h2>
      {#if loading}
        <p class="muted">A carregar...</p>
      {:else if commissions.length === 0}
        <p class="muted">Sem comissoes registadas.</p>
      {:else}
        <table class="list">
          <thead>
            <tr><th>Beneficiario</th><th>Fatura</th><th>Taxa</th><th>Valor</th><th>Estado</th><th></th></tr>
          </thead>
          <tbody>
            {#each commissions as c}
              <tr>
                <td>{c.beneficiary}</td>
                <td class="mono">{c.invoiceNumber ?? "—"}</td>
                <td>{c.ratePct}%</td>
                <td>{money(commissionAmount(c), c.currency)}</td>
                <td><span class="badge {c.status}">{commissionStatusLabel(c.status)}</span></td>
                <td class="actions">
                  {#if c.status === "pending"}
                    <button class="mini" onclick={() => setStatus(c.id, "paid")} disabled={busy}>Liquidar</button>
                    <button class="mini" onclick={() => setStatus(c.id, "void")} disabled={busy}>Anular</button>
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
  .list th, .list td { padding: 0.4rem 0.5rem; border-bottom: 1px solid #262630; text-align: left; }
  .mono { font-family: ui-monospace, monospace; }
  .totals { display: flex; gap: 1.25rem; margin: 0.5rem 0 0.75rem; font-size: 0.85rem; color: #bbb; }
  .actions { display: flex; gap: 0.4rem; }
  .mini { background: #23232c; border: 1px solid #34343f; color: #ddd; border-radius: 5px; padding: 0.2rem 0.5rem; font-size: 0.75rem; cursor: pointer; }
  .ev-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.4rem; }
  .ev-list li { display: flex; justify-content: space-between; align-items: center; padding: 0.4rem; border: 1px solid #2c2c36; border-radius: 5px; font-size: 0.82rem; }
  .ev-list li.suggested { border-color: var(--accent, #4f7cff); }
  .ev-actions { display: flex; gap: 0.3rem; }
  .primary { background: var(--accent, #4f7cff); color: #fff; border: none; border-radius: 6px; padding: 0.45rem 1rem; font-size: 0.85rem; cursor: pointer; margin-top: 0.5rem; }
  .primary:disabled { opacity: 0.6; cursor: default; }
  .badge { padding: 0.1rem 0.45rem; border-radius: 999px; font-size: 0.7rem; text-transform: uppercase; }
  .badge.pending { background: #2b3147; color: #9db4ff; }
  .badge.paid { background: #1f3a2a; color: #8ff0b4; }
  .badge.void { background: #3a2a2a; color: #d99; }
</style>
