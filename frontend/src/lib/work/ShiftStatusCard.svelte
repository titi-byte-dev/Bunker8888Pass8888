<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { fetchShiftStatus } from "$lib/work/access-api";
  import {
    shiftCountdown,
    shiftStatusLabel,
    shiftStatusTone,
  } from "$lib/work/shift-display";

  let within = $state<boolean | null>(null);
  let enabled = $state(false);
  let countdown = $state<string | null>(null);
  let timezone = $state("UTC");
  let busy = $state(true);
  let error = $state("");

  let tick: ReturnType<typeof setInterval> | undefined;

  async function refresh() {
    try {
      const status = await fetchShiftStatus();
      enabled = status.enabled;
      within = status.within_shift;
      timezone = status.timezone || "UTC";
      countdown = shiftCountdown(new Date(status.server_time), status);
      error = "";
    } catch (e) {
      error = e instanceof Error ? e.message : "Indisponível";
      within = null;
    } finally {
      busy = false;
    }
  }

  onMount(() => {
    refresh();
    tick = setInterval(refresh, 30_000);
  });

  onDestroy(() => {
    if (tick) clearInterval(tick);
  });

  const tone = $derived(
    within === null ? "neutral" : shiftStatusTone({ enabled, timezone, schedule: {}, max_clock_skew_seconds: 0 }, within),
  );
  const label = $derived(
    within === null
      ? "Turno"
      : shiftStatusLabel({ enabled, timezone, schedule: {}, max_clock_skew_seconds: 0 }, within),
  );
</script>

{#if busy}
  <p class="card muted">A verificar turno…</p>
{:else if error}
  <p class="card warn">{error}</p>
{:else}
  <a href="/work/shifts" class="card {tone}">
    <div class="card-head">
      <strong>{label}</strong>
      {#if enabled && countdown}
        <span class="badge">Fim em {countdown}</span>
      {/if}
    </div>
    {#if enabled}
      <span class="sub">{timezone}</span>
    {:else}
      <span class="sub">Sem restrição horária activa</span>
    {/if}
  </a>
{/if}

<style>
  .card {
    display: block;
    padding: var(--space-4);
    margin-bottom: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    text-decoration: none;
    color: inherit;
    background: var(--color-bg-surface);
  }

  .card:hover {
    border-color: var(--color-accent);
  }

  .card.ok {
    border-color: color-mix(in srgb, var(--color-success-fg) 40%, var(--color-border));
  }

  .card.warn {
    border-color: color-mix(in srgb, var(--color-danger) 35%, var(--color-border));
  }

  .card-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-2);
    margin-bottom: var(--space-1);
  }

  .badge {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  .sub {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .muted,
  .warn {
    font-size: var(--text-sm);
  }

  .warn {
    color: var(--color-danger);
  }
</style>
