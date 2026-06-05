<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
</script>

<svelte:head>
  <title>CLI mTLS — AegisPass</title>
</svelte:head>

<section class="page">
  <a href="/work" class="back">← Trabalho</a>
  <h1>CLI mTLS</h1>
  <DocHelpLink />
  <p class="lead">
    A CLI <code>aegis</code> injecta segredos em scripts via certificado mTLS — a Master Password
    fica só em memória local (Zero-Knowledge).
  </p>

  <section class="block">
    <h2>1. Registar dispositivo</h2>
    <pre class="code"><code>aegis device register --name "MacBook Pro" --email teu@email.com</code></pre>
    <p class="hint">Pedirá a Master Password e grava o certificado em <code>~/.aegis</code>.</p>
  </section>

  <section class="block">
    <h2>2. Listar itens do cofre</h2>
    <pre class="code"><code>aegis list --type login</code></pre>
  </section>

  <section class="block">
    <h2>3. Injectar segredo num comando</h2>
    <pre class="code"><code>aegis run --item &lt;id&gt; --field password --env SECRET -- ./script.sh</code></pre>
    <p class="hint">A variável <code>SECRET</code> existe só durante a execução do processo filho.</p>
  </section>

  <section class="block">
    <h2>Variáveis de ambiente</h2>
    <dl class="kv">
      <dt><code>AEGIS_API_URL</code></dt>
      <dd>API HTTP (default <code>http://localhost:8080</code>)</dd>
      <dt><code>AEGIS_MTLS_URL</code></dt>
      <dd>API mTLS (default <code>https://localhost:8443</code>)</dd>
      <dt><code>AEGIS_CONFIG_DIR</code></dt>
      <dd>Pasta de certificados (default <code>~/.aegis</code>)</dd>
      <dt><code>AEGIS_MTLS_INSECURE=1</code></dt>
      <dd>Ignora TLS inválido — <strong>só desenvolvimento</strong></dd>
    </dl>
  </section>

  <p class="hint">
    Revoga dispositivos comprometidos em
    <a href="/security/devices">Dispositivos e sessões</a>.
  </p>
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
    line-height: 1.5;
  }

  .block {
    margin-bottom: var(--space-6);
  }

  h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-base);
  }

  .code {
    margin: 0 0 var(--space-2);
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    border: 1px solid var(--color-border);
    overflow-x: auto;
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    line-height: 1.6;
  }

  .kv {
    display: grid;
    grid-template-columns: 11rem 1fr;
    gap: var(--space-2);
    margin: 0;
    font-size: var(--text-sm);
  }

  dt {
    color: var(--color-text-muted);
  }

  dd {
    margin: 0;
    line-height: 1.4;
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: var(--space-2) 0 0;
    line-height: 1.5;
  }

  .hint a {
    color: var(--color-link);
  }

  code {
    font-family: var(--font-mono);
    font-size: 0.9em;
  }
</style>
