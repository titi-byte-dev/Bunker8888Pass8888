/**
 * Guardas de autenticação (UI-003).
 *
 * Didático: separamos **login** (servidor confia na identidade → token) de
 * **unlock** (cliente deriva a Master Key localmente). Após refresh, o token
 * persiste em sessionStorage mas a Master Key volátil desaparece — por isso
 * redireccionamos para `/auth/unlock`.
 */
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { loadSessionToken, loadUserEmail } from "$lib/session";

export type AuthPhase = "guest" | "session" | "unlocked";

/** Fase actual da sessão no browser. */
export function getAuthPhase(): AuthPhase {
  const token = loadSessionToken();
  if (!token) return "guest";
  if (!getMasterKey()) return "session";
  return "unlocked";
}

/** Rotas da app que não exigem unlock (ex.: playground dev autónomo). */
export function isPublicAppPath(pathname: string): boolean {
  return pathname === "/dev" || pathname.startsWith("/dev/");
}

/** Rotas auth onde um guest pode estar sem redirect. */
export function isGuestAuthPath(pathname: string): boolean {
  return (
    pathname === "/auth" ||
    pathname === "/auth/login" ||
    pathname === "/auth/register" ||
    pathname === "/auth/recovery"
  );
}

/** Destino seguro após login/unlock (evita open redirect). */
export function sanitizeRedirect(raw: string | null, fallback = "/vault"): string {
  if (!raw || !raw.startsWith("/") || raw.startsWith("//")) return fallback;
  if (raw.startsWith("/auth")) return fallback;
  return raw;
}

export function resolveAuthRedirect(searchParams: URLSearchParams, fallback = "/vault"): string {
  return sanitizeRedirect(searchParams.get("redirect"), fallback);
}

/** Email para unlock — sessionStorage ou query string de fallback. */
export function resolveUnlockEmail(searchParams: URLSearchParams): string {
  return (loadUserEmail() ?? searchParams.get("email") ?? "").trim().toLowerCase();
}
