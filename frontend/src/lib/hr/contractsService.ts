/**
 * Orquestração de contratos (HR-005) + assinatura digital (HR-006).
 * Exige sessão + Master Key desbloqueada.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { base64ToBytes } from "$lib/auth";
import {
  decryptContractData,
  decryptContractMeta,
  encryptContract,
  MAX_CONTRACT_BYTES,
} from "./contracts";
import { ContractsAPI, type ContractDTO } from "./contractsApi";
import {
  createSigningIdentity,
  importPublicKey,
  sha256hex,
  signBytes,
  unwrapPrivateKey,
  verifyBytes,
} from "./signing";

export { MAX_CONTRACT_BYTES };

export interface DecryptedContract {
  id: string;
  name: string;
  mime: string;
  size: number;
  signed: boolean;
  signedBy: string;
  signedAt: string | null;
  createdAt: string;
}

function api(): ContractsAPI {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessao expirada — inicia sessao de novo.");
  return new ContractsAPI(token);
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para gerir contratos.");
  return mk;
}

/** Lista contratos de uma ficha, decifrando os metadados (nome do ficheiro é ZK). */
export async function listContracts(recordId: string): Promise<DecryptedContract[]> {
  const mk = requireMasterKey();
  const dtos = await api().list(recordId);
  const out: DecryptedContract[] = [];
  for (const c of dtos) {
    let name = "(ilegível)";
    let mime = "";
    let size = c.byte_size;
    try {
      const meta = await decryptContractMeta(mk, base64ToBytes(c.meta_blob), base64ToBytes(c.wrapped_key));
      name = meta.name;
      mime = meta.mime;
      size = meta.size;
    } catch {
      /* Master Key errada — mostra placeholder. */
    }
    out.push({
      id: c.id,
      name,
      mime,
      size,
      signed: c.signed,
      signedBy: c.signed_by ?? "",
      signedAt: c.signed_at ?? null,
      createdAt: c.created_at,
    });
  }
  return out;
}

/** Cifra e carrega um contrato. Recusa ficheiros acima do limite. */
export async function uploadContract(recordId: string, file: File): Promise<void> {
  const mk = requireMasterKey();
  if (file.size > MAX_CONTRACT_BYTES) {
    throw new Error("Ficheiro demasiado grande (máx. ~5 MiB).");
  }
  const { metaBlob, dataBlob, wrappedKey } = await encryptContract(mk, file);
  await api().add(recordId, metaBlob, dataBlob, wrappedKey);
}

/** Remove um contrato da ficha. */
export async function removeContract(recordId: string, contractId: string): Promise<void> {
  await api().remove(recordId, contractId);
}

/** Decifra e descarrega um contrato (devolve Blob + nome do ficheiro). */
export async function downloadContract(
  recordId: string,
  contractId: string,
): Promise<{ blob: Blob; name: string }> {
  const mk = requireMasterKey();
  const dto = await api().get(recordId, contractId);
  const meta = await decryptContractMeta(mk, base64ToBytes(dto.meta_blob), base64ToBytes(dto.wrapped_key));
  const bytes = await decryptContractData(mk, base64ToBytes(dto.data_blob ?? ""), base64ToBytes(dto.wrapped_key));
  return { blob: new Blob([bytes], { type: meta.mime }), name: meta.name };
}

/** Garante que existe identidade de assinatura; cria-a se faltar. */
async function ensurePrivateKey(): Promise<CryptoKey> {
  const mk = requireMasterKey();
  const existing = await api().getSigningIdentity();
  if (existing) {
    return unwrapPrivateKey(mk, existing.wrappedPrivateKey);
  }
  const id = await createSigningIdentity(mk);
  await api().putSigningIdentity(id.publicKey, id.wrappedPrivateKey);
  return id.privateKey;
}

/** Assina digitalmente um contrato (sobre os bytes do ciphertext). */
export async function signContract(recordId: string, contractId: string): Promise<void> {
  const privateKey = await ensurePrivateKey();
  const dto = await api().get(recordId, contractId);
  const dataBlob = base64ToBytes(dto.data_blob ?? "");
  const signature = await signBytes(privateKey, dataBlob);
  const digest = await sha256hex(dataBlob);
  await api().sign(recordId, contractId, digest, signature);
}

export interface VerificationResult {
  valid: boolean;
  reason: string;
}

/** Verifica a assinatura de um contrato com a chave pública do signatário. */
export async function verifyContract(
  recordId: string,
  contractId: string,
): Promise<VerificationResult> {
  const dto: ContractDTO = await api().get(recordId, contractId);
  if (!dto.signed || !dto.signature || !dto.signed_by) {
    return { valid: false, reason: "Contrato não assinado." };
  }
  const dataBlob = base64ToBytes(dto.data_blob ?? "");
  // O digest tem de bater com os bytes actuais (deteta troca do ficheiro).
  const digest = await sha256hex(dataBlob);
  if (dto.content_digest && dto.content_digest !== digest) {
    return { valid: false, reason: "Os bytes do contrato mudaram desde a assinatura." };
  }
  const spki = await api().signerPublicKey(dto.signed_by);
  const publicKey = await importPublicKey(spki);
  const ok = await verifyBytes(publicKey, base64ToBytes(dto.signature), dataBlob);
  return ok
    ? { valid: true, reason: "Assinatura válida." }
    : { valid: false, reason: "Assinatura inválida." };
}
