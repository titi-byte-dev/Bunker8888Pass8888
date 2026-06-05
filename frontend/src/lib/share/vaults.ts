/**
 * Cripto dos Cofres Partilhados (SHARE-002).
 *
 * Didático — o "duplo embrulho" da partilha em grupo:
 *
 *   1. Cada cofre tem UMA chave simétrica própria (AES-256, a "chave do cofre").
 *      Tudo o que está dentro — nome e itens — é cifrado com ela.
 *   2. Para dar acesso a um membro, re-cifra-se essa chave do cofre para a CHAVE
 *      PÚBLICA dele (RSA-OAEP, da SHARE-001). Cada membro guarda a SUA cópia
 *      cifrada; só a sua chave privada a abre.
 *
 *        nome/itens ──AES-GCM(chave_do_cofre)──▶ blobs opacos no servidor
 *        chave_do_cofre ──RSA-OAEP(PK_membro)──▶ wrapped_vault_key por membro
 *
 * Resultado: convidar = embrulhar a chave do cofre para mais alguém; revogar =
 * apagar a cópia dessa pessoa. O servidor nunca vê a chave do cofre em claro.
 */
import { decrypt, encrypt, fromBytes, randomBytes, toBytes, type Bytes } from "../crypto";
import { base64ToBytes, bytesToBase64 } from "../auth";
import { unwrapKeyFromSender, wrapKeyForRecipient } from "./keypair";

/** Marcador de esquema (chave simétrica + embrulho assimétrico), versionado. */
export const SHARED_VAULT_ALGORITHM = "AES-GCM-256+RSA-OAEP-3072";

/** Comprimento da chave do cofre: 32 bytes = AES-256. */
const VAULT_KEY_BYTES = 32;

/** Tipo do único item suportado por agora (nota partilhada). */
export interface VaultItemPayload {
  title: string;
  secret: string;
}

/** Gera uma nova chave de cofre (bytes aleatórios). */
export function generateVaultKey(): Bytes {
  return randomBytes(VAULT_KEY_BYTES);
}

/** Importa a chave do cofre (bytes) como CryptoKey AES-GCM para cifrar/decifrar. */
async function importVaultKey(raw: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM" }, false, [
    "encrypt",
    "decrypt",
  ]);
}

/** Cifra o nome do cofre com a chave do cofre. Devolve base64 (nonce||ct). */
export async function encryptVaultName(vaultKey: Bytes, name: string): Promise<string> {
  const key = await importVaultKey(vaultKey);
  return bytesToBase64(await encrypt(key, toBytes(name)));
}

/** Reverte encryptVaultName. */
export async function decryptVaultName(vaultKey: Bytes, b64: string): Promise<string> {
  const key = await importVaultKey(vaultKey);
  return fromBytes(await decrypt(key, base64ToBytes(b64)));
}

/** Cifra um item (JSON) com a chave do cofre. Devolve base64 (nonce||ct). */
export async function encryptVaultItem(vaultKey: Bytes, payload: VaultItemPayload): Promise<string> {
  const key = await importVaultKey(vaultKey);
  return bytesToBase64(await encrypt(key, toBytes(JSON.stringify(payload))));
}

/** Reverte encryptVaultItem. */
export async function decryptVaultItem(vaultKey: Bytes, b64: string): Promise<VaultItemPayload> {
  const key = await importVaultKey(vaultKey);
  return JSON.parse(fromBytes(await decrypt(key, base64ToBytes(b64)))) as VaultItemPayload;
}

/**
 * Embrulha a chave do cofre para um membro (a sua chave pública RSA-OAEP).
 * `recipientPub` é uma CryptoKey já importada via importPublicKey (keypair.ts).
 */
export async function wrapVaultKeyFor(recipientPub: CryptoKey, vaultKey: Bytes): Promise<string> {
  return wrapKeyForRecipient(recipientPub, vaultKey);
}

/** Reverte wrapVaultKeyFor com a minha chave privada. Devolve os bytes da chave. */
export async function unwrapVaultKey(myPriv: CryptoKey, b64: string): Promise<Bytes> {
  return unwrapKeyFromSender(myPriv, b64);
}
