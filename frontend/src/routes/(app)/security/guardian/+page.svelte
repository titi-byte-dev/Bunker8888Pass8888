<script lang="ts">
  import { onMount } from "svelte";
  import DocHelpLink from "$lib/docs/DocHelpLink.svelte";
  import { listGuardianAudit, type GuardianAuditEntry } from "$lib/agent/auditService";
  import { PageShell, StatusBanner } from "$lib/ui";

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

<PageShell
  title="Auditoria do Guardião"
  taskId="AGENT-010"
  description="Cada execução de tool regista apenas agente, nome da tool e sucesso/erro — nunca passwords, blobs decifrados nem Master Key."
>
  {#snippet actions()}
    <DocHelpLink slug="journey-guardian-audit" label="Como o Guardião protege a Master Key?" />
  {/snippet}

  <p class="hint">
    Em dev, corre prospeção em <a href="/crm">/crm</a> para gerar entradas.
  </p>

  {#if loading}
    <p class="muted">A carregar…</p>
  {:else if error}
    <StatusBanner variant="error">{error}</StatusBanner>
  {:else if entries.length === 0}
    <p class="muted">Sem execuções registadas — simula prospeção ou recrutamento.</p>
  {:else}
    <div class="table-wrap">
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
    </div>
  {/if}
</PageShell>

<style>
  .table-wrap {
    overflow-x: auto;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
  }

  table { width: 100%; border-collapse: collapse; font-size: var(--text-sm); }
  th, td { text-align: left; padding: var(--space-2) var(--space-3); border-bottom: 1px solid var(--color-border); }
  th { color: var(--color-text-muted); font-weight: 500; background: var(--color-bg-surface); }
  .mono { font-family: var(--font-mono); font-size: var(--text-xs); }
  tr.fail td { color: var(--color-danger); }
  .muted { font-size: var(--text-sm); color: var(--color-text-muted); }
  .hint { font-size: var(--text-sm); color: var(--color-text-muted); margin: 0 0 var(--space-4); }
  .hint a { color: var(--color-link); }
</style>
