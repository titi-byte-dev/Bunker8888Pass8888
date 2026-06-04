/**
 * Perfil mínimo da sessão HTTP — email para unlock após refresh.
 */
import { wrapNetworkError } from "$lib/auth/http";
import { loadSessionToken } from "$lib/session";

export type SessionProfile = {
  email: string;
};

/** Email associado ao token Bearer actual (requer sessão válida). */
export async function fetchSessionProfile(apiBase = ""): Promise<SessionProfile> {
  const token = loadSessionToken();
  if (!token) {
    throw new Error("Sessão inválida");
  }
  let res: Response;
  try {
    res = await fetch(`${apiBase}/api/auth/session`, {
      headers: { Authorization: `Bearer ${token}` },
    });
  } catch (err) {
    throw wrapNetworkError(err, "Sessão");
  }
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(body.error ?? "Sessão inválida");
  }
  const data = (await res.json()) as SessionProfile;
  return { email: data.email.trim().toLowerCase() };
}
