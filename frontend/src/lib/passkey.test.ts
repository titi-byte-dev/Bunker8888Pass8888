import { describe, expect, it } from "vitest";
import { _test } from "./passkey";

describe("passkey base64url", () => {
  it("round-trip buffer", () => {
    const raw = new Uint8Array([1, 2, 3, 250]);
    const b64 = _test.bufferToBase64URL(raw.buffer);
    const back = new Uint8Array(_test.base64URLToBuffer(b64));
    expect(back).toEqual(raw);
  });
});
