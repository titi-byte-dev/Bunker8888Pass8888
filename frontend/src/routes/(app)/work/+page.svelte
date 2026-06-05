<script lang="ts">
  import ShiftStatusCard from "$lib/work/ShiftStatusCard.svelte";
</script>

<svelte:head>
  <title>Trabalho — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Trabalho</h1>
  <p class="lead">Turnos, sandbox browser, CLI e ferramentas operacionais BYOD.</p>

  <ShiftStatusCard />

  <ul class="links">
    <li>
      <a href="/work/shifts">
        <strong>Turnos e geofence</strong>
        <span>Estado horário, countdown, desvio NTP e zona GPS (VAULT-010 / VAULT-011)</span>
      </a>
    </li>
    <li>
      <a href="/work/sandbox">
        <strong>Browser sandbox</strong>
        <span>Injeção de credenciais sem revelar a password no painel (VAULT-013)</span>
      </a>
    </li>
    <li>
      <a href="/work/cli">
        <strong>CLI mTLS</strong>
        <span>Registo de dispositivo, listagem e injecção em scripts (VAULT-017)</span>
      </a>
    </li>
    <li>
      <a href="/security/devices">
        <strong>Dispositivos CLI activos</strong>
        <span>Revogar certificados mTLS registados na conta</span>
      </a>
    </li>
  </ul>
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
