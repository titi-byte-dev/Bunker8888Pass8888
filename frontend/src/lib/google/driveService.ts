/**
 * Drive ZK unificado (GOOGLE-002) — backend servidor ou localStorage (offline dev).
 */
import { blobFromBase64, blobToBase64 } from "$lib/vault/items";
import { loadSessionToken } from "$lib/session";
import {
  decryptDriveFile,
  encryptDriveFile,
  opaquePreviewB64,
  type DriveFilePayload,
} from "./driveCrypto";

export type DriveBackend = "server" | "local";

export interface DriveFileView {
  id: string;
  name: string;
  uploadedAt: string;
  cipherB64: string;
  backend: DriveBackend;
}

const LOCAL_KEY = "aegis_google_drive_dev";

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessão expirada — inicia sessão de novo.");
  return t;
}

async function authed(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token()}`,
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return res;
}

function loadLocal(): DriveFileView[] {
  try {
    const raw = localStorage.getItem(LOCAL_KEY);
    if (!raw) return [];
    const parsed = JSON.parse(raw) as { files: DriveFileView[] };
    return parsed.files ?? [];
  } catch {
    return [];
  }
}

function saveLocal(files: DriveFileView[]): void {
  localStorage.setItem(LOCAL_KEY, JSON.stringify({ files }));
}

/** Preferência: servidor se houver sessão; senão localStorage. */
export function preferredDriveBackend(): DriveBackend {
  return loadSessionToken() ? "server" : "local";
}

export async function listDriveFiles(backend: DriveBackend): Promise<DriveFileView[]> {
  if (backend === "local") return loadLocal();
  const j = (await (await authed("/api/work/google/drive/files")).json()) as {
    files: { id: string; blob: string; created_at: string }[];
  };
  return (j.files ?? []).map((f) => ({
    id: f.id,
    name: "(cifrado)",
    uploadedAt: f.created_at,
    cipherB64: f.blob,
    backend: "server" as const,
  }));
}

/** Resolve nomes após decifrar (só para listagem na UI). */
export async function enrichDriveNames(
  masterKey: CryptoKey,
  files: DriveFileView[],
): Promise<DriveFileView[]> {
  const enriched: DriveFileView[] = [];
  for (const f of files) {
    try {
      const payload = await decryptDriveFile(masterKey, blobFromBase64(f.cipherB64));
      enriched.push({ ...f, name: payload.name });
    } catch {
      enriched.push(f);
    }
  }
  return enriched;
}

export async function uploadDriveFile(
  masterKey: CryptoKey,
  name: string,
  content: string,
  backend: DriveBackend,
): Promise<DriveFileView> {
  const blob = await encryptDriveFile(masterKey, name, content);
  const cipherB64 = blobToBase64(blob);
  if (backend === "local") {
    const entry: DriveFileView = {
      id: crypto.randomUUID(),
      name: name.trim() || "ficheiro.txt",
      uploadedAt: new Date().toISOString(),
      cipherB64,
      backend: "local",
    };
    const files = loadLocal();
    files.unshift(entry);
    saveLocal(files);
    return entry;
  }
  const created = (await (
    await authed("/api/work/google/drive/files", {
      method: "POST",
      body: JSON.stringify({ blob: cipherB64 }),
    })
  ).json()) as { id: string; blob: string; created_at: string };
  return {
    id: created.id,
    name: name.trim() || "ficheiro.txt",
    uploadedAt: created.created_at,
    cipherB64: created.blob,
    backend: "server",
  };
}

export async function openDriveFile(masterKey: CryptoKey, file: DriveFileView): Promise<DriveFilePayload> {
  return decryptDriveFile(masterKey, blobFromBase64(file.cipherB64));
}

export function opaquePreview(file: DriveFileView): string {
  return opaquePreviewB64(file.cipherB64);
}

export async function deleteDriveFile(id: string, backend: DriveBackend): Promise<void> {
  if (backend === "local") {
    saveLocal(loadLocal().filter((f) => f.id !== id));
    return;
  }
  const res = await authed(`/api/work/google/drive/files/${id}`, { method: "DELETE" });
  if (res.status !== 204) throw new Error(`delete falhou: ${res.status}`);
}
