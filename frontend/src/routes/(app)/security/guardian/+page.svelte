<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listGuardianAudit, type GuardianAuditEntry } from "$lib/agent/auditService";

  let loading = $state(true);
  let error = $state("");
  let entries = $state<GuardianAuditEntry[]>([]);

  onMount(async () => {
    try {
      entries = await listGuardianAudit();
    } catch (e) {
      error = e instanceof Error ? e.message : "Falha ao carregar auditoria";
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head><title>Guardião — AegisPass</title></svelte:head>

<section class="page">
  <header class="head">
    <div>
      <p class="eyebrow">AGENT-002 · DoD Fase 3</p>
      <h1>Auditoria do Guardião</h1>
      <DocHelpLink slug="journey-guardian-audit" label="Como o Guardião protege a Master Key?" />
    </div>
    <a class="back" href="/security">← Segurança</a>
  </header>

  <p class="lead">
    Cada execução de tool regista <strong>apenas</strong> agente, nome da tool e sucesso/erro —
    nunca passwords, blobs decifrados nem Master Key. Em dev, corre prospeção em <a href="/crm">/crm</a>
    para gerar entradas.
  </p>

  {#if loading}<p class="muted">A carregar…</p>
  {:else if error}<p class="err">{error}</p>
  {:else if entries.length === 0}
    <p class="muted">Sem execuções registadas — simula prospeção ou recrutamento.</p>
  {:else}
    <table>
      <thead><tr><th>Quando</th><th>Agente</th><th>Tool</th><th>Resultado</th></tr></thead>
      <tbody>
        {#each entries as e (e.id)}
          <tr class:fail={!e.success}>
            <td>{new Date(e.createdAt).toLocaleString("pt-PT")}</td>
            <td class="mono">{e.agentId}</td>
            <td class="mono">{e.toolName}</td>
            <td>{e.success ? "✓ OK" : e.errorMsg || "Falhou"}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</section>

<style>
  .page { max-width: 52rem; margin: 0 auto; padding: var(--space-6); }
  .head { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: var(--space-4); }
  .eyebrow { font-size: var(--text-xs); text-transform: uppercase; color: var(--color-text-label); margin: 0; }
  h1 { margin: var(--space-1) 0; }
  .back { font-size: var(--text-sm); color: var(--color-text-muted); }
  .lead { font-size: var(--text-sm); color: var(--color-text-muted); margin-bottom: var(--space-4); }
  table { width: 100%; border-collapse: collapse; font-size: var(--text-sm); }
  th, td { text-align: left; padding: var(--space-2); border-bottom: 1px solid var(--color-border); }
  .mono { font-family: var(--font-mono); font-size: var(--text-xs); }
  tr.fail td { color: var(--color-danger); }
  .muted, .err { font-size: var(--text-sm); }
  .err { color: var(--color-danger); }
</style>
