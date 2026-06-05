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
  decryptVaultItem,
  decryptVaultName,
  encryptVaultItem,
  encryptVaultName,
  generateVaultKey,
  unwrapVaultKey,
  wrapVaultKeyFor,
  type VaultItemPayload,
} from "./vaults";
import type { Bytes } from "../crypto";

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

/** Cofre aberto: o conteúdo decifrado + a chave do cofre em memória para escrever. */
export interface OpenVault {
  vault: DecryptedVault;
  members: VaultMemberDTO[];
  items: DecryptedItem[];
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

  const [members, itemDtos] = await Promise.all([api.members(id), api.items(id)]);
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
