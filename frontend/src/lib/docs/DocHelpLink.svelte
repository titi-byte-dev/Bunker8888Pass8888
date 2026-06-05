<script lang="ts">
  /**
   * Link contextual para a documentação in-app (DOC-013).
   * Didático: mapeia a rota actual (ou slug explícito) para a página didática certa.
   */
  import { page } from "$app/state";
  import { DOC_MANIFEST } from "./loader";
  import { resolveDocHelp } from "./docHelpLinks";

  interface Props {
    /** Slug da doc; omite para resolver pela rota actual */
    slug?: string;
    /** Texto do link; omite para usar o label do mapa de rotas */
    label?: string;
  }

  let { slug, label }: Props = $props();

  const target = $derived(resolveDocHelp(page.url.pathname, slug));
  const linkLabel = $derived(label ?? target?.label ?? "Como funciona?");
  const docTitle = $derived(
    target ? DOC_MANIFEST.docs.find((d) => d.slug === target.slug)?.title : undefined,
  );
</script>

{#if target}
  <a
    class="doc-help-link"
    href="/settings/docs/{target.slug}"
    title={docTitle ? `Documentação: ${docTitle}` : undefined}
  >
    <span class="icon" aria-hidden="true">📖</span>
    <span>{linkLabel}</span>
  </a>
{/if}

<style>
  .doc-help-link {
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
    padding: var(--space-1) var(--space-3);
    border: 1px solid color-mix(in srgb, var(--color-accent) 35%, var(--color-border));
    border-radius: var(--radius-md);
    background: color-mix(in srgb, var(--color-accent) 8%, var(--color-bg-surface));
    color: var(--color-link);
    font-size: var(--text-xs);
    font-weight: 500;
    text-decoration: none;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .doc-help-link:hover {
    border-color: var(--color-accent);
    background: var(--color-accent-muted);
  }

  .icon {
    font-size: 0.95em;
    line-height: 1;
  }
</style>
