<script lang="ts">
  import type { DocFlow, DocFlowStep } from "./types";
  import MermaidBlock from "./MermaidBlock.svelte";
  import { prefersReducedMotion } from "$lib/motion/reduced";

  interface Props {
    flow: DocFlow;
  }

  let { flow }: Props = $props();

  let stepIndex = $state(0);
  let playing = $state(false);

  const messageSteps = $derived(
    flow.steps.filter((s): s is Extract<DocFlowStep, { kind: "message" }> => s.kind === "message"),
  );
  const hasSteps = $derived(messageSteps.length > 0);
  const current = $derived(messageSteps[stepIndex]);
  const total = $derived(messageSteps.length);

  function stepLabel(step: Extract<DocFlowStep, { kind: "message" }>): string {
    return `${step.from} → ${step.to}: ${step.label}`;
  }

  function goTo(idx: number) {
    stepIndex = Math.max(0, Math.min(idx, total - 1));
  }

  function togglePlay() {
    if (prefersReducedMotion()) return;
    playing = !playing;
  }

  // Reprodução automática: avança um passo de cada ~2.8s (respeita pausa).
  $effect(() => {
    if (!playing || total === 0) return;
    const id = setInterval(() => {
      if (stepIndex >= total - 1) {
        playing = false;
        return;
      }
      goTo(stepIndex + 1);
    }, 2800);
    return () => clearInterval(id);
  });
</script>

<div class="flow-player" data-flow-type={flow.type}>
  <MermaidBlock source={flow.source} id={flow.id} title={flow.title || undefined} />

  {#if hasSteps}
    <div class="flow-controls" role="group" aria-label="Percorrer o fluxo passo a passo">
      <div class="control-row">
        <button
          type="button"
          class="ctrl-btn"
          disabled={stepIndex === 0}
          onclick={() => goTo(stepIndex - 1)}
          aria-label="Passo anterior"
        >
          ◀
        </button>
        <span class="step-counter" aria-live="polite">
          Passo {stepIndex + 1} / {total}
        </span>
        <button
          type="button"
          class="ctrl-btn"
          disabled={stepIndex >= total - 1}
          onclick={() => goTo(stepIndex + 1)}
          aria-label="Passo seguinte"
        >
          ▶
        </button>
        {#if !prefersReducedMotion()}
          <button
            type="button"
            class="ctrl-btn play"
            onclick={togglePlay}
            aria-pressed={playing}
          >
            {playing ? "⏸ Pausar" : "▶ Reproduzir"}
          </button>
        {/if}
      </div>

      {#if current}
        <p class="current-step" aria-live="polite">
          <strong>{current.from}</strong>
          <span class="arrow">{current.arrow.includes(">>") ? "⟹" : "→"}</span>
          <strong>{current.to}</strong>
          <span class="msg">{current.label}</span>
        </p>
      {/if}

      <ol class="step-list">
        {#each messageSteps as step, idx (idx)}
          <li class:active={idx === stepIndex} class:done={idx < stepIndex}>
            <button
              type="button"
              class="step-btn"
              onclick={() => goTo(idx)}
              aria-current={idx === stepIndex ? "step" : undefined}
            >
              <span class="step-num">{idx + 1}</span>
              <span class="step-text">{stepLabel(step)}</span>
            </button>
          </li>
        {/each}
      </ol>
    </div>
  {:else if flow.type === "flowchart"}
    <p class="flow-hint">
      Diagrama estático — segue as setas do fluxograma na ordem lógica.
    </p>
  {/if}
</div>

<style>
  .flow-player {
    margin: var(--space-4) 0;
  }

  .flow-controls {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4);
    margin-top: var(--space-2);
  }

  .control-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: var(--space-2);
    margin-bottom: var(--space-3);
  }

  .ctrl-btn {
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-base);
    color: var(--color-text);
    font-size: var(--text-sm);
    cursor: pointer;
    transition: background-color var(--duration-fast) var(--ease-out);
  }

  .ctrl-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .ctrl-btn.play {
    margin-left: auto;
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .step-counter {
    font-size: var(--text-sm);
    font-family: var(--font-mono);
    color: var(--color-text-muted);
    min-width: 6rem;
    text-align: center;
  }

  .current-step {
    margin: 0 0 var(--space-4);
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-sm);
    background: var(--color-accent-muted);
    border-left: 3px solid var(--color-accent);
    font-size: var(--text-sm);
    line-height: 1.5;
    transition: opacity var(--duration-fast) var(--ease-out);
  }

  .arrow {
    margin: 0 var(--space-2);
    color: var(--color-accent);
  }

  .msg {
    display: block;
    margin-top: var(--space-1);
    color: var(--color-text-muted);
  }

  .step-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    max-height: 14rem;
    overflow-y: auto;
  }

  .step-btn {
    display: flex;
    align-items: flex-start;
    gap: var(--space-2);
    width: 100%;
    text-align: left;
    padding: var(--space-2) var(--space-3);
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    line-height: 1.4;
    transition:
      background-color var(--duration-fast) var(--ease-out),
      border-color var(--duration-fast) var(--ease-out);
  }

  .step-list li.active .step-btn {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .step-list li.done .step-btn {
    opacity: 0.75;
  }

  .step-num {
    flex-shrink: 0;
    width: 1.25rem;
    height: 1.25rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    background: var(--color-bg-base);
    font-weight: 600;
    font-size: 0.65rem;
  }

  .step-list li.active .step-num {
    background: var(--color-accent);
    color: var(--color-bg-base);
  }

  .flow-hint {
    margin: var(--space-2) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  @media (prefers-reduced-motion: reduce) {
    .ctrl-btn,
    .current-step,
    .step-btn {
      transition: none;
    }
  }
</style>
