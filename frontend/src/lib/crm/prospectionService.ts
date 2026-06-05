/**
 * Agente de prospeção (AGENT-003) — rascunhos no servidor, cifragem no cliente.
 */
import { loadSessionToken } from "$lib/session";
import { createLead } from "./leadsService";
import type { LeadStage } from "./leads";

export interface ProspectionDraft {
  message_id: string;
  email: string;
  subject: string;
  notes: string;
  suggested_stage: LeadStage;
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

/** Corre o agente de prospeção sobre e-mails pendentes. */
export async function runProspection(): Promise<ProspectionDraft[]> {
  const j = (await (await authedFetch("/api/agent/prospection/run", { method: "POST" })).json()) as {
    drafts: ProspectionDraft[];
  };
  return j.drafts ?? [];
}

/** Marca mensagem de inbox como processada após importar lead. */
export async function markInboxProcessed(messageId: string): Promise<void> {
  await authedFetch(`/api/mail/inbox/${encodeURIComponent(messageId)}/processed`, {
    method: "POST",
  });
}

/** Deriva um nome legível a partir do e-mail ou assunto. */
function suggestName(email: string, subject: string): string {
  const local = email.split("@")[0]?.replace(/[._-]/g, " ").trim();
  if (local && local.length > 1) {
    return local.charAt(0).toUpperCase() + local.slice(1);
  }
  return subject.slice(0, 80) || email;
}

/** Cifra e grava lead; marca inbox como processada (fluxo ZK completo). */
export async function importDraft(draft: ProspectionDraft): Promise<void> {
  await createLead({
    name: suggestName(draft.email, draft.subject),
    email: draft.email,
    stage: draft.suggested_stage ?? "new",
    notes: draft.notes,
    source: draft.source ?? "email",
  });
  await markInboxProcessed(draft.message_id);
}

/** Simula e-mail recebido (stub até MAIL-002). */
export async function seedInboxMessage(fromEmail: string, subject: string, body: string): Promise<void> {
  await authedFetch("/api/mail/inbox", {
    method: "POST",
    body: JSON.stringify({ from_email: fromEmail, subject, body }),
  });
}
