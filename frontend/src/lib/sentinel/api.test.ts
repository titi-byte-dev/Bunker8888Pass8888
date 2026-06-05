import { describe, expect, it } from "vitest";
import { parseSentinelResponse } from "./api";

describe("sentinel api (DW-004)", () => {
  it("parseSentinelResponse detecta code sentinel_step_up", () => {
    const err = parseSentinelResponse(403, {
      code: "sentinel_step_up",
      challenge_id: "abc",
      reason: "impossible_travel",
      detail: "~9500 km",
    });
    expect(err?.stepUp.challengeId).toBe("abc");
  });

  it("parseSentinelResponse ignora outros 403", () => {
    expect(parseSentinelResponse(403, { error: "fora do turno" })).toBeNull();
    expect(parseSentinelResponse(401, {})).toBeNull();
  });
});
