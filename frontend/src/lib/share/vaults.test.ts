import { beforeAll, describe, expect, it } from "vitest";
import { generateSharingKeypair } from "./keypair";
import {
  SHARED_VAULT_ALGORITHM,
  decryptAttachmentMeta,
  decryptFileBytes,
  decryptVaultItem,
  decryptVaultName,
  encryptAttachmentMeta,
  encryptFileBytes,
  encryptVaultItem,
  encryptVaultName,
  generateVaultKey,
  unwrapVaultKey,
  wrapVaultKeyFor,
} from "./vaults";
import type { Bytes } from "../crypto";

// RSA-3072 é pesado — geramos um par uma vez e reutilizamos.
let pair: CryptoKeyPair;

beforeAll(async () => {
  pair = await generateSharingKeypair();
});

describe("SHARE-002 cofres partilhados (cripto)", () => {
  it("expõe um identificador de esquema versionado", () => {
    expect(SHARED_VAULT_ALGORITHM).toBe("AES-GCM-256+RSA-OAEP-3072");
  });

  it("gera uma chave de cofre de 32 bytes (AES-256)", () => {
    expect(generateVaultKey()).toHaveLength(32);
  });

  it("cifra e decifra o nome do cofre (round-trip)", async () => {
    const key = generateVaultKey();
    const blob = await encryptVaultName(key, "Cofre da Equipa de Segurança");
    expect(blob).not.toContain("Cofre"); // está cifrado
    expect(await decryptVaultName(key, blob)).toBe("Cofre da Equipa de Segurança");
  });

  it("cifra e decifra um item (round-trip)", async () => {
    const key = generateVaultKey();
    const payload = { title: "Router admin", secret: "tr0ub4dor&3" };
    const blob = await encryptVaultItem(key, payload);
    expect(await decryptVaultItem(key, blob)).toEqual(payload);
  });

  it("não decifra com a chave de cofre errada (auth GCM)", async () => {
    const blob = await encryptVaultName(generateVaultKey(), "segredo");
    await expect(decryptVaultName(generateVaultKey(), blob)).rejects.toThrow();
  });

  it("embrulha a chave do cofre para um membro e ele abre-a", async () => {
    const vaultKey = generateVaultKey();
    const wrapped = await wrapVaultKeyFor(pair.publicKey, vaultKey);
    const opened = await unwrapVaultKey(pair.privateKey, wrapped);
    expect([...opened]).toEqual([...vaultKey]);

    // E a chave aberta decifra o que o cofre cifrou — prova ponta-a-ponta.
    const blob = await encryptVaultName(vaultKey, "Engenharia");
    expect(await decryptVaultName(opened, blob)).toBe("Engenharia");
  });
});

describe("SHARE-004 anexos cifrados (cripto)", () => {
  it("cifra e decifra os metadados de um anexo (round-trip)", async () => {
    const key = generateVaultKey();
    const meta = { name: "contrato.pem", mime: "application/x-pem-file", size: 2048 };
    const blob = await encryptAttachmentMeta(key, meta);
    expect(blob).not.toContain("contrato"); // nome do ficheiro está cifrado
    expect(await decryptAttachmentMeta(key, blob)).toEqual(meta);
  });

  it("cifra e decifra os bytes de um ficheiro sem os alterar", async () => {
    const key = generateVaultKey();
    const bytes = new Uint8Array([0, 1, 2, 250, 251, 255, 13, 10]) as Bytes;
    const blob = await encryptFileBytes(key, bytes);
    const back = await decryptFileBytes(key, blob);
    expect([...back]).toEqual([...bytes]);
  });

  it("não decifra os bytes com a chave de cofre errada (auth GCM)", async () => {
    const bytes = new Uint8Array([1, 2, 3]) as Bytes;
    const blob = await encryptFileBytes(generateVaultKey(), bytes);
    await expect(decryptFileBytes(generateVaultKey(), blob)).rejects.toThrow();
  });
});
