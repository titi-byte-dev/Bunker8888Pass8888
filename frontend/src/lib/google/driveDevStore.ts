/**
 * Simulador de Google Drive cifrado (GOOGLE-002 dev stub).
 * Ficheiros ficam em localStorage como blobs base64 — nunca saem do browser.
 */
import { bytesToBase64, base64ToBytes, encrypt, decrypt, toBytes, fromBytes } from "$lib/crypto";

const STORAGE_KEY = "aegis_google_drive_dev";

export interface DriveDevFile {
  id: string;
  name: string;
  uploadedAt: string;
  /** Conteúdo cifrado (simula o que a Google veria). */
  cipherB64: string;
}

interface Store {
  files: DriveDevFile[];
}

function load(): Store {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return { files: [] };
    return JSON.parse(raw) as Store;
  } catch {
    return { files: [] };
  }
}

function save(store: Store): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
}

/** «Envia» ficheiro para o mock Drive — cifra com Master Key antes de guardar. */
export async function uploadToMockDrive(masterKey: CryptoKey, name: string, content: string): Promise<DriveDevFile> {
  const blob = await encrypt(masterKey, toBytes(content));
  const entry: DriveDevFile = {
    id: crypto.randomUUID(),
    name,
    uploadedAt: new Date().toISOString(),
    cipherB64: bytesToBase64(blob),
  };
  const store = load();
  store.files.unshift(entry);
  save(store);
  return entry;
}

export function listMockDriveFiles(): DriveDevFile[] {
  return load().files;
}

/** Decifra e devolve texto — só com Master Key (vista AegisPass). */
export async function openMockDriveFile(masterKey: CryptoKey, file: DriveDevFile): Promise<string> {
  return fromBytes(await decrypt(masterKey, base64ToBytes(file.cipherB64)));
}

/** O que o «atacante na Google» veria — bytes opacos. */
export function opaquePreview(file: DriveDevFile): string {
  return file.cipherB64.slice(0, 48) + "…";
}
