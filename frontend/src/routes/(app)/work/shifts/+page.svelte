<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { fetchGeofenceStatus, fetchShiftStatus, type GeofenceStatus, type ShiftStatus } from "$lib/work/access-api";
  import {
    formatClockSkewMs,
    formatWeeklySchedule,
    shiftCountdown,
    shiftStatusLabel,
    shiftStatusTone,
  } from "$lib/work/shift-display";
  import { fetchServerTime, isClockSkewAcceptable } from "$lib/vault/shift";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { getCurrentPosition } from "$lib/vault/geofence";

  let shift = $state<ShiftStatus | null>(null);
  let geofence = $state<GeofenceStatus | null>(null);
  let countdown = $state<string | null>(null);
  let clockSkewOk = $state<boolean | null>(null);
  let clockSkewMs = $state(0);
  let geoBusy = $state(false);
  let busy = $state(true);
  let error = $state("");

  let tick: ReturnType<typeof setInterval> | undefined;

  async function loadGeofence(withGps = false) {
    let pos = null;
    if (withGps) {
      geoBusy = true;
      try {
        pos = await getCurrentPosition();
      } catch (e) {
        error = e instanceof Error ? e.message : "GPS indisponível";
        geoBusy = false;
        return;
      }
      geoBusy = false;
    }
    geofence = await fetchGeofenceStatus(pos);
  }

  async function refresh() {
    error = "";
    try {
      const [shiftStatus, serverTime] = await Promise.all([
        fetchShiftStatus(),
        fetchServerTime(""),
      ]);
      shift = shiftStatus;
      countdown = shiftCountdown(new Date(shiftStatus.server_time), shiftStatus);
      clockSkewMs = Date.now() - serverTime.unix_ms;
      clockSkewOk = isClockSkewAcceptable(
        Date.now(),
        serverTime.unix_ms,
        shiftStatus.max_clock_skew_seconds,
      );
      if (!geofence) {
        await loadGeofence(false);
      }
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar políticas";
    } finally {
      busy = false;
    }
  }

  onMount(() => {
    refresh();
    tick = setInterval(refresh, 15_000);
  });

  onDestroy(() => {
    if (tick) clearInterval(tick);
  });

  const scheduleRows = $derived(shift ? formatWeeklySchedule(shift.schedule) : []);
  const shiftTone = $derived(
    shift ? shiftStatusTone(shift, shift.within_shift) : "neutral",
  );
</script>

<svelte:head>
  <title>Turnos e geofence — AegisPass</title>
</svelte:head>

<section class="page">
  <a href="/work" class="back">← Trabalho</a>
  <h1>Turnos e geofence</h1>
  <DocHelpLink />
  <p class="lead">
    O servidor valida horário (NTP) e zona geográfica antes de permitir acesso ao cofre.
  </p>

  {#if error}
    <p class="error" role="alert">{error}</p>
  {/if}

  {#if busy}
    <p class="muted">A carregar…</p>
  {:else if shift}
    <section class="block">
      <h2>Turno</h2>
      <div class="status-row">
        <span class="pill {shiftTone}">{shiftStatusLabel(shift, shift.within_shift)}</span>
        {#if shift.enabled && countdown}
          <span class="meta">Expira em {countdown}</span>
        {/if}
      </div>
      <dl class="kv">
        <dt>Activo</dt>
        <dd>{shift.enabled ? "Sim" : "Não"}</dd>
        <dt>Fuso</dt>
        <dd>{shift.timezone || "UTC"}</dd>
        <dt>Servidor (UTC)</dt>
        <dd>{new Date(shift.server_time).toLocaleString("pt-PT")}</dd>
        <dt>Desvio relógio</dt>
        <dd class:bad={clockSkewOk === false}>
          {formatClockSkewMs(clockSkewMs)}
          {#if clockSkewOk === false}
            — acima do limite ({shift.max_clock_skew_seconds}s)
          {:else if clockSkewOk}
            — OK
          {/if}
        </dd>
      </dl>
      {#if shift.enabled && scheduleRows.length > 0}
        <table class="schedule">
          <thead>
            <tr><th>Dia</th><th>Janelas</th></tr>
          </thead>
          <tbody>
            {#each scheduleRows as row (row.day)}
              <tr>
                <td>{row.day}</td>
                <td>{row.windows}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      {:else if shift.enabled}
        <p class="muted">Sem janelas configuradas — contacta o administrador.</p>
      {/if}
      <p class="hint">Políticas de turno são definidas pelo admin (UI-008).</p>
    </section>

    <section class="block">
      <h2>Geofence</h2>
      {#if geofence}
        <div class="status-row">
          <span class="pill {geofence.enabled ? (geofence.within_fence ? 'ok' : 'warn') : 'neutral'}">
            {#if !geofence.enabled}
              Desactivado
            {:else if geofence.within_fence}
              Dentro da zona
            {:else}
              Fora da zona
            {/if}
          </span>
        </div>
        <dl class="kv">
          <dt>IP cliente</dt>
          <dd>{geofence.client_ip || "—"}</dd>
          {#if geofence.gps_enabled}
            <dt>GPS</dt>
            <dd>
              Raio {geofence.gps_radius_m}m em
              {geofence.gps_lat?.toFixed(4)}, {geofence.gps_lon?.toFixed(4)}
            </dd>
          {/if}
        </dl>
        {#if geofence.enabled && geofence.gps_enabled}
          <button type="button" disabled={geoBusy} onclick={() => loadGeofence(true)}>
            {geoBusy ? "A obter GPS…" : "Validar com localização actual"}
          </button>
        {/if}
      {:else}
        <p class="muted">Geofence indisponível.</p>
      {/if}
    </section>
  {/if}
</section>

<style>
  .page {
    max-width: 42rem;
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
  }

  .block {
    margin-bottom: var(--space-8);
  }

  h2 {
    margin: 0 0 var(--space-3);
    font-size: var(--text-lg);
  }

  .status-row {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    margin-bottom: var(--space-4);
    flex-wrap: wrap;
  }

  .pill {
    display: inline-block;
    padding: var(--space-1) var(--space-3);
    border-radius: var(--radius-full, 999px);
    font-size: var(--text-sm);
    font-weight: 600;
    border: 1px solid var(--color-border);
  }

  .pill.ok {
    background: var(--color-success-bg);
    color: var(--color-success-fg);
    border-color: transparent;
  }

  .pill.warn {
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    border-color: transparent;
  }

  .pill.neutral {
    background: var(--color-bg-surface);
    color: var(--color-text-muted);
  }

  .meta {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  .kv {
    display: grid;
    grid-template-columns: 9rem 1fr;
    gap: var(--space-2);
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
  }

  dt {
    color: var(--color-text-muted);
  }

  dd {
    margin: 0;
    font-family: var(--font-mono);
    font-size: var(--text-xs);
  }

  dd.bad {
    color: var(--color-danger);
  }

  .schedule {
    width: 100%;
    border-collapse: collapse;
    font-size: var(--text-sm);
    margin-bottom: var(--space-3);
  }

  .schedule th,
  .schedule td {
    text-align: left;
    padding: var(--space-2) var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }

  .schedule th {
    color: var(--color-text-muted);
    font-weight: 500;
  }

  .hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    margin: 0;
  }

  button {
    padding: var(--space-2) var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .error {
    padding: var(--space-3);
    margin-bottom: var(--space-4);
    border-radius: var(--radius-sm);
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
</style>
