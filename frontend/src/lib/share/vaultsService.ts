/**
 * Orquestração dos Cofres Partilhados (SHARE-002).
 *
 * Junta três peças: a cripto do cofre (vaults.ts), o cliente HTTP (vaultsApi.ts)
 * e a identidade de partilha do utilizador (setup.ts / keypair.ts da SHARE-001).
 * Exige sessão + Master Key desbloqueada, tal como o cofre pessoal.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { importPublicKey } from "./keypair";
import { ensureShareIdentity, loadOwnPrivateKey, lookupRecipient } from "./setup";
import { SharedVaultsAPI, type VaultMemberDTO, type VaultRole } from "./vaultsApi";
import {
  decryptAttachmentMeta,
  decryptFileBytes,
  decryptVaultItem,
  decryptVaultName,
  encryptAttachmentMeta,
  encryptFileBytes,
  encryptVaultItem,
  encryptVaultName,
  generateVaultKey,
  unwrapVaultKey,
  wrapVaultKeyFor,
  type VaultItemPayload,
} from "./vaults";
import type { Bytes } from "../crypto";

/**
 * Teto do tamanho do ficheiro (5 MiB). O servidor limita o CIPHERTEXT a 5 MiB;
 * o AES-GCM acrescenta 28 bytes (nonce + tag), por isso recusamos o ficheiro em
 * claro um pouco antes, com uma mensagem amigável.
 */
export const MAX_ATTACHMENT_BYTES = 5 * 1024 * 1024 - 28;

export interface DecryptedVault {
  id: string;
  name: string;
  role: VaultRole;
  ownerId: string;
  createdAt: string;
}

export interface DecryptedItem {
  id: string;
  title: string;
  secret: string;
  createdBy: string;
  createdAt: string;
}

/** Anexo com os metadados já decifrados (o nome do ficheiro também é ZK). */
export interface DecryptedAttachment {
  id: string;
  name: string;
  mime: string;
  size: number; // tamanho do ficheiro original
  createdBy: string;
  createdAt: string;
}

/** Cofre aberto: o conteúdo decifrado + a chave do cofre em memória para escrever. */
export interface OpenVault {
  vault: DecryptedVault;
  members: VaultMemberDTO[];
  items: DecryptedItem[];
  attachments: DecryptedAttachment[];
  /** Chave do cofre em claro — só vive nesta sessão, nunca vai para o servidor. */
  vaultKey: Bytes;
}

function requireVaultAccess(): SharedVaultsAPI {
  const token = loadSessionToken();
  const masterKey = getMasterKey();
  if (!token || !masterKey) {
    throw new Error("Partilha bloqueada — inicia sessão e desbloqueia em /auth/unlock");
  }
  return new SharedVaultsAPI("", token);
}

/** Cria um cofre partilhado: gera a chave, cifra o nome, embrulha a chave p/ mim. */
export async function createSharedVault(name: string): Promise<DecryptedVault> {
  const api = requireVaultAccess();
  const identity = await ensureShareIdentity(); // garante que tenho par de chaves
  const myPub = await importPublicKey(identity.publicKey);

  const vaultKey = generateVaultKey();
  const nameBlob = await encryptVaultName(vaultKey, name.trim());
  const ownerWrapped = await wrapVaultKeyFor(myPub, vaultKey);

  const dto = await api.create(nameBlob, ownerWrapped);
  return {
    id: dto.id,
    name: name.trim(),
    role: dto.role,
    ownerId: dto.owner_id,
    createdAt: dto.created_at,
  };
}

/** Lista os meus cofres, decifrando o nome de cada um com a minha chave privada. */
export async function listSharedVaults(): Promise<DecryptedVault[]> {
  const api = requireVaultAccess();
  const myPriv = await loadOwnPrivateKey();
  const dtos = await api.list();

  const out: DecryptedVault[] = [];
  for (const dto of dtos) {
    const vaultKey = await unwrapVaultKey(myPriv, dto.wrapped_vault_key);
    out.push({
      id: dto.id,
      name: await decryptVaultName(vaultKey, dto.name_blob),
      role: dto.role,
      ownerId: dto.owner_id,
      createdAt: dto.created_at,
    });
  }
  return out;
}

/** Abre um cofre: decifra nome + itens e mantém a chave do cofre em memória. */
export async function openSharedVault(id: string): Promise<OpenVault> {
  const api = requireVaultAccess();
  const myPriv = await loadOwnPrivateKey();
  const dto = await api.get(id);
  const vaultKey = await unwrapVaultKey(myPriv, dto.wrapped_vault_key);

  const [members, itemDtos, attDtos] = await Promise.all([
    api.members(id),
    api.items(id),
    api.attachments(id),
  ]);
  const items: DecryptedItem[] = [];
  for (const it of itemDtos) {
    const payload = await decryptVaultItem(vaultKey, it.blob);
    items.push({
      id: it.id,
      title: payload.title,
      secret: payload.secret,
      createdBy: it.created_by,
      createdAt: it.created_at,
    });
  }
  const attachments: DecryptedAttachment[] = [];
  for (const a of attDtos) {
    const meta = await decryptAttachmentMeta(vaultKey, a.meta_blob);
    attachments.push({
      id: a.id,
      name: meta.name,
      mime: meta.mime,
      size: meta.size,
      createdBy: a.created_by,
      createdAt: a.created_at,
    });
  }

  return {
    vault: {
      id: dto.id,
      name: await decryptVaultName(vaultKey, dto.name_blob),
      role: dto.role,
      ownerId: dto.owner_id,
      createdAt: dto.created_at,
    },
    members,
    items,
    attachments,
    vaultKey,
  };
}

export interface InvitedMember {
  email: string;
  fingerprint: string;
}

/**
 * Convida um colega: procura a chave pública dele, embrulha a chave do cofre
 * para essa chave e envia. Devolve a impressão digital para confirmação fora de
 * banda (anti-MITM). Lança se o colega ainda não activou a partilha.
 */
export async function inviteMember(
  vaultId: string,
  vaultKey: Bytes,
  email: string,
  role: VaultRole,
): Promise<InvitedMember> {
  const api = requireVaultAccess();
  const recipient = await lookupRecipient(email);
  const recipientPub = await importPublicKey(recipient.public_key);
  const wrapped = await wrapVaultKeyFor(recipientPub, vaultKey);
  await api.addMember(vaultId, recipient.user_id, role, wrapped);
  return { email: recipient.email, fingerprint: recipient.fingerprint };
}

/** Adiciona um item cifrado ao cofre. */
export async function addVaultItem(
  vaultId: string,
  vaultKey: Bytes,
  payload: VaultItemPayload,
): Promise<void> {
  const api = requireVaultAccess();
  const blob = await encryptVaultItem(vaultKey, payload);
  await api.addItem(vaultId, "note", blob);
}

/**
 * Carrega um ficheiro como anexo cifrado: cifra metadados e bytes com a chave do
 * cofre antes do upload. O servidor só recebe bytes opacos.
 */
export async function uploadAttachment(
  vaultId: string,
  vaultKey: Bytes,
  file: File,
): Promise<void> {
  if (file.size > MAX_ATTACHMENT_BYTES) {
    throw new Error("Ficheiro demasiado grande (limite 5 MiB).");
  }
  const api = requireVaultAccess();
  const bytes = new Uint8Array(await file.arrayBuffer()) as Bytes;
  const metaBlob = await encryptAttachmentMeta(vaultKey, {
    name: file.name,
    mime: file.type || "application/octet-stream",
    size: file.size,
  });
  const dataBlob = await encryptFileBytes(vaultKey, bytes);
  await api.addAttachment(vaultId, metaBlob, dataBlob);
}

/** Anexo descarregado e decifrado, pronto a guardar no dispositivo. */
export interface DownloadedAttachment {
  name: string;
  mime: string;
  bytes: Bytes;
}

/** Descarrega um anexo e decifra-o no dispositivo (nome + tipo + bytes). */
export async function downloadAttachment(
  vaultId: string,
  vaultKey: Bytes,
  attID: string,
): Promise<DownloadedAttachment> {
  const api = requireVaultAccess();
  const dto = await api.getAttachment(vaultId, attID);
  const meta = await decryptAttachmentMeta(vaultKey, dto.meta_blob);
  const bytes = await decryptFileBytes(vaultKey, dto.data_blob);
  return { name: meta.name, mime: meta.mime, bytes };
}

/** Remove um anexo do cofre. */
export async function removeAttachment(vaultId: string, attID: string): Promise<void> {
  await requireVaultAccess().removeAttachment(vaultId, attID);
}

/** Remove um membro (revogação imediata). */
export async function revokeMember(vaultId: string, userID: string): Promise<void> {
  await requireVaultAccess().removeMember(vaultId, userID);
}

/** Remove um item do cofre. */
export async function removeVaultItem(vaultId: string, itemID: string): Promise<void> {
  await requireVaultAccess().removeItem(vaultId, itemID);
}

/** Apaga um cofre (só owner). */
export async function deleteSharedVault(vaultId: string): Promise<void> {
  await requireVaultAccess().remove(vaultId);
}
