/**
 * Orquestração de alto nível da partilha (SHARE-001) — requer sessão + Master
 * Key desbloqueada, tal como o cofre.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { ShareKeysAPI, type RecipientPublicKey } from "./api";
import {
  exportPublicKey,
  generateSharingKeypair,
  publicKeyFingerprint,
  unwrapPrivateKey,
  wrapPrivateKey,
} from "./keypair";

export interface ShareIdentity {
  /** Chave pública SPKI (base64), partilhável. */
  publicKey: string;
  /** Impressão digital legível para verificação fora de banda. */
  fingerprint: string;
  algorithm: string;
}

function requireShareAccess(): { api: ShareKeysAPI; masterKey: CryptoKey } {
  const token = loadSessionToken();
  const masterKey = getMasterKey();
  if (!token || !masterKey) {
    throw new Error("Partilha bloqueada — inicia sessão e desbloqueia em /auth/unlock");
  }
  return { api: new ShareKeysAPI("", token), masterKey };
}

/** Indica se o utilizador já activou a partilha. */
export async function isShareConfigured(): Promise<boolean> {
  const { api } = requireShareAccess();
  return api.status();
}

/**
 * Garante que existe um par de chaves. Se não existir, gera-o, cifra a privada
 * com a Master Key e envia ambas para o servidor. Devolve a identidade pública.
 */
export async function ensureShareIdentity(): Promise<ShareIdentity> {
  const { api, masterKey } = requireShareAccess();

  if (await api.status()) {
    const stored = await api.fetchOwn();
    return {
      publicKey: stored.public_key,
      fingerprint: await publicKeyFingerprint(stored.public_key),
      algorithm: stored.algorithm,
    };
  }

  const pair = await generateSharingKeypair();
  const publicKey = await exportPublicKey(pair.publicKey);
  const wrappedPrivateKey = await wrapPrivateKey(masterKey, pair.privateKey);
  await api.upload(publicKey, wrappedPrivateKey);

  return {
    publicKey,
    fingerprint: await publicKeyFingerprint(publicKey),
    algorithm: (await api.fetchOwn()).algorithm,
  };
}

/** Carrega a identidade pública existente (ou null se ainda não configurada). */
export async function loadShareIdentity(): Promise<ShareIdentity | null> {
  const { api } = requireShareAccess();
  if (!(await api.status())) return null;
  const stored = await api.fetchOwn();
  return {
    publicKey: stored.public_key,
    fingerprint: await publicKeyFingerprint(stored.public_key),
    algorithm: stored.algorithm,
  };
}

/**
 * Abre a chave privada de partilha desta sessão (decifra com a Master Key).
 * Será usada pela SHARE-002 para abrir itens partilhados comigo.
 */
export async function loadOwnPrivateKey(): Promise<CryptoKey> {
  const { api, masterKey } = requireShareAccess();
  const stored = await api.fetchOwn();
  return unwrapPrivateKey(masterKey, stored.wrapped_private_key);
}

export interface RecipientLookup extends RecipientPublicKey {
  fingerprint: string;
}

/** Procura a chave pública de um colega e calcula a sua impressão digital. */
export async function lookupRecipient(email: string): Promise<RecipientLookup> {
  const { api } = requireShareAccess();
  const recipient = await api.fetchPublicKey(email);
  return { ...recipient, fingerprint: await publicKeyFingerprint(recipient.public_key) };
}
