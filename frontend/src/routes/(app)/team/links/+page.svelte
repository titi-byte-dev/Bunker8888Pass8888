<script lang="ts">
  import { loadSessionToken } from "$lib/session";
  import {
    buildSecretLink,
    encryptSecret,
    generateLinkKey,
  } from "$lib/share/secretLink";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { createSecretLink } from "$lib/share/secretLinkApi";
  import {
    Button,
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
  ];
  const VIEW_OPTIONS = [1, 3, 5];

  let secret = $state("");
  let ttlSeconds = $state(3600);
  let maxViews = $state(1);

  let creating = $state(false);
  let error = $state("");
  let link = $state("");
  let expiresAt = $state("");

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
    try {
      const key = generateLinkKey();
      const ciphertext = await encryptSecret(key, value);
      const created = await createSecretLink(token, ciphertext, ttlSeconds, maxViews);
      link = buildSecretLink(window.location.origin, created.id, key);
      expiresAt = new Date(created.expires_at).toLocaleString("pt-PT");
      secret = "";
      toast.success("Link gerado.");
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao criar link";
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

  function reset() {
    link = "";
    expiresAt = "";
    error = "";
  }
</script>

<svelte:head>
  <title>Secret Links — AegisPass</title>
</svelte:head>

<PageShell
  title="Secret Links"
  taskId="SHARE-003"
  description="Partilha um segredo de uso único com quem não tem conta. O segredo é cifrado no teu dispositivo; o servidor guarda-o só em RAM. A chave vai no fragmento do link — nunca chega ao servidor."
 
>
  {#snippet actions()}
    <DocHelpLink />
    <Button variant="ghost" size="sm" href="/team">← Identidade de partilha</Button>
  {/snippet}

  {#if link}
    <Panel title="Link gerado">
      {#snippet actions()}
        <span class="pill on">Pronto</span>
      {/snippet}
      <div class="link-row">
        <input class="mono" type="text" readonly value={link} aria-label="Link secreto" />
        <Button onclick={copy}>Copiar</Button>
      </div>
      <dl class="props">
        <div class="prop">
          <dt>Expira</dt>
          <dd class="mono">{expiresAt}</dd>
        </div>
        <div class="prop">
          <dt>Visualizações</dt>
          <dd class="mono">{maxViews}</dd>
        </div>
      </dl>
      <p class="panel-foot warn">
        ⚠️ Copia o link agora. A chave de cifra está no próprio link (depois do #) —
        o servidor não a tem e não a consegue recuperar. Se perderes o link, o
        segredo é irrecuperável.
      </p>
      <Button variant="secondary" onclick={reset}>Criar outro link</Button>
    </Panel>
  {:else}
    <Panel title="Novo link">
      <form class="form" onsubmit={create}>
        <Field label="Segredo" required>
          {#snippet control({ id, describedBy })}
            <textarea
              {id}
              aria-describedby={describedBy}
              bind:value={secret}
              rows="3"
              placeholder="Password, token, nota…"
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
          <Field label="Visualizações">
            {#snippet control({ id, describedBy })}
              <select {id} aria-describedby={describedBy} bind:value={maxViews} disabled={creating}>
                {#each VIEW_OPTIONS as v (v)}
                  <option value={v}>{v}{v === 1 ? " (uso único)" : ""}</option>
                {/each}
              </select>
            {/snippet}
          </Field>
        </div>
        {#if error}<StatusBanner variant="error">{error}</StatusBanner>{/if}
        <Button type="submit" disabled={creating || !secret} loading={creating}>Gerar link</Button>
      </form>
    </Panel>

    <Panel title="Como funciona">
      <ol class="steps">
        <li><span>1</span> O segredo é cifrado no teu dispositivo com uma chave aleatória.</li>
        <li><span>2</span> O <em>ciphertext</em> vai para o servidor, que o guarda só em RAM com TTL.</li>
        <li><span>3</span> A chave vai no <em>fragmento</em> do link (#) — nunca chega ao servidor.</li>
        <li><span>4</span> Ao abrir, o segredo é revelado e apagado da RAM. Sem rasto em disco.</li>
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
  .pill {
    flex-shrink: 0;
    font-size: var(--text-xs);
    font-weight: 600;
    letter-spacing: 0.04em;
    padding: 2px var(--space-2);
    border-radius: var(--radius-sm);
    border: 1px solid transparent;
  }
  .pill.on {
    color: var(--color-success-fg);
    background: var(--color-success-bg);
  }
  .panel-foot {
    margin: var(--space-4) 0;
    font-size: var(--text-xs);
    line-height: var(--leading-body);
    color: var(--color-text);
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
