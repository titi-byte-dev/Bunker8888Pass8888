/** Sessão volátil — token HTTP e email nunca em localStorage (só sessionStorage). */
const TOKEN_KEY = "aegis:session-token";
const EMAIL_KEY = "aegis:user-email";

export function saveSessionToken(token: string): void {
  sessionStorage.setItem(TOKEN_KEY, token);
}

export function loadSessionToken(): string | null {
  return sessionStorage.getItem(TOKEN_KEY);
}

export function clearSessionToken(): void {
  sessionStorage.removeItem(TOKEN_KEY);
}

/** Email do utilizador autenticado — necessário para unlock após refresh. */
export function saveUserEmail(email: string): void {
  sessionStorage.setItem(EMAIL_KEY, email.trim().toLowerCase());
}

export function loadUserEmail(): string | null {
  return sessionStorage.getItem(EMAIL_KEY);
}

export function clearUserEmail(): void {
  sessionStorage.removeItem(EMAIL_KEY);
}

/** Limpa token + email (logout). */
export function clearSession(): void {
  clearSessionToken();
  clearUserEmail();
}
