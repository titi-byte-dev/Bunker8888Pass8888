import { describe, expect, it } from "vitest";
import { deriveMasterKeyBytes, DEV_CLIENT_KDF } from "../auth";
import { randomBytes } from "../crypto";
import {
  generateRecoveryCode,
  normalizeRecoveryCode,
  unwrapMasterKeyBytes,
  wrapMasterKeyBytes,
} from "./recovery";

describe("recovery key (VAULT-018)", () => {
  it("generateRecoveryCode produz grupos legíveis", () => {
    const code = generateRecoveryCode();
    expect(code).toMatch(/^[A-Z2-9]{5}(-[A-Z2-9]{5}){3}$/);
  });

  it("wrap/unwrap round-trip da Master Key", async () => {
    const mk = await deriveMasterKeyBytes("master-pw", randomBytes(16), DEV_CLIENT_KDF);
    const code = generateRecoveryCode();
    const envelope = await wrapMasterKeyBytes(mk, code);
    const recovered = await unwrapMasterKeyBytes(envelope, code);
    expect(Array.from(recovered)).toEqual(Array.from(mk));
  });

  it("código errado falha ao descifrar", async () => {
    const mk = await deriveMasterKeyBytes("x", randomBytes(16), DEV_CLIENT_KDF);
    const envelope = await wrapMasterKeyBytes(mk, generateRecoveryCode());
    await expect(unwrapMasterKeyBytes(envelope, "WRONG-CODE-XXXXX-YYYYY-ZZZZZ")).rejects.toThrow();
  });

  it("normalizeRecoveryCode remove hífens", () => {
    expect(normalizeRecoveryCode("abcd-efgh")).toBe("ABCDEFGH");
  });
});
