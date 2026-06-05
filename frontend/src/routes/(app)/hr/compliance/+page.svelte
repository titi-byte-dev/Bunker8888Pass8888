<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { actionLabel, type AuditEntry } from "$lib/hr/audit";
  import {
    fetchAuditLog,
    fetchComplianceReport,
    verifyLocally,
    type ComplianceReport,
  } from "$lib/hr/compliance";

  let locked = $state(false);
  let loading = $state(true);
  let error = $state("");
  let report = $state<ComplianceReport | null>(null);
  let entries = $state<AuditEntry[]>([]);
  let localValid = $state(true);
  let localBroken = $state(0);

  async function load() {
    loading = true;
    error = "";
    try {
      report = await fetchComplianceReport();
      entries = await fetchAuditLog();
      const res = await verifyLocally(entries);
      localValid = res.valid;
      localBroken = res.brokenSeq;
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao gerar relatório";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) load();
    else loading = false;
  });

  function printReport() {
    window.print();
  }

  function fmt(d: string): string {
    return new Date(d).toLocaleString("pt-PT");
  }
</script>

<svelte:head>
  <title>Relatório de Conformidade RGPD — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head no-print">
    <div>
      <p class="eyebrow">HR-008 · Conformidade RGPD</p>
      <h1>Relatório de Conformidade</h1>
      <DocHelpLink />
    </div>
    <div class="head-actions">
      <a class="btn ghost" href="/hr">← Fichas</a>
      <button type="button" class="btn primary" onclick={printReport} disabled={!report}>
        Descarregar PDF
      </button>
    </div>
  </header>

  {#if locked}
    <section class="panel">
      <p>🔒 Cofre bloqueado. Desbloqueia a Master Key para gerar o relatório.</p>
      <a class="btn primary" href="/vault">Ir desbloquear</a>
    </section>
  {:else if loading}
    <p class="muted">A gerar relatório…</p>
  {:else if error}
    <p class="inline-error" role="alert">{error}</p>
  {:else if report}
    <article class="report">
      <header class="report-head">
        <h2>AegisPass · Relatório de Conformidade RGPD</h2>
        <p class="muted">Gerado em {fmt(report.generatedAt)}</p>
      </header>

      <section class="metrics">
        <div class="metric"><span class="n">{report.recordCount}</span> Fichas de empregado</div>
        <div class="metric"><span class="n">{report.activeFieldCount}</span> Campos cifrados activos</div>
        <div class="metric">
          <span class="n">{report.shreddedFieldCount}</span> Campos eliminados (crypto-shred)
        </div>
        <div class="metric">
          <span class="n">{report.certificateCount}</span> Certificados de eliminação
        </div>
        <div class="metric"><span class="n">{report.auditEntryCount}</span> Entradas de auditoria</div>
      </section>

      <section class="integrity">
        <h3>Integridade da cadeia de auditoria (HR-002)</h3>
        <p>
          Servidor:
          <strong class:ok={report.auditChainValid} class:bad={!report.auditChainValid}>
            {report.auditChainValid ? "íntegra" : `partida na entrada #${report.auditBrokenSeq}`}
          </strong>
          · Reverificação no dispositivo:
          <strong class:ok={localValid} class:bad={!localValid}>
            {localValid ? "íntegra ✓" : `partida na entrada #${localBroken}`}
          </strong>
        </p>
        <p class="muted small">
          A cadeia é recalculada localmente (sha256 encadeado) — a conformidade não
          depende da palavra do servidor.
        </p>
      </section>

      {#if entries.length > 0}
        <section class="log">
          <h3>Registo de auditoria (append-only)</h3>
          <table>
            <thead>
              <tr><th>#</th><th>Acção</th><th>Detalhe</th><th>Quando</th><th>entry_hash</th></tr>
            </thead>
            <tbody>
              {#each entries as e (e.seq)}
                <tr>
                  <td>{e.seq}</td>
                  <td>{actionLabel(e.action)}</td>
                  <td class="mono">{e.detail}</td>
                  <td>{fmt(e.occurredAt)}</td>
                  <td class="mono hash">{e.entryHash.slice(0, 12)}…</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </section>
      {/if}

      <footer class="report-foot muted small">
        Documento gerado pelo AegisPass. Os conteúdos das fichas permanecem cifrados
        (Zero-Knowledge); este relatório contém apenas metadados e provas de integridade.
      </footer>
    </article>
  {/if}
</section>

<style>
  .page {
    max-width: 52rem;
  }
  .page-head {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: var(--space-4);
    margin-bottom: var(--space-6);
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
  .head-actions {
    display: inline-flex;
    gap: var(--space-2);
  }
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .small {
    font-size: var(--text-xs);
  }

  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4) var(--space-6);
  }

  .report {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-8);
  }
  .report-head {
    border-bottom: 1px solid var(--color-border);
    padding-bottom: var(--space-4);
    margin-bottom: var(--space-5);
  }
  .report-head h2 {
    margin: 0 0 var(--space-1);
    font-family: var(--font-display);
    font-size: var(--text-lg);
  }

  .metrics {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    gap: var(--space-3);
    margin-bottom: var(--space-6);
  }
  .metric {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: var(--space-3) var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
  .metric .n {
    display: block;
    font-family: var(--font-display);
    font-size: var(--text-xl);
    color: var(--color-text);
  }

  .integrity {
    margin-bottom: var(--space-6);
  }
  .integrity h3,
  .log h3 {
    font-size: var(--text-sm);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
    margin: 0 0 var(--space-2);
  }
  strong.ok {
    color: var(--color-success-fg);
  }
  strong.bad {
    color: var(--color-danger);
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
  .mono {
    font-family: var(--font-mono);
  }
  .hash {
    color: var(--color-text-muted);
  }
  .report-foot {
    margin-top: var(--space-6);
    border-top: 1px solid var(--color-border);
    padding-top: var(--space-3);
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
  .btn.ghost {
    background: none;
    color: var(--color-text-muted);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .inline-error {
    font-size: var(--text-sm);
    color: var(--color-danger);
  }

  /* Impressão → PDF: esconde a navegação e neutraliza o cromado. */
  @media print {
    .no-print {
      display: none !important;
    }
    .report {
      border: none;
      padding: 0;
    }
  }
</style>
