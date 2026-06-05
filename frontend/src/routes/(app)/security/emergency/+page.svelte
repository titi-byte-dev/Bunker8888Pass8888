<script lang="ts">
  import { onMount } from "svelte";
  import { loadUserEmail } from "$lib/session";
  import { getMasterKey } from "$lib/vault/masterKeyStore";
  import { deriveMasterKeyBytes, fetchKdfParams, importAesKeyFromBytes, bytesToBase64, base64ToBytes } from "$lib/auth";
  import {
    generateRecoveryCode,
    wrapMasterKeyBytes,
    unwrapMasterKeyBytes,
  } from "$lib/vault/recovery";
  import {
    approveEmergencyRequest,
    createEmergencyRequest,
    deleteEmergencyConfig,
    fetchEmergencyAccessBlob,
    fetchEmergencyConfig,
    fetchEmergencyRequestStatus,
    listEmergencyRequests,
    rejectEmergencyRequest,
    saveEmergencyConfig,
    secondsUntil,
    type EmergencyRequest,
  } from "$lib/emergency/api";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";

  let heirEmail = $state("");
  let waitDays = $state(7);
  let confirmPassword = $state("");
  let configLoaded = $state(false);
  let hasConfig = $state(false);
  let hasBlob = $state(false);
  let newEmergencyCode = $state("");
  let ownerBusy = $state(false);
  let ownerError = $state("");
  let ownerStatus = $state("");

  let requests = $state<EmergencyRequest[]>([]);

  let ownerEmailForHeir = $state("");
  let heirRequest = $state<EmergencyRequest | null>(null);
  let heirActive = $state(false);
  let heirCode = $state("");
  let heirBusy = $state(false);
  let heirError = $state("");
  let heirStatus = $state("");
  let countdown = $state(0);

  const myEmail = $derived(loadUserEmail() ?? "");

  async function loadOwner() {
    ownerError = "";
    try {
      const cfg = await fetchEmergencyConfig();
      configLoaded = true;
      hasConfig = cfg.configured;
      if (cfg.configured) {
        heirEmail = cfg.heir_email ?? "";
        waitDays = cfg.wait_days ?? 7;
        hasBlob = cfg.has_blob ?? false;
      }
      requests = await listEmergencyRequests();
    } catch (e) {
      ownerError = e instanceof Error ? e.message : "Erro ao carregar";
    }
  }

  async function saveOwnerConfig() {
    ownerBusy = true;
    ownerError = "";
    ownerStatus = "";
    newEmergencyCode = "";
    try {
      const key = getMasterKey();
      if (!key || !confirmPassword || !myEmail) {
        throw new Error("Confirma a Master Password com o cofre desbloqueado.");
      }
      const kdf = await fetchKdfParams("", myEmail);
      const mk = await deriveMasterKeyBytes(confirmPassword, kdf.salt, kdf);
      const code = generateRecoveryCode();
      const envelope = await wrapMasterKeyBytes(mk, code);
      await saveEmergencyConfig({
        heir_email: heirEmail,
        wait_days: waitDays,
        blob: bytesToBase64(new TextEncoder().encode(envelope)),
      });
      newEmergencyCode = code;
      confirmPassword = "";
      ownerStatus = "Herdeiro configurado — partilha o código offline com o herdeiro.";
      await loadOwner();
    } catch (e) {
      ownerError = e instanceof Error ? e.message : "Falha ao guardar";
    } finally {
      ownerBusy = false;
    }
  }

  async function removeConfig() {
    if (!confirm("Remover herdeiro e blob de emergência?")) return;
    ownerBusy = true;
    try {
      await deleteEmergencyConfig();
      hasConfig = false;
      hasBlob = false;
      heirEmail = "";
      ownerStatus = "Configuração removida.";
      await loadOwner();
    } catch (e) {
      ownerError = e instanceof Error ? e.message : "Falha ao remover";
    } finally {
      ownerBusy = false;
    }
  }

  async function refreshHeir() {
    heirError = "";
    if (!ownerEmailForHeir.trim()) {
      heirActive = false;
      heirRequest = null;
      return;
    }
    try {
      const st = await fetchEmergencyRequestStatus(ownerEmailForHeir);
      heirActive = st.active;
      heirRequest = st.request ?? null;
      if (heirRequest?.status === "waiting") {
        countdown = secondsUntil(heirRequest.unlocks_at);
      }
    } catch (e) {
      heirError = e instanceof Error ? e.message : "Erro";
    }
  }

  async function requestAccess() {
    heirBusy = true;
    heirError = "";
    heirStatus = "";
    try {
      heirRequest = await createEmergencyRequest(ownerEmailForHeir);
      heirActive = true;
      countdown = secondsUntil(heirRequest.unlocks_at);
      heirStatus = "Pedido iniciado — período de espera activo.";
      await loadOwner();
    } catch (e) {
      heirError = e instanceof Error ? e.message : "Falha no pedido";
    } finally {
      heirBusy = false;
    }
  }

  async function downloadAccess() {
    heirBusy = true;
    heirError = "";
    try {
      const blobB64 = await fetchEmergencyAccessBlob(ownerEmailForHeir);
      const envelopeJson = new TextDecoder().decode(base64ToBytes(blobB64));
      const mkBytes = await unwrapMasterKeyBytes(envelopeJson, heirCode);
      await importAesKeyFromBytes(mkBytes);
      heirStatus =
        "Master Key recuperada localmente. Inicia sessão como titular se necessário — uso único consumido.";
      await refreshHeir();
    } catch (e) {
      heirError = e instanceof Error ? e.message : "Falha ao obter acesso";
    } finally {
      heirBusy = false;
    }
  }

  function formatCountdown(secs: number): string {
    const d = Math.floor(secs / 86400);
    const h = Math.floor((secs % 86400) / 3600);
    const m = Math.floor((secs % 3600) / 60);
    const s = secs % 60;
    if (d > 0) return `${d}d ${h}h ${m}m`;
    if (h > 0) return `${h}h ${m}m ${s}s`;
    return `${m}m ${s}s`;
  }

  onMount(() => {
    loadOwner();
    const id = setInterval(() => {
      if (heirRequest?.status === "waiting") {
        countdown = secondsUntil(heirRequest.unlocks_at);
        if (countdown === 0) refreshHeir();
      }
    }, 1000);
    return () => clearInterval(id);
  });
</script>

<svelte:head>
  <title>Acesso de emergência — AegisPass</title>
</svelte:head>

<section class="page">
  <header>
    <a href="/security" class="back">← Segurança</a>
    <h1>Acesso de emergência</h1>
    <DocHelpLink />
    <p class="lead">
      Designa um herdeiro digital com período de espera. O servidor orquestra pedidos;
      a Master Key viaja só num blob cifrado (Zero-Knowledge).
    </p>
  </header>

  <article class="card">
    <h2>Como titular</h2>
    {#if ownerError}<p class="error" role="alert">{ownerError}</p>{/if}
    {#if ownerStatus}<p class="ok" role="status">{ownerStatus}</p>{/if}

    <label>Email do herdeiro <input type="email" bind:value={heirEmail} disabled={ownerBusy} /></label>
    <label>
      Período de espera (dias)
      <input type="number" min="1" max="90" bind:value={waitDays} disabled={ownerBusy} />
    </label>
    <label>
      Confirma Master Password (para cifrar blob de emergência)
      <input type="password" bind:value={confirmPassword} disabled={ownerBusy} autocomplete="current-password" />
    </label>

    {#if newEmergencyCode}
      <p class="warn">Código para o herdeiro (offline, uma vez):</p>
      <code class="code">{newEmergencyCode}</code>
    {/if}

    <div class="row">
      <button type="button" onclick={saveOwnerConfig} disabled={ownerBusy || !heirEmail}>Guardar herdeiro</button>
      {#if hasConfig}
        <button type="button" class="secondary" onclick={removeConfig} disabled={ownerBusy}>Remover</button>
      {/if}
    </div>

    {#if configLoaded && hasConfig}
      <p class="muted">Blob de emergência: {hasBlob ? "configurado" : "pendente"}</p>
    {/if}

    <h3>Pedidos recebidos</h3>
    {#if requests.length === 0}
      <p class="muted">Nenhum pedido.</p>
    {:else}
      <ul class="req-list">
        {#each requests as req (req.id)}
          <li>
            <span>{req.heir_email} — <strong>{req.status}</strong></span>
            {#if req.status === "waiting"}
              <span class="muted">desbloqueia {new Date(req.unlocks_at).toLocaleString("pt-PT")}</span>
              <div class="row">
                <button type="button" class="secondary" onclick={() => approveEmergencyRequest(req.id).then(loadOwner)}>
                  Aprovar já
                </button>
                <button type="button" class="secondary" onclick={() => rejectEmergencyRequest(req.id).then(loadOwner)}>
                  Rejeitar
                </button>
              </div>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </article>

  <article class="card">
    <h2>Como herdeiro</h2>
    <p class="muted">Conta autenticada: {myEmail || "—"}</p>
    {#if heirError}<p class="error" role="alert">{heirError}</p>{/if}
    {#if heirStatus}<p class="ok" role="status">{heirStatus}</p>{/if}

    <label>
      Email do titular
      <input
        type="email"
        bind:value={ownerEmailForHeir}
        onchange={refreshHeir}
        disabled={heirBusy}
      />
    </label>

    <div class="row">
      <button type="button" onclick={requestAccess} disabled={heirBusy || !ownerEmailForHeir}>
        Pedir acesso
      </button>
      <button type="button" class="secondary" onclick={refreshHeir} disabled={heirBusy}>Actualizar</button>
    </div>

    {#if heirActive && heirRequest}
      <p>
        Estado: <strong>{heirRequest.status}</strong>
        {#if heirRequest.status === "waiting"}
          — faltam {formatCountdown(countdown)}
        {/if}
      </p>
    {/if}

    {#if heirRequest?.status === "ready"}
      <label>
        Código de emergência (offline)
        <input type="text" bind:value={heirCode} placeholder="XXXXX-XXXXX-XXXXX-XXXXX" disabled={heirBusy} />
      </label>
      <button type="button" onclick={downloadAccess} disabled={heirBusy || !heirCode}>
        Obter Master Key cifrada
      </button>
    {/if}
  </article>
</section>

<style>
  .page {
    max-width: 36rem;
  }

  .back {
    color: var(--color-link);
    text-decoration: none;
    font-size: var(--text-sm);
  }

  h1 {
    margin: var(--space-2) 0;
    font-family: var(--font-display);
    font-size: var(--text-2xl);
  }

  .lead {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    line-height: var(--leading-body);
  }

  .card {
    margin-top: var(--space-6);
    padding: var(--space-4);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    background: var(--color-bg-surface);
  }

  h2 {
    margin: 0 0 var(--space-4);
    font-size: var(--text-lg);
  }

  h3 {
    margin: var(--space-6) 0 var(--space-3);
    font-size: var(--text-base);
  }

  label {
    display: block;
    margin-bottom: var(--space-3);
    font-size: var(--text-sm);
    color: var(--color-text-label);
  }

  input {
    display: block;
    width: 100%;
    margin-top: var(--space-1);
    padding: var(--space-2);
    border: 1px solid var(--color-border-strong);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    box-sizing: border-box;
  }

  .row {
    display: flex;
    gap: var(--space-2);
    flex-wrap: wrap;
    margin-top: var(--space-3);
  }

  button {
    padding: var(--space-2) var(--space-4);
    border: none;
    border-radius: var(--radius-sm);
    background: var(--color-accent);
    color: var(--color-accent-fg);
    cursor: pointer;
    font-size: var(--text-sm);
  }

  button.secondary {
    background: var(--color-accent-muted);
    color: var(--color-text);
  }

  .error {
    color: var(--color-danger);
    font-size: var(--text-sm);
  }

  .ok {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    padding: var(--space-2);
    border-radius: var(--radius-sm);
    font-size: var(--text-sm);
  }

  .warn {
    color: var(--color-warning);
    font-size: var(--text-sm);
  }

  .code {
    display: block;
    font-family: var(--font-mono);
    padding: var(--space-3);
    margin: var(--space-2) 0;
    background: var(--color-bg-inset);
    border-radius: var(--radius-sm);
  }

  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .req-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  .req-list li {
    padding: var(--space-3) 0;
    border-bottom: 1px solid var(--color-border);
  }
</style>
