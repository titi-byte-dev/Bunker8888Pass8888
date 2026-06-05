<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listSubscriptions } from "$lib/fin/subscriptionsService";
  import type { Subscription } from "$lib/fin/subscriptions";
  import {
    connectBank,
    getBankingStatus,
    reportTransactionsSynced,
    syncTransactions,
    type BankConnection,
  } from "$lib/fin/openBankingService";
  import { reconcileTransactions, type BankTransaction, type ReconcileResult } from "$lib/fin/reconcile";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let connection = $state<BankConnection | null>(null);
  let subs = $state<Subscription[]>([]);
  let txs = $state<BankTransaction[]>([]);
  let reconcile = $state<ReconcileResult | null>(null);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  let reconcileApproved = $state(false);

  function money(n: number, cur = "EUR"): string {
    return new Intl.NumberFormat("pt-PT", { style: "currency", currency: cur }).format(Math.abs(n));
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
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
      connection = await getBankingStatus();
      subs = await listSubscriptions();
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao carregar";
    } finally {
      loading = false;
    }
  }

  async function handleConnect() {
    busy = true;
    error = "";
    try {
      connection = await connectBank();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao ligar banco";
    } finally {
      busy = false;
    }
  }

  async function handleSync() {
    busy = true;
    error = "";
    reconcile = null;
    try {
      txs = await syncTransactions();
      reconcile = reconcileTransactions(txs, subs);
      await reportTransactionsSynced(txs.length, reconcile.matched.length, reconcile.unmatched.length);
      connection = await getBankingStatus();
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha na sincronização";
    } finally {
      busy = false;
    }
  }

  async function handleApprove(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
    try {
      const result = await approveSuggestion(ev.id);
      await refreshEvents();
      if (result.action === "reconcile_payments") {
        reconcileApproved = true;
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
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao rejeitar";
    } finally {
      decidingId = null;
    }
  }

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) void refresh();
    else loading = false;
  });
</script>

<svelte:head>
  <title>Open Banking — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">FIN-003 · AGENT-006 · Reconciliação</p>
      <h1>Open Banking</h1>
      <DocHelpLink slug="journey-finance-agent-reconcile" label="Como funciona a reconciliação?" />
    </div>
    <a class="back" href="/fin/costs">← Custos</a>
  </header>

  <p class="lead">
    Liga uma conta bancária (mock em dev), sincroniza movimentos e cruza débitos com
    as tuas subscrições SaaS. O servidor recebe apenas contagens — os movimentos
    ficam no browser.
  </p>

  {#if locked}
    <section class="panel">
      <p class="muted">🔒 Desbloqueia a Master Key para reconciliar.</p>
      <a class="btn primary" href="/vault">Ir desbloquear</a>
    </section>
  {:else}
    {#if error}<p class="inline-error" role="alert">{error}</p>{/if}
    {#if loading}
      <p class="muted">A carregar…</p>
    {:else}
      <section class="panel">
        <div class="panel-head">
          <p class="eyebrow">Ligação bancária</p>
          <span class="badge" class:ok={connection?.status === "connected"}>
            {connection?.status ?? "pending"} · {connection?.provider ?? "mock"}
          </span>
        </div>
        {#if connection?.status !== "connected"}
          <button type="button" class="btn primary" disabled={busy} onclick={handleConnect}>
            Simular consentimento PSD2
          </button>
        {:else}
          <p class="muted">
            Consentimento activo
            {#if connection.consentExpiresAt}
              · expira {new Date(connection.consentExpiresAt).toLocaleDateString("pt-PT")}
            {/if}
            {#if connection.lastSyncAt}
              · última sync {new Date(connection.lastSyncAt).toLocaleString("pt-PT")}
            {/if}
          </p>
          <button type="button" class="btn primary" disabled={busy} onclick={handleSync}>
            Sincronizar movimentos
          </button>
        {/if}
      </section>

      <section class="panel events">
        <h2>Actividade dos agentes</h2>
        {#if agentEvents.length === 0}
          <p class="muted">Sem eventos recentes.</p>
        {:else}
          <ul class="event-list">
            {#each agentEvents.slice(0, 6) as ev (ev.id)}
              <li class:suggested={isPendingSuggestion(ev)}>
                <div class="ev-body">
                  <span class="ev-label">{ev.label}</span>
                  {#if isPendingSuggestion(ev) && ev.payload.action === "reconcile_payments"}
                    <div class="ev-actions">
                      <button type="button" class="btn approve" disabled={decidingId !== null || busy} onclick={() => handleApprove(ev)}>
                        {decidingId === ev.id ? "…" : "Aprovar"}
                      </button>
                      <button type="button" class="btn reject" disabled={decidingId !== null} onclick={() => handleReject(ev)}>
                        Rejeitar
                      </button>
                    </div>
                  {/if}
                </div>
              </li>
            {/each}
          </ul>
        {/if}
      </section>

      {#if reconcileApproved && reconcile}
        <p class="hint" role="status">Reconciliação aprovada — revê os movimentos sem correspondência.</p>
      {/if}

      {#if reconcile}
        <section class="panel">
          <p class="eyebrow">Resultado ({reconcile.matched.length} associados · {reconcile.unmatched.length} por classificar)</p>
          {#if reconcile.matched.length > 0}
            <h3>Associados a subscrições</h3>
            <ul class="match-list">
              {#each reconcile.matched as m (m.transactionId)}
                <li>
                  <span class="name">{m.subscriptionName}</span>
                  <span class="amt">{money(m.amount)}</span>
                </li>
              {/each}
            </ul>
          {/if}
          {#if reconcile.unmatched.length > 0}
            <h3>Sem correspondência</h3>
            <ul class="match-list unmatched">
              {#each reconcile.unmatched as tx (tx.id)}
                <li>
                  <span class="name">{tx.description}</span>
                  <span class="amt">{money(tx.amount, tx.currency)}</span>
                </li>
              {/each}
            </ul>
          {/if}
        </section>
      {/if}
    {/if}
  {/if}
</section>

<style>
  .page { max-width: 52rem; margin: 0 auto; padding: var(--space-6); }
  .page-head { display: flex; justify-content: space-between; align-items: flex-start; gap: var(--space-4); margin-bottom: var(--space-4); }
  .back { font-size: var(--text-sm); color: var(--color-text-muted); text-decoration: none; }
  .lead { color: var(--color-text-muted); font-size: var(--text-sm); margin: 0 0 var(--space-6); }
  .panel { border: 1px solid var(--color-border); border-radius: var(--radius-md); background: var(--color-bg-surface); padding: var(--space-4) var(--space-6); margin-bottom: var(--space-4); }
  .panel-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: var(--space-3); }
  .eyebrow { margin: 0; font-size: var(--text-xs); text-transform: uppercase; letter-spacing: 0.06em; color: var(--color-text-label); }
  .muted { color: var(--color-text-muted); font-size: var(--text-sm); }
  .badge { font-size: var(--text-xs); padding: var(--space-1) var(--space-2); border-radius: var(--radius-sm); background: var(--color-bg-inset); }
  .badge.ok { background: var(--color-accent-muted); color: var(--color-accent); }
  .btn { padding: var(--space-2) var(--space-4); border-radius: var(--radius-sm); border: 1px solid var(--color-border); font-size: var(--text-sm); cursor: pointer; background: var(--color-bg-elevated); color: var(--color-text); text-decoration: none; display: inline-block; }
  .btn.primary { background: var(--color-accent); color: var(--color-accent-fg); border-color: transparent; }
  .btn:disabled { opacity: 0.55; cursor: progress; }
  .inline-error { color: var(--color-danger); font-size: var(--text-sm); }
  .events h2 { margin: 0 0 var(--space-3); font-size: var(--text-base); }
  .event-list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: var(--space-2); }
  .event-list li { padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-sm); background: var(--color-bg-inset); }
  .event-list li.suggested { border-color: var(--color-accent); }
  .ev-body { display: flex; align-items: center; justify-content: space-between; gap: var(--space-3); flex-wrap: wrap; }
  .ev-label { font-size: var(--text-sm); }
  .ev-actions { display: flex; gap: var(--space-2); }
  .btn.approve { background: var(--color-accent); color: var(--color-accent-fg); border-color: transparent; }
  .hint { font-size: var(--text-sm); color: var(--color-accent); margin-bottom: var(--space-4); }
  .match-list { list-style: none; padding: 0; margin: var(--space-2) 0 var(--space-4); display: flex; flex-direction: column; gap: var(--space-1); }
  .match-list li { display: flex; justify-content: space-between; padding: var(--space-2) var(--space-3); background: var(--color-bg-inset); border-radius: var(--radius-sm); font-size: var(--text-sm); }
  .match-list.unmatched li { border-left: 3px solid var(--color-warning, #c79a2e); }
  h3 { margin: var(--space-3) 0 var(--space-2); font-size: var(--text-sm); }
  .amt { font-family: var(--font-mono); }
</style>
