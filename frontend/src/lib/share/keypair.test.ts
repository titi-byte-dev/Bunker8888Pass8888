import { beforeAll, describe, expect, it } from "vitest";
import { deriveMasterKey, randomBytes, type Bytes } from "../crypto";
import {
  SHARE_KEY_ALGORITHM,
  exportPublicKey,
  generateSharingKeypair,
  importPublicKey,
  publicKeyFingerprint,
  unwrapKeyFromSender,
  unwrapPrivateKey,
  wrapKeyForRecipient,
  wrapPrivateKey,
} from "./keypair";

const fastKdf = { iterations: 1_000, hash: "SHA-256" } as const;

// RSA-3072 keygen é pesado — geramos um par uma vez e reutilizamos.
let pair: CryptoKeyPair;
let masterKey: CryptoKey;

beforeAll(async () => {
  pair = await generateSharingKeypair();
  masterKey = await deriveMasterKey("master-pw", randomBytes(16), fastKdf);
});

describe("SHARE-001 keypair (partilha Zero-Knowledge)", () => {
  it("expõe um identificador de algoritmo versionado", () => {
    expect(SHARE_KEY_ALGORITHM).toBe("RSA-OAEP-3072-SHA256");
  });

  it("envolve e abre a chave privada com a Master Key (round-trip)", async () => {
    const wrapped = await wrapPrivateKey(masterKey, pair.privateKey);
    const reopened = await unwrapPrivateKey(masterKey, wrapped);

    // A privada reaberta tem de decifrar o que a pública original cifrou.
    const secret = randomBytes(32);
    const sealed = await wrapKeyForRecipient(pair.publicKey, secret);
    const out = await unwrapKeyFromSender(reopened, sealed);
    expect([...out]).toEqual([...secret]);
  });

  it("não abre a chave privada com a Master Key errada (auth GCM)", async () => {
    const wrapped = await wrapPrivateKey(masterKey, pair.privateKey);
    const wrongKey = await deriveMasterKey("outra-pw", randomBytes(16), fastKdf);
    await expect(unwrapPrivateKey(wrongKey, wrapped)).rejects.toThrow();
  });

  it("envolve uma chave de item para o destinatário e este abre-a", async () => {
    const itemKey = randomBytes(32) as Bytes;
    const sealed = await wrapKeyForRecipient(pair.publicKey, itemKey);
    expect(sealed).not.toContain(",");
    const opened = await unwrapKeyFromSender(pair.privateKey, sealed);
    expect([...opened]).toEqual([...itemKey]);
  });

  it("exporta/importa a chave pública (SPKI) preservando a função", async () => {
    const spki = await exportPublicKey(pair.publicKey);
    const reimported = await importPublicKey(spki);

    const secret = randomBytes(32);
    const sealed = await wrapKeyForRecipient(reimported, secret);
    const out = await unwrapKeyFromSender(pair.privateKey, sealed);
    expect([...out]).toEqual([...secret]);
  });

  it("calcula uma impressão digital determinista e legível", async () => {
    const spki = await exportPublicKey(pair.publicKey);
    const fp1 = await publicKeyFingerprint(spki);
    const fp2 = await publicKeyFingerprint(spki);
    expect(fp1).toBe(fp2);
    // Hex em grupos de 4, maiúsculas (SHA-256 = 64 hex chars = 16 grupos).
    expect(fp1).toMatch(/^[0-9A-F]{4}( [0-9A-F]{4})+$/);
    expect(fp1.replace(/ /g, "")).toHaveLength(64);
  });
});
