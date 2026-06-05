<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listSentinelEvents, reasonLabel, type LoginEvent } from "$lib/sentinel/api";
  import { PageShell, StatusBanner } from "$lib/ui";

  let events = $state<LoginEvent[]>([]);
  let alerts24h = $state(0);
  let busy = $state(true);
  let error = $state("");

  async function load() {
    busy = true;
    error = "";
    try {
      const data = await listSentinelEvents();
      events = data.events;
      alerts24h = data.suspiciousLast24h;
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      busy = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Sentinel Mode — AegisPass</title>
</svelte:head>

<PageShell
  title="Sentinel Mode"
  taskId="VAULT-014"
  description="Deteção de logins geograficamente impossíveis (DW-004). Compara GPS entre sessões; se a velocidade implícita exceder ~900 km/h, exige passkey antes de emitir token."
>
  {#snippet actions()}
    <DocHelpLink />
  {/snippet}

  {#if alerts24h > 0}
    <StatusBanner variant="warning">
      {alerts24h} alerta(s) suspeito(s) nas últimas 24 horas.
    </StatusBanner>
  {/if}

  {#if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {/if}

  {#if busy}
    <p class="muted">A carregar…</p>
  {:else if events.length === 0}
    <p class="muted">Nenhum evento de login registado ainda.</p>
  {:else}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Data</th>
            <th>IP</th>
            <th>GPS</th>
            <th>Estado</th>
            <th>Notas</th>
          </tr>
        </thead>
        <tbody>
          {#each events as ev (ev.id)}
            <tr class:suspicious={ev.suspicious}>
              <td>{new Date(ev.created_at).toLocaleString("pt-PT")}</td>
              <td class="mono">{ev.client_ip || "—"}</td>
              <td class="mono">
                {#if ev.geo_lat != null && ev.geo_lon != null}
                  {ev.geo_lat.toFixed(2)}, {ev.geo_lon.toFixed(2)}
                {:else}
                  —
                {/if}
              </td>
              <td>
                {#if ev.success}
                  <span class="ok">OK</span>
                {:else if ev.step_up_required}
                  <span class="warn">Step-up</span>
                {:else}
                  <span class="fail">Falhou</span>
                {/if}
              </td>
              <td>
                {#if ev.suspicious}
                  {reasonLabel(ev.reason)}
                {:else}
                  —
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  <p class="hint">
    GPS no browser pode ser falsificado; combina Sentinel com turnos, geofence e passkeys.
  </p>
</PageShell>

<style>
  .table-wrap {
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    margin-bottom: var(--space-4);
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

  tr.suspicious td {
    background: color-mix(in srgb, var(--color-danger) 6%, transparent);
  }

  th {
    color: var(--color-text-muted);
    font-weight: 500;
    background: var(--color-bg-surface);
  }

  .mono {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
  }

  .ok {
    color: var(--color-success-fg);
  }

  .warn {
    color: var(--color-danger);
  }

  .fail {
    color: var(--color-text-muted);
  }

  .hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    line-height: 1.5;
  }

  .muted {
    font-size: var(--text-sm);
  }

  .muted {
    color: var(--color-text-muted);
  }
</style>
