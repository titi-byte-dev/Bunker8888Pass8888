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
  import {
    Button,
    confirmDialog,
    EmptyState,
    Field,
    PageShell,
    Panel,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";

  let loading = $state(true);
  let inboxLoading = $state(true);
  let busy = $state(false);
  let error = $state("");
  let aliases = $state<EmailAlias[]>([]);
  let inbox = $state<InboxMessage[]>([]);
  let destination = $state("");
  let label = $state("");
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
      toast.success("Alias criado.");
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
      toast.info(a.active ? "Alias desligado." : "Alias activado.");
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao alterar alias";
    } finally {
      busy = false;
    }
  }

  async function onDelete(a: EmailAlias) {
    const ok = await confirmDialog({
      title: "Apagar alias?",
      message: `Remove «${a.aliasAddress}». O reencaminhamento deixa de funcionar de imediato.`,
      confirmLabel: "Apagar",
      variant: "danger",
    });
    if (!ok) return;
    busy = true;
    try {
      await deleteAlias(a.id);
      await refresh();
      toast.success("Alias apagado.");
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
      toast.success("E-mail enviado (Mailpit em dev).");
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao enviar";
    } finally {
      busy = false;
    }
  }

  async function copy(addr: string) {
    try {
      await navigator.clipboard.writeText(addr);
      toast.success("Endereço copiado.");
    } catch {
      toast.error("Não foi possível copiar.");
    }
  }
</script>

<svelte:head>
  <title>Aliases de E-mail — AegisPass</title>
</svelte:head>

<PageShell
  title="Aliases de E-mail"
  taskId="MAIL-001/002"
  description="Endereços descartáveis que reencaminham para o teu e-mail real. Envio em dev via Mailpit (SMTP localhost:1025). Envia para o teu alias @aegis.email e corre a prospeção no CRM."
  width="narrow"
>
  {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

  <Panel title="Novo alias">
    <form onsubmit={onCreate} class="row-form">
      <Field label="Reencaminhar para (e-mail real)" required>
        {#snippet control({ id, describedBy })}
          <input
            {id}
            aria-describedby={describedBy}
            type="email"
            bind:value={destination}
            placeholder="tu@empresa.pt"
            disabled={busy}
            required
          />
        {/snippet}
      </Field>
      <Field label="Rótulo (opcional)">
        {#snippet control({ id, describedBy })}
          <input
            {id}
            aria-describedby={describedBy}
            type="text"
            bind:value={label}
            placeholder="ex.: Netflix"
            disabled={busy}
          />
        {/snippet}
      </Field>
      <Button type="submit" disabled={busy || !destination.trim()} loading={busy}>
        Gerar alias
      </Button>
    </form>
  </Panel>

  <Panel title="Os meus aliases">
    {#if loading}
      <Skeleton variant="row" />
      <Skeleton variant="row" />
    {:else if aliases.length === 0}
      <EmptyState title="Sem aliases" description="Gera o primeiro alias acima." />
    {:else}
      <ul class="list">
        {#each aliases as a (a.id)}
          <li class="alias" class:off={!a.active}>
            <div class="a-main">
              <button type="button" class="addr mono" onclick={() => copy(a.aliasAddress)}>
                {a.aliasAddress}
              </button>
              <span class="muted sm">→ {a.destination}</span>
              {#if a.label}<span class="tag">{a.label}</span>{/if}
            </div>
            <div class="a-actions">
              <span class="state" class:on={a.active}>{a.active ? "activo" : "desligado"}</span>
              <Button variant="ghost" size="sm" onclick={() => onToggle(a)} disabled={busy}>
                {a.active ? "desligar" : "ligar"}
              </Button>
              <Button variant="ghost" size="sm" onclick={() => onDelete(a)} disabled={busy}>
                apagar
              </Button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </Panel>

  <Panel title="Compor e-mail (MAIL-004)">
    <form class="compose-form" onsubmit={onCompose}>
      <Field label="Alias (remetente)" required>
        {#snippet control({ id, describedBy })}
          <select {id} aria-describedby={describedBy} bind:value={composeAliasId} disabled={busy || aliases.length === 0}>
            <option value="">— escolhe —</option>
            {#each aliases.filter((a) => a.active) as a (a.id)}
              <option value={a.id}>{a.aliasAddress}</option>
            {/each}
          </select>
        {/snippet}
      </Field>
      <Field label="Para" required>
        {#snippet control({ id, describedBy })}
          <input {id} aria-describedby={describedBy} type="email" bind:value={composeTo} required disabled={busy} />
        {/snippet}
      </Field>
      <Field label="Assunto" required>
        {#snippet control({ id, describedBy })}
          <input {id} aria-describedby={describedBy} bind:value={composeSubject} required disabled={busy} />
        {/snippet}
      </Field>
      <Field label="Mensagem" required>
        {#snippet control({ id, describedBy })}
          <textarea {id} aria-describedby={describedBy} bind:value={composeBody} rows="3" required disabled={busy}></textarea>
        {/snippet}
      </Field>
      <Button type="submit" disabled={busy || !composeAliasId} loading={busy}>Enviar</Button>
    </form>
  </Panel>

  <Panel title="Caixa de entrada (MAIL-002)">
    {#snippet actions()}
      <Button variant="ghost" size="sm" href="/crm">Prospeção no CRM →</Button>
    {/snippet}
    {#if inboxLoading}
      <Skeleton variant="row" />
      <Skeleton variant="row" />
    {:else if inbox.length === 0}
      <EmptyState
        title="Inbox vazia"
        description="Envia SMTP para um alias activo (Mailpit UI ou swaks --to TEU_ALIAS@aegis.email --server localhost:1025)."
      />
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
  </Panel>
</PageShell>

<style>
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
    gap: var(--space-3);
  }
  .row-form {
    flex-direction: row;
    flex-wrap: wrap;
    align-items: flex-end;
  }
  .row-form :global(.field) {
    flex: 1;
    min-width: 10rem;
  }
  .compose-form select,
  .compose-form textarea,
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
  input:focus-visible,
  select:focus-visible,
  textarea:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  @media (max-width: 560px) {
    .row-form {
      flex-direction: column;
      align-items: stretch;
    }
  }
  .list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
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
    gap: var(--space-2);
    flex-shrink: 0;
  }
  .state {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
  .state.on {
    color: var(--color-success-fg);
  }
  .inbox-list .inbox-item {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: var(--space-2);
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
  }
  .inbox-item.done {
    opacity: 0.65;
  }
  .snippet {
    margin: var(--space-1) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
