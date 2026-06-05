<script lang="ts">
</script>

<svelte:head>
  <title>Segurança — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Segurança</h1>
  <p class="lead">Higiene, dispositivos, sessões e confiança operacional.</p>

  <ul class="links">
    <li>
      <a href="/security/hygiene">
        <strong>Higiene de passwords</strong>
        <span>Score no cliente + verificação de fugas k-anonymity (UI-007 / DW-001)</span>
      </a>
    </li>
    <li>
      <a href="/security/emergency">
        <strong>Acesso de emergência</strong>
        <span>Herdeiro digital, countdown, aprovar/rejeitar (VAULT-016)</span>
      </a>
    </li>
  </ul>

  <p class="muted">Dispositivos e sessões — em breve.</p>
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
    margin: 0 0 var(--space-6);
  }

  .links a {
    display: block;
    padding: var(--space-4);
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

  .muted {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
