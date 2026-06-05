/**
 * Par de chaves assimétricas por utilizador (SHARE-001) — fundação da partilha
 * em Zero-Knowledge.
 *
 * Didático: para partilhar um segredo sem o servidor o ver, re-cifra-se a chave
 * do item para a **chave pública** do destinatário (RSA-OAEP). Só ele, com a sua
 * **chave privada**, a consegue abrir. O servidor encaminha bytes opacos.
 *
 * Modelo de guarda:
 *   • Chave pública  → exportada em SPKI, guardada em CLARO (partilhável).
 *   • Chave privada  → exportada em PKCS8 e cifrada com a Master Key do dono
 *                      (AES-GCM, via `encrypt`), antes de ir para o servidor.
 *
 * Tudo corre NO DISPOSITIVO. A chave privada só existe em claro durante a sessão
 * desbloqueada; o servidor nunca a vê.
 */
import { decrypt, encrypt, type Bytes } from "../crypto";
import { base64ToBytes, bytesToBase64 } from "../auth";

/** Identificador do esquema — versionado para futura rotação de parâmetros. */
export const SHARE_KEY_ALGORITHM = "RSA-OAEP-3072-SHA256";

/** Parâmetros RSA-OAEP. 3072 bits = equilíbrio segurança/desempenho no browser. */
const RSA_PARAMS: RsaHashedKeyGenParams = {
  name: "RSA-OAEP",
  modulusLength: 3072,
  publicExponent: new Uint8Array([0x01, 0x00, 0x01]), // 65537
  hash: "SHA-256",
};

function toBytes(buf: ArrayBuffer): Bytes {
  return new Uint8Array(buf) as Bytes;
}

/**
 * Gera um novo par de chaves de partilha. `extractable: true` é necessário para
 * exportar (a pública para partilhar, a privada para a cifrar com a Master Key).
 */
export async function generateSharingKeypair(): Promise<CryptoKeyPair> {
  return crypto.subtle.generateKey(RSA_PARAMS, true, ["encrypt", "decrypt"]);
}

/** Exporta a chave pública em SPKI (base64) — formato partilhável e portável. */
export async function exportPublicKey(pub: CryptoKey): Promise<string> {
  return bytesToBase64(toBytes(await crypto.subtle.exportKey("spki", pub)));
}

/** Importa uma chave pública SPKI (base64). Só serve para cifrar (`encrypt`). */
export async function importPublicKey(b64: string): Promise<CryptoKey> {
  return crypto.subtle.importKey(
    "spki",
    base64ToBytes(b64),
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["encrypt"],
  );
}

/**
 * Cifra (envolve) a chave privada com a Master Key e devolve base64.
 * Resultado = AES-GCM(Master Key, PKCS8) no formato nonce||ciphertext.
 */
export async function wrapPrivateKey(masterKey: CryptoKey, priv: CryptoKey): Promise<string> {
  const pkcs8 = toBytes(await crypto.subtle.exportKey("pkcs8", priv));
  const sealed = await encrypt(masterKey, pkcs8);
  return bytesToBase64(sealed);
}

/** Reverte `wrapPrivateKey`: decifra com a Master Key e importa a chave privada. */
export async function unwrapPrivateKey(masterKey: CryptoKey, b64: string): Promise<CryptoKey> {
  const pkcs8 = await decrypt(masterKey, base64ToBytes(b64));
  return crypto.subtle.importKey(
    "pkcs8",
    pkcs8 as Bytes,
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["decrypt"],
  );
}

/**
 * Envolve uma chave simétrica (ex.: a chave de um item do cofre) para um
 * destinatário, usando a sua chave pública. É a primitiva que a SHARE-002 usa
 * para partilhar itens. Devolve base64.
 */
export async function wrapKeyForRecipient(recipientPub: CryptoKey, rawKey: Bytes): Promise<string> {
  const sealed = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, recipientPub, rawKey);
  return bytesToBase64(toBytes(sealed));
}

/** Reverte `wrapKeyForRecipient` com a minha chave privada. Devolve os bytes. */
export async function unwrapKeyFromSender(myPriv: CryptoKey, b64: string): Promise<Bytes> {
  const raw = await crypto.subtle.decrypt({ name: "RSA-OAEP" }, myPriv, base64ToBytes(b64));
  return toBytes(raw);
}

/**
 * Impressão digital (SHA-256) de uma chave pública SPKI, em hex agrupado.
 *
 * Didático: como as "safety numbers" do Signal ou o fingerprint de uma chave
 * SSH — permite confirmar fora de banda que a chave pública é mesmo da pessoa
 * certa, defendendo contra troca maliciosa de chaves (MITM).
 */
export async function publicKeyFingerprint(b64Spki: string): Promise<string> {
  const digest = toBytes(await crypto.subtle.digest("SHA-256", base64ToBytes(b64Spki)));
  const hex = Array.from(digest, (b) => b.toString(16).padStart(2, "0")).join("");
  return (hex.match(/.{4}/g) ?? [hex]).join(" ").toUpperCase();
}
