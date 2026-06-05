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

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");

  let subs = $state<Subscription[]>([]);
  let logins = $state<VaultLoginRef[]>([]);

  // Formulário (criar/editar).
  let editingId = $state<string | null>(null);
  let fName = $state("");
  let fCost = $state(0);
  let fCurrency = $state("EUR");
  let fCycle = $state<BillingCycle>("monthly");
  let fCategory = $state("");
  let fVaultItemId = $state("");
  let fLastUsedAt = $state("");
  let fActive = $state(true);

  const summary = $derived(costSummary(subs));
  const alerts = $derived(detectAlerts(subs));

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
    if (!locked) refresh();
    else loading = false;
  });

  function resetForm() {
    editingId = null;
    fName = "";
    fCost = 0;
    fCurrency = "EUR";
    fCycle = "monthly";
    fCategory = "";
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
      vaultItemId: fVaultItemId || undefined,
      vaultItemTitle: login?.title,
      lastUsedAt: fLastUsedAt ? new Date(fLastUsedAt).toISOString() : undefined,
      active: fActive,
    };
    try {
      if (editingId) await updateSubscription(editingId, payload);
      else await createSubscription(payload);
      resetForm();
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao gravar";
    } finally {
      busy = false;
    }
  }

  async function remove(id: string) {
    busy = true;
    try {
      await deleteSubscription(id);
      if (editingId === id) resetForm();
      await refresh();
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

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">FIN-001/002 · Custos SaaS</p>
      <h1>Monitorização de Custos</h1>
      <DocHelpLink />
    </div>
    <p class="lead">
      As subscrições são cifradas com a tua Master Key — só tu vês os custos. O
      dashboard cruza cada subscrição com o login do cofre e assinala licenças
      esquecidas (sem uso) ou sem credencial associada.
    </p>
  </header>

  {#if locked}
    <section class="panel">
      <p class="muted">🔒 Desbloqueia a Master Key para ver os custos.</p>
      <a class="btn primary" href="/vault">Ir desbloquear</a>
    </section>
  {:else}
    {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

    <section class="metrics">
      <div class="metric"><span class="n">{money(summary.monthly)}</span> por mês</div>
      <div class="metric"><span class="n">{money(summary.yearly)}</span> por ano</div>
      <div class="metric"><span class="n">{summary.activeCount}</span> activas</div>
      <div class="metric warn">
        <span class="n">{money(summary.potentialMonthlySaving)}</span> poupança potencial/mês
      </div>
    </section>

    {#if alerts.length > 0}
      <section class="panel alerts">
        <div class="panel-head"><p class="eyebrow">⚠️ Alertas ({alerts.length})</p></div>
        <ul class="alert-list">
          {#each alerts as a (a.subscriptionId + a.kind)}
            <li class="alert" class:stale={a.kind === "stale"} class:orphan={a.kind === "orphan"}>
              <span class="a-name">{a.name}</span>
              <span class="a-reason">{a.reason}</span>
              {#if a.kind === "stale"}<span class="a-save">−{money(a.monthlySaving)}/mês</span>{/if}
            </li>
          {/each}
        </ul>
      </section>
    {/if}

    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">{editingId ? "Editar subscrição" : "Nova subscrição"}</p>
        {#if editingId}
          <button type="button" class="link-btn" onclick={resetForm}>cancelar edição</button>
        {/if}
      </div>
      <form onsubmit={save}>
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
          <label class="field grow"><span>Login no cofre (cruza com vault)</span>
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
        <button type="submit" class="btn primary" disabled={busy || !fName.trim()}>
          {editingId ? "Guardar alterações" : "Adicionar"}
        </button>
      </form>
    </section>

    <section class="panel">
      <div class="panel-head"><p class="eyebrow">Subscrições</p></div>
      {#if loading}
        <p class="muted">A carregar…</p>
      {:else if subs.length === 0}
        <p class="muted">Sem subscrições. Adiciona a primeira acima.</p>
      {:else}
        <table>
          <thead>
            <tr><th>Serviço</th><th>Mensal</th><th>Ciclo</th><th>Cofre</th><th>Estado</th><th></th></tr>
          </thead>
          <tbody>
            {#each subs as s (s.id)}
              <tr class:off={!s.active}>
                <td>
                  {s.name}
                  {#if alertsFor(s.id).length > 0}<span class="flag">⚠</span>{/if}
                </td>
                <td class="mono">{money(monthlyCost(s), s.currency)}</td>
                <td>{s.cycle === "yearly" ? "Anual" : "Mensal"}</td>
                <td class="muted sm">{s.vaultItemTitle ?? "—"}</td>
                <td>{s.active ? "activa" : "inactiva"}</td>
                <td class="actions">
                  <button type="button" class="link-btn" onclick={() => edit(s)}>editar</button>
                  <button type="button" class="link-btn" onclick={() => remove(s.id)} disabled={busy}>apagar</button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </section>
  {/if}
</section>

<style>
  .page {
    max-width: 56rem;
  }
  .page-head {
    margin-bottom: var(--space-5);
  }
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }
  .lead {
    margin: var(--space-3) 0 0;
    max-width: 42rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
    gap: var(--space-3);
    margin-bottom: var(--space-4);
  }
  .metric {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
  .metric.warn .n {
    color: var(--color-danger);
  }
  .metric .n {
    display: block;
    font-family: var(--font-display);
    font-size: var(--text-xl);
    color: var(--color-text);
  }
  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4) var(--space-6);
    margin-bottom: var(--space-4);
  }
  .panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-3);
  }
  .panel-head .eyebrow {
    margin: 0;
  }
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .sm {
    font-size: var(--text-xs);
  }
  .mono {
    font-family: var(--font-mono);
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
  .alert.stale {
    border-left-color: var(--color-danger);
  }
  .alert.orphan {
    border-left-color: var(--color-warning, #c79a2e);
  }
  .a-name {
    font-weight: 500;
  }
  .a-reason {
    color: var(--color-text-muted);
    flex: 1;
  }
  .a-save {
    color: var(--color-danger);
    font-family: var(--font-mono);
  }
  .row {
    display: flex;
    gap: var(--space-2);
    flex-wrap: wrap;
  }
  .field {
    display: block;
    margin-bottom: var(--space-3);
  }
  .field.grow {
    flex: 1;
    min-width: 9rem;
  }
  .field.check {
    align-self: center;
  }
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
  table {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--text-sm);
  }
  th,
  td {
    text-align: left;
    padding: var(--space-2) var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }
  th {
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  tr.off {
    opacity: 0.55;
  }
  .flag {
    color: var(--color-danger);
  }
  td.actions {
    display: flex;
    gap: var(--space-3);
  }
  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    color: var(--color-text);
  }
  .btn.primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
    border-color: transparent;
  }
  .btn.primary:hover:not(:disabled) {
    filter: brightness(1.08);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .link-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    padding: 0;
  }
  .link-btn:hover {
    color: var(--color-danger);
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
</style>
