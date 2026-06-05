<script lang="ts">
</script>

<svelte:head>
  <title>Trabalho — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Trabalho</h1>
  <p class="lead">Turnos, sandbox browser, CLI e passkeys.</p>

  <ul class="links">
    <li>
      <a href="/work/sandbox">
        <strong>Browser sandbox</strong>
        <span>Injeção de credenciais sem revelar a password no painel (VAULT-013)</span>
      </a>
    </li>
  </ul>

  <p class="muted">Turnos, passkeys e dispositivos CLI — em breve.</p>
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
