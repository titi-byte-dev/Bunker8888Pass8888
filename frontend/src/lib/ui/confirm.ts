/**
 * Diálogo de confirmação global (UI-017).
 *
 * Didático: `window.confirm()` bloqueia o thread e não permite estilo acessível.
 * Este módulo expõe uma Promise — o componente ConfirmDialog no layout resolve-a.
 */

import { get, writable } from "svelte/store";

export type ConfirmVariant = "default" | "danger";

export type ConfirmOptions = {
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: ConfirmVariant;
};

type ConfirmState = {
  open: boolean;
  options: ConfirmOptions | null;
  resolve: ((value: boolean) => void) | null;
};

const initial: ConfirmState = {
  open: false,
  options: null,
  resolve: null,
};

export const confirmStore = writable<ConfirmState>(initial);

export function confirmState(): ConfirmState {
  return get(confirmStore);
}

export function confirmDialog(options: ConfirmOptions): Promise<boolean> {
  const current = get(confirmStore);
  if (current.open && current.resolve) {
    current.resolve(false);
  }
  return new Promise((resolve) => {
    confirmStore.set({ open: true, options, resolve });
  });
}

export function closeConfirm(result: boolean): void {
  const current = get(confirmStore);
  const r = current.resolve;
  confirmStore.set(initial);
  r?.(result);
}
