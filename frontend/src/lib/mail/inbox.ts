/**
 * Caixa de entrada (MAIL-002 / AGENT-003) — mensagens recebidas via alias.
 */
import { loadSessionToken } from "$lib/session";

export interface InboxMessage {
  id: string;
  fromEmail: string;
  subject: string;
  body: string;
  receivedAt: string;
  processedAt?: string;
}

interface InboxDTO {
  id: string;
  from_email: string;
  subject: string;
  body: string;
  received_at: string;
  processed_at?: string;
}

function mapInbox(d: InboxDTO): InboxMessage {
  return {
    id: d.id,
    fromEmail: d.from_email,
    subject: d.subject,
    body: d.body,
    receivedAt: d.received_at,
    processedAt: d.processed_at,
  };
}

async function authedFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessao expirada — inicia sessao de novo.");
  const res = await fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token}`,
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

/** Lista mensagens da inbox (opcional: só pendentes). */
export async function listInbox(unprocessedOnly = false): Promise<InboxMessage[]> {
  const q = unprocessedOnly ? "?unprocessed=1" : "";
  const j = (await (await authedFetch(`/api/mail/inbox${q}`)).json()) as {
    messages: InboxDTO[];
  };
  return (j.messages ?? []).map(mapInbox);
}
