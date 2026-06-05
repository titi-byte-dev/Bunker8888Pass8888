/**
 * Aliases de e-mail (MAIL-001) — cliente HTTP + orquestração.
 *
 * Ao contrário do resto do AegisPass, o destino é visível ao servidor: é ele
 * que (no MAIL-002, com SMTP) fará o reencaminhamento. Aqui gerimos só a
 * configuração dos aliases.
 */
import { loadSessionToken } from "$lib/session";

export interface EmailAlias {
  id: string;
  aliasAddress: string;
  destination: string;
  label: string;
  active: boolean;
  createdAt: string;
}

interface AliasDTO {
  id: string;
  alias_address: string;
  destination: string;
  label: string;
  active: boolean;
  created_at: string;
}

function mapAlias(d: AliasDTO): EmailAlias {
  return {
    id: d.id,
    aliasAddress: d.alias_address,
    destination: d.destination,
    label: d.label,
    active: d.active,
    createdAt: d.created_at,
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

export async function listAliases(): Promise<EmailAlias[]> {
  const j = (await (await authedFetch("/api/mail/aliases")).json()) as { aliases: AliasDTO[] };
  return (j.aliases ?? []).map(mapAlias);
}

export async function createAlias(destination: string, label: string): Promise<EmailAlias> {
  const dto = (await (
    await authedFetch("/api/mail/aliases", {
      method: "POST",
      body: JSON.stringify({ destination, label }),
    })
  ).json()) as AliasDTO;
  return mapAlias(dto);
}

export async function setAliasActive(id: string, active: boolean): Promise<void> {
  await authedFetch(`/api/mail/aliases/${encodeURIComponent(id)}`, {
    method: "PATCH",
    body: JSON.stringify({ active }),
  });
}

export async function deleteAlias(id: string): Promise<void> {
  await authedFetch(`/api/mail/aliases/${encodeURIComponent(id)}`, { method: "DELETE" });
}

/** Envia e-mail com remetente = alias (MAIL-004 compose). */
export async function composeFromAlias(
  aliasId: string,
  to: string,
  subject: string,
  body: string,
): Promise<void> {
  await authedFetch("/api/mail/compose", {
    method: "POST",
    body: JSON.stringify({ alias_id: aliasId, to, subject, body }),
  });
}
