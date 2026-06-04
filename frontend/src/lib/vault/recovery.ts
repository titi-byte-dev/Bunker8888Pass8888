/**
 * Chave de recuperação (VAULT-018).
 *
 * Didático: se perderes a Master Password, a chave de recuperação (código offline)
 * permite obter de volta a Master Key a partir do blob guardado no servidor.
 * O blob = AES-GCM(Master Key, derivada da chave de recuperação).
 *
 * ⚠️ Segurança: guarda o código offline (papel/cofre). Mostramo-lo UMA vez na
 * criação — o servidor nunca o vê em claro.
 */
import { encrypt, decrypt, randomBytes, type Bytes } from "../crypto";
import {
  deriveMasterKeyBytes,
  importAesKeyFromBytes,
  bytesToBase64,
  base64ToBytes,
  type ClientKdfParams,
} from "../auth";

const RECOVERY_VERSION = 1;

/** KDF para a chave de recuperação (custo moderado — uso raro). */
export const RECOVERY_KDF: ClientKdfParams = {
  time: 2,
  memory: 16384,
  threads: 1,
};

/** Alfabeto sem caracteres ambíguos (0/O, 1/I/l). */
const CODE_ALPHABET = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";

export interface RecoveryEnvelope {
  v: number;
  salt: string;
  kdf: ClientKdfParams;
  blob: string;
}

/** Gera código legível (4×5 chars, ~100 bits de entropia). */
export function generateRecoveryCode(): string {
  const groups: string[] = [];
  for (let g = 0; g < 4; g++) {
    let chunk = "";
    for (let i = 0; i < 5; i++) {
      chunk += CODE_ALPHABET[randomBytes(1)[0]! % CODE_ALPHABET.length];
    }
    groups.push(chunk);
  }
  return groups.join("-");
}

/** Normaliza input do utilizador (remove espaços e hífens). */
export function normalizeRecoveryCode(code: string): string {
  return code.toUpperCase().replace(/[\s-]/g, "");
}

/** Cifra bytes da Master Key com a chave de recuperação. */
export async function wrapMasterKeyBytes(
  masterKeyBytes: Uint8Array,
  recoveryCode: string,
): Promise<string> {
  const salt = randomBytes(16);
  const recoveryKeyBytes = await deriveMasterKeyBytes(
    normalizeRecoveryCode(recoveryCode),
    salt,
    RECOVERY_KDF,
  );
  const aesKey = await importAesKeyFromBytes(recoveryKeyBytes);
  const ciphertext = await encrypt(aesKey, masterKeyBytes as Bytes);
  const envelope: RecoveryEnvelope = {
    v: RECOVERY_VERSION,
    salt: bytesToBase64(salt),
    kdf: RECOVERY_KDF,
    blob: bytesToBase64(ciphertext),
  };
  return JSON.stringify(envelope);
}

/** Descifra o blob de recuperação e devolve bytes da Master Key. */
export async function unwrapMasterKeyBytes(
  envelopeJson: string,
  recoveryCode: string,
): Promise<Uint8Array> {
  const envelope = JSON.parse(envelopeJson) as RecoveryEnvelope;
  if (envelope.v !== RECOVERY_VERSION) {
    throw new Error(`versão de recuperação não suportada: ${envelope.v}`);
  }
  const salt = base64ToBytes(envelope.salt);
  const recoveryKeyBytes = await deriveMasterKeyBytes(
    normalizeRecoveryCode(recoveryCode),
    salt,
    envelope.kdf,
  );
  const aesKey = await importAesKeyFromBytes(recoveryKeyBytes);
  const plaintext = await decrypt(aesKey, base64ToBytes(envelope.blob));
  if (plaintext.length !== 32) {
    throw new Error("Master Key inválida após recuperação");
  }
  return plaintext;
}

export async function uploadRecoveryBackup(
  baseURL: string,
  token: string,
  envelopeJson: string,
): Promise<void> {
  const res = await fetch(`${baseURL}/api/vault/recovery-backup`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ blob: bytesToBase64(new TextEncoder().encode(envelopeJson)) }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
}

export async function fetchRecoveryBackupByEmail(
  baseURL: string,
  email: string,
): Promise<string> {
  const res = await fetch(
    `${baseURL}/api/vault/recovery-backup/lookup?email=${encodeURIComponent(email)}`,
  );
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? "Recuperação indisponível");
  }
  const { blob } = (await res.json()) as { blob: string };
  return new TextDecoder().decode(base64ToBytes(blob));
}

export async function fetchRecoveryBackupStatus(
  baseURL: string,
  token: string,
): Promise<boolean> {
  const res = await fetch(`${baseURL}/api/vault/recovery-backup/status`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) return false;
  const { configured } = (await res.json()) as { configured: boolean };
  return configured;
}

/** Recupera Master Key a partir de email + código (sem sessão). */
export async function recoverMasterKeyFromEmail(
  baseURL: string,
  email: string,
  recoveryCode: string,
): Promise<CryptoKey> {
  const envelopeJson = await fetchRecoveryBackupByEmail(baseURL, email);
  const mkBytes = await unwrapMasterKeyBytes(envelopeJson, recoveryCode);
  return importAesKeyFromBytes(mkBytes);
}
