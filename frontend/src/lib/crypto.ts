/**
 * Cripto client-side (Zero-Knowledge) do AegisPass.
 *
 * Tudo aqui corre NO DISPOSITIVO do utilizador. A Master Key derivada nunca sai
 * do cliente — o servidor só vê blobs cifrados e o auth hash.
 *
 * Usamos a WebCrypto API nativa (crypto.subtle), nunca cripto "caseira".
 * Para a KDF usamos PBKDF2 (nativo na WebCrypto). Argon2id (via WASM) é uma
 * melhoria futura — mais resistente, mas requer uma biblioteca extra.
 */

const encoder = new TextEncoder();

export interface KdfParams {
  /** Nº de iterações do PBKDF2. Mais = mais lento = mais seguro. */
  iterations: number;
  hash: "SHA-256";
}

/** Parâmetros por omissão. Rever periodicamente à medida que o hardware evolui. */
export const DEFAULT_KDF: KdfParams = { iterations: 600_000, hash: "SHA-256" };

/** NONCE_BYTES = 12 é o tamanho recomendado para o nonce (IV) do AES-GCM. */
const NONCE_BYTES = 12;

/** Gera `n` bytes aleatórios criptograficamente seguros. */
export function randomBytes(n: number): Uint8Array {
  return crypto.getRandomValues(new Uint8Array(n));
}

/**
 * Deriva a Master Key (CryptoKey AES-GCM 256) a partir da Master Password.
 *
 * `extractable: false` impede que a chave seja exportada do objeto CryptoKey,
 * reduzindo o risco de ela acabar acidentalmente em logs, rede ou disco.
 */
export async function deriveMasterKey(
  password: string,
  salt: Uint8Array,
  params: KdfParams = DEFAULT_KDF,
): Promise<CryptoKey> {
  // Importamos a password como material base para a KDF (só permite deriveKey).
  const baseKey = await crypto.subtle.importKey(
    "raw",
    encoder.encode(password),
    "PBKDF2",
    false,
    ["deriveKey"],
  );

  return crypto.subtle.deriveKey(
    { name: "PBKDF2", salt, iterations: params.iterations, hash: params.hash },
    baseKey,
    { name: "AES-GCM", length: 256 },
    false, // não exportável
    ["encrypt", "decrypt"],
  );
}

/**
 * Cifra `plaintext` com AES-GCM. Devolve nonce||ciphertext (com a tag incluída
 * no fim do ciphertext pela WebCrypto), pronto a guardar/enviar.
 */
export async function encrypt(
  key: CryptoKey,
  plaintext: Uint8Array,
  aad?: Uint8Array,
): Promise<Uint8Array> {
  const nonce = randomBytes(NONCE_BYTES);
  const ciphertext = new Uint8Array(
    await crypto.subtle.encrypt(
      { name: "AES-GCM", iv: nonce, additionalData: aad },
      key,
      plaintext,
    ),
  );

  // Prefixamos o nonce ao ciphertext: nonce || ciphertext+tag.
  const out = new Uint8Array(nonce.length + ciphertext.length);
  out.set(nonce, 0);
  out.set(ciphertext, nonce.length);
  return out;
}

/** Reverte encrypt(). Lança se a autenticação falhar (dados adulterados). */
export async function decrypt(
  key: CryptoKey,
  blob: Uint8Array,
  aad?: Uint8Array,
): Promise<Uint8Array> {
  const nonce = blob.subarray(0, NONCE_BYTES);
  const ciphertext = blob.subarray(NONCE_BYTES);
  const plaintext = await crypto.subtle.decrypt(
    { name: "AES-GCM", iv: nonce, additionalData: aad },
    key,
    ciphertext,
  );
  return new Uint8Array(plaintext);
}

/** Converte texto em bytes (UTF-8) — atalho conveniente para o chamador. */
export function toBytes(s: string): Uint8Array {
  return encoder.encode(s);
}

/** Converte bytes (UTF-8) em texto. */
export function fromBytes(b: Uint8Array): string {
  return new TextDecoder().decode(b);
}
