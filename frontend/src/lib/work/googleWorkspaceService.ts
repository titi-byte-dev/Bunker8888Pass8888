/**
 * Google Workspace (GOOGLE-001) — estado do provider no servidor.
 */
import { loadSessionToken } from "$lib/session";

export interface GoogleWorkspaceStatus {
  provider: string;
  enabled: boolean;
  delegated_user?: string;
  ready: boolean;
  scopes: string[];
  message?: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

export async function getGoogleWorkspaceStatus(): Promise<GoogleWorkspaceStatus> {
  const res = await fetch("/api/work/google/status", {
    headers: { Authorization: `Bearer ${token()}` },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  const j = (await res.json()) as { status: GoogleWorkspaceStatus };
  return j.status;
}
