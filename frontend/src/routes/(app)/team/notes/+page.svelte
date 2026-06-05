<script lang="ts">
  import { loadSessionToken } from "$lib/session";
  import { buildNoteLink, encryptNoteContent, generateNoteKey } from "$lib/share/burnNote";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { burnNoteManually, createBurnNote } from "$lib/share/burnNoteApi";
  import {
    Button,
    confirmDialog,
    Field,
    PageShell,
    Panel,
    StatusBanner,
    toast,
  } from "$lib/ui";

  const TTL_OPTIONS = [
    { label: "10 minutos", seconds: 600 },
    { label: "1 hora", seconds: 3600 },
    { label: "24 horas", seconds: 86400 },
    { label: "7 dias", seconds: 604800 },
  ];

  let secret = $state("");
  let passphrase = $state("");
  let ttlSeconds = $state(86400);

  let creating = $state(false);
  let error = $state("");
  let link = $state("");
  let noteId = $state("");
  let burnToken = $state("");
  let expiresAt = $state("");
  let burned = $state(false);
  let burnError = $state("");

  async function create(event: SubmitEvent) {
    event.preventDefault();
    const value = secret;
    if (!value) return;
    const token = loadSessionToken();
    if (!token) {
      error = "Sessao expirada — inicia sessao de novo.";
      return;
    }
    creating = true;
    error = "";
    link = "";
    burned = false;
    burnError = "";
    try {
      const key = generateNoteKey();
      const { ciphertext, salt } = await encryptNoteContent(key, value, passphrase || undefined);
      const created = await createBurnNote(token, ciphertext, ttlSeconds);
      noteId = created.id;
      burnToken = created.burn_token;
      link = buildNoteLink(window.location.origin, created.id, key, salt);
      expiresAt = new Date(created.expires_at).toLocaleString("pt-PT");
      secret = "";
      passphrase = "";
      toast.success("Nota criada.");
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao criar nota";
    } finally {
      creating = false;
    }
  }

  async function copy() {
    try {
      await navigator.clipboard.writeText(link);
      toast.success("Link copiado.");
    } catch {
      toast.error("Não foi possível copiar.");
    }
  }

  async function burnNow() {
    const ok = await confirmDialog({
      title: "Destruir nota?",
      message: "O link deixa de funcionar de imediato — mesmo que ainda não tenha sido lida.",
      confirmLabel: "Destruir",
      variant: "danger",
    });
    if (!ok) return;
    burnError = "";
    try {
      await burnNoteManually(noteId, burnToken);
      burned = true;
      toast.info("Nota destruída.");
    } catch (e) {
      burned = true;
      burnError = e instanceof Error ? e.message : "";
    }
  }

  function reset() {
    link = "";
    noteId = "";
    burnToken = "";
    expiresAt = "";
    error = "";
    burned = false;
    burnError = "";
  }
</script>

<svelte:head>
  <title>Notas Auto-Destrutivas — AegisPass</title>
</svelte:head>

<PageShell
  title="Notas Auto-Destrutivas"
  taskId="SHARE-005"
  description="Uma nota cifrada que se lê uma única vez e arde a seguir. A chave vai no fragmento do link (nunca chega ao servidor). Passphrase opcional como 2.ª camada."
 
>
  {#snippet actions()}
    <DocHelpLink />
    <Button variant="ghost" size="sm" href="/team">← Identidade de partilha</Button>
  {/snippet}

  {#if link}
    <Panel title="Nota gerada">
      {#snippet actions()}
        <span class="pill" class:on={!burned} class:gone={burned}>{burned ? "Destruída" : "Ativa"}</span>
      {/snippet}

      {#if burned}
        <p class="panel-body">A nota foi destruída. O link já não abre nada.</p>
        <Button variant="secondary" onclick={reset}>Criar outra nota</Button>
      {:else}
        <div class="link-row">
          <input class="mono" type="text" readonly value={link} aria-label="Link da nota" />
          <Button onclick={copy}>Copiar</Button>
        </div>
        <dl class="props">
          <div class="prop">
            <dt>Expira</dt>
            <dd class="mono">{expiresAt}</dd>
          </div>
        </dl>
        <p class="panel-foot warn">
          ⚠️ Copia o link agora. A chave de cifra está no próprio link (depois do #)
          — o servidor não a tem. A nota arde após a primeira leitura.
        </p>
        <div class="actions">
          <Button variant="ghost" onclick={burnNow}>Destruir já</Button>
          <Button variant="secondary" onclick={reset}>Criar outra nota</Button>
        </div>
        {#if burnError}<StatusBanner variant="info">{burnError}</StatusBanner>{/if}
      {/if}
    </Panel>
  {:else}
    <Panel title="Nova nota">
      <form class="form" onsubmit={create}>
        <Field label="Nota" required>
          {#snippet control({ id, describedBy })}
            <textarea
              {id}
              aria-describedby={describedBy}
              bind:value={secret}
              rows="3"
              placeholder="Mensagem, password, token…"
              disabled={creating}
              required
            ></textarea>
          {/snippet}
        </Field>
        <div class="row">
          <Field label="Expira após">
            {#snippet control({ id, describedBy })}
              <select {id} aria-describedby={describedBy} bind:value={ttlSeconds} disabled={creating}>
                {#each TTL_OPTIONS as opt (opt.seconds)}
                  <option value={opt.seconds}>{opt.label}</option>
                {/each}
              </select>
            {/snippet}
          </Field>
          <Field label="Passphrase (opcional)" hint="2.ª camada — combina por um canal à parte.">
            {#snippet control({ id, describedBy })}
              <input
                {id}
                aria-describedby={describedBy}
                type="text"
                bind:value={passphrase}
                placeholder="Dita à parte"
                autocomplete="off"
                spellcheck="false"
                disabled={creating}
              />
            {/snippet}
          </Field>
        </div>
        {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}
        <Button type="submit" disabled={creating || !secret} loading={creating}>Gerar nota</Button>
      </form>
    </Panel>

    <Panel title="Como funciona">
      <ol class="steps">
        <li><span>1</span> A nota é cifrada no teu dispositivo com uma chave aleatória.</li>
        <li><span>2</span> Com passphrase, ciframos <em>outra vez</em> — nem o link sozinho a abre.</li>
        <li><span>3</span> O servidor guarda o <em>ciphertext</em> só em RAM, com TTL.</li>
        <li><span>4</span> Ao ler, a nota é revelada <em>uma vez</em> e apagada. Sem rasto.</li>
      </ol>
    </Panel>
  {/if}
</PageShell>

<style>
  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }
  textarea,
  select,
  input {
    width: 100%;
    padding: var(--space-2) var(--space-3);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    background: var(--color-bg-inset);
    color: var(--color-text);
    font-family: var(--font-ui);
    font-size: var(--text-sm);
    box-sizing: border-box;
  }
  textarea {
    resize: vertical;
  }
  textarea:focus-visible,
  select:focus-visible,
  input:focus-visible {
    outline: none;
    border-color: var(--color-accent);
  }
  .row {
    display: flex;
    gap: var(--space-3);
  }
  .row :global(.field) {
    flex: 1;
  }
  .link-row {
    display: flex;
    gap: var(--space-2);
    margin-bottom: var(--space-3);
  }
  .link-row input {
    flex: 1;
    min-width: 0;
    font-size: var(--text-xs);
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
    font-size: var(--text-xs);
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-label);
  }
  .prop dd {
    margin: 0;
  }
  .mono {
    font-family: var(--font-mono);
  }
  .panel-body {
    margin: 0 0 var(--space-3);
    font-size: var(--text-sm);
  }
  .panel-foot {
    margin: var(--space-4) 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text);
  }
  .pill {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
    color: var(--color-text-muted);
  }
  .pill.on {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
  }
  .pill.gone {
    color: var(--color-danger);
    border-color: var(--color-danger);
  }
  .actions {
    display: flex;
    gap: var(--space-2);
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
