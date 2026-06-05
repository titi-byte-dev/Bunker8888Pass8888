/**
 * Cripto das Notas Auto-Destrutivas (SHARE-005).
 *
 * Assenta no padrão da SHARE-003 (chave no fragmento do URL), com uma 2.ª camada
 * OPCIONAL por passphrase:
 *
 *   1. Geramos uma chave aleatória (noteKey) e ciframos a nota com ela (AES-GCM).
 *   2. Se houver passphrase, ciframos OUTRA vez com uma chave derivada dela
 *      (PBKDF2). Assim, mesmo quem intercete o link inteiro não abre a nota sem
 *      a passphrase — que se combina por um canal à parte (ex.: dito ao telefone).
 *   3. A noteKey (e o salt, se aplicável) viajam no FRAGMENTO do URL (#), que o
 *      browser nunca envia. O servidor guarda bytes opacos e não sabe sequer se
 *      existe passphrase.
 *
 *        link = /n/{id}#k=<noteKey>            (sem passphrase)
 *        link = /n/{id}#k=<noteKey>&p=1&s=<salt>   (com passphrase)
 *                       └──────────────┬───────────┘
 *                        fica no browser; o servidor recebe só {id}
 */
import { decrypt, encrypt, fromBytes, randomBytes, toBytes, type Bytes } from "../crypto";
import { base64ToBytes, bytesToBase64 } from "../auth";
import { fromBase64Url, toBase64Url } from "./secretLink";

/** Comprimento da chave da nota: 32 bytes = AES-256. */
const NOTE_KEY_BYTES = 32;
/** Bytes de salt da derivação por passphrase. */
const SALT_BYTES = 16;
/** Iterações do PBKDF2 da passphrase (defesa contra força bruta offline). */
const PASSPHRASE_ITERATIONS = 200_000;

/** Gera uma chave de nota aleatória. */
export function generateNoteKey(): Bytes {
  return randomBytes(NOTE_KEY_BYTES);
}

async function importNoteKey(raw: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM" }, false, ["encrypt", "decrypt"]);
}

/** Deriva uma chave AES-GCM a partir da passphrase + salt (PBKDF2-SHA256). */
async function derivePassphraseKey(passphrase: string, salt: Bytes): Promise<CryptoKey> {
  const base = await crypto.subtle.importKey("raw", toBytes(passphrase), "PBKDF2", false, [
    "deriveKey",
  ]);
  return crypto.subtle.deriveKey(
    { name: "PBKDF2", salt, iterations: PASSPHRASE_ITERATIONS, hash: "SHA-256" },
    base,
    { name: "AES-GCM", length: 256 },
    false,
    ["encrypt", "decrypt"],
  );
}

/** Resultado da cifra de uma nota: ciphertext (base64) + salt (se houve passphrase). */
export interface EncryptedNote {
  ciphertext: string;
  salt: Bytes | null;
}

/**
 * Cifra o conteúdo da nota. Com passphrase, aplica a 2.ª camada e devolve também
 * o salt (que vai no fragmento, não para o servidor).
 */
export async function encryptNoteContent(
  noteKey: Bytes,
  secret: string,
  passphrase?: string,
): Promise<EncryptedNote> {
  const inner = await encrypt(await importNoteKey(noteKey), toBytes(secret));
  if (!passphrase) {
    return { ciphertext: bytesToBase64(inner), salt: null };
  }
  const salt = randomBytes(SALT_BYTES);
  const outer = await encrypt(await derivePassphraseKey(passphrase, salt), inner);
  return { ciphertext: bytesToBase64(outer), salt };
}

/**
 * Reverte encryptNoteContent. Se a nota tinha passphrase, `passphrase`+`salt` são
 * obrigatórios. Lança se a passphrase ou a chave estiverem erradas (auth GCM).
 */
export async function decryptNoteContent(
  noteKey: Bytes,
  ciphertext: string,
  passphrase: string | null,
  salt: Bytes | null,
): Promise<string> {
  let bytes = base64ToBytes(ciphertext);
  if (salt) {
    if (!passphrase) throw new Error("Esta nota precisa de passphrase.");
    bytes = await decrypt(await derivePassphraseKey(passphrase, salt), bytes);
  }
  return fromBytes(await decrypt(await importNoteKey(noteKey), bytes));
}

/** Constrói o link completo, com a chave (e salt, se houver) no fragmento. */
export function buildNoteLink(origin: string, id: string, key: Bytes, salt: Bytes | null): string {
  let link = `${origin}/n/${id}#k=${toBase64Url(key)}`;
  if (salt) link += `&p=1&s=${toBase64Url(salt)}`;
  return link;
}

/** Dados extraídos do fragmento de um link de nota. */
export interface NoteFragment {
  key: Bytes;
  requiresPassphrase: boolean;
  salt: Bytes | null;
}

/** Extrai a chave (e salt) do fragmento. Devolve null se ausente/inválida. */
export function parseNoteFragment(hash: string): NoteFragment | null {
  const km = /[#&]k=([A-Za-z0-9\-_]+)/.exec(hash);
  if (!km) return null;
  let key: Bytes;
  try {
    key = fromBase64Url(km[1]!);
  } catch {
    return null;
  }
  if (key.length !== NOTE_KEY_BYTES) return null;

  const sm = /[#&]s=([A-Za-z0-9\-_]+)/.exec(hash);
  const requiresPassphrase = /[#&]p=1/.test(hash) || sm !== null;
  let salt: Bytes | null = null;
  if (sm) {
    try {
      salt = fromBase64Url(sm[1]!);
    } catch {
      return null;
    }
  }
  return { key, requiresPassphrase, salt };
}
