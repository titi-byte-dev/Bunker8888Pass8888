<script lang="ts">
  import { onMount } from "svelte";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { maskSensitiveText, unmaskText } from "$lib/google/masking";
  import {
    listMockDriveFiles,
    openMockDriveFile,
    opaquePreview,
    uploadToMockDrive,
    type DriveDevFile,
  } from "$lib/google/driveDevStore";

  let locked = $state(false);
  let tab = $state<"drive" | "sheets">("drive");

  let fileName = $state("contrato.txt");
  let fileContent = $state("Documento confidencial da empresa.");
  let driveFiles = $state<DriveDevFile[]>([]);
  let opened = $state("");

  let sheetRaw = $state("Cliente: ACME\nNIF: 123456789\nIBAN: PT50002600009987654321076\nSalário: 2500");
  let maskedView = $state("");
  let tokenMap = $state<Record<string, string>>({});

  async function refreshDrive() {
    driveFiles = listMockDriveFiles();
  }

  async function handleUpload() {
    const mk = getMasterKey();
    if (!mk) return;
    await uploadToMockDrive(mk, fileName.trim() || "ficheiro.txt", fileContent);
    fileContent = "";
    await refreshDrive();
  }

  async function handleOpen(f: DriveDevFile) {
    const mk = getMasterKey();
    if (!mk) return;
    opened = await openMockDriveFile(mk, f);
  }

  function applyMask() {
    const r = maskSensitiveText(sheetRaw);
    maskedView = r.masked;
    tokenMap = r.tokens;
  }

  onMount(() => {
    locked = !getMasterKey();
    if (!locked) void refreshDrive();
    applyMask();
  });
</script>

<svelte:head><title>Google proxy (dev) — AegisPass</title></svelte:head>

<section class="page">
  <header class="head">
    <div>
      <p class="eyebrow">GOOGLE-002/003 · Simulação dev · DoD Fase 2</p>
      <h1>Google Workspace (stub)</h1>
      <DocHelpLink slug="journey-google-dev-stub" label="Como funciona sem VPS?" />
    </div>
    <a class="back" href="/work">← Trabalho</a>
  </header>

  <p class="lead">
    Substitui OAuth/Google real em desenvolvimento: Drive cifrado no browser e Sheets com
    mascaramento NIF/IBAN. Em produção, ver tasks <code>GOOGLE-*</code> + <code>INFRA-001</code>.
  </p>

  <div class="tabs">
    <button type="button" class:active={tab === "drive"} onclick={() => (tab = "drive")}>Drive</button>
    <button type="button" class:active={tab === "sheets"} onclick={() => (tab = "sheets")}>Sheets</button>
  </div>

  {#if locked}
    <p class="muted">🔒 Desbloqueia a Master Key para usar o stub.</p>
  {:else if tab === "drive"}
    <section class="panel">
      <h2>Upload simulado (ZK)</h2>
      <label>Nome<input bind:value={fileName} /></label>
      <label>Conteúdo<textarea bind:value={fileContent} rows="3"></textarea></label>
      <button type="button" class="btn primary" onclick={handleUpload}>Enviar para mock Drive</button>
    </section>
    <section class="panel">
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
              <button type="button" class="btn sm" onclick={() => handleOpen(f)}>Abrir no AegisPass</button>
            </li>
          {/each}
        </ul>
      {/if}
      {#if opened}<pre class="opened">{opened}</pre>{/if}
    </section>
  {:else}
    <section class="panel">
      <h2>Mascaramento (o que a Google veria)</h2>
      <label>Dados sensíveis<textarea bind:value={sheetRaw} rows="5" oninput={applyMask}></textarea></label>
      <h3>Vista Sheet (tokens)</h3>
      <pre class="masked">{maskedView}</pre>
      <h3>Vista AegisPass (desmascarado)</h3>
      <pre class="opened">{unmaskText(maskedView, tokenMap)}</pre>
    </section>
  {/if}
</section>

<style>
  .page { max-width: 48rem; margin: 0 auto; padding: var(--space-6); }
  .head { display: flex; justify-content: space-between; margin-bottom: var(--space-4); }
  .eyebrow { font-size: var(--text-xs); text-transform: uppercase; color: var(--color-text-label); margin: 0; }
  .lead { font-size: var(--text-sm); color: var(--color-text-muted); margin-bottom: var(--space-4); }
  .back { font-size: var(--text-sm); color: var(--color-text-muted); }
  .tabs { display: flex; gap: var(--space-2); margin-bottom: var(--space-4); }
  .tabs button { padding: var(--space-2) var(--space-3); border: 1px solid var(--color-border); border-radius: var(--radius-sm); background: var(--color-bg-inset); cursor: pointer; }
  .tabs button.active { border-color: var(--color-accent); color: var(--color-accent); }
  .panel { border: 1px solid var(--color-border); border-radius: var(--radius-md); padding: var(--space-4); margin-bottom: var(--space-4); }
  label { display: block; margin-bottom: var(--space-3); font-size: var(--text-sm); }
  input, textarea { width: 100%; margin-top: var(--space-1); padding: var(--space-2); border: 1px solid var(--color-border); border-radius: var(--radius-sm); box-sizing: border-box; }
  .btn { padding: var(--space-2) var(--space-4); border-radius: var(--radius-sm); border: 1px solid var(--color-border); cursor: pointer; }
  .btn.primary { background: var(--color-accent); color: var(--color-accent-fg); border-color: transparent; }
  .btn.sm { font-size: var(--text-xs); }
  .files { list-style: none; padding: 0; }
  .files li { display: flex; justify-content: space-between; align-items: center; padding: var(--space-2); border-bottom: 1px solid var(--color-border); }
  .opaque { display: block; font-size: var(--text-xs); color: var(--color-text-muted); }
  .mono { font-family: var(--font-mono); }
  pre { background: var(--color-bg-inset); padding: var(--space-3); border-radius: var(--radius-sm); font-size: var(--text-xs); overflow-x: auto; }
  .muted { color: var(--color-text-muted); font-size: var(--text-sm); }
</style>
