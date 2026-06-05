<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import AdminGate from "$lib/admin/AdminGate.svelte";
  import { hasAdminKey } from "$lib/admin/adminKey";
  import { listWipeAuditEvents, type WipeAuditEvent } from "$lib/admin/api";

  let unlocked = $state(hasAdminKey());
  let events = $state<WipeAuditEvent[]>([]);
  let busy = $state(false);
  let error = $state("");

  async function load() {
    if (!hasAdminKey()) return;
    busy = true;
    error = "";
    try {
      events = await listWipeAuditEvents();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar";
    } finally {
      busy = false;
    }
  }

  onMount(load);

  function onGateChange() {
    unlocked = hasAdminKey();
    if (unlocked) load();
  }
</script>

<svelte:head>
  <title>Auditoria — Admin</title>
</svelte:head>

<section class="page">
  <a href="/admin" class="back">← Administração</a>
  <h1>Auditoria</h1>
  <DocHelpLink />
  <p class="lead">Registo append-only de remote wipe (VAULT-012).</p>

  <AdminGate onUnlocked={onGateChange} />

  {#if unlocked}
    {#if error}
      <p class="error" role="alert">{error}</p>
    {/if}
    {#if busy}
      <p class="muted">A carregar…</p>
    {:else if events.length === 0}
      <p class="muted">Nenhum evento registado.</p>
    {:else}
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Data</th>
              <th>Alvo</th>
              <th>Iniciado por</th>
              <th>Motivo</th>
              <th>Dispositivos</th>
            </tr>
          </thead>
          <tbody>
            {#each events as ev (ev.id)}
              <tr>
                <td>{new Date(ev.created_at).toLocaleString("pt-PT")}</td>
                <td>{ev.target_email || ev.target_user_id.slice(0, 8)}</td>
                <td>{ev.initiated_by}</td>
                <td>{ev.reason || "—"}</td>
                <td>{ev.devices_notified}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  {/if}
</section>

<style>
  .page {
    max-width: 52rem;
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

  .table-wrap {
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
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
    color: var(--color-text-muted);
    font-weight: 500;
    background: var(--color-bg-surface);
  }

  tr:last-child td {
    border-bottom: none;
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
