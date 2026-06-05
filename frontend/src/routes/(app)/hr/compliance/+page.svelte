<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    Button,
    DataTable,
    EmptyState,
    MetricCard,
    PageShell,
    Skeleton,
    StatusBanner,
    type DataColumn,
  } from "$lib/ui";
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

  const auditColumns: DataColumn<AuditEntry>[] = [
    { id: "seq", label: "#", mono: true, accessor: (e) => String(e.seq) },
    { id: "action", label: "Acção", accessor: (e) => actionLabel(e.action) },
    { id: "detail", label: "Detalhe", mono: true, accessor: (e) => e.detail },
    { id: "when", label: "Quando", accessor: (e) => fmt(e.occurredAt) },
    { id: "hash", label: "entry_hash", mono: true, muted: true, accessor: (e) => `${e.entryHash.slice(0, 12)}…` },
  ];
</script>

<svelte:head>
  <title>Relatório de Conformidade RGPD — AegisPass</title>
</svelte:head>

<div class="compliance-page">
<PageShell
  title="Relatório de Conformidade"
  taskId="HR-008"
  description="Metadados RGPD e prova de integridade da cadeia de auditoria — conteúdos das fichas permanecem cifrados."
  width="wide"
>
  {#snippet actions()}
    <span class="no-print">
      <DocHelpLink />
      <Button variant="ghost" size="sm" href="/hr">← Fichas</Button>
      <Button onclick={printReport} disabled={!report}>Descarregar PDF</Button>
    </span>
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para gerar o relatório.">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else if loading}
    <Skeleton variant="block" height="8rem" />
    <Skeleton variant="table" rows={5} cols={5} />
  {:else if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {:else if report}
    <article class="report">
      <header class="report-head">
        <h2>AegisPass · Relatório de Conformidade RGPD</h2>
        <p class="muted">Gerado em {fmt(report.generatedAt)}</p>
      </header>

      <section class="metrics">
        <MetricCard label="fichas de empregado" value={String(report.recordCount)} />
        <MetricCard label="campos cifrados activos" value={String(report.activeFieldCount)} />
        <MetricCard label="campos eliminados" value={String(report.shreddedFieldCount)} variant="warning" />
        <MetricCard label="certificados de eliminação" value={String(report.certificateCount)} />
        <MetricCard label="entradas de auditoria" value={String(report.auditEntryCount)} />
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
          <DataTable
            columns={auditColumns}
            rows={entries}
            keyFn={(e) => String(e.seq)}
            dense
          />
        </section>
      {/if}

      <footer class="report-foot muted small">
        Documento gerado pelo AegisPass. Os conteúdos das fichas permanecem cifrados
        (Zero-Knowledge); este relatório contém apenas metadados e provas de integridade.
      </footer>
    </article>
  {/if}
</PageShell>
</div>

<style>
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .small {
    font-size: var(--text-xs);
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
    grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
    gap: var(--space-3);
    margin-bottom: var(--space-6);
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

  .report-foot {
    margin-top: var(--space-6);
    border-top: 1px solid var(--color-border);
    padding-top: var(--space-3);
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
