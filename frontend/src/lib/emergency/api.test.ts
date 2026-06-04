import { describe, expect, it } from "vitest";
import { secondsUntil } from "./api";

describe("emergency api", () => {
  it("secondsUntil devolve 0 após unlocks_at", () => {
    const past = new Date(Date.now() - 60_000).toISOString();
    expect(secondsUntil(past)).toBe(0);
  });

  it("secondsUntil arredonda para cima", () => {
    const future = new Date(Date.now() + 2500).toISOString();
    expect(secondsUntil(future)).toBe(3);
  });
});
