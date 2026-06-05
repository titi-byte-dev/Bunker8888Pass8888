<script lang="ts">
  /**
   * MetricCard (UI-015) — KPI tipo Stripe Dashboard.
   * Usado em /fin/costs, /crm (funil), admin.
   */
  type Variant = "default" | "warning" | "success";
  type TrendDir = "up" | "down" | "flat";

  interface Props {
    label: string;
    value: string;
    variant?: Variant;
    trend?: { direction: TrendDir; label: string };
  }

  let { label, value, variant = "default", trend }: Props = $props();

  const trendSymbol = $derived(
    trend?.direction === "up" ? "↑" : trend?.direction === "down" ? "↓" : "→",
  );
</script>

<div class="metric {variant}">
  <span class="value">{value}</span>
  <span class="label">{label}</span>
  {#if trend}
    <span class="trend {trend.direction}" aria-label={trend.label}>
      <span aria-hidden="true">{trendSymbol}</span>
      {trend.label}
    </span>
  {/if}
</div>

<style>
  .metric {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    padding: var(--space-3) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    min-width: 0;
  }

  .value {
    font-family: var(--font-display);
    font-size: var(--text-xl);
    font-weight: 600;
    line-height: var(--leading-tight);
    color: var(--color-text);
  }

  .warning .value {
    color: var(--color-danger);
  }

  .success .value {
    color: var(--color-success-fg);
  }

  .label {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .trend {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    margin-top: var(--space-1);
  }

  .trend.up {
    color: var(--color-success-fg);
  }

  .trend.down {
    color: var(--color-danger);
  }
</style>
