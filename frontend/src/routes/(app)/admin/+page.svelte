<script lang="ts">
  import AdminGate from "$lib/admin/AdminGate.svelte";
  import { hasAdminKey } from "$lib/admin/adminKey";

  let unlocked = $state(hasAdminKey());

  function refreshGate() {
    unlocked = hasAdminKey();
  }
</script>

<svelte:head>
  <title>Administração — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Administração</h1>
  <p class="lead">Utilizadores, políticas de acesso, auditoria e remote wipe (UI-008).</p>

  <AdminGate onUnlocked={refreshGate} />

  {#if unlocked}
    <ul class="links">
      <li>
        <a href="/admin/users">
          <strong>Utilizadores</strong>
          <span>Lista, turnos, geofence e remote wipe por colaborador</span>
        </a>
      </li>
      <li>
        <a href="/admin/audit">
          <strong>Auditoria</strong>
          <span>Eventos de remote wipe (append-only)</span>
        </a>
      </li>
    </ul>
  {/if}
</section>

<style>
  .page {
    max-width: 36rem;
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

  .links {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .links a {
    display: block;
    padding: var(--space-4);
    margin-bottom: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    text-decoration: none;
    color: inherit;
    background: var(--color-bg-surface);
  }

  .links a:hover {
    border-color: var(--color-accent);
  }

  .links strong {
    display: block;
    margin-bottom: var(--space-1);
  }

  .links span {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
