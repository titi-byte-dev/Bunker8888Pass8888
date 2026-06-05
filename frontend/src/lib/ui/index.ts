/**
 * lib/ui — biblioteca de componentes reutilizaveis (UI-012).
 *
 * Importar daqui (nao copiar CSS para as paginas):
 *   import { PageShell, Panel, Button } from "$lib/ui";
 *
 * Antes: ~85 blocos `.page-head`/`.panel`/`.eyebrow` re-definidos em 32 paginas.
 * Depois: um unico componente por padrao, tokens via var(--*).
 */
export { default as PageShell } from "./PageShell.svelte";
export { default as Panel } from "./Panel.svelte";
export { default as Button } from "./Button.svelte";
export { default as Field } from "./Field.svelte";
export { default as Eyebrow } from "./Eyebrow.svelte";
export { default as EmptyState } from "./EmptyState.svelte";
export { default as StatusBanner } from "./StatusBanner.svelte";
export { default as HubLinks } from "./HubLinks.svelte";
export { default as Breadcrumbs } from "./Breadcrumbs.svelte";

export type { HubLinkItem } from "./HubLinks.svelte";
