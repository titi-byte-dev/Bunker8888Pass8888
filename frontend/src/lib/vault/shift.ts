/**
 * Políticas de acesso por turnos (VAULT-010).
 *
 * Didático: o servidor valida o turno com relógio UTC (NTP do SO). O cliente
 * espelha a mesma lógica para expurgar a Master Key ao fim da janela, sem
 * depender só do servidor para "lembrete".
 */

export type ShiftDay = "mon" | "tue" | "wed" | "thu" | "fri" | "sat" | "sun";

export interface ShiftWindow {
  start: string; // "HH:MM" 24h
  end: string;
}

export type WeeklyShiftSchedule = Partial<Record<ShiftDay, ShiftWindow[]>>;

export interface ShiftPolicy {
  enabled: boolean;
  timezone: string;
  schedule: WeeklyShiftSchedule;
  max_clock_skew_seconds: number;
}

function dayKey(date: Date, timeZone: string): ShiftDay {
  const weekday = new Intl.DateTimeFormat("en-US", { timeZone, weekday: "short" })
    .format(date)
    .toLowerCase()
    .slice(0, 3);
  const map: Record<string, ShiftDay> = {
    mon: "mon",
    tue: "tue",
    wed: "wed",
    thu: "thu",
    fri: "fri",
    sat: "sat",
    sun: "sun",
  };
  return map[weekday] ?? "mon";
}

function localMinutes(date: Date, timeZone: string): number {
  const parts = new Intl.DateTimeFormat("en-GB", {
    timeZone,
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).formatToParts(date);
  const h = Number(parts.find((p) => p.type === "hour")?.value ?? "0");
  const m = Number(parts.find((p) => p.type === "minute")?.value ?? "0");
  return h * 60 + m;
}

function parseHM(s: string): number {
  const [h, m] = s.trim().split(":").map(Number);
  if (Number.isNaN(h) || Number.isNaN(m) || h < 0 || h > 23 || m < 0 || m > 59) {
    throw new Error(`hora inválida: ${s}`);
  }
  return h * 60 + m;
}

function inWindow(minutes: number, start: number, end: number): boolean {
  if (start <= end) return minutes >= start && minutes < end;
  return minutes >= start || minutes < end;
}

/** Devolve true se `now` está dentro de alguma janela do turno. */
export function isWithinShift(now: Date, policy: ShiftPolicy): boolean {
  if (!policy.enabled) return true;
  const tz = policy.timezone || "UTC";
  const key = dayKey(now, tz);
  const windows = policy.schedule[key] ?? [];
  if (windows.length === 0) return false;
  const minutes = localMinutes(now, tz);
  return windows.some((w) => inWindow(minutes, parseHM(w.start), parseHM(w.end)));
}

/**
 * Milissegundos até ao fim da janela actual (para agendar expurgo da Master Key).
 * Devolve null se não estiver dentro de um turno.
 */
export function msUntilShiftEnd(now: Date, policy: ShiftPolicy): number | null {
  if (!policy.enabled || !isWithinShift(now, policy)) return null;
  const tz = policy.timezone || "UTC";
  const key = dayKey(now, tz);
  const windows = policy.schedule[key] ?? [];
  const minutes = localMinutes(now, tz);

  let best: number | null = null;
  for (const w of windows) {
    const start = parseHM(w.start);
    const end = parseHM(w.end);
    if (!inWindow(minutes, start, end)) continue;
    const minsLeft = start <= end ? end - minutes : minutes >= start ? 24 * 60 - minutes + end : end - minutes;
    if (best === null || minsLeft < best) best = minsLeft;
  }
  return best === null ? null : best * 60_000;
}

/** Valida desvio entre relógio local e servidor (mitiga manipulação de hora). */
export function isClockSkewAcceptable(
  clientNowMs: number,
  serverUnixMs: number,
  maxSkewSeconds: number,
): boolean {
  const skewMs = Math.abs(clientNowMs - serverUnixMs);
  return skewMs <= maxSkewSeconds * 1000;
}

export interface ServerTimeResponse {
  server_time: string;
  unix_ms: number;
}

/** Obtém a hora do servidor para sincronização NTP. */
export async function fetchServerTime(baseURL: string): Promise<ServerTimeResponse> {
  const res = await globalThis.fetch(`${baseURL}/api/time`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return res.json() as Promise<ServerTimeResponse>;
}
