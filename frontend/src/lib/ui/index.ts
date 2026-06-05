/**
 * lib/ui — biblioteca de componentes reutilizaveis (UI-012, UI-015, UI-017).
 *
 * Importar daqui (nao copiar CSS para as paginas):
 *   import { PageShell, Panel, Button, toast, confirmDialog } from "$lib/ui";
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

/* UI-015 — dados e densidade */
export { default as MetricCard } from "./MetricCard.svelte";
export { default as ListRow } from "./ListRow.svelte";
export { default as DataTable } from "./DataTable.svelte";
export type { DataColumn } from "./data-table";

/* UI-017 — feedback global */
export { default as Skeleton } from "./Skeleton.svelte";
export { default as ToastHost } from "./ToastHost.svelte";
export { default as ConfirmDialog } from "./ConfirmDialog.svelte";
export { toast, pushToast, dismissToast, toastQueue, toastStore } from "./toast";
export type { ToastVariant, ToastItem } from "./toast";
export { confirmDialog, closeConfirm, confirmState, confirmStore } from "./confirm";
export type { ConfirmOptions, ConfirmVariant } from "./confirm";

export type { HubLinkItem } from "./HubLinks.svelte";
