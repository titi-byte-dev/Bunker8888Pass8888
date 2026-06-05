<script lang="ts">
  import ArgonProgress from "$lib/auth/ArgonProgress.svelte";
  import {
    CATALOG_SECTIONS,
    MOCK_HEALTH_REPORT,
    SEMANTIC_COLORS,
    TYPE_SCALE,
  } from "$lib/design/catalog";
  import { themeModeLabel, type ThemeMode } from "$lib/design";
  import { loadThemePreference, setThemePreference } from "$lib/design";
  import SecurityHealthCard from "$lib/security/SecurityHealthCard.svelte";
  import CopyField from "$lib/vault/CopyField.svelte";

  let themeMode = $state<ThemeMode>(loadThemePreference());

  function setTheme(mode: ThemeMode) {
    themeMode = mode;
    setThemePreference(mode);
  }
</script>

<svelte:head>
  <title>Catálogo de componentes — AegisPass</title>
</svelte:head>

<section class="catalog">
  <header class="hero">
    <a href="/dev" class="back">← Playground</a>
    <h1>Catálogo de componentes</h1>
    <p class="lead">
      Design tokens (UI-001) e peças reutilizáveis da app — só visível em desenvolvimento.
    </p>
    <nav class="toc" aria-label="Secções">
      {#each CATALOG_SECTIONS as sec (sec.id)}
        <a href="#{sec.id}">{sec.title}</a>
      {/each}
    </nav>
    <div class="theme-row">
      {#each (["light", "dark", "system"] as const) as mode (mode)}
        <button
          type="button"
          class:active={themeMode === mode}
          onclick={() => setTheme(mode)}
        >
          {themeModeLabel(mode)}
        </button>
      {/each}
    </div>
  </header>

  <section id="typography" class="block">
    <h2>Tipografia</h2>
    <p class="hint">Fontes: <code>--font-ui</code>, <code>--font-display</code>, <code>--font-mono</code></p>
    <div class="type-grid">
      {#each TYPE_SCALE as row (row.token)}
        <div class="type-row" style="font-size: var({row.token})">
          <span class="token">{row.token}</span>
          <span>{row.sample}</span>
        </div>
      {/each}
    </div>
    <p class="display-sample">AegisPass — cofre suíço no teu dispositivo</p>
  </section>

  <section id="colors" class="block">
    <h2>Cores semânticas</h2>
    <div class="swatches">
      {#each SEMANTIC_COLORS as c (c.name)}
        <div class="swatch">
          <div
            class="chip"
            style="background: var({c.bg}); color: var({c.var})"
          ></div>
          <span>{c.name}</span>
        </div>
      {/each}
    </div>
    <div class="surfaces">
      <div class="surface base">base</div>
      <div class="surface elevated">elevated</div>
      <div class="surface inset">inset</div>
    </div>
  </section>

  <section id="spacing" class="block">
    <h2>Spacing & radius</h2>
    <div class="space-demo">
      {#each [
        { label: "1", v: "var(--space-1)" },
        { label: "2", v: "var(--space-2)" },
        { label: "3", v: "var(--space-3)" },
        { label: "4", v: "var(--space-4)" },
        { label: "6", v: "var(--space-6)" },
        { label: "8", v: "var(--space-8)" },
      ] as sp (sp.label)}
        <div class="space-bar" style="width: {sp.v}" title="--space-{sp.label}"></div>
      {/each}
    </div>
    <div class="radius-demo">
      <div class="r sm">sm</div>
      <div class="r md">md</div>
      <div class="r lg">lg</div>
    </div>
  </section>

  <section id="buttons" class="block">
    <h2>Botões</h2>
    <div class="btn-row">
      <button type="button" class="btn primary">Primário</button>
      <button type="button" class="btn secondary">Secundário</button>
      <button type="button" class="btn danger">Perigo</button>
      <button type="button" class="btn" disabled>Desactivado</button>
    </div>
  </section>

  <section id="forms" class="block">
    <h2>Formulários</h2>
    <label class="field">
      Email
      <input type="email" placeholder="tu@empresa.com" />
    </label>
    <label class="field check">
      <input type="checkbox" checked />
      Opção activa
    </label>
    <CopyField label="Password (secret)" value="hunter2-demo" secret copyTimeoutMs={8000} />
    <CopyField label="Nota pública" value="Sem segredo — clipboard 8s" />
  </section>

  <section id="components" class="block">
    <h2>Componentes vivos</h2>
    <h3>SecurityHealthCard</h3>
    <SecurityHealthCard report={MOCK_HEALTH_REPORT} />
    <h3>ArgonProgress</h3>
    <ArgonProgress active label="A derivar chaves (demo Argon2id)…" />
    <h3>Estados</h3>
    <p class="status ok">Sucesso operacional</p>
    <p class="status warn">Atenção — rever política</p>
    <p class="status err">Erro ou acção destrutiva</p>
  </section>
</section>

<style>
  .catalog {
    max-width: 48rem;
    padding-bottom: var(--space-12);
  }

  .hero {
    margin-bottom: var(--space-8);
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
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
  }

  .toc {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    margin-bottom: var(--space-4);
  }

  .toc a {
    font-size: var(--text-xs);
    color: var(--color-link);
    text-decoration: none;
    padding: var(--space-1) var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
  }

  .theme-row {
    display: flex;
    gap: var(--space-2);
  }

  .theme-row button {
    padding: var(--space-1) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-surface);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  .theme-row button.active {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .block {
    margin-bottom: var(--space-10);
    padding-top: var(--space-2);
  }

  h2 {
    margin: 0 0 var(--space-4);
    font-size: var(--text-lg);
    font-family: var(--font-display);
  }

  h3 {
    margin: var(--space-4) 0 var(--space-2);
    font-size: var(--text-base);
  }

  .hint {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0 0 var(--space-3);
  }

  .type-grid {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    overflow: hidden;
  }

  .type-row {
    display: flex;
    justify-content: space-between;
    padding: var(--space-2) var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }

  .type-row:last-child {
    border-bottom: none;
  }

  .token {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .display-sample {
    margin-top: var(--space-4);
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .swatches {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-4);
    margin-bottom: var(--space-4);
  }

  .swatch {
    text-align: center;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .chip {
    width: 4rem;
    height: 4rem;
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
    margin-bottom: var(--space-1);
  }

  .surfaces {
    display: flex;
    gap: var(--space-3);
    flex-wrap: wrap;
  }

  .surface {
    padding: var(--space-4);
    border-radius: var(--radius-md);
    font-size: var(--text-xs);
    border: 1px solid var(--color-border);
    min-width: 6rem;
    text-align: center;
  }

  .surface.base {
    background: var(--color-bg-base);
  }

  .surface.elevated {
    background: var(--color-bg-elevated);
  }

  .surface.inset {
    background: var(--color-bg-inset);
  }

  .space-demo {
    display: flex;
    align-items: flex-end;
    gap: var(--space-2);
    margin-bottom: var(--space-4);
  }

  .space-bar {
    height: var(--space-8);
    min-width: 4px;
    background: var(--color-accent);
    border-radius: var(--radius-sm);
  }

  .radius-demo {
    display: flex;
    gap: var(--space-3);
  }

  .r {
    padding: var(--space-3) var(--space-4);
    background: var(--color-bg-surface);
    border: 1px solid var(--color-border);
    font-size: var(--text-xs);
  }

  .r.sm {
    border-radius: var(--radius-sm);
  }

  .r.md {
    border-radius: var(--radius-md);
  }

  .r.lg {
    border-radius: var(--radius-lg);
  }

  .btn-row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
  }

  .btn {
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    cursor: pointer;
    border: 1px solid var(--color-border);
  }

  .btn.primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
    border-color: transparent;
  }

  .btn.secondary {
    background: var(--color-bg-surface);
    color: var(--color-text);
  }

  .btn.danger {
    background: transparent;
    color: var(--color-danger);
    border-color: var(--color-danger);
  }

  .btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .field {
    display: block;
    font-size: var(--text-sm);
    margin-bottom: var(--space-3);
  }

  .field input[type="email"] {
    display: block;
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    color: var(--color-text);
    box-sizing: border-box;
  }

  .field.check {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }

  .status {
    padding: var(--space-3);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
    margin-bottom: var(--space-2);
  }

  .status.ok {
    background: var(--color-success-bg);
    color: var(--color-success-fg);
  }

  .status.warn {
    background: color-mix(in srgb, var(--color-warning) 15%, transparent);
    color: var(--color-warning);
  }

  .status.err {
    background: color-mix(in srgb, var(--color-danger) 12%, transparent);
    color: var(--color-danger);
  }

  code {
    font-family: var(--font-mono);
    font-size: 0.9em;
  }
</style>
