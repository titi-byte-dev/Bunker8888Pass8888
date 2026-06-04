/**
 * Helpers partilhados para pedidos HTTP ao backend em dev.
 */
export function normalizeEmail(email: string): string {
  return email.trim().toLowerCase();
}

/** Mensagem legível quando o fetch falha (backend parado, proxy, etc.). */
export function wrapNetworkError(err: unknown, context: string): Error {
  if (err instanceof TypeError && /fetch|network|failed/i.test(err.message)) {
    return new Error(
      `${context}: servidor inacessível. Corre «docker compose up» e confirma http://localhost:8080/healthz`,
    );
  }
  if (err instanceof Error) return err;
  return new Error(`${context}: erro desconhecido`);
}
