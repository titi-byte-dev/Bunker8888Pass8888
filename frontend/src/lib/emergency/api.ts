/**
 * API de acesso de emergência (VAULT-016).
 */
import { loadSessionToken } from "$lib/session";

const API = "";

function authHeaders(): HeadersInit {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão inválida");
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export type EmergencyConfig = {
  configured: boolean;
  heir_email?: string;
  wait_days?: number;
  has_blob?: boolean;
  updated_at?: string;
};

export type EmergencyRequest = {
  id: string;
  heir_email: string;
  status: "waiting" | "rejected" | "ready" | "consumed";
  requested_at: string;
  unlocks_at: string;
  rejected_at?: string;
  consumed_at?: string;
};

export async function fetchEmergencyConfig(): Promise<EmergencyConfig> {
  const res = await fetch(`${API}/api/emergency/config`, { headers: authHeaders() });
  if (!res.ok) throw new Error("Falha ao carregar configuração");
  return res.json() as Promise<EmergencyConfig>;
}

export async function saveEmergencyConfig(body: {
  heir_email: string;
  wait_days: number;
  blob?: string;
}): Promise<void> {
  const res = await fetch(`${API}/api/emergency/config`, {
    method: "PUT",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Falha ao guardar");
  }
}

export async function deleteEmergencyConfig(): Promise<void> {
  const res = await fetch(`${API}/api/emergency/config`, {
    method: "DELETE",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao remover herdeiro");
}

export async function listEmergencyRequests(): Promise<EmergencyRequest[]> {
  const res = await fetch(`${API}/api/emergency/requests`, { headers: authHeaders() });
  if (!res.ok) throw new Error("Falha ao listar pedidos");
  const j = (await res.json()) as { requests: EmergencyRequest[] };
  return j.requests ?? [];
}

export async function approveEmergencyRequest(id: string): Promise<void> {
  const res = await fetch(`${API}/api/emergency/requests/${id}/approve`, {
    method: "POST",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao aprovar");
}

export async function rejectEmergencyRequest(id: string): Promise<void> {
  const res = await fetch(`${API}/api/emergency/requests/${id}/reject`, {
    method: "POST",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao rejeitar");
}

export async function createEmergencyRequest(ownerEmail: string): Promise<EmergencyRequest> {
  const res = await fetch(`${API}/api/emergency/request`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ owner_email: ownerEmail }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Falha ao pedir acesso");
  }
  const j = (await res.json()) as { request: EmergencyRequest };
  return j.request;
}

export async function fetchEmergencyRequestStatus(
  ownerEmail: string,
): Promise<{ active: boolean; request?: EmergencyRequest }> {
  const q = encodeURIComponent(ownerEmail.trim().toLowerCase());
  const res = await fetch(`${API}/api/emergency/request/status?owner_email=${q}`, {
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao consultar estado");
  return res.json() as Promise<{ active: boolean; request?: EmergencyRequest }>;
}

export async function fetchEmergencyAccessBlob(ownerEmail: string): Promise<string> {
  const q = encodeURIComponent(ownerEmail.trim().toLowerCase());
  const res = await fetch(`${API}/api/emergency/access?owner_email=${q}`, {
    headers: authHeaders(),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Acesso indisponível");
  }
  const j = (await res.json()) as { blob: string };
  return j.blob;
}

/** Segundos restantes até unlocks_at (0 se já passou). */
export function secondsUntil(iso: string, now = Date.now()): number {
  const t = new Date(iso).getTime();
  return Math.max(0, Math.ceil((t - now) / 1000));
}
