<script lang="ts">
  import { onMount } from "svelte";
  import {
    ensureShareIdentity,
    loadShareIdentity,
    lookupRecipient,
    type RecipientLookup,
    type ShareIdentity,
  } from "$lib/share/setup";
  import { PageShell, Panel, Button, StatusBanner, HubLinks } from "$lib/ui";
  import type { HubLinkItem } from "$lib/ui";
  import { routeChildren } from "$lib/shell/routes";

  type Status = "loading" | "locked" | "ready" | "error";

  let status = $state<Status>("loading");
  let identity = $state<ShareIdentity | null>(null);
  let loadError = $state("");

  let activating = $state(false);
  let activateError = $state("");

  let email = $state("");
  let looking = $state(false);
  let lookupError = $state("");
  let recipient = $state<RecipientLookup | null>(null);

  const descriptions: Record<string, string> = {
    "/team/vaults": "Cofres por equipa com permissoes e revogacao imediata.",
    "/team/notes": "Notas que ardem apos a primeira leitura, passphrase opcional.",
    "/team/links": "Segredo de uso unico para quem nao tem conta — chave no fragmento.",
  };
  const hubItems: HubLinkItem[] = routeChildren("/team").map((c) => ({
    href: c.href,
    title: c.label,
    description: descriptions[c.href],
    taskId: c.taskId,
    comingSoon: c.comingSoon,
  }));

  function isLocked(message: string): boolean {
    return message.toLowerCase().includes("bloquead");
  }

  async function load() {
    status = "loading";
    loadError = "";
    try {
      identity = await loadShareIdentity();
      status = "ready";
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha ao carregar partilha";
      if (isLocked(msg)) {
        status = "locked";
      } else {
        loadError = msg;
        status = "error";
      }
    }
  }

  async function activate() {
    activating = true;
    activateError = "";
    try {
      identity = await ensureShareIdentity();
    } catch (e) {
      activateError = e instanceof Error ? e.message : "Falha ao activar partilha";
    } finally {
      activating = false;
    }
  }

  async function lookup(event: SubmitEvent) {
    event.preventDefault();
    const target = email.trim();
    if (!target) return;
    looking = true;
    lookupError = "";
    recipient = null;
    try {
      recipient = await lookupRecipient(target);
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Falha na procura";
      lookupError = msg.includes("não encontrado")
        ? "Sem chave pública para este email (o colega ainda não activou a partilha)."
        : msg;
    } finally {
      looking = false;
    }
  }

  onMount(load);
</script>

<svelte:head>
  <title>Equipa — AegisPass</title>
</svelte:head>

<PageShell
  title="Equipa"
  taskId="SHARE"
  description="Chaves assimétricas por utilizador — a base para partilhar segredos sem o servidor alguma vez os ver."
>
  {#if status === "loading"}
    <Panel><p class="muted">A carregar identidade de partilha…</p></Panel>
  {:else if status === "locked"}
    <Panel title="Cofre bloqueado">
      <p class="body">
        A chave de partilha é protegida pela tua Master Password. Desbloqueia o
        cofre para a gerir.
      </p>
      <Button href="/auth/unlock">Desbloquear</Button>
    </Panel>
  {:else if status === "error"}
    <StatusBanner variant="error">{loadError}</StatusBanner>
    <div class="retry"><Button variant="secondary" onclick={load}>Tentar de novo</Button></div>
  {:else}
    <!-- Painel A — a minha chave -->
    <Panel title="A minha chave">
      {#snippet actions()}
        <span class="pill" class:on={!!identity} class:off={!identity}>
          {identity ? "Activa" : "Inactiva"}
        </span>
      {/snippet}

      {#if identity}
        <dl class="props">
          <div class="prop">
            <dt>Algoritmo</dt>
            <dd class="mono">{identity.algorithm}</dd>
          </div>
          <div class="prop">
            <dt>Impressão digital</dt>
            <dd class="mono fingerprint">{identity.fingerprint}</dd>
          </div>
        </dl>
        <p class="foot">
          A chave privada está cifrada com a tua Master Password antes de chegar
          ao servidor. Partilha esta impressão digital por um canal de confiança
          para os colegas confirmarem que é mesmo a tua.
        </p>
      {:else}
        <p class="body">
          Ainda não activaste a partilha. Vamos gerar um par de chaves no teu
          dispositivo: a pública fica partilhável, a privada é cifrada com a tua
          Master Password.
        </p>
        {#if activateError}<StatusBanner variant="error">{activateError}</StatusBanner>{/if}
        <Button onclick={activate} loading={activating}>
          {activating ? "A gerar par de chaves…" : "Activar partilha"}
        </Button>
      {/if}
    </Panel>

    <!-- Painel B — verificar chave de colega -->
    <Panel title="Verificar colega">
      <p class="body">Procura a chave pública de um colega para lhe partilhares um segredo.</p>
      <form class="lookup" onsubmit={lookup}>
        <input
          type="email"
          bind:value={email}
          placeholder="colega@empresa.pt"
          autocomplete="off"
          spellcheck="false"
          aria-label="Email do colega"
          disabled={looking}
        />
        <Button variant="secondary" type="submit" loading={looking} disabled={!email.trim()}>
          {looking ? "A procurar…" : "Procurar"}
        </Button>
      </form>

      {#if lookupError}
        <p class="note" role="status">{lookupError}</p>
      {:else if recipient}
        <dl class="props">
          <div class="prop">
            <dt>Email</dt>
            <dd class="mono">{recipient.email}</dd>
          </div>
          <div class="prop">
            <dt>Impressão digital</dt>
            <dd class="mono fingerprint">{recipient.fingerprint}</dd>
          </div>
        </dl>
        <p class="foot">
          Confirma esta impressão digital com o colega por um canal à parte antes
          de partilhar — defende contra troca maliciosa de chaves.
        </p>
      {/if}
    </Panel>

    <!-- Hub — destinos de partilha (derivados da ROUTE_TREE) -->
    <h2 class="hub-title">Partilhar</h2>
    <HubLinks items={hubItems} />

    <!-- Como funciona -->
    <Panel title="Como funciona">
      <ol class="steps">
        <li><span>1</span> Cada utilizador gera um par de chaves no dispositivo.</li>
        <li><span>2</span> Partilhar = re-cifrar a chave do item para a <em>chave pública</em> do destinatário.</li>
        <li><span>3</span> Só a <em>chave privada</em> dele a abre. O servidor encaminha bytes opacos.</li>
      </ol>
    </Panel>
  {/if}
</PageShell>

<style>
  .body {
    margin: 0 0 var(--space-4);
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .foot {
    margin: var(--space-4) 0 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text-muted);
  }
  .muted {
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    margin: 0;
  }
  .retry {
    margin-top: var(--space-3);
  }
  .hub-title {
    margin: var(--space-6) 0 var(--space-3);
    font-size: var(--text-base);
    font-weight: 600;
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
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  .prop dd {
    margin: 0;
    text-align: right;
    word-break: break-word;
  }
  .mono {
    font-family: var(--font-mono);
    font-size: var(--text-sm);
  }
  .fingerprint {
    font-size: var(--text-xs);
    letter-spacing: 0.04em;
    color: var(--color-text);
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
  }
  .pill.on {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
    border-color: transparent;
  }
  .pill.off {
    color: var(--color-text-muted);
  }

  .lookup {
    display: flex;
    gap: var(--space-2);
  }
  .lookup input {
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
  .lookup input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }

  .note {
    margin: var(--space-4) 0 0;
    font-size: var(--text-sm);
    color: var(--color-text-muted);
  }

  .steps {
    margin: 0;
    padding: 0;
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }
  .steps li {
    display: flex;
    gap: var(--space-3);
    align-items: baseline;
    font-size: var(--text-sm);
    color: var(--color-text);
  }
  .steps span {
    flex-shrink: 0;
    width: 1.25rem;
    height: 1.25rem;
    display: inline-grid;
    place-items: center;
    border-radius: 50%;
    background: var(--color-accent-muted);
    color: var(--color-accent);
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-mono);
  }
  .steps em {
    font-style: normal;
    color: var(--color-accent);
  }
</style>
