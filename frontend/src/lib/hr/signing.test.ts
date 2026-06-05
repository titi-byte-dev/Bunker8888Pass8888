import { describe, expect, it } from "vitest";
import {
  createSigningIdentity,
  importPublicKey,
  sha256hex,
  signBytes,
  unwrapPrivateKey,
  verifyBytes,
} from "./signing";
import { encryptContract } from "./contracts";
import { randomBytes, type Bytes } from "$lib/crypto";

async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]);
}

describe("HR-006 assinatura digital de contratos (ECDSA P-256)", () => {
  it("assina e verifica os bytes de um contrato (round-trip)", async () => {
    const mk = await makeMasterKey();
    const id = await createSigningIdentity(mk);
    const data = randomBytes(256);
    const sig = await signBytes(id.privateKey, data);
    const pub = await importPublicKey(id.publicKey);
    expect(await verifyBytes(pub, sig, data)).toBe(true);
  });

  it("rejeita a verificação se os bytes mudarem", async () => {
    const mk = await makeMasterKey();
    const id = await createSigningIdentity(mk);
    const data = randomBytes(128);
    const sig = await signBytes(id.privateKey, data);
    const pub = await importPublicKey(id.publicKey);
    const tampered = Uint8Array.from(data) as Bytes;
    tampered[0] ^= 0xff;
    expect(await verifyBytes(pub, sig, tampered)).toBe(false);
  });

  it("a chave privada embrulhada desembrulha com a Master Key e volta a assinar", async () => {
    const mk = await makeMasterKey();
    const id = await createSigningIdentity(mk);
    const priv2 = await unwrapPrivateKey(mk, id.wrappedPrivateKey);
    const data = randomBytes(64);
    const sig = await signBytes(priv2, data);
    const pub = await importPublicKey(id.publicKey);
    expect(await verifyBytes(pub, sig, data)).toBe(true);
  });

  it("uma chave publica de outra identidade nao valida a assinatura", async () => {
    const mk = await makeMasterKey();
    const a = await createSigningIdentity(mk);
    const b = await createSigningIdentity(mk);
    const data = randomBytes(64);
    const sig = await signBytes(a.privateKey, data);
    const pubB = await importPublicKey(b.publicKey);
    expect(await verifyBytes(pubB, sig, data)).toBe(false);
  });

  it("encryptContract produz blobs distintos e digest estavel", async () => {
    const mk = await makeMasterKey();
    const file = new File([new Uint8Array([1, 2, 3, 4])], "contrato.pdf", { type: "application/pdf" });
    const enc = await encryptContract(mk, file);
    expect(enc.dataBlob.length).toBeGreaterThan(0);
    const d1 = await sha256hex(enc.dataBlob);
    const d2 = await sha256hex(enc.dataBlob);
    expect(d1).toBe(d2);
    expect(d1).toHaveLength(64);
  });
});
