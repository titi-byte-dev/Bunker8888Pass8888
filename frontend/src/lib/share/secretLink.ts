/**
 * Cripto dos Secret Links efémeros (SHARE-003).
 *
 * Didático — a chave NUNCA chega ao servidor:
 *
 *   1. Geramos uma chave aleatória e ciframos o segredo com ela (AES-GCM).
 *   2. O ciphertext vai para o servidor (que o guarda só em RAM, com TTL).
 *   3. A chave vai no FRAGMENTO do URL (depois do #). O browser não envia o
 *      fragmento nos pedidos HTTP — por isso o servidor nunca a vê.
 *
 *        link = https://host/s/{id}#k=<chave-base64url>
 *                                  └──────────┬─────────┘
 *                         fica no browser; o servidor recebe só {id}
 *
 * Quem abre o link lê a chave do fragmento, pede o ciphertext (uso único) e
 * decifra no dispositivo. Zero-Knowledge ponta-a-ponta.
 */
import { decrypt, encrypt, fromBytes, randomBytes, toBytes, type Bytes } from "../crypto";
import { base64ToBytes, bytesToBase64 } from "../auth";

/** Comprimento da chave do link: 32 bytes = AES-256. */
const LINK_KEY_BYTES = 32;

/** Gera uma chave de link aleatória. */
export function generateLinkKey(): Bytes {
  return randomBytes(LINK_KEY_BYTES);
}

async function importLinkKey(raw: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM" }, false, ["encrypt", "decrypt"]);
}

/** Cifra o segredo (texto) com a chave do link. Devolve base64 (nonce||ct). */
export async function encryptSecret(key: Bytes, secret: string): Promise<string> {
  const k = await importLinkKey(key);
  return bytesToBase64(await encrypt(k, toBytes(secret)));
}

/** Reverte encryptSecret. Lança se a chave estiver errada (auth GCM). */
export async function decryptSecret(key: Bytes, ciphertextB64: string): Promise<string> {
  const k = await importLinkKey(key);
  return fromBytes(await decrypt(k, base64ToBytes(ciphertextB64)));
}

/** Codifica bytes em base64url (sem padding) — seguro em fragmentos de URL. */
export function toBase64Url(b: Bytes): string {
  return bytesToBase64(b).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

/** Descodifica base64url de volta a bytes. */
export function fromBase64Url(s: string): Bytes {
  const pad = s.length % 4 === 0 ? "" : "=".repeat(4 - (s.length % 4));
  return base64ToBytes(s.replace(/-/g, "+").replace(/_/g, "/") + pad);
}

/** Constrói o link completo, com a chave no fragmento. */
export function buildSecretLink(origin: string, id: string, key: Bytes): string {
  return `${origin}/s/${id}#k=${toBase64Url(key)}`;
}

/** Extrai a chave do fragmento (`#k=...`). Devolve null se ausente/inválida. */
export function keyFromFragment(hash: string): Bytes | null {
  const match = /[#&]k=([A-Za-z0-9\-_]+)/.exec(hash);
  if (!match) return null;
  try {
    const key = fromBase64Url(match[1]!);
    return key.length === LINK_KEY_BYTES ? key : null;
  } catch {
    return null;
  }
}
