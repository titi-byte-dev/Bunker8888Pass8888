<script lang="ts">
  import { loadSessionToken } from "$lib/session";
  import { buildNoteLink, encryptNoteContent, generateNoteKey } from "$lib/share/burnNote";
  import { burnNoteManually, createBurnNote } from "$lib/share/burnNoteApi";

  const TTL_OPTIONS = [
    { label: "10 minutos", seconds: 600 },
    { label: "1 hora", seconds: 3600 },
    { label: "24 horas", seconds: 86400 },
    { label: "7 dias", seconds: 604800 },
  ];

  let secret = $state("");
  let passphrase = $state("");
  let ttlSeconds = $state(86400);

  let creating = $state(false);
  let error = $state("");
  let link = $state("");
  let noteId = $state("");
  let burnToken = $state("");
  let expiresAt = $state("");
  let copied = $state(false);
  let burned = $state(false);
  let burnError = $state("");

  async function create(event: SubmitEvent) {
    event.preventDefault();
    const value = secret;
    if (!value) return;
    const token = loadSessionToken();
    if (!token) {
      error = "Sessao expirada — inicia sessao de novo.";
      return;
    }
    creating = true;
    error = "";
    link = "";
    copied = false;
    burned = false;
    burnError = "";
    try {
      const key = generateNoteKey();
      const { ciphertext, salt } = await encryptNoteContent(key, value, passphrase || undefined);
      const created = await createBurnNote(token, ciphertext, ttlSeconds);
      noteId = created.id;
      burnToken = created.burn_token;
      link = buildNoteLink(window.location.origin, created.id, key, salt);
      expiresAt = new Date(created.expires_at).toLocaleString("pt-PT");
      secret = "";
      passphrase = "";
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao criar nota";
    } finally {
      creating = false;
    }
  }

  async function copy() {
    try {
      await navigator.clipboard.writeText(link);
      copied = true;
      setTimeout(() => (copied = false), 2000);
    } catch {
      copied = false;
    }
  }

  async function burnNow() {
    burnError = "";
    try {
      await burnNoteManually(noteId, burnToken);
      burned = true;
    } catch (e) {
      // Se ja foi lida/destruida, tambem conta como queimada.
      burned = true;
      burnError = e instanceof Error ? e.message : "";
    }
  }

  function reset() {
    link = "";
    noteId = "";
    burnToken = "";
    expiresAt = "";
    error = "";
    burned = false;
    burnError = "";
  }
</script>

<svelte:head>
  <title>Notas Auto-Destrutivas — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">SHARE-005 · Notas que ardem após leitura</p>
      <h1>Notas Auto-Destrutivas</h1>
    </div>
    <p class="lead">
      Uma nota cifrada que se lê <strong>uma única vez</strong> e arde a seguir. A
      chave vai no fragmento do link (nunca chega ao servidor). Podes juntar uma
      passphrase (2.ª camada) combinada por um canal à parte e destruir a nota a
      qualquer momento antes de ser lida.
    </p>
    <a class="back" href="/team">← Identidade de partilha</a>
  </header>

  {#if link}
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">Nota gerada</p>
        <span class="pill" class:on={!burned} class:gone={burned}>{burned ? "Destruída" : "Ativa"}</span>
      </div>

      {#if burned}
        <p class="panel-body">A nota foi destruída. O link já não abre nada.</p>
      {:else}
        <div class="link-row">
          <input class="mono" type="text" readonly value={link} aria-label="Link da nota" />
          <button type="button" class="btn primary" onclick={copy}>
            {copied ? "Copiado!" : "Copiar"}
          </button>
        </div>
        <dl class="props">
          <div class="prop">
            <dt>Expira</dt>
            <dd class="mono">{expiresAt}</dd>
          </div>
        </dl>
        <p class="panel-foot warn">
          ⚠️ Copia o link agora. A chave de cifra está no próprio link (depois do #)
          — o servidor não a tem. A nota arde após a primeira leitura.
        </p>
        <div class="actions">
          <button type="button" class="btn danger-btn" onclick={burnNow}>Destruir já</button>
          <button type="button" class="btn secondary" onclick={reset}>Criar outra nota</button>
        </div>
        {#if burnError}<p class="inline-note" role="status">{burnError}</p>{/if}
      {/if}

      {#if burned}
        <button type="button" class="btn secondary" onclick={reset}>Criar outra nota</button>
      {/if}
    </section>
  {:else}
    <section class="panel">
      <div class="panel-head"><p class="eyebrow">Nova nota</p></div>
      <form onsubmit={create}>
        <label class="field">
          <span>Nota</span>
          <textarea
            bind:value={secret}
            rows="3"
            placeholder="Mensagem, password, token…"
            disabled={creating}
          ></textarea>
        </label>
        <div class="row">
          <label class="field">
            <span>Expira após</span>
            <select bind:value={ttlSeconds} disabled={creating}>
              {#each TTL_OPTIONS as opt (opt.seconds)}
                <option value={opt.seconds}>{opt.label}</option>
              {/each}
            </select>
          </label>
          <label class="field">
            <span>Passphrase (opcional)</span>
            <input
              type="text"
              bind:value={passphrase}
              placeholder="2.ª camada, dita à parte"
              autocomplete="off"
              spellcheck="false"
              disabled={creating}
            />
          </label>
        </div>
        {#if error}<p class="inline-error" role="alert">{error}</p>{/if}
        <button type="submit" class="btn primary" disabled={creating || !secret}>
          {creating ? "A cifrar…" : "Gerar nota"}
        </button>
      </form>
    </section>

    <section class="panel">
      <div class="panel-head"><p class="eyebrow">Como funciona</p></div>
      <ol class="steps">
        <li><span>1</span> A nota é cifrada no teu dispositivo com uma chave aleatória.</li>
        <li><span>2</span> Com passphrase, ciframos <em>outra vez</em> — nem o link sozinho a abre.</li>
        <li><span>3</span> O servidor guarda o <em>ciphertext</em> só em RAM, com TTL.</li>
        <li><span>4</span> Ao ler, a nota é revelada <em>uma vez</em> e apagada. Sem rasto.</li>
      </ol>
    </section>
  {/if}
</section>

<style>
  .page {
    max-width: 44rem;
  }
  .page-head {
    margin-bottom: var(--space-8);
  }
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
  }
  h1 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    line-height: var(--leading-tight);
  }
  .lead {
    margin: var(--space-3) 0 0;
    max-width: 38rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .back {
    display: inline-block;
    margin-top: var(--space-3);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    text-decoration: none;
  }
  .back:hover {
    color: var(--color-text);
  }

  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4) var(--space-6);
    margin-bottom: var(--space-4);
  }
  .panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    margin-bottom: var(--space-3);
  }
  .panel-head .eyebrow {
    margin: 0;
  }
  .panel-body {
    margin: 0 0 var(--space-2);
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .panel-foot {
    margin: var(--space-4) 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text-muted);
  }
  .panel-foot.warn {
    color: var(--color-text);
  }

  .field {
    display: block;
    margin-bottom: var(--space-3);
  }
  .field > span {
    display: block;
    margin-bottom: var(--space-1);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  textarea,
  select,
  input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }
  textarea {
    resize: vertical;
  }
  textarea:focus-visible,
  select:focus-visible,
  input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  .row {
    display: flex;
    gap: var(--space-3);
  }
  .row .field {
    flex: 1;
  }

  .link-row {
    display: flex;
    gap: var(--space-2);
    margin-bottom: var(--space-3);
  }
  .link-row input {
    flex: 1;
    min-width: 0;
    font-size: var(--text-xs);
  }

  .props {
    margin: 0;
    border-top: 1px solid var(--color-border);
  }
  .prop {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-3) 0;
    border-bottom: 1px solid var(--color-border);
  }
  .prop dt {
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  .prop dd {
    margin: 0;
  }
  .mono {
    font-family: var(--font-mono);
  }

  .pill {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
    color: var(--color-text-muted);
  }
  .pill.on {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
  }
  .pill.gone {
    color: var(--color-danger);
    border-color: var(--color-danger);
  }

  .actions {
    display: flex;
    gap: var(--space-2);
  }

  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
  }
  .btn.primary {
    background: var(--color-accent);
    color: var(--color-accent-fg);
    border-color: transparent;
  }
  .btn.primary:hover:not(:disabled) {
    filter: brightness(1.08);
  }
  .btn.secondary {
    background: var(--color-bg-elevated);
    color: var(--color-text);
  }
  .btn.secondary:hover:not(:disabled) {
    background: var(--color-accent-muted);
  }
  .btn.danger-btn {
    color: var(--color-danger);
    border-color: var(--color-danger);
    background: none;
  }
  .btn.danger-btn:hover {
    background: var(--color-danger);
    color: var(--color-accent-fg);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }

  .inline-error {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
  .inline-note {
    margin: var(--space-2) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .steps {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }
  .steps li {
    display: flex;
    gap: var(--space-3);
    align-items: baseline;
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .steps span {
    flex-shrink: 0;
    width: 1.25rem;
    height: 1.25rem;
    display: inline-grid;
    place-items: center;
    border-radius: 50%;
    background: var(--color-accent-muted);
    color: var(--color-accent);
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-mono);
  }
  .steps em {
    font-style: normal;
    color: var(--color-accent);
  }
</style>
