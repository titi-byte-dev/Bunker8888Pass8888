<script lang="ts">
  import { onMount } from "svelte";
  import { analyzeVaultHygiene, type HygieneSummary, type ItemHygieneResult } from "$lib/vault/hygiene";
  import {
    buildHealthReport,
    checkPasswordBreached,
    itemsRequiringPasswordChange,
    remediationEditUrl,
    saveHealthSnapshot,
    type BreachCheckResult,
  } from "$lib/darkweb";
  import SecurityHealthCard from "$lib/security/SecurityHealthCard.svelte";
  import { loadDecodedLogins } from "$lib/vault/ui";

  let summary = $state<HygieneSummary | null>(null);
  let passwordsByItem = $state<Map<string, string>>(new Map());
  let breachByItem = $state<Map<string, BreachCheckResult>>(new Map());
  let breachBusy = $state(false);
  let breachError = $state("");
  let busy = $state(true);
  let error = $state("");

  const healthReport = $derived(
    summary ? buildHealthReport(summary, breachByItem) : null,
  );

  const remediation = $derived(
    summary ? itemsRequiringPasswordChange(summary, breachByItem) : [],
  );

  async function loadHygiene() {
    busy = true;
    error = "";
    summary = null;
    breachByItem = new Map();
    try {
      const logins = await loadDecodedLogins();
      passwordsByItem = new Map(logins.map((l) => [l.meta.id, l.login.password]));
      summary = await analyzeVaultHygiene(
        logins.map(({ meta, login }) => ({
          itemId: meta.id,
          title: login.title,
          password: login.password,
        })),
      );
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao analisar higiene";
    } finally {
      busy = false;
    }
  }

  async function scanBreaches() {
    if (!summary) return;
    breachBusy = true;
    breachError = "";
    const next = new Map(breachByItem);
    try {
      for (const item of summary.items) {
        const password = passwordsByItem.get(item.itemId);
        if (!password) continue;
        const result = await checkPasswordBreached(password);
        next.set(item.itemId, result);
        breachByItem = new Map(next);
        await new Promise((r) => setTimeout(r, 400));
      }
      if (summary) {
        const report = buildHealthReport(summary, next);
        saveHealthSnapshot(report);
      }
    } catch (e) {
      breachError = e instanceof Error ? e.message : "Verificação falhou";
    } finally {
      breachBusy = false;
    }
  }

  function issueLabel(issues: ItemHygieneResult["issues"]): string {
    const labels: Record<string, string> = {
      weak: "Fraca",
      reused: "Reutilizada",
      common: "Comum",
      short: "Curta",
    };
    return issues.map((i) => labels[i] ?? i).join(", ") || "OK";
  }

  function scoreClass(score: number): string {
    if (score >= 75) return "good";
    if (score >= 50) return "warn";
    return "bad";
  }

  onMount(loadHygiene);
</script>

<svelte:head>
  <title>Saúde de segurança — AegisPass</title>
</svelte:head>

<section class="page">
  <a href="/security" class="back">← Segurança</a>
  <h1>Saúde de segurança</h1>
  <p class="lead">
    Higiene + fugas (DW-003). Análise no cliente; k-anonymity para breach check (DW-001).
  </p>

  {#if busy}
    <p class="muted">A analisar cofre…</p>
  {:else if error}
    <p class="error" role="alert">{error}</p>
  {:else if summary && healthReport}
    <SecurityHealthCard report={healthReport} />

    {#if remediation.length > 0}
      <div class="remediation-banner" role="alert">
        <strong>Acção necessária (DW-002)</strong>
        <p>{remediation.length} password(s) exposta(s) em fugas — altera-as já.</p>
        <ul>
          {#each remediation as item (item.itemId)}
            <li>
              <a href={remediationEditUrl(item.itemId, item.reason)}>{item.title} — Alterar password</a>
            </li>
          {/each}
        </ul>
      </div>
    {/if}

    <div class="breach-panel">
      <h2>Verificação de fugas</h2>
      <p class="hint">
        API HIBP com k-anonymity — só os primeiros 5 caracteres do hash SHA-1 saem do dispositivo.
      </p>
      <button type="button" onclick={scanBreaches} disabled={breachBusy || summary.items.length === 0}>
        {breachBusy ? "A verificar…" : "Verificar fugas"}
      </button>
      {#if breachError}
        <p class="error" role="alert">{breachError}</p>
      {/if}
    </div>

    <h2 class="section-title">Detalhe por login</h2>
    <ul class="items">
      {#each summary.items as item (item.itemId)}
        {@const breach = breachByItem.get(item.itemId)}
        {@const needsFix = breach?.breached}
        <li class:urgent={needsFix}>
          <div class="row">
            <a href="/vault/{item.itemId}">{item.title}</a>
            <span class="badge {scoreClass(item.score)}">{item.score}</span>
          </div>
          <p class="issues">{issueLabel(item.issues)}</p>
          {#if breach}
            <p class="breach" class:exposed={breach.breached}>
              {#if breach.breached}
                ⚠ Exposta ({breach.exposureCount.toLocaleString("pt-PT")}×)
              {:else}
                ✓ Sem ocorrências conhecidas
              {/if}
            </p>
          {/if}
          {#if needsFix}
            <a class="fix-btn" href={remediationEditUrl(item.itemId, item.issues.includes("weak") ? "weak_and_breach" : "breach")}>
              Alterar password
            </a>
          {/if}
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .page {
    max-width: 40rem;
  }

  .back {
    display: inline-block;
    margin-bottom: var(--space-4);
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  h1 {
    margin: 0 0 var(--space-2);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .lead {
    color: var(--color-text-muted);
    margin: 0 0 var(--space-6);
    font-size: var(--text-sm);
    line-height: var(--leading-body);
  }

  .section-title {
    font-size: var(--text-lg);
    margin: 0 0 var(--space-3);
  }

  .remediation-banner {
    padding: var(--space-4);
    margin-bottom: var(--space-6);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-warning);
    background: color-mix(in srgb, var(--color-warning) 10%, transparent);
  }

  .remediation-banner p {
    margin: var(--space-2) 0;
    font-size: var(--text-sm);
  }

  .remediation-banner ul {
    margin: 0;
    padding-left: var(--space-4);
    font-size: var(--text-sm);
  }

  .remediation-banner a {
    color: var(--color-link);
  }

  .breach-panel {
    margin-bottom: var(--space-6);
    padding: var(--space-4);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-md);
  }

  .breach-panel h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-lg);
  }

  .hint {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .breach-panel button {
    padding: var(--space-2) var(--space-4);
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    font-weight: 600;
    cursor: pointer;
  }

  .breach-panel button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .items {
    list-style: none;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .items li {
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border);
  }

  .items li.urgent {
    background: color-mix(in srgb, var(--color-warning) 6%, transparent);
  }

  .items li:last-child {
    border-bottom: none;
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-3);
  }

  .row a {
    color: inherit;
    font-weight: 500;
    text-decoration: none;
  }

  .row a:hover {
    color: var(--color-link);
  }

  .badge {
    font-size: var(--text-xs);
    font-weight: 700;
    padding: 2px 8px;
    border-radius: 999px;
    background: var(--color-bg-inset);
  }

  .badge.good {
    color: var(--color-success-fg);
  }

  .badge.warn {
    color: var(--color-warning);
  }

  .badge.bad {
    color: var(--color-danger);
  }

  .issues {
    margin: var(--space-1) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .breach {
    margin: var(--space-1) 0 0;
    font-size: var(--text-sm);
    color: var(--color-success-fg);
  }

  .breach.exposed {
    color: var(--color-warning);
  }

  .fix-btn {
    display: inline-block;
    margin-top: var(--space-2);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-accent);
    text-decoration: none;
  }

  .error {
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
  }
</style>
