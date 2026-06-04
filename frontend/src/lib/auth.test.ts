import { describe, expect, it } from "vitest";
import { deriveAuthHashBytes, deriveMasterKeyBytes, DEV_CLIENT_KDF } from "./auth";
import { randomBytes } from "./crypto";

describe("auth Argon2id (alinhado com backend Go)", () => {
  it("Master Key e Auth Hash são derivações distintas", async () => {
    const password = "test-master-password";
    const salt = randomBytes(16);
    const mk = await deriveMasterKeyBytes(password, salt, DEV_CLIENT_KDF);
    const ah = await deriveAuthHashBytes(mk, password, DEV_CLIENT_KDF);
    expect(mk.length).toBe(32);
    expect(ah.length).toBe(32);
    expect(Array.from(mk).join()).not.toBe(Array.from(ah).join());
  });
});
