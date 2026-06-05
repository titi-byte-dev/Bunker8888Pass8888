/**
 * Compatibilidade com código que usava localStorage apenas (GOOGLE-002 dev).
 * Novos fluxos devem usar `driveService.ts` directamente.
 */
import {
  deleteDriveFile,
  enrichDriveNames,
  listDriveFiles,
  openDriveFile,
  opaquePreview,
  uploadDriveFile,
  type DriveFileView,
} from "./driveService";

export type DriveDevFile = DriveFileView;

export async function uploadToMockDrive(
  masterKey: CryptoKey,
  name: string,
  content: string,
): Promise<DriveDevFile> {
  return uploadDriveFile(masterKey, name, content, "local");
}

export function listMockDriveFiles(): DriveDevFile[] {
  // Síncrono legado — preferir listDriveFiles async no UI novo.
  try {
    const raw = localStorage.getItem("aegis_google_drive_dev");
    if (!raw) return [];
    return (JSON.parse(raw) as { files: DriveDevFile[] }).files ?? [];
  } catch {
    return [];
  }
}

export async function openMockDriveFile(masterKey: CryptoKey, file: DriveDevFile): Promise<string> {
  const payload = await openDriveFile(masterKey, file);
  return payload.content;
}

export { opaquePreview, enrichDriveNames, deleteDriveFile };
