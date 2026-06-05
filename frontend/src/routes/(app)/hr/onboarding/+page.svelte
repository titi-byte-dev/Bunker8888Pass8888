<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { onboardEmployee, type OnboardingResult, type OnboardingStep } from "$lib/hr/onboarding";

  let locked = $state(false);
  let busy = $state(false);
  let error = $state("");

  let fullName = $state("");
  let email = $state("");
  let role = $state("");

  let steps = $state<OnboardingStep[]>([]);
  let result = $state<OnboardingResult | null>(null);
  let copied = $state(false);

  onMount(() => {
    locked = !getMasterKey();
  });

  async function run(e: SubmitEvent) {
    e.preventDefault();
    if (!fullName.trim() || !email.trim()) return;
    busy = true;
    error = "";
    result = null;
    steps = [];
    try {
      result = await onboardEmployee(
        { fullName: fullName.trim(), email: email.trim(), role: role.trim() },
        (s) => (steps = s),
      );
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
  }

  async function copyAlias() {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result.alias);
      copied = true;
      setTimeout(() => (copied = false), 1500);
    } catch {
      copied = false;
    }
  }
</script>

<svelte:head>
  <title>Onboarding — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">HR-007 · Onboarding em 1 clique</p>
      <h1>Onboarding de Empregado</h1>
    </div>
    <a class="back" href="/hr">← Fichas</a>
  </header>

  {#if locked}
    <section class="panel">
      <p class="muted">🔒 Desbloqueia a Master Key para fazer onboarding (cifra os campos).</p>
      <a class="btn primary" href="/vault">Ir desbloquear</a>
    </section>
  {:else}
    <p class="lead">
      Num clique: cria a ficha cifrada (HR-001), guarda os campos iniciais
      campo-a-campo e gera um <strong>alias de e-mail</strong> (MAIL-001) que
      reencaminha para o e-mail pessoal. Cada passo fica no registo imutável (HR-002).
    </p>

    {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

    {#if !result}
      <section class="panel">
        <form onsubmit={run}>
          <label class="field">
            <span>Nome completo</span>
            <input bind:value={fullName} placeholder="Joana Silva" disabled={busy} />
          </label>
          <label class="field">
            <span>E-mail pessoal (destino do alias)</span>
            <input type="email" bind:value={email} placeholder="joana@gmail.com" disabled={busy} />
          </label>
          <label class="field">
            <span>Função (opcional)</span>
            <input bind:value={role} placeholder="Engenheira de software" disabled={busy} />
          </label>
          <button type="submit" class="btn primary" disabled={busy || !fullName.trim() || !email.trim()}>
            {busy ? "A fazer onboarding…" : "Onboard em 1 clique"}
          </button>
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
      </section>
    {:else}
      <section class="panel done-panel">
        <p class="eyebrow">Onboarding concluído</p>
        <ul class="steps">
          {#each result.steps as s (s.label)}
            <li class="done"><span class="mark">✓</span> {s.label}</li>
          {/each}
        </ul>
        <div class="alias-box">
          <span class="muted sm">Alias gerado</span>
          <div class="alias-row">
            <span class="mono">{result.alias}</span>
            <button type="button" class="btn sm" onclick={copyAlias}>{copied ? "Copiado!" : "Copiar"}</button>
          </div>
        </div>
        <div class="actions">
          <a class="btn primary" href="/hr">Abrir fichas</a>
          <button type="button" class="btn secondary" onclick={reset}>Onboard outro</button>
        </div>
      </section>
    {/if}
  {/if}
</section>

<style>
  .page {
    max-width: 40rem;
  }
  .page-head {
    margin-bottom: var(--space-4);
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
  .back {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    text-decoration: none;
  }
  .back:hover {
    color: var(--color-text);
  }
  .lead {
    margin: 0 0 var(--space-4);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-5) var(--space-6);
    margin-bottom: var(--space-4);
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
  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    color: var(--color-text);
  }
  .btn.sm {
    padding: var(--space-1) var(--space-3);
    font-size: var(--text-xs);
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
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
</style>
