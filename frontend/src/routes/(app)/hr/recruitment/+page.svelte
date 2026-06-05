<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";
  import {
    markInboxProcessed,
    runRecruitment,
    seedRecruitmentEmail,
    type CandidateDraft,
  } from "$lib/hr/recruitmentService";

  let busy = $state(false);
  let error = $state("");
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  let drafts = $state<CandidateDraft[]>([]);

  let seedFrom = $state("candidato@example.com");
  let seedSubject = $state("Candidatura — Engenheira de software");
  let seedBody = $state("Género: Feminino\nEtnia: —\nExperiência: 5 anos Go e PostgreSQL.");

  onMount(() => void refreshEvents());

  async function refreshEvents() {
    try {
      agentEvents = await listAgentEvents();
    } catch {
      agentEvents = [];
    }
  }

  function isPendingSuggestion(ev: AgentEvent): boolean {
    return ev.type === "orchestrator.action.suggested" && (ev.approvalStatus ?? "pending") === "pending";
  }

  async function handleApprove(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
    try {
      const result = await approveSuggestion(ev.id);
      await refreshEvents();
      if (result.action === "screen_candidate") {
        await handleScreening();
      }
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao aprovar";
    } finally {
      decidingId = null;
    }
  }

  async function handleRejectSuggestion(ev: AgentEvent) {
    decidingId = ev.id;
    try {
      await rejectSuggestion(ev.id);
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao rejeitar";
    } finally {
      decidingId = null;
    }
  }

  async function handleScreening() {
    busy = true;
    error = "";
    try {
      drafts = await runRecruitment();
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha na triagem";
    } finally {
      busy = false;
    }
  }

  async function handleSeed(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    error = "";
    try {
      await seedRecruitmentEmail(seedFrom.trim(), seedSubject.trim(), seedBody.trim());
      await refreshEvents();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao simular e-mail";
    } finally {
      busy = false;
    }
  }

  async function dismissDraft(d: CandidateDraft) {
    busy = true;
    try {
      await markInboxProcessed(d.message_id);
      drafts = drafts.filter((x) => x.message_id !== d.message_id);
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Recrutamento — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">AGENT-007 · Triagem às cegas</p>
      <h1>Recrutamento</h1>
      <DocHelpLink slug="journey-hr-agent-recruitment" label="Como funciona a triagem às cegas?" />
    </div>
    <a class="back" href="/hr">← Fichas</a>
  </header>

  <p class="lead">
    Candidaturas por e-mail são analisadas com <strong>triagem às cegas</strong> — género, etnia e
    idade ficam ocultos antes de chegares ao resumo (CT-RGPD-04).
  </p>

  {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

  <section class="panel events">
    <h2>Sugestões do orquestrador</h2>
    {#if agentEvents.length === 0}
      <p class="muted">Sem eventos recentes.</p>
    {:else}
      <ul class="event-list">
        {#each agentEvents.slice(0, 8) as ev (ev.id)}
          <li class:suggested={isPendingSuggestion(ev)}>
            <div class="ev-body">
              <span class="ev-label">{ev.label}</span>
              {#if isPendingSuggestion(ev) && ev.payload.action === "screen_candidate"}
                <div class="ev-actions">
                  <button type="button" class="btn approve" disabled={decidingId !== null || busy} onclick={() => handleApprove(ev)}>
                    {decidingId === ev.id ? "…" : "Aprovar"}
                  </button>
                  <button type="button" class="btn reject" disabled={decidingId !== null} onclick={() => handleRejectSuggestion(ev)}>
                    Rejeitar
                  </button>
                </div>
              {/if}
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>

  <section class="panel">
    <h2>Simular candidatura</h2>
    <form class="seed" onsubmit={handleSeed}>
      <label>
        De
        <input type="email" bind:value={seedFrom} disabled={busy} />
      </label>
      <label>
        Assunto
        <input bind:value={seedSubject} disabled={busy} />
      </label>
      <label>
        Corpo (com campos a ocultar)
        <textarea bind:value={seedBody} rows="3" disabled={busy}></textarea>
      </label>
      <button type="submit" class="btn secondary" disabled={busy}>Simular</button>
    </form>
    <button type="button" class="btn primary" disabled={busy} onclick={handleScreening}>
      {busy ? "A processar…" : "Correr triagem"}
    </button>
  </section>

  {#if drafts.length > 0}
    <section class="panel">
      <h2>Candidatos (triagem às cegas)</h2>
      <ul class="drafts">
        {#each drafts as d (d.message_id)}
          <li>
            <p class="draft-head"><strong>{d.email}</strong> — {d.subject}</p>
            <pre class="summary">{d.summary}</pre>
            {#if d.blind}<span class="badge">Triagem às cegas</span>{/if}
            <button type="button" class="btn sm" disabled={busy} onclick={() => dismissDraft(d)}>Marcar processada</button>
          </li>
        {/each}
      </ul>
    </section>
  {/if}
</section>

<style>
  .page { max-width: 42rem; }
  .eyebrow { font-size: var(--text-xs); text-transform: uppercase; color: var(--color-text-muted); }
  h1 { margin: 0; font-family: var(--font-display); }
  h2 { margin: 0 0 var(--space-3); font-size: var(--text-base); }
  .lead { color: var(--color-text-muted); font-size: var(--text-sm); }
  .panel { border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: var(--space-5); margin-bottom: var(--space-4); background: var(--color-bg-surface); }
  .muted { color: var(--color-text-muted); font-size: var(--text-sm); }
  .inline-error { color: var(--color-danger); font-size: var(--text-sm); }
  .back { font-size: var(--text-xs); color: var(--color-text-muted); }
  label { display: block; margin-bottom: var(--space-2); font-size: var(--text-sm); }
  input, textarea { width: 100%; padding: var(--space-2); border: 1px solid var(--color-border); border-radius: var(--radius-sm); box-sizing: border-box; }
  .btn { padding: var(--space-2) var(--space-4); border-radius: var(--radius-sm); border: 1px solid var(--color-border); cursor: pointer; margin-top: var(--space-2); }
  .btn.primary { background: var(--color-accent); color: var(--color-accent-fg); border: none; }
  .btn.approve { background: var(--color-success-bg); color: var(--color-success-fg); border: none; }
  .btn.sm { font-size: var(--text-xs); padding: var(--space-1) var(--space-2); }
  .event-list, .drafts { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: var(--space-2); }
  .event-list li, .drafts li { padding: var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-sm); background: var(--color-bg-inset); }
  .event-list li.suggested { border-color: var(--color-accent); }
  .ev-actions { display: flex; gap: var(--space-2); margin-top: var(--space-2); }
  .summary { font-size: var(--text-xs); white-space: pre-wrap; overflow-x: auto; max-height: 12rem; }
  .badge { font-size: var(--text-xs); color: var(--color-accent); }
  .seed { margin-bottom: var(--space-3); }
</style>
