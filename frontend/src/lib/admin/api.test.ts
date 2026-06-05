import { describe, expect, it } from "vitest";
import { formatSimpleSchedule, parseSimpleSchedule } from "./api";

describe("admin api helpers (UI-008)", () => {
  it("parseSimpleSchedule lê linhas por dia", () => {
    const text = "09:00-17:00\n09:00-17:00\n-\n10:00-14:00";
    const s = parseSimpleSchedule(text);
    expect(s.mon).toEqual([{ start: "09:00", end: "17:00" }]);
    expect(s.wed).toBeUndefined();
    expect(s.thu).toEqual([{ start: "10:00", end: "14:00" }]);
  });

  it("formatSimpleSchedule gera texto editável", () => {
    const text = formatSimpleSchedule({
      mon: [{ start: "08:00", end: "16:00" }],
    });
    expect(text.startsWith("08:00-16:00")).toBe(true);
  });
});
