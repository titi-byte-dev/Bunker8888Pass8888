<script lang="ts">
  import type { Snippet } from "svelte";

  /**
   * Field (UI-012) — label + input + hint/erro acessivel.
   * Para inputs nativos usa o snippet `control` (recebe o id a ligar).
   * Substitui labels+inputs inline espalhados pelos formularios.
   */
  interface Props {
    label: string;
    /** Texto de ajuda neutro por baixo do campo. */
    hint?: string;
    /** Mensagem de erro — quando presente, marca aria-invalid. */
    error?: string;
    required?: boolean;
    /** Recebe { id, describedBy } para ligar ao input nativo. */
    control: Snippet<[{ id: string; describedBy: string | undefined }]>;
  }
  let { label, hint, error, required = false, control }: Props = $props();

  const id = `f-${Math.random().toString(36).slice(2, 9)}`;
  const msgId = `${id}-msg`;
  const describedBy = $derived(hint || error ? msgId : undefined);
</script>

<div class="field" class:has-error={!!error}>
  <label for={id}>
    {label}{#if required}<span class="req" aria-hidden="true"> *</span>{/if}
  </label>
  {@render control({ id, describedBy })}
  {#if error}
    <p class="msg err" id={msgId} role="alert">{error}</p>
  {:else if hint}
    <p class="msg hint" id={msgId}>{hint}</p>
  {/if}
</div>

<style>
  .field { display: flex; flex-direction: column; gap: var(--space-1); }
  label {
    font-size: var(--text-xs);
    font-weight: 600;
    color: var(--color-text-label);
  }
  .req { color: var(--color-danger); }
  .msg { margin: 0; font-size: var(--text-xs); }
  .hint { color: var(--color-text-muted); }
  .err { color: var(--color-danger); }
</style>
