<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listSubscriptions, updateSubscription } from "$lib/fin/subscriptionsService";
  import type { Subscription } from "$lib/fin/subscriptions";
  import {
    applyFiscalSuggestions,
    fiscalLabel,
    fiscalSummary,
    FISCAL_CATEGORIES,
    type FiscalCode,
  } from "$lib/fin/fiscal";

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");
  let subs = $state<Subscription[]>([]);

  const summary = $derived(fiscalSummary(subs));

  async function refresh() {
    loading = true;
    error = "";
    try {
      const mk = getMasterKey();
      if (!mk) {
        locked = true;
        subs = [];
        return;
      }
      locked = false;
      subs = await listSubscriptions();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      loading = false;
    }
  }

  onMount(() => void refresh());

  async function autoClassify() {
    busy = true;
    error = "";
    try {
      const suggestions = applyFiscalSuggestions(subs);
      for (const s of suggestions) {
        const sub = subs.find((x) => x.id === s.id);
        if (!sub) continue;
        await updateSubscription(s.id, { ...sub, fiscalCode: s.fiscalCode });
      }
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha na classificação";
    } finally {
      busy = false;
    }
  }

  async function setCode(id: string, code: FiscalCode) {
    const sub = subs.find((s) => s.id === id);
    if (!sub) return;
    busy = true;
    try {
      await updateSubscription(id, { ...sub, fiscalCode: code });
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao gravar";
    } finally {
      busy = false;
    }
  }

  function money(n: number, cur = "EUR"): string {
    return new Intl.NumberFormat("pt-PT", { style: "currency", currency: cur }).format(n);
  }
</script>

<svelte:head><title>Fiscal — AegisPass</title></svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">FIN-005</p>
      <h1>Categorização fiscal</h1>
      <DocHelpLink slug="journey-fiscal-categorization" label="Como funciona?" />
    </div>
    <a class="back" href="/fin">← Custos SaaS</a>
  </header>

  <p class="lead">
    Classificação IRC calculada no cliente — o servidor só guarda blobs cifrados. Não
    substitui aconselhamento de contabilista certificado.
  </p>

  {#if locked}
    <p class="muted">🔒 Desbloqueia a Master Key em <a href="/vault">/vault</a>.</p>
  {:else}
    {#if error}<p class="error" role="alert">{error}</p>{/if}

    <section class="metrics">
      <div class="metric">
        <span class="n">{money(summary.totalMonthly, summary.currency)}</span> gasto/mês
      </div>
      <div class="metric">
        <span class="n">{money(summary.totalDeductibleMonthly, summary.currency)}</span>
        dedutível estimado/mês
      </div>
      <div class="metric warn">
        <span class="n">{summary.pendingCount}</span> por classificar
      </div>
    </section>

    <div class="actions">
      <button type="button" class="btn primary" disabled={busy || summary.pendingCount === 0} onclick={autoClassify}>
        Sugerir automaticamente
      </button>
    </div>

    <section class="panel">
      {#if loading}
        <p class="muted">A carregar…</p>
      {:else if summary.lines.length === 0}
        <p class="muted">Sem subscrições activas.</p>
      {:else}
        <table>
          <thead>
            <tr><th>Serviço</th><th>Mensal</th><th>Classificação</th><th>Dedutível</th></tr>
          </thead>
          <tbody>
            {#each summary.lines as line (line.subscriptionId)}
              <tr>
                <td>{line.name}</td>
                <td class="mono">{money(line.monthly, line.currency)}</td>
                <td>
                  <select
                    value={line.fiscalCode}
                    disabled={busy}
                    onchange={(e) => setCode(line.subscriptionId, (e.currentTarget as HTMLSelectElement).value as FiscalCode)}
                  >
                    {#each FISCAL_CATEGORIES as cat}
                      <option value={cat.code}>{cat.label}</option>
                    {/each}
                  </select>
                </td>
                <td class="mono">{money(line.deductibleMonthly, line.currency)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </section>
  {/if}
</section>

<style>
  .page { max-width: 56rem; }
  .page-head {
    display: flex;
    justify-content: space-between;
    margin-bottom: var(--space-4);
  }
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  h1 { margin: 0; font-family: var(--font-display); }
  .lead { color: var(--color-text-muted); font-size: var(--text-sm); max-width: 40rem; }
  .back { font-size: var(--text-sm); color: var(--color-text-muted); }
  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
    gap: var(--space-3);
    margin: var(--space-4) 0;
  }
  .metric {
    padding: var(--space-3);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    font-size: var(--text-sm);
  }
  .metric .n { display: block; font-size: var(--text-lg); font-weight: 600; }
  .metric.warn .n { color: var(--color-warning, #b45309); }
  .actions { margin-bottom: var(--space-4); }
  .btn.primary {
    padding: var(--space-2) var(--space-4);
    background: var(--color-accent);
    color: var(--color-on-accent);
    border: none;
    border-radius: var(--radius-sm);
    font-weight: 600;
    cursor: pointer;
  }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: var(--space-4);
    overflow-x: auto;
  }
  table { width: 100%; border-collapse: collapse; font-size: var(--text-sm); }
  th, td { text-align: left; padding: var(--space-2); border-bottom: 1px solid var(--color-border); }
  .mono { font-family: var(--font-mono); }
  select { max-width: 16rem; }
  .muted, .error { font-size: var(--text-sm); }
  .error { color: var(--color-danger); }
</style>
