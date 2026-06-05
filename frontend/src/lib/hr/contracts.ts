/**
 * Cifragem de contratos ficheiro-a-ficheiro (HR-005).
 *
 * Cada contrato tem a SUA chave aleatoria (file_key): cifra os metadados e os
 * bytes do ficheiro; a file_key e embrulhada com a Master Key. Mesmo padrao da
 * HR-001, mas para ficheiros. O servidor so ve blobs opacos.
 */
import { decrypt, encrypt, fromBytes, randomBytes, toBytes, type Bytes } from "$lib/crypto";

/** Limite do ficheiro em claro (~5 MiB menos o overhead do AES-GCM). */
export const MAX_CONTRACT_BYTES = 5 * 1024 * 1024 - 28;

export interface ContractMeta {
  name: string;
  mime: string;
  size: number;
}

export interface EncryptedContract {
  metaBlob: Bytes;
  dataBlob: Bytes;
  wrappedKey: Bytes;
}

async function importFileKey(raw: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM", length: 256 }, false, [
    "encrypt",
    "decrypt",
  ]);
}

/** Cifra um ficheiro: gera file_key, cifra meta+bytes, embrulha a chave. */
export async function encryptContract(
  masterKey: CryptoKey,
  file: File,
): Promise<EncryptedContract> {
  const bytes = new Uint8Array(await file.arrayBuffer()) as Bytes;
  const meta: ContractMeta = { name: file.name, mime: file.type || "application/octet-stream", size: bytes.length };

  const rawFileKey = randomBytes(32);
  const fileKey = await importFileKey(rawFileKey);
  const metaBlob = await encrypt(fileKey, toBytes(JSON.stringify(meta)));
  const dataBlob = await encrypt(fileKey, bytes);
  const wrappedKey = await encrypt(masterKey, rawFileKey);
  return { metaBlob, dataBlob, wrappedKey };
}

/** Decifra os metadados de um contrato (nome/tipo/tamanho). */
export async function decryptContractMeta(
  masterKey: CryptoKey,
  metaBlob: Bytes,
  wrappedKey: Bytes,
): Promise<ContractMeta> {
  const rawFileKey = (await decrypt(masterKey, wrappedKey)) as Bytes;
  const fileKey = await importFileKey(rawFileKey);
  return JSON.parse(fromBytes(await decrypt(fileKey, metaBlob))) as ContractMeta;
}

/** Decifra os bytes de um contrato para download. */
export async function decryptContractData(
  masterKey: CryptoKey,
  dataBlob: Bytes,
  wrappedKey: Bytes,
): Promise<Bytes> {
  const rawFileKey = (await decrypt(masterKey, wrappedKey)) as Bytes;
  const fileKey = await importFileKey(rawFileKey);
  return (await decrypt(fileKey, dataBlob)) as Bytes;
}
