/**
 * API de dispositivos e sessões (Pacote A — hub Segurança).
 */
import { loadSessionToken } from "$lib/session";
import type { PasskeyMeta } from "$lib/passkey";
import { listPasskeys } from "$lib/passkey";

const API = "";

function authHeaders(): HeadersInit {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão inválida");
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
}

export type HttpSession = {
  id: string;
  created_at: string;
  expires_at: string;
  current: boolean;
};

export type CliDevice = {
  id: string;
  name: string;
  created_at: string;
};

export async function listHttpSessions(): Promise<HttpSession[]> {
  const res = await fetch(`${API}/api/auth/sessions`, { headers: authHeaders() });
  if (!res.ok) throw new Error("Falha ao listar sessões");
  const data = (await res.json()) as { sessions: HttpSession[] };
  return data.sessions ?? [];
}

export async function revokeHttpSession(sessionId: string): Promise<void> {
  const res = await fetch(`${API}/api/auth/sessions/${sessionId}`, {
    method: "DELETE",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao revogar sessão");
}

export async function revokeOtherHttpSessions(): Promise<number> {
  const res = await fetch(`${API}/api/auth/sessions/revoke-others`, {
    method: "POST",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao revogar outras sessões");
  const data = (await res.json()) as { revoked: number };
  return data.revoked ?? 0;
}

export async function listCliDevices(): Promise<CliDevice[]> {
  const res = await fetch(`${API}/api/cli/devices`, { headers: authHeaders() });
  if (!res.ok) throw new Error("Falha ao listar dispositivos CLI");
  const data = (await res.json()) as { devices: CliDevice[] };
  return data.devices ?? [];
}

export async function revokeCliDevice(deviceId: string): Promise<void> {
  const res = await fetch(`${API}/api/cli/devices/${deviceId}`, {
    method: "DELETE",
    headers: authHeaders(),
  });
  if (!res.ok) throw new Error("Falha ao revogar dispositivo");
}

export async function listRegisteredPasskeys(): Promise<PasskeyMeta[]> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão inválida");
  return listPasskeys(API, token);
}
