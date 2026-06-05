<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import {
    createSubscription,
    deleteSubscription,
    listSubscriptions,
    listVaultLogins,
    updateSubscription,
    type VaultLoginRef,
  } from "$lib/fin/subscriptionsService";
  import type { BillingCycle, Subscription } from "$lib/fin/subscriptions";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { costSummary, detectAlerts, monthlyCost, type Alert } from "$lib/fin/alerts";
  import {
    FISCAL_CATEGORIES,
    fiscalLabel,
    suggestFiscalCode,
    type FiscalCode,
  } from "$lib/fin/fiscal";
  import { reportStaleAlerts } from "$lib/fin/financeAgentService";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";
  import {
    Button,
    confirmDialog,
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

  let subs = $state<Subscription[]>([]);
  let logins = $state<VaultLoginRef[]>([]);

  let editingId = $state<string | null>(null);
  let fName = $state("");
  let fCost = $state(0);
  let fCurrency = $state("EUR");
  let fCycle = $state<BillingCycle>("monthly");
  let fCategory = $state("");
  let fFiscalCode = $state<FiscalCode>("pendente");
  let fVaultItemId = $state("");
  let fLastUsedAt = $state("");
  let fActive = $state(true);

  const summary = $derived(costSummary(subs));
  const alerts = $derived(detectAlerts(subs));

  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  let reviewApproved = $state(false);

  const subColumns = $derived<DataColumn<Subscription>[]>([
    {
      id: "name",
      label: "Serviço",
      accessor: (s) => (alertsFor(s.id).length > 0 ? `${s.name} ⚠` : s.name),
    },
    {
      id: "monthly",
      label: "Mensal",
      align: "right",
      mono: true,
      accessor: (s) => money(monthlyCost(s), s.currency),
    },
    {
      id: "fiscal",
      label: "Fiscal",
      muted: true,
      accessor: (s) => fiscalLabel(s.fiscalCode),
    },
    {
      id: "cycle",
      label: "Ciclo",
      accessor: (s) => (s.cycle === "yearly" ? "Anual" : "Mensal"),
    },
    {
      id: "vault",
      label: "Cofre",
      muted: true,
      accessor: (s) => s.vaultItemTitle ?? "—",
    },
    {
      id: "status",
      label: "Estado",
      accessor: (s) => (s.active ? "activa" : "inactiva"),
    },
  ]);

  async function refreshEvents() {
    try {
      agentEvents = await listAgentEvents();
    } catch {
      agentEvents = [];
    }
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
  }

  async function handleReportToAgent() {
    if (alerts.length === 0) return;
    busy = true;
    error = "";
    try {
      await reportStaleAlerts(alerts, subs);
      await refreshEvents();
      toast.success("Pedido de revisão enviado ao agente.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao reportar alertas";
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
      if (result.action === "review_saas_licenses") {
        reviewApproved = true;
        toast.success("Revisão aprovada — revê os alertas abaixo.");
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

  async function refresh() {
    loading = true;
    error = "";
    try {
      subs = await listSubscriptions();
      logins = await listVaultLogins();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar subscrições";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) {
      refresh();
      void refreshEvents();
    } else loading = false;
  });

  function resetForm() {
    editingId = null;
    fName = "";
    fCost = 0;
    fCurrency = "EUR";
    fCycle = "monthly";
    fCategory = "";
    fFiscalCode = "pendente";
    fVaultItemId = "";
    fLastUsedAt = "";
    fActive = true;
  }

  function edit(s: Subscription) {
    editingId = s.id;
    fName = s.name;
    fCost = s.cost;
    fCurrency = s.currency;
    fCycle = s.cycle;
    fCategory = s.category ?? "";
    fFiscalCode = s.fiscalCode ?? "pendente";
    fVaultItemId = s.vaultItemId ?? "";
    fLastUsedAt = s.lastUsedAt ? s.lastUsedAt.slice(0, 10) : "";
    fActive = s.active;
  }

  async function save(e: SubmitEvent) {
    e.preventDefault();
    if (!fName.trim()) return;
    busy = true;
    error = "";
    const login = logins.find((l) => l.id === fVaultItemId);
    const payload = {
      name: fName.trim(),
      cost: Number(fCost) || 0,
      currency: fCurrency,
      cycle: fCycle,
      category: fCategory.trim() || undefined,
      fiscalCode: fFiscalCode === "pendente" ? undefined : fFiscalCode,
      vaultItemId: fVaultItemId || undefined,
      vaultItemTitle: login?.title,
      lastUsedAt: fLastUsedAt ? new Date(fLastUsedAt).toISOString() : undefined,
      active: fActive,
    };
    const wasEdit = !!editingId;
    try {
      if (editingId) await updateSubscription(editingId, payload);
      else await createSubscription(payload);
      resetForm();
      await refresh();
      toast.success(wasEdit ? "Subscrição actualizada." : "Subscrição adicionada.");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao gravar";
    } finally {
      busy = false;
    }
  }

  async function remove(s: Subscription) {
    const ok = await confirmDialog({
      title: "Apagar subscrição?",
      message: `Remove «${s.name}» da lista de custos. Os dados cifrados são eliminados localmente.`,
      variant: "danger",
      confirmLabel: "Apagar",
    });
    if (!ok) return;
    busy = true;
    try {
      await deleteSubscription(s.id);
      if (editingId === s.id) resetForm();
      await refresh();
      toast.success(`«${s.name}» removido.`);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao apagar";
    } finally {
      busy = false;
    }
  }

  function money(n: number, currency = "EUR"): string {
    return new Intl.NumberFormat("pt-PT", { style: "currency", currency }).format(n);
  }

  function alertsFor(id: string): Alert[] {
    return alerts.filter((a) => a.subscriptionId === id);
  }
</script>

<svelte:head>
  <title>Custos SaaS — AegisPass</title>
</svelte:head>

<PageShell
  title="Monitorização de Custos"
  taskId="FIN-001/002 · AGENT-006"
  description="As subscrições são cifradas com a tua Master Key — só tu vês os custos. O dashboard cruza cada subscrição com o login do cofre e assinala licenças esquecidas ou sem credencial associada."
 
>
  {#snippet actions()}
    <DocHelpLink slug="journey-finance-agent-saas" label="Como funciona o agente financeiro?" />
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para ver e gerir custos SaaS.">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else}
    {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

    <section class="metrics">
      <MetricCard label="por mês" value={money(summary.monthly)} />
      <MetricCard label="por ano" value={money(summary.yearly)} />
      <MetricCard label="activas" value={String(summary.activeCount)} />
      <MetricCard
        label="poupança potencial/mês"
        value={money(summary.potentialMonthlySaving)}
        variant="warning"
      />
    </section>

    <Panel title="Actividade dos agentes">
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos recentes.</p>
      {:else}
        <ul class="event-list">
          {#each agentEvents.slice(0, 6) as ev (ev.id)}
            <li class:suggested={isPendingSuggestion(ev)}>
              <div class="ev-body">
                <span class="ev-label">{ev.label}</span>
                {#if isPendingSuggestion(ev) && ev.payload.action === "review_saas_licenses"}
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

    {#if reviewApproved}
      <StatusBanner variant="info">Revisão aprovada — revê os alertas abaixo e cancela licenças inactivas.</StatusBanner>
    {/if}

    {#if alerts.length > 0}
      <Panel title="Alertas ({alerts.length})">
        {#snippet actions()}
          <Button variant="secondary" size="sm" disabled={busy} onclick={handleReportToAgent}>
            Pedir revisão ao agente
          </Button>
        {/snippet}
        <ul class="alert-list">
          {#each alerts as a (a.subscriptionId + a.kind)}
            <li class="alert" class:stale={a.kind === "stale"} class:orphan={a.kind === "orphan"}>
              <span class="a-name">{a.name}</span>
              <span class="a-reason">{a.reason}</span>
              {#if a.kind === "stale"}<span class="a-save">−{money(a.monthlySaving)}/mês</span>{/if}
            </li>
          {/each}
        </ul>
      </Panel>
    {/if}

    <Panel title={editingId ? "Editar subscrição" : "Nova subscrição"}>
      {#snippet actions()}
        {#if editingId}
          <Button variant="ghost" size="sm" onclick={resetForm}>Cancelar edição</Button>
        {/if}
      {/snippet}
      <form class="sub-form" onsubmit={save}>
        <div class="row">
          <label class="field grow"><span>Serviço</span>
            <input bind:value={fName} placeholder="Netflix" disabled={busy} /></label>
          <label class="field"><span>Custo</span>
            <input type="number" min="0" step="0.01" bind:value={fCost} disabled={busy} /></label>
          <label class="field"><span>Moeda</span>
            <input bind:value={fCurrency} disabled={busy} /></label>
          <label class="field"><span>Ciclo</span>
            <select bind:value={fCycle} disabled={busy}>
              <option value="monthly">Mensal</option>
              <option value="yearly">Anual</option>
            </select></label>
        </div>
        <div class="row">
          <label class="field"><span>Categoria</span>
            <input bind:value={fCategory} placeholder="Entretenimento" disabled={busy} /></label>
          <label class="field grow"><span>Fiscal (FIN-005)</span>
            <select bind:value={fFiscalCode} disabled={busy}>
              {#each FISCAL_CATEGORIES as cat}
                <option value={cat.code}>{cat.label}</option>
              {/each}
            </select></label>
          <button
            type="button"
            class="link-btn"
            disabled={busy || !fName.trim()}
            onclick={() => { fFiscalCode = suggestFiscalCode({ name: fName, category: fCategory }); }}
          >sugerir</button>
          <label class="field grow"><span>Login no cofre</span>
            <select bind:value={fVaultItemId} disabled={busy}>
              <option value="">— sem associação —</option>
              {#each logins as l (l.id)}
                <option value={l.id}>{l.title}</option>
              {/each}
            </select></label>
          <label class="field"><span>Último uso</span>
            <input type="date" bind:value={fLastUsedAt} disabled={busy} /></label>
          <label class="field check"><span>Activa</span>
            <input type="checkbox" bind:checked={fActive} disabled={busy} /></label>
        </div>
        <Button type="submit" disabled={busy || !fName.trim()} loading={busy}>
          {editingId ? "Guardar alterações" : "Adicionar"}
        </Button>
      </form>
    </Panel>

    <Panel title="Subscrições" padding="none">
      <DataTable
        columns={subColumns}
        rows={subs}
        keyFn={(s) => s.id}
        loading={loading}
        dense
        rowClass={(s) => (!s.active ? "off" : undefined)}
        emptyTitle="Sem subscrições"
        emptyDescription="Adiciona a primeira subscrição no formulário acima."
      >
        {#snippet actions(row)}
          <Button variant="ghost" size="sm" onclick={() => edit(row)}>editar</Button>
          <Button variant="ghost" size="sm" disabled={busy} onclick={() => remove(row)}>apagar</Button>
        {/snippet}
      </DataTable>
    </Panel>
  {/if}
</PageShell>

<style>
  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
    gap: var(--space-3);
  }

  .muted {
    margin: 0;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .alert-list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }

  .alert {
    display: flex;
    align-items: baseline;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    border-left: 3px solid var(--color-border);
    font-size: var(--text-sm);
  }

  .alert.stale { border-left-color: var(--color-danger); }
  .alert.orphan { border-left-color: var(--color-warning); }
  .a-name { font-weight: 500; }
  .a-reason { color: var(--color-text-muted); flex: 1; }
  .a-save { color: var(--color-danger); font-family: var(--font-mono); }

  .sub-form .row {
    display: flex;
    gap: var(--space-2);
    flex-wrap: wrap;
    margin-bottom: var(--space-2);
  }

  .field {
    display: block;
    margin-bottom: var(--space-3);
  }

  .field.grow { flex: 1; min-width: 9rem; }
  .field.check { align-self: center; }

  .field > span {
    display: block;
    margin-bottom: var(--space-1);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }

  input,
  select {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    box-sizing: border-box;
    max-width: 9rem;
  }

  .field.grow input,
  .field.grow select {
    max-width: none;
    width: 100%;
  }

  input:focus-visible,
  select:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }

  .link-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    padding: 0;
    align-self: flex-end;
    margin-bottom: var(--space-3);
  }

  .link-btn:hover:not(:disabled) { color: var(--color-accent); }
  .link-btn:disabled { opacity: 0.5; cursor: not-allowed; }

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
</style>
