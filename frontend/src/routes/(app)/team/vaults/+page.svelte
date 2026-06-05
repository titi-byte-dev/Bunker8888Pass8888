<script lang="ts">
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import {
    Button,
    confirmDialog,
    EmptyState,
    PageShell,
    Panel,
    Skeleton,
    StatusBanner,
    toast,
  } from "$lib/ui";
  import { onMount } from "svelte";
  import {
    addVaultItem,
    createSharedVault,
    deleteSharedVault,
    downloadAttachment,
    inviteMember,
    listSharedVaults,
    openSharedVault,
    removeAttachment,
    removeVaultItem,
    revokeMember,
    uploadAttachment,
    type DecryptedAttachment,
    type DecryptedVault,
    type InvitedMember,
    type OpenVault,
  } from "$lib/share/vaultsService";
  import type { VaultRole } from "$lib/share/vaultsApi";

  type Status = "loading" | "locked" | "ready" | "error";

  let status = $state<Status>("loading");
  let loadError = $state("");
  let vaults = $state<DecryptedVault[]>([]);

  // Criar cofre
  let newName = $state("");
  let creating = $state(false);
  let createError = $state("");

  // Cofre aberto
  let open = $state<OpenVault | null>(null);
  let opening = $state(false);
  let openError = $state("");

  // Convidar membro
  let inviteEmail = $state("");
  let inviteRole = $state<VaultRole>("member");
  let inviting = $state(false);
  let inviteError = $state("");
  let lastInvited = $state<InvitedMember | null>(null);

  // Adicionar item
  let itemTitle = $state("");
  let itemSecret = $state("");
  let addingItem = $state(false);
  let itemError = $state("");
  let revealed = $state<Set<string>>(new Set());

  // Anexos cifrados
  let attachFile = $state<File | null>(null);
  let uploading = $state(false);
  let attachError = $state("");
  let fileInput = $state<HTMLInputElement | null>(null);

  const canManage = $derived(open?.vault.role === "owner" || open?.vault.role === "admin");
  const canWrite = $derived(
    open?.vault.role === "owner" || open?.vault.role === "admin" || open?.vault.role === "member",
  );
  const isOwner = $derived(open?.vault.role === "owner");

  function isLocked(message: string): boolean {
    return message.toLowerCase().includes("bloquead");
  }

  async function load() {
    status = "loading";
    loadError = "";
    try {
      vaults = await listSharedVaults();
      status = "ready";
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha ao carregar cofres";
      if (isLocked(msg)) status = "locked";
      else {
        loadError = msg;
        status = "error";
      }
    }
  }

  async function create(event: SubmitEvent) {
    event.preventDefault();
    const name = newName.trim();
    if (!name) return;
    creating = true;
    createError = "";
    try {
      const v = await createSharedVault(name);
      vaults = [v, ...vaults];
      newName = "";
      await openVault(v.id);
    } catch (e) {
      createError = e instanceof Error ? e.message : "Falha ao criar cofre";
    } finally {
      creating = false;
    }
  }

  async function openVault(id: string) {
    opening = true;
    openError = "";
    lastInvited = null;
    revealed = new Set();
    try {
      open = await openSharedVault(id);
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao abrir cofre";
    } finally {
      opening = false;
    }
  }

  function closeVault() {
    open = null;
    inviteEmail = "";
    itemTitle = "";
    itemSecret = "";
  }

  async function invite(event: SubmitEvent) {
    event.preventDefault();
    if (!open) return;
    const email = inviteEmail.trim();
    if (!email) return;
    inviting = true;
    inviteError = "";
    lastInvited = null;
    try {
      lastInvited = await inviteMember(open.vault.id, open.vaultKey, email, inviteRole);
      inviteEmail = "";
      open = await openSharedVault(open.vault.id);
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha ao convidar";
      inviteError = msg.includes("não encontrado")
        ? "Sem chave pública para este email (o colega ainda não activou a partilha)."
        : msg;
    } finally {
      inviting = false;
    }
  }

  async function removeMember(userID: string) {
    if (!open) return;
    const ok = await confirmDialog({
      title: "Remover membro?",
      message: "Revoga o acesso ao cofre — a cópia da chave dele é apagada.",
      confirmLabel: "Remover",
      variant: "danger",
    });
    if (!ok) return;
    try {
      await revokeMember(open.vault.id, userID);
      open = await openSharedVault(open.vault.id);
      toast.success("Membro removido.");
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao remover membro";
    }
  }

  async function addItem(event: SubmitEvent) {
    event.preventDefault();
    if (!open) return;
    const title = itemTitle.trim();
    const secret = itemSecret;
    if (!title || !secret) return;
    addingItem = true;
    itemError = "";
    try {
      await addVaultItem(open.vault.id, open.vaultKey, { title, secret });
      itemTitle = "";
      itemSecret = "";
      open = await openSharedVault(open.vault.id);
    } catch (e) {
      itemError = e instanceof Error ? e.message : "Falha ao adicionar item";
    } finally {
      addingItem = false;
    }
  }

  async function deleteItem(itemID: string) {
    if (!open) return;
    const ok = await confirmDialog({
      title: "Apagar item?",
      message: "Remove o segredo cifrado deste cofre partilhado.",
      confirmLabel: "Apagar",
      variant: "danger",
    });
    if (!ok) return;
    try {
      await removeVaultItem(open.vault.id, itemID);
      open = await openSharedVault(open.vault.id);
      toast.success("Item apagado.");
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao remover item";
    }
  }

  async function destroyVault() {
    if (!open) return;
    const ok = await confirmDialog({
      title: "Apagar cofre partilhado?",
      message: `Remove «${open.vault.name}» para todos os membros. Irreversível.`,
      confirmLabel: "Apagar cofre",
      variant: "danger",
    });
    if (!ok) return;
    try {
      const id = open.vault.id;
      await deleteSharedVault(id);
      vaults = vaults.filter((v) => v.id !== id);
      closeVault();
      toast.success("Cofre apagado.");
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao apagar cofre";
    }
  }

  function onFilePick(event: Event) {
    const input = event.currentTarget as HTMLInputElement;
    attachFile = input.files && input.files.length > 0 ? input.files[0] : null;
    attachError = "";
  }

  async function upload(event: SubmitEvent) {
    event.preventDefault();
    if (!open || !attachFile) return;
    uploading = true;
    attachError = "";
    try {
      await uploadAttachment(open.vault.id, open.vaultKey, attachFile);
      attachFile = null;
      if (fileInput) fileInput.value = "";
      open = await openSharedVault(open.vault.id);
    } catch (e) {
      attachError = e instanceof Error ? e.message : "Falha ao carregar anexo";
    } finally {
      uploading = false;
    }
  }

  async function download(att: DecryptedAttachment) {
    if (!open) return;
    try {
      const file = await downloadAttachment(open.vault.id, open.vaultKey, att.id);
      const blob = new Blob([file.bytes], { type: file.mime });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = file.name;
      document.body.appendChild(a);
      a.click();
      a.remove();
      URL.revokeObjectURL(url);
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao descarregar anexo";
    }
  }

  async function deleteAttachment(attID: string) {
    if (!open) return;
    const ok = await confirmDialog({
      title: "Apagar anexo?",
      message: "Remove o ficheiro cifrado deste cofre.",
      confirmLabel: "Apagar",
      variant: "danger",
    });
    if (!ok) return;
    try {
      await removeAttachment(open.vault.id, attID);
      open = await openSharedVault(open.vault.id);
      toast.success("Anexo apagado.");
    } catch (e) {
      openError = e instanceof Error ? e.message : "Falha ao remover anexo";
    }
  }

  function formatBytes(n: number): string {
    if (n < 1024) return `${n} B`;
    if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KiB`;
    return `${(n / (1024 * 1024)).toFixed(1)} MiB`;
  }

  function toggleReveal(id: string) {
    const next = new Set(revealed);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    revealed = next;
  }

  function roleLabel(role: VaultRole): string {
    return { owner: "Dono", admin: "Admin", member: "Membro", viewer: "Leitor" }[role];
  }

  onMount(load);
</script>

<svelte:head>
  <title>Cofres Partilhados — AegisPass</title>
</svelte:head>

<PageShell
  title="Cofres Partilhados"
  taskId="SHARE-002"
  description="Coleções cifradas sob uma chave de cofre própria. Convidar um colega = re-cifrar essa chave para a chave pública dele. O servidor nunca vê o conteúdo."
 
>
  {#snippet actions()}
    <DocHelpLink />
    <Button variant="ghost" size="sm" href="/team">← Identidade de partilha</Button>
  {/snippet}

  {#if status === "loading"}
    <Skeleton variant="row" />
    <Skeleton variant="row" />
  {:else if status === "locked"}
    <EmptyState
      title="Cofre bloqueado"
      description="Os cofres partilhados são protegidos pela tua Master Password. Desbloqueia para os gerir."
    >
      {#snippet action()}
        <Button href="/auth/unlock">Desbloquear</Button>
      {/snippet}
    </EmptyState>
  {:else if status === "error"}
    <StatusBanner variant="error">{loadError}</StatusBanner>
    <Button variant="secondary" onclick={load}>Tentar de novo</Button>
  {:else}
    <!-- Criar -->
    <section class="panel">
      <div class="panel-head"><p class="eyebrow">Novo cofre</p></div>
      <form class="lookup" onsubmit={create}>
        <input
          type="text"
          bind:value={newName}
          placeholder="Ex.: Credenciais de Infraestrutura"
          autocomplete="off"
          disabled={creating}
        />
        <button type="submit" class="btn primary" disabled={creating || !newName.trim()}>
          {creating ? "A criar…" : "Criar cofre"}
        </button>
      </form>
      {#if createError}<p class="inline-error" role="alert">{createError}</p>{/if}
    </section>

    <!-- Lista -->
    <section class="panel">
      <div class="panel-head">
        <p class="eyebrow">Os meus cofres</p>
        <span class="pill off">{vaults.length}</span>
      </div>
      {#if vaults.length === 0}
        <p class="panel-body">Ainda não tens cofres partilhados. Cria o primeiro acima.</p>
      {:else}
        <ul class="vault-list">
          {#each vaults as v (v.id)}
            <li>
              <button
                type="button"
                class="vault-row"
                class:active={open?.vault.id === v.id}
                onclick={() => openVault(v.id)}
              >
                <span class="vault-name">{v.name}</span>
                <span class="pill role-{v.role}">{roleLabel(v.role)}</span>
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    {#if opening}
      <div class="panel muted-panel">A abrir cofre…</div>
    {:else if open}
      <!-- Detalhe do cofre aberto -->
      <section class="panel detail">
        <div class="panel-head">
          <p class="eyebrow">Cofre · {open.vault.name}</p>
          <button type="button" class="link-btn" onclick={closeVault}>Fechar</button>
        </div>
        {#if openError}<StatusBanner variant="error">{openError}</StatusBanner>{/if}

        <!-- Membros -->
        <div class="subhead">Membros · permissões</div>
        <dl class="props">
          {#each open.members as m (m.user_id)}
            <div class="prop member">
              <dt class="mono">{m.email}</dt>
              <dd>
                <span class="pill role-{m.role}">{roleLabel(m.role)}</span>
                {#if canManage && m.role !== "owner"}
                  <button type="button" class="link-btn danger-link" onclick={() => removeMember(m.user_id)}>
                    Remover
                  </button>
                {/if}
              </dd>
            </div>
          {/each}
        </dl>

        {#if canManage}
          <form class="invite" onsubmit={invite}>
            <input
              type="email"
              bind:value={inviteEmail}
              placeholder="colega@empresa.pt"
              autocomplete="off"
              spellcheck="false"
              disabled={inviting}
            />
            <select bind:value={inviteRole} disabled={inviting}>
              <option value="admin">Admin</option>
              <option value="member">Membro</option>
              <option value="viewer">Leitor</option>
            </select>
            <button type="submit" class="btn secondary" disabled={inviting || !inviteEmail.trim()}>
              {inviting ? "A convidar…" : "Convidar"}
            </button>
          </form>
          {#if inviteError}<p class="inline-note" role="status">{inviteError}</p>{/if}
          {#if lastInvited}
            <p class="panel-foot">
              Convite enviado a <span class="mono">{lastInvited.email}</span>. Confirma a
              impressão digital por um canal à parte (anti-MITM):
              <br /><span class="mono fingerprint">{lastInvited.fingerprint}</span>
            </p>
          {/if}
        {/if}

        <!-- Itens -->
        <div class="subhead">Segredos do cofre</div>
        {#if open.items.length === 0}
          <p class="panel-body">Sem segredos neste cofre ainda.</p>
        {:else}
          <ul class="item-list">
            {#each open.items as it (it.id)}
              <li class="item">
                <div class="item-main">
                  <span class="item-title">{it.title}</span>
                  <code class="item-secret">{revealed.has(it.id) ? it.secret : "••••••••••"}</code>
                </div>
                <div class="item-actions">
                  <button type="button" class="link-btn" onclick={() => toggleReveal(it.id)}>
                    {revealed.has(it.id) ? "Ocultar" : "Revelar"}
                  </button>
                  {#if canWrite}
                    <button type="button" class="link-btn danger-link" onclick={() => deleteItem(it.id)}>
                      Apagar
                    </button>
                  {/if}
                </div>
              </li>
            {/each}
          </ul>
        {/if}

        {#if canWrite}
          <form class="add-item" onsubmit={addItem}>
            <input type="text" bind:value={itemTitle} placeholder="Título" disabled={addingItem} />
            <input type="text" bind:value={itemSecret} placeholder="Segredo" disabled={addingItem} />
            <button type="submit" class="btn secondary" disabled={addingItem || !itemTitle.trim() || !itemSecret}>
              {addingItem ? "A guardar…" : "Adicionar"}
            </button>
          </form>
          {#if itemError}<p class="inline-error" role="alert">{itemError}</p>{/if}
        {:else}
          <p class="panel-foot">Tens acesso de leitor — podes ver, mas não escrever.</p>
        {/if}

        <!-- Anexos cifrados -->
        <div class="subhead">Anexos cifrados</div>
        {#if open.attachments.length === 0}
          <p class="panel-body">Sem anexos neste cofre ainda.</p>
        {:else}
          <ul class="item-list">
            {#each open.attachments as att (att.id)}
              <li class="item">
                <div class="item-main">
                  <span class="item-title">{att.name}</span>
                  <span class="item-secret">{att.mime} · {formatBytes(att.size)}</span>
                </div>
                <div class="item-actions">
                  <button type="button" class="link-btn" onclick={() => download(att)}>Descarregar</button>
                  {#if canWrite}
                    <button type="button" class="link-btn danger-link" onclick={() => deleteAttachment(att.id)}>
                      Apagar
                    </button>
                  {/if}
                </div>
              </li>
            {/each}
          </ul>
        {/if}

        {#if canWrite}
          <form class="add-item" onsubmit={upload}>
            <input
              type="file"
              bind:this={fileInput}
              onchange={onFilePick}
              disabled={uploading}
            />
            <button type="submit" class="btn secondary" disabled={uploading || !attachFile}>
              {uploading ? "A cifrar…" : "Carregar anexo"}
            </button>
          </form>
          {#if attachError}<p class="inline-error" role="alert">{attachError}</p>{/if}
          <p class="panel-foot">
            Cada ficheiro é cifrado no teu dispositivo com a chave do cofre (máx. 5 MiB).
            O servidor guarda só bytes opacos — nunca vê o nome nem o conteúdo.
          </p>
        {/if}

        {#if isOwner}
          <div class="danger-zone">
            <button type="button" class="btn danger-btn" onclick={destroyVault}>Apagar cofre</button>
          </div>
        {/if}
      </section>
    {/if}
  {/if}
</PageShell>

<style>
  .eyebrow {
    margin: 0 0 var(--space-1);
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--color-text-muted);
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
    margin: 0 0 var(--space-2);
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .panel-title {
    margin: 0 0 var(--space-2);
    font-weight: 600;
    font-size: var(--text-sm);
  }
  .panel-foot {
    margin: var(--space-3) 0 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text-muted);
  }
  .muted-panel {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }
  .subhead {
    margin: var(--space-5) 0 var(--space-2);
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }

  .props {
    margin: 0;
    border-top: 1px solid var(--color-border);
  }
  .prop {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-3) 0;
    border-bottom: 1px solid var(--color-border);
  }
  .prop dt {
    flex-shrink: 0;
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .prop dd {
    margin: 0;
    display: flex;
    align-items: center;
    gap: var(--space-3);
  }

  .mono {
    font-family: var(--font-mono);
    font-size: var(--text-sm);
  }
  .fingerprint {
    font-size: var(--text-xs);
    letter-spacing: 0.04em;
    line-height: var(--leading-body);
  }

  .pill {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid var(--color-border);
    color: var(--color-text-muted);
  }
  .pill.off {
    font-family: var(--font-mono);
  }
  .pill.role-owner {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    border-color: transparent;
  }
  .pill.role-admin {
    color: var(--color-accent);
    border-color: var(--color-accent);
  }
  .pill.role-member {
    color: var(--color-text);
  }
  .pill.role-viewer {
    color: var(--color-text-muted);
  }

  .lookup,
  .invite,
  .add-item {
    display: flex;
    gap: var(--space-2);
    margin-top: var(--space-2);
  }
  input,
  select {
    flex: 1;
    min-width: 0;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
  }
  select {
    flex: 0 0 auto;
  }
  input:focus-visible,
  select:focus-visible {
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
    text-decoration: none;
    cursor: pointer;
    white-space: nowrap;
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
    color: var(--color-text);
  }
  .btn.secondary:hover:not(:disabled) {
    background: var(--color-accent-muted);
  }
  .btn:disabled {
    opacity: 0.55;
    cursor: progress;
  }

  .link-btn {
    background: none;
    border: none;
    padding: 0;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    cursor: pointer;
  }
  .link-btn:hover {
    color: var(--color-text);
  }
  .danger-link:hover {
    color: var(--color-danger);
  }

  .vault-list,
  .item-list {
    margin: 0;
    padding: 0;
    list-style: none;
  }
  .vault-row {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-3) 0;
    border: none;
    border-bottom: 1px solid var(--color-border);
    background: none;
    color: var(--color-text);
    font-size: var(--text-sm);
    text-align: left;
    cursor: pointer;
  }
  .vault-row:hover,
  .vault-row.active {
    color: var(--color-accent);
  }
  .vault-name {
    font-weight: 500;
  }

  .item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: var(--space-3) 0;
    border-bottom: 1px solid var(--color-border);
  }
  .item-main {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }
  .item-title {
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .item-secret {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    word-break: break-all;
  }
  .item-actions {
    display: flex;
    gap: var(--space-3);
    flex-shrink: 0;
  }

  .inline-error {
    margin: var(--space-2) 0 0;
    font-size: var(--text-sm);
    color: var(--color-danger);
  }
  .inline-note {
    margin: var(--space-2) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .danger-zone {
    margin-top: var(--space-5);
    padding-top: var(--space-4);
    border-top: 1px solid var(--color-border);
  }
  .danger-btn {
    color: var(--color-danger);
    border-color: var(--color-danger);
    background: none;
  }
  .danger-btn:hover {
    background: var(--color-danger);
    color: var(--color-accent-fg);
  }
</style>
