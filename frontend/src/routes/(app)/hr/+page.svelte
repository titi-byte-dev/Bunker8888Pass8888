<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { fieldLabel, SUGGESTED_FIELDS } from "$lib/hr/employees";
  import {
    createRecord,
    deleteRecord,
    listCertificates,
    listRecords,
    openRecord,
    removeField,
    saveField,
    shredField,
    shredRecord,
    type OpenRecord,
    type RecordSummary,
    type VerifiedCertificate,
  } from "$lib/hr/employeesService";
  import { certificateToJSON } from "$lib/hr/erasure";
  import {
    downloadContract,
    listContracts,
    removeContract,
    signContract,
    uploadContract,
    verifyContract,
    type DecryptedContract,
  } from "$lib/hr/contractsService";

  let locked = $state(false);
  let loading = $state(true);
  let error = $state("");
  let records = $state<RecordSummary[]>([]);
  let open = $state<OpenRecord | null>(null);
  let busy = $state(false);

  // Formulário de novo campo.
  let newFieldName = $state("");
  let newFieldValue = $state("");

  // Certificados de eliminação (HR-004).
  let certs = $state<VerifiedCertificate[]>([]);

  // Contratos (HR-005/006).
  let contracts = $state<DecryptedContract[]>([]);
  let contractInput = $state<HTMLInputElement | null>(null);
  let verifyMsg = $state<Record<string, string>>({});

  async function refreshContracts() {
    if (!open) return;
    try {
      contracts = await listContracts(open.id);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar contratos";
    }
  }

  async function onUploadContract(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file || !open) return;
    busy = true;
    error = "";
    try {
      await uploadContract(open.id, file);
      await refreshContracts();
    } catch (err) {
      error = err instanceof Error ? err.message : "Falha ao carregar contrato";
    } finally {
      busy = false;
      if (contractInput) contractInput.value = "";
    }
  }

  async function onDownloadContract(id: string) {
    if (!open) return;
    try {
      const { blob, name } = await downloadContract(open.id, id);
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = name;
      a.click();
      URL.revokeObjectURL(url);
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao descarregar contrato";
    }
  }

  async function onSignContract(id: string) {
    if (!open) return;
    busy = true;
    error = "";
    try {
      await signContract(open.id, id);
      await refreshContracts();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao assinar contrato";
    } finally {
      busy = false;
    }
  }

  async function onVerifyContract(id: string) {
    if (!open) return;
    try {
      const res = await verifyContract(open.id, id);
      verifyMsg = { ...verifyMsg, [id]: res.reason };
    } catch (e) {
      verifyMsg = { ...verifyMsg, [id]: e instanceof Error ? e.message : "Falha" };
    }
  }

  async function onRemoveContract(id: string) {
    if (!open) return;
    busy = true;
    try {
      await removeContract(open.id, id);
      await refreshContracts();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao remover contrato";
    } finally {
      busy = false;
    }
  }

  function fmtBytes(n: number): string {
    if (n < 1024) return `${n} B`;
    if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`;
    return `${(n / (1024 * 1024)).toFixed(1)} MB`;
  }

  async function refreshCerts() {
    try {
      certs = await listCertificates();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar certificados";
    }
  }

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
    if (!locked) {
      refresh();
      refreshCerts();
    } else loading = false;
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
    verifyMsg = {};
    try {
      open = await openRecord(id);
      await refreshContracts();
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

  async function onShredField(name: string) {
    if (!open) return;
    if (
      !confirm(
        `Crypto-shred do campo "${fieldLabel(name)}"? A chave é destruída e o valor fica irrecuperável (RGPD Art. 17). É emitido um certificado.`,
      )
    )
      return;
    busy = true;
    error = "";
    try {
      await shredField(open.id, name);
      open = await openRecord(open.id);
      await refreshCerts();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao eliminar campo";
    } finally {
      busy = false;
    }
  }

  async function onShredRecord() {
    if (!open) return;
    if (!confirm("Crypto-shred de TODA a ficha? Cada campo emite um certificado.")) return;
    busy = true;
    error = "";
    try {
      await shredRecord(open.id);
      open = await openRecord(open.id);
      await refreshCerts();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao eliminar ficha";
    } finally {
      busy = false;
    }
  }

  function downloadCert(cert: VerifiedCertificate) {
    const blob = new Blob([certificateToJSON(cert)], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `erasure-${cert.fieldName}-${cert.id.slice(0, 8)}.json`;
    a.click();
    URL.revokeObjectURL(url);
  }

  async function onDeleteRecord() {
    if (!open) return;
    if (!confirm("Apagar a ficha inteira? Os campos cifrados vão com ela.")) return;
    busy = true;
    try {
      await deleteRecord(open.id);
      open = null;
      await refresh();
      await refreshCerts();
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
    <span class="report-link">
      <a class="btn ghost" href="/hr/onboarding">Onboarding →</a>
      <a class="btn ghost" href="/hr/compliance">Relatório RGPD →</a>
    </span>
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
            <span class="head-actions">
              <button type="button" class="btn ghost sm" onclick={onShredRecord} disabled={busy}>
                🔥 Shred ficha
              </button>
              <button type="button" class="btn danger-btn sm" onclick={onDeleteRecord} disabled={busy}>
                Apagar ficha
              </button>
            </span>
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
                    <span class="field-actions">
                      {#if !f.shredded}
                        <button
                          type="button"
                          class="link-btn shred"
                          onclick={() => onShredField(f.name)}
                          disabled={busy}
                          title="Crypto-shred (RGPD Art. 17)"
                        >
                          shred
                        </button>
                      {/if}
                      <button
                        type="button"
                        class="link-btn"
                        onclick={() => onRemoveField(f.name)}
                        disabled={busy}
                      >
                        remover
                      </button>
                    </span>
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

          <!-- Contratos cifrados (HR-005) + assinatura digital (HR-006) -->
          <div class="contracts">
            <div class="contracts-head">
              <p class="eyebrow">Contratos cifrados</p>
              <label class="btn primary sm upload">
                + Carregar
                <input
                  type="file"
                  bind:this={contractInput}
                  onchange={onUploadContract}
                  disabled={busy}
                  hidden
                />
              </label>
            </div>
            {#if contracts.length === 0}
              <p class="muted sm">Sem contratos. Carrega um ficheiro (máx. ~5 MiB).</p>
            {:else}
              <ul class="contract-list">
                {#each contracts as c (c.id)}
                  <li class="contract">
                    <div class="c-main">
                      <span class="c-name">{c.name}</span>
                      <span class="muted sm">{fmtBytes(c.size)}</span>
                      {#if c.signed}
                        <span class="c-badge">✍ assinado</span>
                      {/if}
                    </div>
                    <div class="c-actions">
                      <button type="button" class="link-btn" onclick={() => onDownloadContract(c.id)}>
                        descarregar
                      </button>
                      {#if c.signed}
                        <button type="button" class="link-btn" onclick={() => onVerifyContract(c.id)}>
                          verificar
                        </button>
                      {:else}
                        <button
                          type="button"
                          class="link-btn sign"
                          onclick={() => onSignContract(c.id)}
                          disabled={busy}
                        >
                          assinar
                        </button>
                      {/if}
                      <button type="button" class="link-btn" onclick={() => onRemoveContract(c.id)} disabled={busy}>
                        remover
                      </button>
                    </div>
                    {#if verifyMsg[c.id]}
                      <p class="c-verify" role="status">{verifyMsg[c.id]}</p>
                    {/if}
                  </li>
                {/each}
              </ul>
            {/if}
          </div>
        {/if}
      </section>
    </div>

    <!-- Certificados de eliminação (HR-004) -->
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">Certificados de eliminação · RGPD Art. 17</p>
        <span class="muted sm">{certs.length} emitido(s)</span>
      </div>
      <p class="muted">
        Cada crypto-shred emite uma prova verificável: <span class="mono">sha256</span> sobre os
        factos da eliminação. O selo ✓ confirma que o <span class="mono">cert_hash</span> foi
        recalculado no teu dispositivo e bate certo.
      </p>
      {#if certs.length > 0}
        <ul class="cert-list">
          {#each certs as c (c.id)}
            <li class="cert">
              <span class="cert-badge" class:ok={c.valid} class:bad={!c.valid}>
                {c.valid ? "✓ íntegro" : "✗ inválido"}
              </span>
              <span class="cert-field">{fieldLabel(c.fieldName)}</span>
              <span class="muted sm mono">{c.certHash.slice(0, 16)}…</span>
              <span class="muted sm">{new Date(c.shreddedAt).toLocaleString("pt-PT")}</span>
              <button type="button" class="link-btn" onclick={() => downloadCert(c)}>
                descarregar
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>
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
  .link-btn.shred:hover {
    color: var(--color-danger);
  }
  .field-actions {
    display: inline-flex;
    gap: var(--space-2);
    flex-shrink: 0;
  }
  .head-actions {
    display: inline-flex;
    gap: var(--space-2);
  }
  .report-link {
    display: inline-flex;
    gap: var(--space-2);
    margin-top: var(--space-3);
  }
  .report-link a {
    text-decoration: none;
  }
  .btn.ghost {
    background: none;
    color: var(--color-text-muted);
  }
  .btn.ghost:hover:not(:disabled) {
    color: var(--color-danger);
    border-color: var(--color-danger);
  }

  .cert-list {
    margin: var(--space-3) 0 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .cert {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    font-size: var(--text-sm);
  }
  .cert-field {
    font-weight: 500;
  }
  .cert-badge {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
  }
  .cert-badge.ok {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
  }
  .cert-badge.bad {
    color: var(--color-danger);
    border: 1px solid var(--color-danger);
  }
  .cert .link-btn {
    margin-left: auto;
    color: var(--color-accent);
  }
  .inline-error {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-danger);
  }

  .contracts {
    border-top: 1px solid var(--color-border);
    margin-top: var(--space-4);
    padding-top: var(--space-4);
  }
  .contracts-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: var(--space-2);
  }
  .contracts-head .eyebrow {
    margin: 0;
  }
  .upload {
    cursor: pointer;
  }
  .contract-list {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .contract {
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    padding: var(--space-2) var(--space-3);
  }
  .c-main {
    display: flex;
    align-items: center;
    gap: var(--space-3);
  }
  .c-name {
    font-size: var(--text-sm);
    font-weight: 500;
    word-break: break-word;
  }
  .c-badge {
    font-size: var(--text-xs);
    font-weight: 600;
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    padding: 1px var(--space-2);
    border-radius: var(--radius-sm);
  }
  .c-actions {
    display: flex;
    gap: var(--space-3);
    margin-top: var(--space-1);
  }
  .link-btn.sign:hover {
    color: var(--color-accent);
  }
  .c-verify {
    margin: var(--space-2) 0 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }
</style>
