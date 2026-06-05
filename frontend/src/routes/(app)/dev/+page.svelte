<script lang="ts">
</script>

<svelte:head>
  <title>Playground — AegisPass</title>
</svelte:head>

<section class="page">
  <h1>Desenvolvimento</h1>
  <p class="lead">Ferramentas de teste e referência visual — só em build dev.</p>

  <ul class="links">
    <li>
      <a href="/dev/playground">
        <strong>Cofre playground</strong>
        <span>CRUD end-to-end legacy (VAULT-019)</span>
      </a>
    </li>
    <li>
      <a href="/dev/components">
        <strong>Catálogo de componentes</strong>
        <span>Tokens, botões, forms e previews (UI-010)</span>
      </a>
    </li>
  </ul>

  <p class="hint">Motion demo e ArgonProgress estão no catálogo de componentes.</p>
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

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
</style>
