/**
 * Store global de toasts (UI-017).
 *
 * Didático: em vez de cada página inventar um `<p class="status">`, centralizamos
 * feedback efémero numa fila. O ToastHost no layout subscreve `toastStore`.
 *
 * Usamos `writable` (svelte/store) em vez de `$state` aqui para o módulo ser
 * testável no Vitest sem compilar runes.
 *
 * > ⚠️ **Segurança:** nunca passar passwords, tokens ou IBAN na mensagem.
 */

import { get, writable } from "svelte/store";

export type ToastVariant = "info" | "success" | "warning" | "error";

export type ToastItem = {
  id: string;
  message: string;
  variant: ToastVariant;
  /** ms até auto-dismiss; 0 = persistir até o utilizador fechar */
  durationMs: number;
};

const DEFAULT_DURATION = 4_000;

export const toastStore = writable<ToastItem[]>([]);

const timers = new Map<string, ReturnType<typeof setTimeout>>();

/** Snapshot actual — útil em testes */
export function toastQueue(): ToastItem[] {
  return get(toastStore);
}

export function dismissToast(id: string): void {
  const t = timers.get(id);
  if (t) {
    clearTimeout(t);
    timers.delete(id);
  }
  toastStore.update((q) => q.filter((item) => item.id !== id));
}

export function pushToast(
  message: string,
  opts?: { variant?: ToastVariant; durationMs?: number },
): string {
  const id =
    typeof crypto !== "undefined" && crypto.randomUUID
      ? crypto.randomUUID()
      : `toast-${Date.now()}-${Math.random().toString(36).slice(2)}`;

  const durationMs = opts?.durationMs ?? DEFAULT_DURATION;
  const item: ToastItem = {
    id,
    message,
    variant: opts?.variant ?? "info",
    durationMs,
  };

  toastStore.update((q) => [...q, item]);

  if (durationMs > 0) {
    const timer = setTimeout(() => dismissToast(id), durationMs);
    timers.set(id, timer);
  }

  return id;
}

/** Atalhos por variante semântica */
export const toast = {
  info: (message: string, durationMs?: number) =>
    pushToast(message, { variant: "info", durationMs }),
  success: (message: string, durationMs?: number) =>
    pushToast(message, { variant: "success", durationMs }),
  warning: (message: string, durationMs?: number) =>
    pushToast(message, { variant: "warning", durationMs }),
  error: (message: string, durationMs?: number) =>
    pushToast(message, { variant: "error", durationMs }),
};
