import { describe, expect, it, vi, afterEach } from "vitest";
import {
  getMasterKey,
  purgeMasterKey,
  scheduleShiftPurge,
  setMasterKey,
} from "./masterKeyStore";
import type { ShiftPolicy } from "./shift";

const fakeKey = {} as CryptoKey;

const shortPolicy: ShiftPolicy = {
  enabled: true,
  timezone: "UTC",
  max_clock_skew_seconds: 300,
  schedule: {
    wed: [{ start: "09:00", end: "17:00" }],
  },
};

describe("masterKeyStore", () => {
  afterEach(() => {
    purgeMasterKey();
    vi.useRealTimers();
  });

  it("set/get/purge", () => {
    setMasterKey(fakeKey);
    expect(getMasterKey()).toBe(fakeKey);
    purgeMasterKey();
    expect(getMasterKey()).toBeNull();
  });

  it("expurga imediatamente fora do turno", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-06-03T20:00:00.000Z"));
    setMasterKey(fakeKey);
    const onPurged = vi.fn();
    scheduleShiftPurge(shortPolicy, onPurged);
    expect(getMasterKey()).toBeNull();
    expect(onPurged).toHaveBeenCalledOnce();
  });
});
