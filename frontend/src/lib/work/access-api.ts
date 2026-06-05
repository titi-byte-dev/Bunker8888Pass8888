/**
 * Estado de políticas de acesso (turnos + geofence) para a área Trabalho.
 *
 * Didático: o servidor é a fonte de verdade para `within_shift` / `within_fence`;
 * o cliente espelha a lógica em `shift.ts` para countdown e expurgo da Master Key.
 */
import { loadSessionToken } from "$lib/session";
import type { ShiftPolicy } from "$lib/vault/shift";
import type { GeofencePolicy } from "$lib/vault/geofence";
import { geoHeaders, type GeoPosition } from "$lib/vault/geofence";

const API = "";

function authHeaders(extra: Record<string, string> = {}): HeadersInit {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão inválida");
  return {
    Authorization: `Bearer ${token}`,
    ...extra,
  };
}

export type ShiftStatus = ShiftPolicy & {
  within_shift: boolean;
  server_time: string;
};

export type GeofenceStatus = GeofencePolicy & {
  within_fence: boolean;
  client_ip: string;
};

/** Política de turno do utilizador + estado actual. */
export async function fetchShiftStatus(): Promise<ShiftStatus> {
  const res = await fetch(`${API}/api/access/shift`, { headers: authHeaders() });
  if (!res.ok) throw new Error("Falha ao obter turno");
  return res.json() as Promise<ShiftStatus>;
}

/**
 * Política de geofence + validação.
 * Se GPS estiver activo no servidor, passa coordenadas nos headers X-Geo-*.
 */
export async function fetchGeofenceStatus(position?: GeoPosition | null): Promise<GeofenceStatus> {
  const res = await fetch(`${API}/api/access/geofence`, {
    headers: authHeaders(geoHeaders(position ?? null)),
  });
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(body.error ?? "Falha ao obter geofence");
  }
  return res.json() as Promise<GeofenceStatus>;
}
