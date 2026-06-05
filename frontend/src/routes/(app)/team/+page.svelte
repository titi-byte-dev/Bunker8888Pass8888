<script lang="ts">
  import { onMount } from "svelte";
  import {
    ensureShareIdentity,
    loadShareIdentity,
    lookupRecipient,
    type RecipientLookup,
    type ShareIdentity,
  } from "$lib/share/setup";

  type Status = "loading" | "locked" | "ready" | "error";

  let status = $state<Status>("loading");
  let identity = $state<ShareIdentity | null>(null);
  let loadError = $state("");

  let activating = $state(false);
  let activateError = $state("");

  let email = $state("");
  let looking = $state(false);
  let lookupError = $state("");
  let recipient = $state<RecipientLookup | null>(null);

  function isLocked(message: string): boolean {
    return message.toLowerCase().includes("bloquead");
  }

  async function load() {
    status = "loading";
    loadError = "";
    try {
      identity = await loadShareIdentity();
      status = "ready";
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha ao carregar partilha";
      if (isLocked(msg)) {
        status = "locked";
      } else {
        loadError = msg;
        status = "error";
      }
    }
  }

  async function activate() {
    activating = true;
    activateError = "";
    try {
      identity = await ensureShareIdentity();
    } catch (e) {
      activateError = e instanceof Error ? e.message : "Falha ao activar partilha";
    } finally {
      activating = false;
    }
  }

  async function lookup(event: SubmitEvent) {
    event.preventDefault();
    const target = email.trim();
    if (!target) return;
    looking = true;
    lookupError = "";
    recipient = null;
    try {
      recipient = await lookupRecipient(target);
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha na procura";
      lookupError = msg.includes("não encontrado")
        ? "Sem chave pública para este email (o colega ainda não activou a partilha)."
        : msg;
    } finally {
      looking = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Equipa — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">SHARE · Partilha blindada</p>
      <h1>Equipa</h1>
    </div>
    <p class="lead">
      Chaves assimétricas por utilizador — a base para partilhar segredos sem o
      servidor alguma vez os ver.
    </p>
  </header>

  {#if status === "loading"}
    <div class="panel muted-panel">A carregar identidade de partilha…</div>
  {:else if status === "locked"}
    <div class="panel locked">
      <p class="panel-title">Cofre bloqueado</p>
      <p class="panel-body">
        A chave de partilha é protegida pela tua Master Password. Desbloqueia o
        cofre para a gerir.
      </p>
      <a class="btn primary" href="/auth/unlock">Desbloquear</a>
    </div>
  {:else if status === "error"}
    <div class="panel danger">
      <p class="panel-title">Erro</p>
      <p class="panel-body">{loadError}</p>
      <button type="button" class="btn secondary" onclick={load}>Tentar de novo</button>
    </div>
  {:else}
    <!-- Painel A — a minha chave -->
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">A minha chave</p>
        <span class="pill" class:on={!!identity} class:off={!identity}>
          {identity ? "Activa" : "Inactiva"}
        </span>
      </div>

      {#if identity}
        <dl class="props">
          <div class="prop">
            <dt>Algoritmo</dt>
            <dd class="mono">{identity.algorithm}</dd>
          </div>
          <div class="prop">
            <dt>Impressão digital</dt>
            <dd class="mono fingerprint">{identity.fingerprint}</dd>
          </div>
        </dl>
        <p class="panel-foot">
          A chave privada está cifrada com a tua Master Password antes de chegar
          ao servidor. Partilha esta impressão digital por um canal de confiança
          para os colegas confirmarem que é mesmo a tua.
        </p>
      {:else}
        <p class="panel-body">
          Ainda não activaste a partilha. Vamos gerar um par de chaves no teu
          dispositivo: a pública fica partilhável, a privada é cifrada com a tua
          Master Password.
        </p>
        {#if activateError}
          <p class="inline-error" role="alert">{activateError}</p>
        {/if}
        <button type="button" class="btn primary" onclick={activate} disabled={activating}>
          {activating ? "A gerar par de chaves…" : "Activar partilha"}
        </button>
      {/if}
    </section>

    <!-- Painel B — verificar chave de colega -->
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">Verificar colega</p>
      </div>
      <p class="panel-body">
        Procura a chave pública de um colega para lhe partilhares um segredo.
      </p>
      <form class="lookup" onsubmit={lookup}>
        <input
          type="email"
          bind:value={email}
          placeholder="colega@empresa.pt"
          autocomplete="off"
          spellcheck="false"
          disabled={looking}
        />
        <button type="submit" class="btn secondary" disabled={looking || !email.trim()}>
          {looking ? "A procurar…" : "Procurar"}
        </button>
      </form>

      {#if lookupError}
        <p class="inline-note" role="status">{lookupError}</p>
      {:else if recipient}
        <dl class="props">
          <div class="prop">
            <dt>Email</dt>
            <dd class="mono">{recipient.email}</dd>
          </div>
          <div class="prop">
            <dt>Impressão digital</dt>
            <dd class="mono fingerprint">{recipient.fingerprint}</dd>
          </div>
        </dl>
        <p class="panel-foot">
          Confirma esta impressão digital com o colega por um canal à parte antes
          de partilhar — defende contra troca maliciosa de chaves.
        </p>
      {/if}
    </section>

    <!-- Painel C — como funciona -->
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">Como funciona</p>
      </div>
      <ol class="steps">
        <li><span>1</span> Cada utilizador gera um par de chaves no dispositivo.</li>
        <li><span>2</span> Partilhar = re-cifrar a chave do item para a <em>chave pública</em> do destinatário.</li>
        <li><span>3</span> Só a <em>chave privada</em> dele a abre. O servidor encaminha bytes opacos.</li>
      </ol>
    </section>

    <!-- Painel D — a seguir -->
    <section class="panel next">
      <div class="panel-head">
        <p class="eyebrow">A seguir</p>
        <span class="pill soon">SHARE-002</span>
      </div>
      <p class="panel-body">
        Cofres partilhados por departamento, com permissões e revogação imediata —
        construídos sobre estas chaves.
      </p>
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
    max-width: 36rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
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
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text);
  }

  .panel-foot {
    margin: var(--space-4) 0 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text-muted);
  }

  .muted-panel,
  .locked .panel-body,
  .danger .panel-body {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .panel-title {
    margin: 0 0 var(--space-2);
    font-weight: 600;
    font-size: var(--text-sm);
  }

  /* Property rows — densidade tipo Stripe/Linear */
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
    flex-shrink: 0;
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }

  .prop dd {
    margin: 0;
    text-align: right;
    word-break: break-word;
  }

  .mono {
    font-family: var(--font-mono);
    font-size: var(--text-sm);
  }

  .fingerprint {
    font-size: var(--text-xs);
    letter-spacing: 0.04em;
    color: var(--color-text);
    line-height: var(--leading-body);
  }

  /* Status pills */
  .pill {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
  }

  .pill.on {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    border-color: transparent;
  }

  .pill.off {
    color: var(--color-text-muted);
  }

  .pill.soon {
    color: var(--color-text-muted);
    font-family: var(--font-mono);
  }

  /* Lookup form */
  .lookup {
    display: flex;
    gap: var(--space-2);
  }

  .lookup input {
    flex: 1;
    min-width: 0;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
  }

  .lookup input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }

  /* Botões */
  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    text-decoration: none;
    cursor: pointer;
    transition: background-color var(--duration-fast) var(--ease-out);
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
    margin: var(--space-4) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  /* Steps */
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

  .next {
    border-style: dashed;
    background: transparent;
  }

  .danger .panel-title {
    color: var(--color-danger);
  }

  @media (prefers-reduced-motion: reduce) {
    .btn {
      transition: none;
    }
  }
</style>
