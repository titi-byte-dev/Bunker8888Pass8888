<script lang="ts">
  import { onMount } from "svelte";
  import { analyzeVaultHygiene, type HygieneSummary, type ItemHygieneResult } from "$lib/vault/hygiene";
  import { checkPasswordBreached, type BreachCheckResult } from "$lib/darkweb/breach";
  import { loadDecodedLogins } from "$lib/vault/ui";

  let summary = $state<HygieneSummary | null>(null);
  let passwordsByItem = $state<Map<string, string>>(new Map());
  let breachByItem = $state<Map<string, BreachCheckResult>>(new Map());
  let breachBusy = $state(false);
  let breachError = $state("");
  let busy = $state(true);
  let error = $state("");

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

  /** Verificação k-anonymity (DW-001) — só corre quando o utilizador pede. */
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
  <title>Higiene — AegisPass</title>
</svelte:head>

<section class="page">
  <a href="/security" class="back">← Segurança</a>
  <h1>Higiene de passwords</h1>
  <p class="lead">
    Análise 100% no cliente — o servidor nunca vê as tuas passwords. A verificação
    de fugas (Dark Web) usa <strong>k-anonymity</strong>: só saem 5 caracteres do hash.
  </p>

  {#if busy}
    <p class="muted">A analisar cofre…</p>
  {:else if error}
    <p class="error" role="alert">{error}</p>
  {:else if summary}
    <div class="score-card {scoreClass(summary.overallScore)}">
      <span class="score-value">{summary.overallScore}</span>
      <div>
        <strong>Score global</strong>
        <p>{summary.totalLogins} login{summary.totalLogins === 1 ? "" : "s"} · {summary.weakCount} fraca(s) · {summary.reusedCount} reutilizada(s)</p>
      </div>
    </div>

    <div class="breach-panel">
      <h2>Preview Dark Web (DW-001)</h2>
      <p class="hint">
        Clica para verificar fugas conhecidas via API pública (HIBP). Comparação final
        local — hash completo nunca sai do dispositivo.
      </p>
      <button type="button" onclick={scanBreaches} disabled={breachBusy || summary.items.length === 0}>
        {breachBusy ? "A verificar…" : "Verificar fugas"}
      </button>
      {#if breachError}
        <p class="error" role="alert">{breachError}</p>
      {/if}
    </div>

    <ul class="items">
      {#each summary.items as item (item.itemId)}
        {@const breach = breachByItem.get(item.itemId)}
        <li>
          <div class="row">
            <a href="/vault/{item.itemId}">{item.title}</a>
            <span class="badge {scoreClass(item.score)}">{item.score}</span>
          </div>
          <p class="issues">{issueLabel(item.issues)}</p>
          {#if breach}
            <p class="breach" class:exposed={breach.breached}>
              {#if breach.breached}
                ⚠ Exposta em fugas ({breach.exposureCount.toLocaleString("pt-PT")}×)
              {:else}
                ✓ Sem ocorrências conhecidas
              {/if}
            </p>
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

  .score-card {
    display: flex;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-4);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    margin-bottom: var(--space-6);
  }

  .score-card.good {
    background: var(--color-success-bg);
  }

  .score-card.warn {
    background: color-mix(in srgb, var(--color-warning) 12%, var(--color-bg-surface));
  }

  .score-card.bad {
    background: color-mix(in srgb, var(--color-danger) 10%, var(--color-bg-surface));
  }

  .score-value {
    font-family: var(--font-display);
    font-size: var(--text-3xl);
    font-weight: 700;
    min-width: 3rem;
    text-align: center;
  }

  .score-card p {
    margin: var(--space-1) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
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
