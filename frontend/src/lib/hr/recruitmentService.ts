/**
 * Agente de recrutamento às cegas (AGENT-007).
 */
import { loadSessionToken } from "$lib/session";
import { blindScreenCV } from "./blindScreen";

export interface CandidateDraft {
  message_id: string;
  email: string;
  subject: string;
  summary: string;
  blind: boolean;
  source: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

async function authedFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return res;
}

/** Corre triagem às cegas sobre candidaturas pendentes na inbox. */
export async function runRecruitment(): Promise<CandidateDraft[]> {
  const j = (await (await authedFetch("/api/agent/recruitment/run", { method: "POST" })).json()) as {
    drafts: CandidateDraft[];
  };
  return (j.drafts ?? []).map((d) => ({
    ...d,
    summary: blindScreenCV(d.summary),
  }));
}

export async function seedRecruitmentEmail(fromEmail: string, subject: string, body: string): Promise<void> {
  await authedFetch("/api/mail/inbox", {
    method: "POST",
    body: JSON.stringify({ from_email: fromEmail, subject, body }),
  });
}

export async function markInboxProcessed(messageId: string): Promise<void> {
  await authedFetch(`/api/mail/inbox/${encodeURIComponent(messageId)}/processed`, { method: "POST" });
}
