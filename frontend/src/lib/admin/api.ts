/**
 * API administrativa (UI-008) — requer header X-Admin-Key.
 */
import { loadAdminKey } from "./adminKey";
import type { ShiftPolicy, WeeklyShiftSchedule } from "$lib/vault/shift";
import type { GeofencePolicy } from "$lib/vault/geofence";

const API = "";

export type AdminUser = {
  ID: string;
  Email: string;
  CreatedAt: string;
};

export type WipeAuditEvent = {
  id: string;
  target_user_id: string;
  target_email: string;
  initiated_by: string;
  reason: string;
  devices_notified: number;
  created_at: string;
};

export type RemoteWipeResult = {
  devices_notified: number;
  sessions_revoked: number;
};

function adminHeaders(): HeadersInit {
  const key = loadAdminKey();
  if (!key) throw new Error("Chave admin em falta");
  return {
    "X-Admin-Key": key,
    "Content-Type": "application/json",
  };
}

async function adminFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(`${API}${path}`, {
    ...init,
    headers: { ...adminHeaders(), ...init.headers },
  });
  if (res.status === 403) throw new Error("Chave admin inválida");
  if (res.status === 503) throw new Error("Admin desactivado no servidor (AEGIS_ADMIN_KEY)");
  return res;
}

/** Valida a chave tentando listar utilizadores. */
export async function verifyAdminKey(key: string): Promise<boolean> {
  const res = await fetch(`${API}/api/admin/users`, {
    headers: { "X-Admin-Key": key.trim() },
  });
  if (res.status === 403 || res.status === 503) return false;
  return res.ok;
}

export async function listAdminUsers(): Promise<AdminUser[]> {
  const res = await adminFetch("/api/admin/users");
  if (!res.ok) throw new Error("Falha ao listar utilizadores");
  const data = (await res.json()) as {
    users: { id: string; email: string; created_at: string }[];
  };
  return (data.users ?? []).map((u) => ({
    ID: u.id,
    Email: u.email,
    CreatedAt: u.created_at,
  }));
}

export async function getAdminUser(id: string): Promise<{ id: string; email: string }> {
  const res = await adminFetch(`/api/admin/users/${id}`);
  if (!res.ok) throw new Error("Utilizador não encontrado");
  return res.json() as Promise<{ id: string; email: string }>;
}

export async function getAdminShiftPolicy(userId: string): Promise<ShiftPolicy> {
  const res = await adminFetch(`/api/admin/users/${userId}/access-shift`);
  if (!res.ok) throw new Error("Falha ao obter turno");
  const data = await res.json();
  return {
    enabled: data.enabled,
    timezone: data.timezone,
    schedule: data.schedule as WeeklyShiftSchedule,
    max_clock_skew_seconds: data.max_clock_skew_seconds,
  };
}

export async function setAdminShiftPolicy(userId: string, policy: ShiftPolicy): Promise<void> {
  const res = await adminFetch(`/api/admin/users/${userId}/access-shift`, {
    method: "PUT",
    body: JSON.stringify({
      enabled: policy.enabled,
      timezone: policy.timezone,
      schedule: policy.schedule,
      max_clock_skew_seconds: policy.max_clock_skew_seconds,
    }),
  });
  if (!res.ok) throw new Error("Falha ao gravar turno");
}

export async function getAdminGeofencePolicy(userId: string): Promise<GeofencePolicy> {
  const res = await adminFetch(`/api/admin/users/${userId}/access-geofence`);
  if (!res.ok) throw new Error("Falha ao obter geofence");
  return res.json() as Promise<GeofencePolicy>;
}

export async function setAdminGeofencePolicy(userId: string, policy: GeofencePolicy): Promise<void> {
  const res = await adminFetch(`/api/admin/users/${userId}/access-geofence`, {
    method: "PUT",
    body: JSON.stringify({
      enabled: policy.enabled,
      allowed_cidrs: policy.allowed_cidrs,
      gps_enabled: policy.gps_enabled,
      gps_lat: policy.gps_lat,
      gps_lon: policy.gps_lon,
      gps_radius_m: policy.gps_radius_m,
    }),
  });
  if (!res.ok) throw new Error("Falha ao gravar geofence");
}

export async function triggerRemoteWipe(userId: string, reason: string): Promise<RemoteWipeResult> {
  const res = await adminFetch(`/api/admin/users/${userId}/remote-wipe`, {
    method: "POST",
    body: JSON.stringify({ reason }),
  });
  if (!res.ok) throw new Error("Remote wipe falhou");
  return res.json() as Promise<RemoteWipeResult>;
}

export async function listWipeAuditEvents(): Promise<WipeAuditEvent[]> {
  const res = await adminFetch("/api/admin/audit/wipe-events");
  if (!res.ok) throw new Error("Falha ao carregar auditoria");
  const data = (await res.json()) as { events: WipeAuditEvent[] };
  return data.events ?? [];
}

/** Converte texto «09:00-17:00» por linha num schedule mon-fri simples. */
export function parseSimpleSchedule(text: string): WeeklyShiftSchedule {
  const lines = text
    .split("\n")
    .map((l) => l.trim())
    .filter(Boolean);
  const schedule: WeeklyShiftSchedule = {};
  const days = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"] as const;
  for (let i = 0; i < lines.length && i < days.length; i++) {
    const line = lines[i]!;
    if (line === "-" || line === "—") continue;
    const [start, end] = line.split(/[-–]/).map((s) => s.trim());
    if (start && end) {
      schedule[days[i]!] = [{ start, end }];
    }
  }
  return schedule;
}

export function formatSimpleSchedule(schedule: WeeklyShiftSchedule): string {
  const days = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"] as const;
  return days
    .map((d) => {
      const w = schedule[d];
      if (!w?.length) return "-";
      return `${w[0]!.start}-${w[0]!.end}`;
    })
    .join("\n");
}
