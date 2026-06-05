<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    Button,
    EmptyState,
    Field,
    PageShell,
    Panel,
    StatusBanner,
    toast,
  } from "$lib/ui";
  import { onboardEmployee, type OnboardingResult, type OnboardingStep } from "$lib/hr/onboarding";
  import { listAgentEvents, type AgentEvent } from "$lib/agent/eventsService";
  import { approveSuggestion, rejectSuggestion } from "$lib/agent/approvalService";

  let locked = $state(false);
  let busy = $state(false);
  let error = $state("");

  let fullName = $state("");
  let email = $state("");
  let role = $state("");

  let steps = $state<OnboardingStep[]>([]);
  let result = $state<OnboardingResult | null>(null);
  let agentEvents = $state<AgentEvent[]>([]);
  let decidingId = $state<string | null>(null);
  /** Ficha aprovada via AGENT-007 — o wizard reutiliza-a em vez de criar outra. */
  let approvedRecordId = $state<string | null>(null);

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) void refreshEvents();
  });

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

  async function handleApproveSuggestion(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
    try {
      const result = await approveSuggestion(ev.id);
      await refreshEvents();
      toast.success("Sugestão aprovada.");
      if (result.action === "run_onboarding" && typeof ev.payload.record_id === "string") {
        approvedRecordId = ev.payload.record_id;
      }
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao aprovar";
    } finally {
      decidingId = null;
    }
  }

  async function handleRejectSuggestion(ev: AgentEvent) {
    decidingId = ev.id;
    error = "";
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

  async function run(e: SubmitEvent) {
    e.preventDefault();
    if (!fullName.trim() || !email.trim()) return;
    busy = true;
    error = "";
    result = null;
    steps = [];
    try {
      result = await onboardEmployee(
        {
          fullName: fullName.trim(),
          email: email.trim(),
          role: role.trim(),
          recordId: approvedRecordId ?? undefined,
        },
        (s) => (steps = s),
      );
      approvedRecordId = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha no onboarding";
    } finally {
      busy = false;
    }
  }

  function reset() {
    fullName = "";
    email = "";
    role = "";
    steps = [];
    result = null;
    error = "";
    approvedRecordId = null;
  }

  async function copyAlias() {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result.alias);
      toast.success("Alias copiado.");
    } catch {
      toast.error("Não foi possível copiar.");
    }
  }
</script>

<svelte:head>
  <title>Onboarding — AegisPass</title>
</svelte:head>

<PageShell
  title="Onboarding de Empregado"
  taskId="HR-007 · AGENT-007"
  description="Num clique: cria a ficha cifrada (HR-001), guarda campos iniciais e gera um alias de e-mail (MAIL-001). Cada passo fica no registo imutável (HR-002)."
 
>
  {#snippet actions()}
    <DocHelpLink slug="journey-hr-agent-onboarding" label="Como funciona o agente RH?" />
    <Button variant="ghost" size="sm" href="/hr">← Fichas</Button>
  {/snippet}

  {#if locked}
    <EmptyState title="Cofre bloqueado" description="Desbloqueia a Master Key para fazer onboarding (cifra os campos).">
      {#snippet action()}
        <Button href="/vault">Ir desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else}
    {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}

    <Panel title="Sugestões do orquestrador">
      <p class="muted sm">
        Quando crias uma ficha vazia em Fichas, o agente RH sugere completar o onboarding
        (AGENT-009 — aprova antes de cifrar dados).
      </p>
      {#if agentEvents.length === 0}
        <p class="muted">Sem eventos recentes.</p>
      {:else}
        <ul class="event-list">
          {#each agentEvents.slice(0, 6) as ev (ev.id)}
            <li
              class:suggested={isPendingSuggestion(ev)}
              class:approved={ev.approvalStatus === "approved"}
              class:rejected={ev.approvalStatus === "rejected"}
            >
              <div class="ev-body">
                <span class="ev-label">{ev.label}</span>
                {#if isPendingSuggestion(ev)}
                  <div class="ev-actions">
                    <Button
                      size="sm"
                      disabled={decidingId !== null || busy}
                      loading={decidingId === ev.id}
                      onclick={() => handleApproveSuggestion(ev)}
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
              <span class="ev-meta">{new Date(ev.createdAt).toLocaleString("pt-PT")}</span>
            </li>
          {/each}
        </ul>
      {/if}
    </Panel>

    {#if approvedRecordId}
      <p class="hint" role="status">
        Ficha <code class="mono">{approvedRecordId}</code> aprovada — preenche os dados abaixo.
      </p>
    {/if}

    {#if !result}
      <Panel title="Dados do empregado">
        <form onsubmit={run} class="form">
          <Field label="Nome completo" required>
            {#snippet control({ id, describedBy })}
              <input {id} aria-describedby={describedBy} bind:value={fullName} placeholder="Joana Silva" disabled={busy} required />
            {/snippet}
          </Field>
          <Field label="E-mail pessoal (destino do alias)" required>
            {#snippet control({ id, describedBy })}
              <input {id} aria-describedby={describedBy} type="email" bind:value={email} placeholder="joana@gmail.com" disabled={busy} required />
            {/snippet}
          </Field>
          <Field label="Função (opcional)">
            {#snippet control({ id, describedBy })}
              <input {id} aria-describedby={describedBy} bind:value={role} placeholder="Engenheira de software" disabled={busy} />
            {/snippet}
          </Field>
          <Button type="submit" disabled={busy || !fullName.trim() || !email.trim()} loading={busy}>
            Onboard em 1 clique
          </Button>
        </form>

        {#if steps.length > 0}
          <ul class="steps">
            {#each steps as s (s.label)}
              <li class:done={s.done}>
                <span class="mark">{s.done ? "✓" : "○"}</span>
                {s.label}
              </li>
            {/each}
          </ul>
        {/if}
      </Panel>
    {:else}
      <Panel title="Onboarding concluído">
        <ul class="steps">
          {#each result.steps as s (s.label)}
            <li class="done"><span class="mark">✓</span> {s.label}</li>
          {/each}
        </ul>
        <div class="alias-box">
          <span class="muted sm">Alias gerado</span>
          <div class="alias-row">
            <span class="mono">{result.alias}</span>
            <Button variant="secondary" size="sm" onclick={copyAlias}>Copiar</Button>
          </div>
        </div>
        <div class="actions">
          <Button href="/hr">Abrir fichas</Button>
          <Button variant="secondary" onclick={reset}>Onboard outro</Button>
        </div>
      </Panel>
    {/if}
  {/if}
</PageShell>

<style>
  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }
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
  input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  .steps {
    margin: var(--space-4) 0 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .steps li {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }
  .steps li.done {
    color: var(--color-text);
  }
  .mark {
    display: inline-block;
    width: 1.2rem;
    color: var(--color-success-fg);
    font-family: var(--font-mono);
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
  .alias-box {
    margin: var(--space-4) 0;
    padding: var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
  }
  .alias-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    margin-top: var(--space-1);
  }
  .actions {
    display: flex;
    gap: var(--space-2);
  }
  .event-list {
    margin: var(--space-3) 0 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .event-list li {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    font-size: var(--text-sm);
  }
  .event-list li.suggested {
    border-color: var(--color-accent);
    background: color-mix(in srgb, var(--color-accent-muted) 40%, var(--color-bg-inset));
  }
  .ev-body {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .ev-label {
    font-weight: 500;
  }
  .ev-meta {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    white-space: nowrap;
  }
  .ev-actions {
    display: flex;
    gap: var(--space-2);
  }
  .hint {
    margin: 0 0 var(--space-3);
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    font-size: var(--text-sm);
  }
</style>
