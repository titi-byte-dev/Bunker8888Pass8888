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
  import { Button, PageShell, Panel, StatusBanner, toast } from "$lib/ui";

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
      toast.info("Sugestão rejeitada.");
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
      toast.success(`${drafts.length} candidato(s) após triagem às cegas.`);
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
      toast.success("Candidatura simulada na inbox.");
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
      toast.info("Candidatura marcada como processada.");
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Recrutamento — AegisPass</title>
</svelte:head>

<PageShell
  title="Recrutamento"
  taskId="AGENT-007"
  description="Candidaturas por e-mail com triagem às cegas — género, etnia e idade ficam ocultos antes do resumo (CT-RGPD-04)."
 
>
  {#snippet actions()}
    <DocHelpLink slug="journey-hr-agent-recruitment" label="Como funciona a triagem às cegas?" />
    <Button variant="ghost" size="sm" href="/hr">← Fichas</Button>
  {/snippet}

  {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

  <Panel title="Sugestões do orquestrador">
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
                  <Button
                    variant="primary"
                    size="sm"
                    loading={decidingId === ev.id}
                    disabled={decidingId !== null || busy}
                    onclick={() => handleApprove(ev)}
                  >
                    Aprovar
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    disabled={decidingId !== null}
                    onclick={() => handleRejectSuggestion(ev)}
                  >
                    Rejeitar
                  </Button>
                </div>
              {/if}
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </Panel>

  <Panel title="Simular candidatura">
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
      <Button type="submit" variant="secondary" disabled={busy}>Simular</Button>
    </form>
    <Button disabled={busy} loading={busy} onclick={handleScreening}>Correr triagem</Button>
  </Panel>

  {#if drafts.length > 0}
    <Panel title="Candidatos (triagem às cegas)">
      <ul class="drafts">
        {#each drafts as d (d.message_id)}
          <li>
            <p class="draft-head"><strong>{d.email}</strong> — {d.subject}</p>
            <pre class="summary">{d.summary}</pre>
            {#if d.blind}<span class="badge">Triagem às cegas</span>{/if}
            <Button variant="ghost" size="sm" disabled={busy} onclick={() => dismissDraft(d)}>
              Marcar processada
            </Button>
          </li>
        {/each}
      </ul>
    </Panel>
  {/if}
</PageShell>

<style>
  .muted { margin: 0; color: var(--color-text-muted); font-size: var(--text-sm); }

  label {
    display: block;
    margin-bottom: var(--space-2);
    font-size: var(--text-sm);
    color: var(--color-text-label);
  }

  input, textarea {
    width: 100%;
    padding: var(--space-2);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    box-sizing: border-box;
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-size: var(--text-sm);
  }

  .event-list, .drafts {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .event-list li, .drafts li {
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
  }

  .event-list li.suggested { border-color: var(--color-accent); }

  .ev-body { display: flex; flex-direction: column; gap: var(--space-2); }
  .ev-label { font-size: var(--text-sm); }
  .ev-actions { display: flex; gap: var(--space-2); }

  .summary {
    font-size: var(--text-xs);
    white-space: pre-wrap;
    overflow-x: auto;
    max-height: 12rem;
    margin: var(--space-2) 0;
  }

  .badge { font-size: var(--text-xs); color: var(--color-accent); }
  .seed { margin-bottom: var(--space-3); }
  .draft-head { margin: 0 0 var(--space-1); font-size: var(--text-sm); }
</style>
