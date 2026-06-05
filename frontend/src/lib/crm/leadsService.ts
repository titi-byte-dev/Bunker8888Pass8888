/**
 * API CRM — leads cifrados (CRM-001).
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { base64ToBytes, bytesToBase64 } from "$lib/auth";
import { decryptLead, encryptLead, type Lead, type LeadPayload } from "./leads";

interface LeadDTO {
  id: string;
  blob: string;
  created_at: string;
  updated_at: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para ver leads.");
  return mk;
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

export async function listLeads(): Promise<Lead[]> {
  const mk = requireMasterKey();
  const j = (await (await authedFetch("/api/crm/leads")).json()) as { leads: LeadDTO[] };
  const out: Lead[] = [];
  for (const dto of j.leads ?? []) {
    try {
      const payload = await decryptLead(mk, base64ToBytes(dto.blob));
      out.push({ id: dto.id, ...payload });
    } catch {
      /* blob corrompido — ignorar */
    }
  }
  return out;
}

export async function createLead(payload: LeadPayload): Promise<Lead> {
  const mk = requireMasterKey();
  const blob = await encryptLead(mk, payload);
  const dto = (await (
    await authedFetch("/api/crm/leads", {
      method: "POST",
      body: JSON.stringify({ blob: bytesToBase64(blob) }),
    })
  ).json()) as LeadDTO;
  return { id: dto.id, ...payload };
}

export async function updateLead(id: string, payload: LeadPayload): Promise<Lead> {
  const mk = requireMasterKey();
  const blob = await encryptLead(mk, payload);
  await authedFetch(`/api/crm/leads/${id}`, {
    method: "PUT",
    body: JSON.stringify({ blob: bytesToBase64(blob) }),
  });
  return { id, ...payload };
}

export async function deleteLead(id: string): Promise<void> {
  await authedFetch(`/api/crm/leads/${id}`, { method: "DELETE" });
}
