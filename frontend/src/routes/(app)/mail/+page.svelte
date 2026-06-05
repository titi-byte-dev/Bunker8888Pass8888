<script lang="ts">
  import { onMount } from "svelte";
  import {
    composeFromAlias,
    createAlias,
    deleteAlias,
    listAliases,
    setAliasActive,
    type EmailAlias,
  } from "$lib/mail/aliases";
  import { listInbox, type InboxMessage } from "$lib/mail/inbox";

  let loading = $state(true);
  let inboxLoading = $state(true);
  let busy = $state(false);
  let error = $state("");
  let aliases = $state<EmailAlias[]>([]);
  let inbox = $state<InboxMessage[]>([]);
  let destination = $state("");
  let label = $state("");
  let copied = $state("");
  let composeAliasId = $state("");
  let composeTo = $state("");
  let composeSubject = $state("");
  let composeBody = $state("");

  async function refresh() {
    loading = true;
    error = "";
    try {
      aliases = await listAliases();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar aliases";
    } finally {
      loading = false;
    }
  }

  async function refreshInbox() {
    inboxLoading = true;
    try {
      inbox = await listInbox(false);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar inbox";
    } finally {
      inboxLoading = false;
    }
  }

  onMount(() => {
    refresh();
    refreshInbox();
  });

  async function onCreate(e: SubmitEvent) {
    e.preventDefault();
    if (!destination.trim()) return;
    busy = true;
    error = "";
    try {
      await createAlias(destination.trim(), label.trim());
      destination = "";
      label = "";
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao criar alias";
    } finally {
      busy = false;
    }
  }

  async function onToggle(a: EmailAlias) {
    busy = true;
    try {
      await setAliasActive(a.id, !a.active);
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao alterar alias";
    } finally {
      busy = false;
    }
  }

  async function onDelete(id: string) {
    busy = true;
    try {
      await deleteAlias(id);
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao apagar alias";
    } finally {
      busy = false;
    }
  }

  async function onCompose(e: SubmitEvent) {
    e.preventDefault();
    if (!composeAliasId || !composeTo.trim()) return;
    busy = true;
    error = "";
    try {
      await composeFromAlias(
        composeAliasId,
        composeTo.trim(),
        composeSubject.trim(),
        composeBody.trim(),
      );
      composeTo = "";
      composeSubject = "";
      composeBody = "";
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao enviar";
    } finally {
      busy = false;
    }
  }

  async function copy(addr: string) {
    try {
      await navigator.clipboard.writeText(addr);
      copied = addr;
      setTimeout(() => (copied = ""), 1500);
    } catch {
      copied = "";
    }
  }
</script>

<svelte:head>
  <title>Aliases de E-mail — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">MAIL-001/002 · Aliases + SMTP dev (Mailpit)</p>
      <h1>Aliases de E-mail</h1>
    </div>
    <p class="lead">
      Endereços descartáveis que reencaminham para o teu e-mail real. Dá um alias
      a cada serviço; se vazar ou fizer spam, desligas só esse alias. O envio
      efectivo em dev usa Mailpit (SMTP <code>localhost:1025</code>). Envia para o teu
      alias <code>@aegis.email</code> e a mensagem aparece na inbox abaixo — depois
      corre a prospeção em <a href="/crm">CRM</a>.
    </p>
  </header>

  {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

  <section class="panel">
    <div class="panel-head"><p class="eyebrow">Novo alias</p></div>
    <form onsubmit={onCreate} class="row-form">
      <input type="email" bind:value={destination} placeholder="Reencaminhar para (e-mail real)" disabled={busy} />
      <input type="text" bind:value={label} placeholder="Rótulo (ex.: Netflix)" disabled={busy} />
      <button type="submit" class="btn primary" disabled={busy || !destination.trim()}>
        Gerar alias
      </button>
    </form>
  </section>

  <section class="panel">
    <div class="panel-head"><p class="eyebrow">Os meus aliases</p></div>
    {#if loading}
      <p class="muted">A carregar…</p>
    {:else if aliases.length === 0}
      <p class="muted">Sem aliases. Gera o primeiro acima.</p>
    {:else}
      <ul class="list">
        {#each aliases as a (a.id)}
          <li class="alias" class:off={!a.active}>
            <div class="a-main">
              <button type="button" class="addr mono" onclick={() => copy(a.aliasAddress)}>
                {a.aliasAddress}
              </button>
              {#if copied === a.aliasAddress}<span class="copied">copiado!</span>{/if}
              <span class="muted sm">→ {a.destination}</span>
              {#if a.label}<span class="tag">{a.label}</span>{/if}
            </div>
            <div class="a-actions">
              <span class="state" class:on={a.active}>{a.active ? "activo" : "desligado"}</span>
              <button type="button" class="link-btn" onclick={() => onToggle(a)} disabled={busy}>
                {a.active ? "desligar" : "ligar"}
              </button>
              <button type="button" class="link-btn" onclick={() => onDelete(a.id)} disabled={busy}>
                apagar
              </button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>

  <section class="panel">
    <div class="panel-head"><p class="eyebrow">Compor e-mail (MAIL-004)</p></div>
    <form class="compose-form" onsubmit={onCompose}>
      <label>
        Alias (remetente)
        <select bind:value={composeAliasId} disabled={busy || aliases.length === 0}>
          <option value="">— escolhe —</option>
          {#each aliases.filter((a) => a.active) as a (a.id)}
            <option value={a.id}>{a.aliasAddress}</option>
          {/each}
        </select>
      </label>
      <label>
        Para
        <input type="email" bind:value={composeTo} required disabled={busy} />
      </label>
      <label>
        Assunto
        <input bind:value={composeSubject} required disabled={busy} />
      </label>
      <label>
        Mensagem
        <textarea bind:value={composeBody} rows="3" required disabled={busy}></textarea>
      </label>
      <button type="submit" class="btn primary" disabled={busy || !composeAliasId}>Enviar</button>
    </form>
  </section>

  <section class="panel">
    <div class="panel-head">
      <p class="eyebrow">Caixa de entrada (MAIL-002)</p>
      <a class="crm-link" href="/crm">Prospeção no CRM →</a>
    </div>
    {#if inboxLoading}
      <p class="muted">A carregar inbox…</p>
    {:else if inbox.length === 0}
      <p class="muted">
        Sem mensagens. Envia SMTP para um alias activo (ex.: via Mailpit UI ou
        <code>swaks --to TEU_ALIAS@aegis.email --server localhost:1025</code>).
      </p>
    {:else}
      <ul class="list inbox-list">
        {#each inbox as m (m.id)}
          <li class="inbox-item" class:done={!!m.processedAt}>
            <div>
              <strong>{m.fromEmail}</strong>
              <span class="muted sm">{m.subject}</span>
              <p class="snippet">{m.body.slice(0, 120)}{m.body.length > 120 ? "…" : ""}</p>
            </div>
            <span class="state" class:on={!m.processedAt}>
              {m.processedAt ? "processada" : "pendente"}
            </span>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</section>

<style>
  .page {
    max-width: 48rem;
  }
  .page-head {
    margin-bottom: var(--space-6);
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
  }
  .lead {
    margin: var(--space-3) 0 0;
    max-width: 40rem;
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
    margin-bottom: var(--space-3);
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: var(--space-2);
  }
  .crm-link {
    font-size: var(--text-xs);
    color: var(--color-accent);
    text-decoration: none;
  }
  .inbox-list .inbox-item {
    flex-direction: column;
    align-items: stretch;
    gap: var(--space-2);
  }
  .inbox-item.done {
    opacity: 0.65;
  }
  .snippet {
    margin: var(--space-1) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
  .panel-head .eyebrow {
    margin: 0;
  }
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .sm {
    font-size: var(--text-xs);
  }
  .mono {
    font-family: var(--font-mono);
  }
  .row-form,
  .compose-form {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .row-form {
    flex-direction: row;
  }
  .compose-form label {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-sm);
  }
  .compose-form select,
  .compose-form textarea {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
  }
  @media (max-width: 560px) {
    .row-form {
      flex-direction: column;
    }
  }
  input {
    flex: 1;
    min-width: 0;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }
  input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  .list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .alias {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
  }
  .alias.off {
    opacity: 0.6;
  }
  .a-main {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    flex-wrap: wrap;
  }
  .addr {
    background: none;
    border: none;
    color: var(--color-accent);
    cursor: pointer;
    padding: 0;
    font-size: var(--text-sm);
  }
  .copied {
    font-size: var(--text-xs);
    color: var(--color-success-fg);
  }
  .tag {
    font-size: var(--text-xs);
    padding: 1px var(--space-2);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    color: var(--color-accent);
  }
  .a-actions {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    flex-shrink: 0;
  }
  .state {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
  .state.on {
    color: var(--color-success-fg);
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
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .link-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    padding: 0;
  }
  .link-btn:hover {
    color: var(--color-danger);
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
</style>
