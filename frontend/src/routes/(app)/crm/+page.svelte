<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { createLead, deleteLead, listLeads, updateLead } from "$lib/crm/leadsService";
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

  onMount(() => {
    locked = !getMasterKey();
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
      <p class="eyebrow">CRM-001/002 · Funil de vendas</p>
      <h1>Leads</h1>
      <DocHelpLink slug="journey-admin-onboarding" label="Como funciona o onboarding?" />
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
  .btn.primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
