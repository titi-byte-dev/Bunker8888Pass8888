<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { maskSensitiveText, unmaskText } from "$lib/google/masking";
  import {
    deleteDriveFile,
    enrichDriveNames,
    listDriveFiles,
    openDriveFile,
    opaquePreview,
    preferredDriveBackend,
    uploadDriveFile,
    type DriveBackend,
    type DriveFileView,
  } from "$lib/google/driveService";
  import { Button, PageShell, Panel, StatusBanner } from "$lib/ui";

  let locked = $state(false);
  let tab = $state<"drive" | "sheets">("drive");
  let backend = $state<DriveBackend>("server");

  let fileName = $state("contrato.txt");
  let fileContent = $state("Documento confidencial da empresa.");
  let driveFiles = $state<DriveFileView[]>([]);
  let opened = $state("");

  let sheetRaw = $state("Cliente: ACME\nNIF: 123456789\nIBAN: PT50002600009987654321076\nSalário: 2500");
  let maskedView = $state("");
  let tokenMap = $state<Record<string, string>>({});

  async function refreshDrive() {
    const mk = getMasterKey();
    if (!mk) return;
    const raw = await listDriveFiles(backend);
    driveFiles = await enrichDriveNames(mk, raw);
  }

  async function handleUpload() {
    const mk = getMasterKey();
    if (!mk) return;
    await uploadDriveFile(mk, fileName.trim() || "ficheiro.txt", fileContent, backend);
    fileContent = "";
    await refreshDrive();
  }

  async function handleOpen(f: DriveFileView) {
    const mk = getMasterKey();
    if (!mk) return;
    const payload = await openDriveFile(mk, f);
    opened = payload.content;
  }

  async function handleDelete(id: string) {
    await deleteDriveFile(id, backend);
    await refreshDrive();
  }

  function applyMask() {
    const r = maskSensitiveText(sheetRaw);
    maskedView = r.masked;
    tokenMap = r.tokens;
  }

  onMount(() => {
    locked = !getMasterKey();
    backend = preferredDriveBackend();
    if (!locked) void refreshDrive();
    applyMask();
  });
</script>

<svelte:head><title>Google proxy (dev) — AegisPass</title></svelte:head>

<PageShell
  title="Google Workspace (stub)"
  taskId="GOOGLE-002"
  description="Drive com cifragem ZK: servidor guarda blobs opacos (GOOGLE-002); local usa só localStorage. Sheets com mascaramento NIF/IBAN. Simulação dev — DoD Fase 2."
>
  {#snippet actions()}
    <DocHelpLink slug="journey-google-dev-stub" label="Como funciona sem VPS?" />
    <Button variant="ghost" size="sm" href="/work/google">Estado OAuth →</Button>
  {/snippet}

  <div class="tabs">
    <button type="button" class:active={tab === "drive"} onclick={() => (tab = "drive")}>Drive</button>
    <button type="button" class:active={tab === "sheets"} onclick={() => (tab = "sheets")}>Sheets</button>
  </div>

  {#if locked}
    <StatusBanner variant="warning">
      Desbloqueia a Master Key em <a href="/vault">/vault</a> para usar o stub.
    </StatusBanner>
  {:else if tab === "drive"}
    <Panel>
      <div class="row">
        <label>Armazenamento
          <select bind:value={backend} onchange={() => refreshDrive()}>
            <option value="server">Servidor (PostgreSQL opaco)</option>
            <option value="local">localStorage (offline)</option>
          </select>
        </label>
      </div>
      <h2>Upload ZK</h2>
      <label>Nome<input bind:value={fileName} /></label>
      <label>Conteúdo<textarea bind:value={fileContent} rows="3"></textarea></label>
      <button type="button" class="btn primary" onclick={handleUpload}>Enviar</button>
    </Panel>
    <Panel>
      <h2>Ficheiros (vista Google = opaco)</h2>
      {#if driveFiles.length === 0}<p class="muted">Vazio.</p>
      {:else}
        <ul class="files">
          {#each driveFiles as f (f.id)}
            <li>
              <div>
                <strong>{f.name}</strong>
                <span class="opaque mono">{opaquePreview(f)}</span>
              </div>
              <div class="actions">
                <button type="button" class="btn sm" onclick={() => handleOpen(f)}>Abrir</button>
                <button type="button" class="btn sm danger" onclick={() => handleDelete(f.id)}>Apagar</button>
              </div>
            </li>
          {/each}
        </ul>
      {/if}
      {#if opened}<pre class="opened">{opened}</pre>{/if}
    </Panel>
  {:else}
    <Panel>
      <h2>Mascaramento (o que a Google veria)</h2>
      <label>Dados sensíveis<textarea bind:value={sheetRaw} rows="5" oninput={applyMask}></textarea></label>
      <h3>Vista Sheet (tokens)</h3>
      <pre class="masked">{maskedView}</pre>
      <h3>Vista AegisPass (desmascarado)</h3>
      <pre class="opened">{unmaskText(maskedView, tokenMap)}</pre>
    </Panel>
  {/if}
</PageShell>

<style>
  .tabs { display: flex; gap: var(--space-2); margin-bottom: var(--space-4); }
  .tabs button { padding: var(--space-2) var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-sm); background: var(--color-bg-inset); cursor: pointer; }
  .tabs button.active { border-color: var(--color-accent); color: var(--color-accent); }
  .row { margin-bottom: var(--space-3); }
  label { display: block; margin-bottom: var(--space-3); font-size: var(--text-sm); }
  input, textarea, select { width: 100%; margin-top: var(--space-1); padding: var(--space-2); border: 1px solid var(--color-border); border-radius: var(--radius-sm); box-sizing: border-box; }
  .btn { padding: var(--space-2) var(--space-4); border-radius: var(--radius-sm); border: 1px solid var(--color-border); cursor: pointer; }
  .btn.primary { background: var(--color-accent); color: var(--color-accent-fg); border-color: transparent; }
  .btn.sm { font-size: var(--text-xs); }
  .btn.danger { color: var(--color-danger, #b91c1c); }
  .files { list-style: none; padding: 0; }
  .files li { display: flex; justify-content: space-between; align-items: center; padding: var(--space-2); border-bottom: 1px solid var(--color-border); gap: var(--space-2); }
  .actions { display: flex; gap: var(--space-2); flex-shrink: 0; }
  .opaque { display: block; font-size: var(--text-xs); color: var(--color-text-muted); }
  .mono { font-family: var(--font-mono); }
  pre { background: var(--color-bg-inset); padding: var(--space-3); border-radius: var(--radius-sm); font-size: var(--text-xs); overflow-x: auto; }
  .muted { color: var(--color-text-muted); font-size: var(--text-sm); }
  h2 { margin: 0 0 var(--space-3); font-size: var(--text-lg); }
  h3 { margin: var(--space-4) 0 var(--space-2); font-size: var(--text-base); }
</style>
