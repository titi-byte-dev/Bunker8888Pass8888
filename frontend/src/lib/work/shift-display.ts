/**
 * Formatação legível de turnos para UI (Pacote B — Trabalho).
 */
import type { ShiftDay, ShiftPolicy, ShiftWindow, WeeklyShiftSchedule } from "$lib/vault/shift";
import { isWithinShift, msUntilShiftEnd } from "$lib/vault/shift";

export const DAY_LABELS: Record<ShiftDay, string> = {
  mon: "Segunda",
  tue: "Terça",
  wed: "Quarta",
  thu: "Quinta",
  fri: "Sexta",
  sat: "Sábado",
  sun: "Domingo",
};

const DAY_ORDER: ShiftDay[] = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"];

export function formatWindow(w: ShiftWindow): string {
  return `${w.start} – ${w.end}`;
}

/** Lista janelas por dia para exibição na UI. */
export function formatWeeklySchedule(schedule: WeeklyShiftSchedule): { day: string; windows: string }[] {
  return DAY_ORDER.flatMap((key) => {
    const windows = schedule[key];
    if (!windows?.length) return [];
    return [{ day: DAY_LABELS[key], windows: windows.map(formatWindow).join(", ") }];
  });
}

/** Rótulo curto do estado do turno. */
export function shiftStatusLabel(policy: ShiftPolicy, within: boolean): string {
  if (!policy.enabled) return "Turnos desactivados";
  return within ? "Dentro do turno" : "Fora do turno";
}

/** Classe CSS semântica para badges de estado. */
export function shiftStatusTone(policy: ShiftPolicy, within: boolean): "ok" | "warn" | "neutral" {
  if (!policy.enabled) return "neutral";
  return within ? "ok" : "warn";
}

/** Formata milissegundos como «Xh Ym» ou «Ym Zs». */
export function formatCountdown(ms: number): string {
  if (ms <= 0) return "0s";
  const totalSec = Math.ceil(ms / 1000);
  const h = Math.floor(totalSec / 3600);
  const m = Math.floor((totalSec % 3600) / 60);
  const s = totalSec % 60;
  if (h > 0) return `${h}h ${m}m`;
  if (m > 0) return `${m}m ${s}s`;
  return `${s}s`;
}

/** Countdown até fim do turno (null se fora ou desactivado). */
export function shiftCountdown(now: Date, policy: ShiftPolicy): string | null {
  const ms = msUntilShiftEnd(now, policy);
  if (ms === null) return null;
  return formatCountdown(ms);
}

/** Espelha servidor: dentro do turno com base na política e hora actual. */
export function clientWithinShift(now: Date, policy: ShiftPolicy): boolean {
  return isWithinShift(now, policy);
}

/** Desvio de relógio legível (segundos). */
export function formatClockSkewMs(skewMs: number): string {
  const sec = Math.round(Math.abs(skewMs) / 1000);
  return `${sec}s`;
}
