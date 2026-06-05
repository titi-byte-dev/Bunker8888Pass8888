/**
 * Chave admin volátil — só sessionStorage (UI-008).
 *
 * ⚠️ Segurança: nunca persistir em localStorage; expira ao fechar o separador.
 */
const ADMIN_KEY = "aegis:admin-key";

export function saveAdminKey(key: string): void {
  sessionStorage.setItem(ADMIN_KEY, key.trim());
}

export function loadAdminKey(): string | null {
  return sessionStorage.getItem(ADMIN_KEY);
}

export function clearAdminKey(): void {
  sessionStorage.removeItem(ADMIN_KEY);
}

export function hasAdminKey(): boolean {
  return !!loadAdminKey();
}
