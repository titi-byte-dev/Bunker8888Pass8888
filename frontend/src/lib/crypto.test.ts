import { describe, it, expect } from "vitest";
import {
  deriveMasterKey,
  encrypt,
  decrypt,
  randomBytes,
  toBytes,
  fromBytes,
} from "./crypto";

// Parâmetros leves para os testes correrem depressa (em produção usa-se DEFAULT_KDF).
const fastKdf = { iterations: 1_000, hash: "SHA-256" } as const;

describe("crypto (Zero-Knowledge client-side)", () => {
  it("deriva a mesma chave para a mesma password+salt (determinista)", async () => {
    const salt = toBytes("salt-fixo-16-byte");
    const k1 = await deriveMasterKey("password", salt, fastKdf);
    const k2 = await deriveMasterKey("password", salt, fastKdf);

    // CryptoKey não é exportável; verificamos a igualdade pela via funcional:
    // o que k1 cifra, k2 tem de conseguir decifrar.
    const blob = await encrypt(k1, toBytes("segredo"));
    const out = await decrypt(k2, blob);
    expect(fromBytes(out)).toBe("segredo");
  });

  it("faz round-trip encrypt -> decrypt", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const blob = await encrypt(key, toBytes("ola mundo"));
    const out = await decrypt(key, blob);
    expect(fromBytes(out)).toBe("ola mundo");
  });

  it("deteta adulteração (autenticação GCM)", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const blob = await encrypt(key, toBytes("mensagem"));
    blob[blob.length - 1] ^= 0x01; // altera um byte
    await expect(decrypt(key, blob)).rejects.toThrow();
  });

  it("usa nonce diferente a cada cifragem", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const a = await encrypt(key, toBytes("igual"));
    const b = await encrypt(key, toBytes("igual"));
    expect(Buffer.from(a).equals(Buffer.from(b))).toBe(false);
  });

  it("respeita o AAD (additional authenticated data)", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const blob = await encrypt(key, toBytes("x"), toBytes("aad-1"));
    await expect(decrypt(key, blob, toBytes("aad-2"))).rejects.toThrow();
  });
});
