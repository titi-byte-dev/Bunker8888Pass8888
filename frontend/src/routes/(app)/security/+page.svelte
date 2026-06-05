<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
</script>

<svelte:head>
  <title>Segurança — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Segurança</h1>
  <p class="lead">Higiene, dispositivos, sessões e confiança operacional.</p>
  <DocHelpLink />

  <ul class="links">
    <li>
      <a href="/security/hygiene">
        <strong>Saúde de segurança</strong>
        <span>Score composto, fugas k-anonymity, alteração forçada (UI-007 / DW-002 / DW-003)</span>
      </a>
    </li>
    <li>
      <a href="/security/devices">
        <strong>Dispositivos e sessões</strong>
        <span>Sessões HTTP, passkeys e certificados CLI</span>
      </a>
    </li>
    <li>
      <a href="/security/sentinel">
        <strong>Sentinel Mode</strong>
        <span>Logins impossíveis, step-up passkey e histórico (DW-004)</span>
      </a>
    </li>
    <li>
      <a href="/security/emergency">
        <strong>Acesso de emergência</strong>
        <span>Herdeiro digital, countdown, aprovar/rejeitar (VAULT-016)</span>
      </a>
    </li>
    <li>
      <a href="/security/guardian">
        <strong>Auditoria do Guardião</strong>
        <span>Tools de agentes sem exposição da Master Key (DoD Fase 3)</span>
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
