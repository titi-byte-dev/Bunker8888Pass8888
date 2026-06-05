/**
 * Assinatura digital de contratos (HR-006) — ECDSA P-256.
 *
 * O utilizador tem uma IDENTIDADE DE ASSINATURA: um par de chaves ECDSA. A
 * chave publica fica em claro no servidor (qualquer um verifica); a privada e
 * exportada em PKCS8 e embrulhada com a Master Key (so o dono assina).
 *
 * Assina-se diretamente os bytes do CIPHERTEXT do contrato (data_blob): a
 * verificacao nao precisa de decifrar nada (Zero-Knowledge) e prova que o
 * signatario se comprometeu com exactamente aqueles bytes.
 */
import { decrypt, encrypt, type Bytes } from "$lib/crypto";

const ALGO = { name: "ECDSA", namedCurve: "P-256" } as const;
const SIGN_ALGO = { name: "ECDSA", hash: "SHA-256" } as const;

/** Gera um novo par de chaves de assinatura (extraivel para exportar/embrulhar). */
async function generateKeyPair(): Promise<CryptoKeyPair> {
  return crypto.subtle.generateKey(ALGO, true, ["sign", "verify"]);
}

export interface NewIdentity {
  publicKey: Bytes; // SPKI
  wrappedPrivateKey: Bytes; // PKCS8 cifrado com a Master Key
  privateKey: CryptoKey;
}

/** Cria uma identidade nova e embrulha a chave privada com a Master Key. */
export async function createSigningIdentity(masterKey: CryptoKey): Promise<NewIdentity> {
  const pair = await generateKeyPair();
  const publicKey = new Uint8Array(await crypto.subtle.exportKey("spki", pair.publicKey)) as Bytes;
  const pkcs8 = new Uint8Array(await crypto.subtle.exportKey("pkcs8", pair.privateKey)) as Bytes;
  const wrappedPrivateKey = await encrypt(masterKey, pkcs8);
  return { publicKey, wrappedPrivateKey, privateKey: pair.privateKey };
}

/** Desembrulha a chave privada (PKCS8 cifrado) com a Master Key e importa-a. */
export async function unwrapPrivateKey(
  masterKey: CryptoKey,
  wrappedPrivateKey: Bytes,
): Promise<CryptoKey> {
  const pkcs8 = (await decrypt(masterKey, wrappedPrivateKey)) as Bytes;
  return crypto.subtle.importKey("pkcs8", pkcs8, ALGO, false, ["sign"]);
}

/** Importa uma chave publica (SPKI) para verificacao. */
export async function importPublicKey(spki: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("spki", spki, ALGO, false, ["verify"]);
}

/** Assina os bytes (data_blob do contrato) com a chave privada. */
export async function signBytes(privateKey: CryptoKey, data: Bytes): Promise<Bytes> {
  return new Uint8Array(await crypto.subtle.sign(SIGN_ALGO, privateKey, data)) as Bytes;
}

/** Verifica uma assinatura sobre os bytes, com a chave publica do signatario. */
export async function verifyBytes(
  publicKey: CryptoKey,
  signature: Bytes,
  data: Bytes,
): Promise<boolean> {
  return crypto.subtle.verify(SIGN_ALGO, publicKey, signature, data);
}

/** sha256(bytes) em hex — usado como content_digest legivel do contrato. */
export async function sha256hex(data: Bytes): Promise<string> {
  const d = await crypto.subtle.digest("SHA-256", data);
  return Array.from(new Uint8Array(d))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}
