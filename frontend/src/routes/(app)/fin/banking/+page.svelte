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
  import {
    Button,
    EmptyState,
    PageShell,
    Panel,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";

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
      toast.success("Consentimento bancário simulado (PSD2).");
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
      toast.success(
        `${reconcile.matched.length} associados · ${reconcile.unmatched.length} por classificar.`,
      );
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
        toast.success("Reconciliação aprovada.");
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

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) void refresh();
    else loading = false;
  });
</script>

<svelte:head>
  <title>Reconciliação bancária — AegisPass</title>
</svelte:head>

<PageShell
  title="Open Banking"
  taskId="FIN-003 · AGENT-006"
  description="Liga uma conta bancária (mock em dev), sincroniza movimentos e cruza débitos com as tuas subscrições SaaS. O servidor recebe apenas contagens — os movimentos ficam no browser."
>
  {#snippet actions()}
    <DocHelpLink slug="journey-finance-agent-reconcile" label="Como funciona a reconciliação?" />
    <Button variant="ghost" size="sm" href="/fin/costs">← Custos</Button>
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para reconciliar movimentos.">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else if loading}
    <Skeleton variant="block" height="6rem" />
    <Skeleton variant="block" height="8rem" />
  {:else}
    {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

    <Panel title="Ligação bancária">
      {#snippet actions()}
        <span class="badge" class:ok={connection?.status === "connected"}>
          {connection?.status ?? "pending"} · {connection?.provider ?? "mock"}
        </span>
      {/snippet}
      {#if connection?.status !== "connected"}
        <Button disabled={busy} loading={busy} onclick={handleConnect}>
          Simular consentimento PSD2
        </Button>
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
        <Button disabled={busy} loading={busy} onclick={handleSync}>
          Sincronizar movimentos
        </Button>
      {/if}
    </Panel>

    <Panel title="Actividade dos agentes">
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
                    <Button
                      variant="primary"
                      size="sm"
                      loading={decidingId === ev.id}
                      disabled={decidingId !== null || busy}
                      onclick={() => handleApprove(ev)}
                    >
                      Aprovar
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      disabled={decidingId !== null}
                      onclick={() => handleReject(ev)}
                    >
                      Rejeitar
                    </Button>
                  </div>
                {/if}
              </div>
            </li>
          {/each}
        </ul>
      {/if}
    </Panel>

    {#if reconcileApproved && reconcile}
      <StatusBanner variant="info">Reconciliação aprovada — revê os movimentos sem correspondência.</StatusBanner>
    {/if}

    {#if reconcile}
      <Panel title="Resultado ({reconcile.matched.length} associados · {reconcile.unmatched.length} por classificar)">
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
      </Panel>
    {/if}
  {/if}
</PageShell>

<style>
  .badge {
    font-size: var(--text-xs);
    padding: var(--space-1) var(--space-2);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text-muted);
  }

  .badge.ok {
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }

  .muted {
    margin: 0 0 var(--space-3);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .event-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .event-list li {
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
  }

  .event-list li.suggested { border-color: var(--color-accent); }

  .ev-body {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    flex-wrap: wrap;
  }

  .ev-label { font-size: var(--text-sm); }
  .ev-actions { display: flex; gap: var(--space-2); }

  .match-list {
    list-style: none;
    padding: 0;
    margin: var(--space-2) 0 var(--space-4);
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .match-list li {
    display: flex;
    justify-content: space-between;
    padding: var(--space-2) var(--space-3);
    background: var(--color-bg-inset);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
  }

  .match-list.unmatched li {
    border-left: 3px solid var(--color-warning);
  }

  h3 {
    margin: var(--space-3) 0 var(--space-2);
    font-size: var(--text-sm);
    font-weight: 600;
  }

  .amt { font-family: var(--font-mono); }
</style>
