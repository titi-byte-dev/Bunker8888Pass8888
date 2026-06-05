<script lang="ts">
  import type { SecurityHealthReport } from "$lib/darkweb/health";

  interface Props {
    report: SecurityHealthReport;
  }

  let { report }: Props = $props();

  const grade = $derived(
    report.compositeScore >= 75 ? "good" : report.compositeScore >= 50 ? "warn" : "bad",
  );

  function trendLabel(): string {
    if (report.trend === "unknown") return "Primeira medição";
    if (report.trend === "flat") return "Estável";
    const sign = report.trendDelta > 0 ? "+" : "";
    return `${sign}${report.trendDelta} vs última análise`;
  }
</script>

<div class="health-card {grade}">
  <div class="score-block">
    <span class="score">{report.compositeScore}</span>
    <span class="label">Saúde de segurança</span>
  </div>
  <div class="metrics">
    <p><strong>{report.totalLogins}</strong> logins · <strong>{report.scannedCount}</strong> verificados em fugas</p>
    <p>{report.exposedCount} exposta(s) · {report.weakCount} fraca(s) · {report.reusedCount} reutilizada(s)</p>
    <p class="trend" class:up={report.trend === "up"} class:down={report.trend === "down"}>{trendLabel()}</p>
  </div>
</div>

<style>
  .health-card {
    display: flex;
    gap: var(--space-6);
    padding: var(--space-4);
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    margin-bottom: var(--space-6);
  }

  .health-card.good {
    background: var(--color-success-bg);
  }

  .health-card.warn {
    background: color-mix(in srgb, var(--color-warning) 12%, var(--color-bg-surface));
  }

  .health-card.bad {
    background: color-mix(in srgb, var(--color-danger) 10%, var(--color-bg-surface));
  }

  .score-block {
    text-align: center;
    min-width: 5rem;
  }

  .score {
    display: block;
    font-family: var(--font-display);
    font-size: var(--text-3xl);
    font-weight: 700;
    line-height: 1;
  }

  .label {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .metrics p {
    margin: 0 0 var(--space-1);
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .trend {
    margin-top: var(--space-2) !important;
    font-weight: 500;
  }

  .trend.up {
    color: var(--color-success-fg);
  }

  .trend.down {
    color: var(--color-warning);
  }
</style>
