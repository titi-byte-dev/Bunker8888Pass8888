<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { createLead, deleteLead, listLeads, updateLead } from "$lib/crm/leadsService";
  import {
    importDraft,
    runProspection,
    seedInboxMessage,
    type ProspectionDraft,
  } from "$lib/crm/prospectionService";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { LEAD_STAGES, type Lead, type LeadStage } from "$lib/crm/leads";

  let locked = $state(false);
  let loading = $state(true);
  let busy = $state(false);
  let error = $state("");
  let leads = $state<Lead[]>([]);

  let fName = $state("");
  let fEmail = $state("");
  let fCompany = $state("");
  let fStage = $state<LeadStage>("new");
  let fNotes = $state("");

  let drafts = $state<ProspectionDraft[]>([]);
  let prospectionBusy = $state(false);
  let seedFrom = $state("lead@empresa.pt");
  let seedSubject = $state("Pedido de demonstração");
  let seedBody = $state("Olá, gostávamos de agendar uma demo do produto.");
  let agentEvents = $state<AgentEvent[]>([]);

  const byStage = $derived(
    LEAD_STAGES.map((s) => ({
      ...s,
      items: leads.filter((l) => l.stage === s.id),
    })),
  );

  async function refresh() {
    loading = true;
    error = "";
    try {
      leads = await listLeads();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar leads";
    } finally {
      loading = false;
    }
  }

  async function refreshEvents() {
    try {
      agentEvents = await listAgentEvents();
    } catch {
      agentEvents = [];
    }
  }

  onMount(() => {
    locked = !getMasterKey();
    refreshEvents();
    if (!locked) refresh();
    else loading = false;
  });

  async function handleCreate(e: SubmitEvent) {
    e.preventDefault();
    if (!fName.trim() || !fEmail.trim()) return;
    busy = true;
    error = "";
    try {
      await createLead({
        name: fName.trim(),
        email: fEmail.trim(),
        company: fCompany.trim() || undefined,
        stage: fStage,
        notes: fNotes.trim() || undefined,
        source: "manual",
      });
      fName = "";
      fEmail = "";
      fCompany = "";
      fNotes = "";
      fStage = "new";
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao criar lead";
    } finally {
      busy = false;
    }
  }

  async function moveStage(lead: Lead, stage: LeadStage) {
    busy = true;
    error = "";
    try {
      await updateLead(lead.id, { ...lead, stage });
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao actualizar";
    } finally {
      busy = false;
    }
  }

  async function handleProspection() {
    prospectionBusy = true;
    error = "";
    try {
      drafts = await runProspection();
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha na prospeção";
    } finally {
      prospectionBusy = false;
    }
  }

  async function handleImportDraft(draft: ProspectionDraft) {
    busy = true;
    error = "";
    try {
      await importDraft(draft);
      drafts = drafts.filter((d) => d.message_id !== draft.message_id);
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao importar";
    } finally {
      busy = false;
    }
  }

  async function handleSeedInbox(e: SubmitEvent) {
    e.preventDefault();
    prospectionBusy = true;
    error = "";
    try {
      await seedInboxMessage(seedFrom.trim(), seedSubject.trim(), seedBody.trim());
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao simular e-mail";
    } finally {
      prospectionBusy = false;
    }
  }

  async function remove(lead: Lead) {
    if (!confirm(`Apagar lead «${lead.name}»?`)) return;
    busy = true;
    try {
      await deleteLead(lead.id);
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao apagar";
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>CRM — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">CRM-001/002 · AGENT-003 · Funil de vendas</p>
      <h1>Leads</h1>
      <DocHelpLink slug="journey-crm-prospection" label="Como funciona a prospeção?" />
    </div>
    <p class="lead">
      Contactos cifrados com a Master Key — o servidor só vê blobs opacos. Arrasta
      mentalmente pelas colunas do funil (estágio no payload decifrado no cliente).
    </p>
  </header>

  {#if locked}
    <section class="panel">
      <p class="muted">🔒 Desbloqueia o cofre para gerir leads.</p>
      <a class="btn primary" href="/auth/unlock">Desbloquear</a>
    </section>
  {:else}
    {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

    <section class="panel events">
      <h2>Actividade dos agentes</h2>
      <p class="muted small">Event Bus AGENT-004 — feed auditável de mail e prospeção.</p>
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos recentes.</p>
      {:else}
        <ul class="event-list">
          {#each agentEvents.slice(0, 8) as ev (ev.id)}
            <li class:suggested={ev.type === "orchestrator.action.suggested"}>
              <span class="ev-label">{ev.label}</span>
              <span class="ev-meta">{new Date(ev.createdAt).toLocaleString("pt-PT")}</span>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section class="panel prospection">
      <h2>Prospeção automática</h2>
      <p class="muted small">
        O agente lê e-mails pendentes (stub até MAIL-002), gera rascunhos e tu
        importas com cifragem local — Zero-Knowledge mantido.
      </p>
      <form class="seed" onsubmit={handleSeedInbox}>
        <p class="seed-label">Simular e-mail recebido</p>
        <div class="grid">
          <label>
            De
            <input type="email" bind:value={seedFrom} disabled={prospectionBusy} />
          </label>
          <label>
            Assunto
            <input bind:value={seedSubject} disabled={prospectionBusy} />
          </label>
        </div>
        <label>
          Corpo
          <textarea bind:value={seedBody} rows="2" disabled={prospectionBusy}></textarea>
        </label>
        <button type="submit" class="btn secondary" disabled={prospectionBusy}>Simular</button>
      </form>
      <button
        type="button"
        class="btn primary"
        disabled={prospectionBusy || busy}
        onclick={handleProspection}
      >
        {prospectionBusy ? "A processar…" : "Correr prospeção"}
      </button>
      {#if drafts.length > 0}
        <ul class="draft-list">
          {#each drafts as draft (draft.message_id)}
            <li class="draft-card">
              <strong>{draft.email}</strong>
              <span class="email">{draft.subject}</span>
              <button
                type="button"
                class="btn secondary"
                disabled={busy}
                onclick={() => handleImportDraft(draft)}
              >
                Importar para o funil
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <form class="panel form" onsubmit={handleCreate}>
      <h2>Novo lead</h2>
      <div class="grid">
        <label>
          Nome
          <input bind:value={fName} required disabled={busy} />
        </label>
        <label>
          E-mail
          <input type="email" bind:value={fEmail} required disabled={busy} />
        </label>
        <label>
          Empresa
          <input bind:value={fCompany} disabled={busy} />
        </label>
        <label>
          Estágio
          <select bind:value={fStage} disabled={busy}>
            {#each LEAD_STAGES as s (s.id)}
              <option value={s.id}>{s.label}</option>
            {/each}
          </select>
        </label>
      </div>
      <label>
        Notas
        <textarea bind:value={fNotes} rows="2" disabled={busy}></textarea>
      </label>
      <button type="submit" class="btn primary" disabled={busy}>Adicionar</button>
    </form>

    {#if loading}
      <p class="muted">A carregar funil…</p>
    {:else}
      <div class="board">
        {#each byStage as col (col.id)}
          <section class="column">
            <h3>{col.label} <span class="count">{col.items.length}</span></h3>
            <ul>
              {#each col.items as lead (lead.id)}
                <li class="card">
                  <strong>{lead.name}</strong>
                  <span class="email">{lead.email}</span>
                  {#if lead.company}<span class="co">{lead.company}</span>{/if}
                  <div class="actions">
                    <select
                      value={lead.stage}
                      onchange={(e) => moveStage(lead, (e.currentTarget as HTMLSelectElement).value as LeadStage)}
                      disabled={busy}
                    >
                      {#each LEAD_STAGES as s (s.id)}
                        <option value={s.id}>{s.label}</option>
                      {/each}
                    </select>
                    <button type="button" class="linkish" disabled={busy} onclick={() => remove(lead)}>
                      Apagar
                    </button>
                  </div>
                </li>
              {/each}
            </ul>
          </section>
        {/each}
      </div>
    {/if}
  {/if}
</section>

<style>
  .page {
    max-width: 72rem;
  }
  .page-head {
    margin-bottom: var(--space-5);
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
    max-width: 42rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .panel {
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    margin-bottom: var(--space-4);
  }
  .form h2 {
    margin: 0 0 var(--space-3);
    font-size: var(--text-base);
  }
  .grid {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
    margin-bottom: var(--space-3);
  }
  label {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-sm);
  }
  input,
  select,
  textarea {
    padding: var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    color: inherit;
    font-family: var(--font-ui);
  }
  .board {
    display: grid;
    gap: var(--space-3);
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
    align-items: start;
  }
  .column {
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    min-height: 8rem;
  }
  .column h3 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-sm);
    display: flex;
    justify-content: space-between;
  }
  .count {
    color: var(--color-text-muted);
    font-weight: 400;
  }
  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .card {
    padding: var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    font-size: var(--text-xs);
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .email,
  .co {
    color: var(--color-text-muted);
  }
  .actions {
    display: flex;
    gap: var(--space-2);
    margin-top: var(--space-1);
    align-items: center;
  }
  .actions select {
    flex: 1;
    font-size: var(--text-xs);
  }
  .linkish {
    border: none;
    background: none;
    color: var(--color-link);
    font-size: var(--text-xs);
    cursor: pointer;
    padding: 0;
  }
  .inline-error {
    color: var(--color-danger, #c44);
    font-size: var(--text-sm);
  }
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .btn.primary {
    display: inline-block;
    margin-top: var(--space-2);
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-md);
    background: var(--color-accent);
    color: var(--color-bg-base);
    text-decoration: none;
    border: none;
    cursor: pointer;
    font-size: var(--text-sm);
  }
  .btn.primary:disabled,
  .btn.secondary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  .btn.secondary {
    display: inline-block;
    margin-top: var(--space-2);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-md);
    background: transparent;
    color: var(--color-accent);
    border: 1px solid var(--color-accent);
    cursor: pointer;
    font-size: var(--text-sm);
  }
  .events h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-base);
  }
  .event-list {
    list-style: none;
    margin: var(--space-2) 0 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .event-list li {
    display: flex;
    justify-content: space-between;
    gap: var(--space-2);
    font-size: var(--text-xs);
    padding: var(--space-1) 0;
    border-bottom: 1px solid var(--color-border);
  }
  .ev-label {
    font-weight: 500;
  }
  .ev-meta {
    color: var(--color-text-muted);
    flex-shrink: 0;
  }
  .event-list li.suggested {
    background: var(--color-accent-muted, rgba(99, 102, 241, 0.08));
    border-radius: var(--radius-sm);
    padding: var(--space-1) var(--space-2);
  }
  .prospection h2 {
    margin: 0 0 var(--space-2);
    font-size: var(--text-base);
  }
  .small {
    margin: 0 0 var(--space-3);
  }
  .seed {
    margin-bottom: var(--space-3);
    padding-bottom: var(--space-3);
    border-bottom: 1px solid var(--color-border);
  }
  .seed-label {
    margin: 0 0 var(--space-2);
    font-size: var(--text-sm);
    font-weight: 600;
  }
  .draft-list {
    list-style: none;
    margin: var(--space-3) 0 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .draft-card {
    padding: var(--space-2);
    border: 1px dashed var(--color-border);
    border-radius: var(--radius-sm);
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    font-size: var(--text-sm);
  }
</style>
