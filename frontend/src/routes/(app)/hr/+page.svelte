<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { fieldLabel, SUGGESTED_FIELDS } from "$lib/hr/employees";
  import {
    createRecord,
    deleteRecord,
    listRecords,
    openRecord,
    removeField,
    saveField,
    type OpenRecord,
    type RecordSummary,
  } from "$lib/hr/employeesService";

  let locked = $state(false);
  let loading = $state(true);
  let error = $state("");
  let records = $state<RecordSummary[]>([]);
  let open = $state<OpenRecord | null>(null);
  let busy = $state(false);

  // Formulário de novo campo.
  let newFieldName = $state("");
  let newFieldValue = $state("");

  async function refresh() {
    loading = true;
    error = "";
    try {
      records = await listRecords();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar fichas";
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) refresh();
    else loading = false;
  });

  async function onCreate() {
    busy = true;
    error = "";
    try {
      const id = await createRecord();
      await refresh();
      await onOpen(id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao criar ficha";
    } finally {
      busy = false;
    }
  }

  async function onOpen(id: string) {
    busy = true;
    error = "";
    try {
      open = await openRecord(id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao abrir ficha";
    } finally {
      busy = false;
    }
  }

  async function onSaveField() {
    if (!open || !newFieldName.trim() || !newFieldValue) return;
    busy = true;
    error = "";
    try {
      await saveField(open.id, newFieldName.trim(), newFieldValue);
      newFieldName = "";
      newFieldValue = "";
      open = await openRecord(open.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao gravar campo";
    } finally {
      busy = false;
    }
  }

  async function onRemoveField(name: string) {
    if (!open) return;
    busy = true;
    try {
      await removeField(open.id, name);
      open = await openRecord(open.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao remover campo";
    } finally {
      busy = false;
    }
  }

  async function onDeleteRecord() {
    if (!open) return;
    if (!confirm("Apagar a ficha inteira? Os campos cifrados vão com ela.")) return;
    busy = true;
    try {
      await deleteRecord(open.id);
      open = null;
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao apagar ficha";
    } finally {
      busy = false;
    }
  }

  function shortId(id: string): string {
    return id.slice(0, 8);
  }
</script>

<svelte:head>
  <title>Fichas de Empregado — AegisPass</title>
</svelte:head>

<section class="page">
  <header class="page-head">
    <div>
      <p class="eyebrow">HR-001 · Cifragem campo-a-campo</p>
      <h1>Fichas de Empregado</h1>
    </div>
    <p class="lead">
      Cada campo da ficha é cifrado <strong>de forma independente</strong>, com a
      sua própria chave. Essa chave é embrulhada com a tua Master Key — o servidor
      só vê blobs opacos. Apagar a chave de um campo torna-o irrecuperável
      (fundação do crypto-shredding RGPD, HR-003).
    </p>
  </header>

  {#if locked}
    <section class="panel">
      <p class="panel-body">
        🔒 Cofre bloqueado. Desbloqueia a tua Master Key para gerir fichas — é ela
        que embrulha as chaves de cada campo.
      </p>
      <a class="btn primary" href="/vault">Ir desbloquear</a>
    </section>
  {:else}
    {#if error}<p class="inline-error" role="alert">{error}</p>{/if}

    <div class="grid">
      <!-- Coluna esquerda: lista de fichas -->
      <section class="panel">
        <div class="panel-head">
          <p class="eyebrow">Fichas</p>
          <button type="button" class="btn primary sm" onclick={onCreate} disabled={busy}>
            + Nova
          </button>
        </div>
        {#if loading}
          <p class="muted">A carregar…</p>
        {:else if records.length === 0}
          <p class="muted">Sem fichas. Cria a primeira.</p>
        {:else}
          <ul class="list">
            {#each records as r (r.id)}
              <li>
                <button
                  type="button"
                  class="row"
                  class:active={open?.id === r.id}
                  onclick={() => onOpen(r.id)}
                >
                  <span class="mono">{shortId(r.id)}</span>
                  <span class="muted sm">{new Date(r.updatedAt).toLocaleDateString("pt-PT")}</span>
                </button>
              </li>
            {/each}
          </ul>
        {/if}
      </section>

      <!-- Coluna direita: detalhe da ficha -->
      <section class="panel">
        {#if !open}
          <p class="muted">Seleciona ou cria uma ficha para ver os campos.</p>
        {:else}
          <div class="panel-head">
            <p class="eyebrow">Ficha · <span class="mono">{shortId(open.id)}</span></p>
            <button type="button" class="btn danger-btn sm" onclick={onDeleteRecord} disabled={busy}>
              Apagar ficha
            </button>
          </div>

          {#if open.fields.length === 0}
            <p class="muted">Ficha vazia. Adiciona o primeiro campo abaixo.</p>
          {:else}
            <dl class="fields">
              {#each open.fields as f (f.name)}
                <div class="field-row">
                  <dt>{fieldLabel(f.name)}</dt>
                  <dd>
                    {#if f.shredded}
                      <span class="shredded">🔥 campo destruído (sem chave)</span>
                    {:else}
                      <span class="value">{f.value}</span>
                    {/if}
                    <button
                      type="button"
                      class="link-btn"
                      onclick={() => onRemoveField(f.name)}
                      disabled={busy}
                    >
                      remover
                    </button>
                  </dd>
                </div>
              {/each}
            </dl>
          {/if}

          <div class="add-field">
            <p class="eyebrow">Adicionar / actualizar campo</p>
            <div class="row-form">
              <input
                list="suggested-fields"
                bind:value={newFieldName}
                placeholder="Campo (ex.: salary)"
                disabled={busy}
              />
              <datalist id="suggested-fields">
                {#each SUGGESTED_FIELDS as s (s)}
                  <option value={s}>{fieldLabel(s)}</option>
                {/each}
              </datalist>
              <input
                bind:value={newFieldValue}
                placeholder="Valor (cifrado no teu dispositivo)"
                disabled={busy}
              />
              <button
                type="button"
                class="btn primary"
                onclick={onSaveField}
                disabled={busy || !newFieldName.trim() || !newFieldValue}
              >
                Gravar
              </button>
            </div>
          </div>
        {/if}
      </section>
    </div>
  {/if}
</section>

<style>
  .page {
    max-width: 56rem;
  }
  .page-head {
    margin-bottom: var(--space-6);
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
    line-height: var(--leading-tight);
  }
  .lead {
    margin: var(--space-3) 0 0;
    max-width: 42rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .grid {
    display: grid;
    grid-template-columns: 16rem 1fr;
    gap: var(--space-4);
    align-items: start;
  }
  @media (max-width: 720px) {
    .grid {
      grid-template-columns: 1fr;
    }
  }

  .panel {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
    padding: var(--space-4) var(--space-6);
    margin-bottom: var(--space-4);
  }
  .panel-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    margin-bottom: var(--space-3);
  }
  .panel-head .eyebrow {
    margin: 0;
  }
  .panel-body {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-text);
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

  .list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .row {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
    padding: var(--space-2) var(--space-3);
    border: 1px solid transparent;
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    cursor: pointer;
    font-size: var(--text-sm);
  }
  .row:hover {
    border-color: var(--color-border);
  }
  .row.active {
    border-color: var(--color-accent);
  }

  .fields {
    margin: 0 0 var(--space-4);
  }
  .field-row {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-3) 0;
    border-bottom: 1px solid var(--color-border);
  }
  .field-row dt {
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
    flex-shrink: 0;
    width: 9rem;
  }
  .field-row dd {
    margin: 0;
    flex: 1;
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: var(--space-3);
  }
  .value {
    font-size: var(--text-sm);
    color: var(--color-text);
    word-break: break-word;
  }
  .shredded {
    font-size: var(--text-sm);
    color: var(--color-danger);
  }

  .add-field {
    border-top: 1px solid var(--color-border);
    padding-top: var(--space-4);
  }
  .row-form {
    display: flex;
    gap: var(--space-2);
    margin-top: var(--space-2);
  }
  @media (max-width: 560px) {
    .row-form {
      flex-direction: column;
    }
  }
  input {
    flex: 1;
    min-width: 0;
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

  .btn {
    display: inline-block;
    padding: var(--space-2) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    text-decoration: none;
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
  .btn.danger-btn {
    color: var(--color-danger);
    border-color: var(--color-danger);
    background: none;
  }
  .btn.danger-btn:hover {
    background: var(--color-danger);
    color: var(--color-accent-fg);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }
  .link-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    font-size: var(--text-xs);
    cursor: pointer;
    padding: 0;
    flex-shrink: 0;
  }
  .link-btn:hover {
    color: var(--color-danger);
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
</style>
