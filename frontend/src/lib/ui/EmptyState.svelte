<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * EmptyState (UI-012) — estado vazio consistente com CTA opcional.
   * Substitui blocos `.empty` variados. CTA pela snippet `action`.
   */
  interface Props {
    title: string;
    description?: string;
    /** Glifo/emoji decorativo opcional. */
    icon?: string;
    action?: Snippet;
  }
  let { title, description, icon, action }: Props = $props();
</script>

<div class="empty">
  {#if icon}<div class="icon" aria-hidden="true">{icon}</div>{/if}
  <h3>{title}</h3>
  {#if description}<p>{description}</p>{/if}
  {#if action}<div class="cta">{@render action()}</div>{/if}
</div>

<style>
  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: var(--space-2);
    padding: var(--space-12) var(--space-4);
    color: var(--color-text-muted);
  }
  .icon { font-size: var(--text-3xl); opacity: 0.7; }
  h3 {
    margin: 0;
    font-family: var(--font-display);
    font-size: var(--text-lg);
    color: var(--color-text);
  }
  p { margin: 0; font-size: var(--text-sm); max-width: 32rem; }
  .cta { margin-top: var(--space-2); }
</style>
