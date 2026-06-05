/**
 * Cifragem ZK de ficheiros Drive (GOOGLE-002).
 * O blob enviado ao servidor (ou à Google) é opaco — nome e conteúdo só no cliente.
 */
import { decrypt, encrypt, fromBytes, toBytes, type Bytes } from "$lib/crypto";

export interface DriveFilePayload {
  name: string;
  content: string;
}

/** Cifra nome + conteúdo num único blob AES-GCM. */
export async function encryptDriveFile(
  masterKey: CryptoKey,
  name: string,
  content: string,
): Promise<Bytes> {
  const payload: DriveFilePayload = { name: name.trim() || "ficheiro.txt", content };
  return encrypt(masterKey, toBytes(JSON.stringify(payload)));
}

export async function decryptDriveFile(masterKey: CryptoKey, blob: Bytes): Promise<DriveFilePayload> {
  return JSON.parse(fromBytes(await decrypt(masterKey, blob))) as DriveFilePayload;
}

/** Pré-visualização opaca (o que um atacante na Google veria). */
export function opaquePreviewB64(b64: string, max = 48): string {
  return b64.length <= max ? b64 : `${b64.slice(0, max)}…`;
}
