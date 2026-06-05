<script lang="ts">
  import type { DocFlowStep } from "./types";
  import { prefersReducedMotion } from "$lib/motion/reduced";

  interface Props {
    messageSteps: Extract<DocFlowStep, { kind: "message" }>[];
    stepIndex: number;
    playing: boolean;
    onStep: (idx: number) => void;
    onTogglePlay: () => void;
  }

  let { messageSteps, stepIndex, playing, onStep, onTogglePlay }: Props = $props();

  const total = $derived(messageSteps.length);
  const current = $derived(messageSteps[stepIndex]);

  function stepLabel(step: Extract<DocFlowStep, { kind: "message" }>): string {
    return `${step.from} → ${step.to}: ${step.label}`;
  }
</script>

<div class="flow-controls" role="group" aria-label="Percorrer o fluxo passo a passo">
  <div class="control-row">
    <button
      type="button"
      class="ctrl-btn"
      disabled={stepIndex === 0}
      onclick={() => onStep(stepIndex - 1)}
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
      onclick={() => onStep(stepIndex + 1)}
      aria-label="Passo seguinte"
    >
      ▶
    </button>
    {#if !prefersReducedMotion()}
      <button type="button" class="ctrl-btn play" onclick={onTogglePlay} aria-pressed={playing}>
        {playing ? "⏸ Pausar" : "▶ Reproduzir"}
      </button>
    {/if}
  </div>

  <p class="flow-legend" aria-hidden="true">
    <span class="leg current">Passo actual</span>
    <span class="leg done">Concluído</span>
    <span class="leg from">Origem</span>
    <span class="leg to">Destino</span>
  </p>

  {#if current}
    <p class="current-step" aria-live="polite">
      <strong class="from-tag">{current.from}</strong>
      <span class="arrow">{current.arrow.includes("--") ? "⟹" : "→"}</span>
      <strong class="to-tag">{current.to}</strong>
      <span class="msg">{current.label}</span>
    </p>
  {/if}

  <ol class="step-list">
    {#each messageSteps as step, idx (idx)}
      <li class:active={idx === stepIndex} class:done={idx < stepIndex}>
        <button
          type="button"
          class="step-btn"
          onclick={() => onStep(idx)}
          aria-current={idx === stepIndex ? "step" : undefined}
        >
          <span class="step-num">{idx + 1}</span>
          <span class="step-text">{stepLabel(step)}</span>
        </button>
      </li>
    {/each}
  </ol>
</div>

<style>
  .flow-controls {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4);
    margin-top: 0;
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
    max-height: min(22rem, 50vh);
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
  }

  .flow-legend {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-2);
    margin: 0 0 var(--space-3);
    font-size: 0.65rem;
    color: var(--color-text-muted);
  }

  .leg {
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
  }

  .leg.current {
    border-color: var(--color-accent);
    background: color-mix(in srgb, var(--color-accent) 12%, transparent);
  }

  .leg.done {
    border-color: var(--color-success-fg);
    color: var(--color-success-fg);
  }

  .leg.from {
    border-color: var(--color-warning);
    color: var(--color-warning);
  }

  .leg.to {
    border-color: var(--color-success-fg);
    background: var(--color-success-bg);
    color: var(--color-success-fg);
  }

  .from-tag {
    color: var(--color-warning);
  }

  .to-tag {
    color: var(--color-success-fg);
  }

  .step-list li.done .step-btn {
    color: var(--color-success-fg);
    border-color: color-mix(in srgb, var(--color-success-fg) 35%, transparent);
    background: color-mix(in srgb, var(--color-success-fg) 8%, transparent);
  }

  .step-list li.active .step-btn {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .step-list li.done .step-num {
    background: var(--color-success-bg);
    color: var(--color-success-fg);
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
</style>
