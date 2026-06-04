import { describe, expect, it } from "vitest";
import {
  SANDBOX_FILL_MESSAGE,
  isSandboxFillPayload,
  isSandboxReadyPayload,
} from "./protocol";

describe("sandbox protocol (VAULT-013)", () => {
  it("valida payload de fill", () => {
    expect(
      isSandboxFillPayload({
        type: SANDBOX_FILL_MESSAGE,
        username: "dev",
        password: "secret",
      }),
    ).toBe(true);
    expect(isSandboxFillPayload({ type: "evil" })).toBe(false);
  });

  it("valida payload ready", () => {
    expect(isSandboxReadyPayload({ type: "aegis:sandbox:ready" })).toBe(true);
  });
});
