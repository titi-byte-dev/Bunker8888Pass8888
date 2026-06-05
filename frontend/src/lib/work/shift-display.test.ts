import { describe, expect, it } from "vitest";
import {
  formatCountdown,
  formatWeeklySchedule,
  shiftCountdown,
  shiftStatusLabel,
  shiftStatusTone,
} from "./shift-display";
import type { ShiftPolicy } from "$lib/vault/shift";

const policy: ShiftPolicy = {
  enabled: true,
  timezone: "UTC",
  max_clock_skew_seconds: 120,
  schedule: {
    mon: [{ start: "09:00", end: "17:00" }],
    wed: [{ start: "10:00", end: "14:00" }],
  },
};

describe("shift-display (Pacote B)", () => {
  it("formatWeeklySchedule omite dias sem janelas", () => {
    const rows = formatWeeklySchedule(policy.schedule);
    expect(rows).toHaveLength(2);
    expect(rows[0]?.day).toBe("Segunda");
    expect(rows[0]?.windows).toBe("09:00 – 17:00");
  });

  it("shiftStatusLabel reflecte enabled e within", () => {
    expect(shiftStatusLabel({ ...policy, enabled: false }, false)).toBe("Turnos desactivados");
    expect(shiftStatusLabel(policy, true)).toBe("Dentro do turno");
    expect(shiftStatusLabel(policy, false)).toBe("Fora do turno");
  });

  it("shiftStatusTone mapeia estados", () => {
    expect(shiftStatusTone({ ...policy, enabled: false }, false)).toBe("neutral");
    expect(shiftStatusTone(policy, true)).toBe("ok");
    expect(shiftStatusTone(policy, false)).toBe("warn");
  });

  it("formatCountdown formata horas e minutos", () => {
    expect(formatCountdown(3_661_000)).toBe("1h 1m");
    expect(formatCountdown(90_000)).toBe("1m 30s");
    expect(formatCountdown(5_000)).toBe("5s");
  });

  it("shiftCountdown devolve texto dentro do turno", () => {
    // Segunda 2024-06-03 12:00 UTC — dentro de 09:00-17:00
    const now = new Date("2024-06-03T12:00:00.000Z");
    const cd = shiftCountdown(now, policy);
    expect(cd).toMatch(/5h/);
  });
});
