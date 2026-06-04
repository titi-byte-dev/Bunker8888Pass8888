/** Sessão volátil — token HTTP nunca em localStorage (só sessionStorage). */
const TOKEN_KEY = "aegis:session-token";

export function saveSessionToken(token: string): void {
  sessionStorage.setItem(TOKEN_KEY, token);
}

export function loadSessionToken(): string | null {
  return sessionStorage.getItem(TOKEN_KEY);
}

export function clearSessionToken(): void {
  sessionStorage.removeItem(TOKEN_KEY);
}
