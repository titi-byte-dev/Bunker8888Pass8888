/**
 * Cliente da API de contratos + identidade de assinatura (HR-005/006).
 * Só transporta bytes (base64); a cifra e a assinatura acontecem no cliente.
 */
import { base64ToBytes, bytesToBase64 } from "$lib/auth";
import type { Bytes } from "$lib/crypto";

export interface ContractDTO {
  id: string;
  record_id: string;
  meta_blob: string;
  wrapped_key: string;
  byte_size: number;
  signed: boolean;
  content_digest?: string;
  signature?: string;
  signed_by?: string;
  signed_at?: string;
  created_by: string;
  created_at: string;
  data_blob?: string;
}

export class ContractsAPI {
  constructor(
    private token: string,
  ) {}

  private async fetch(path: string, init: RequestInit = {}): Promise<Response> {
    const res = await fetch(path, {
      ...init,
      headers: {
        Authorization: `Bearer ${this.token}`,
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

  /** Devolve a identidade de assinatura, ou null se ainda não existir (404). */
  async getSigningIdentity(): Promise<{ publicKey: Bytes; wrappedPrivateKey: Bytes } | null> {
    const res = await fetch("/api/hr/signing-identity", {
      headers: { Authorization: `Bearer ${this.token}` },
    });
    if (res.status === 404) return null;
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const j = (await res.json()) as { public_key: string; wrapped_private_key: string };
    return {
      publicKey: base64ToBytes(j.public_key),
      wrappedPrivateKey: base64ToBytes(j.wrapped_private_key),
    };
  }

  async putSigningIdentity(publicKey: Bytes, wrappedPrivateKey: Bytes): Promise<void> {
    await this.fetch("/api/hr/signing-identity", {
      method: "PUT",
      body: JSON.stringify({
        public_key: bytesToBase64(publicKey),
        wrapped_private_key: bytesToBase64(wrappedPrivateKey),
      }),
    });
  }

  async signerPublicKey(userId: string): Promise<Bytes> {
    const res = await this.fetch(`/api/hr/signers/${encodeURIComponent(userId)}/public-key`);
    return base64ToBytes(((await res.json()) as { public_key: string }).public_key);
  }

  async list(recordId: string): Promise<ContractDTO[]> {
    const res = await this.fetch(`/api/hr/employees/${encodeURIComponent(recordId)}/contracts`);
    return ((await res.json()) as { contracts: ContractDTO[] }).contracts ?? [];
  }

  async add(recordId: string, metaBlob: Bytes, dataBlob: Bytes, wrappedKey: Bytes): Promise<ContractDTO> {
    const res = await this.fetch(`/api/hr/employees/${encodeURIComponent(recordId)}/contracts`, {
      method: "POST",
      body: JSON.stringify({
        meta_blob: bytesToBase64(metaBlob),
        data_blob: bytesToBase64(dataBlob),
        wrapped_key: bytesToBase64(wrappedKey),
      }),
    });
    return (await res.json()) as ContractDTO;
  }

  async get(recordId: string, contractId: string): Promise<ContractDTO> {
    const res = await this.fetch(
      `/api/hr/employees/${encodeURIComponent(recordId)}/contracts/${encodeURIComponent(contractId)}`,
    );
    return (await res.json()) as ContractDTO;
  }

  async sign(recordId: string, contractId: string, contentDigest: string, signature: Bytes): Promise<ContractDTO> {
    const res = await this.fetch(
      `/api/hr/employees/${encodeURIComponent(recordId)}/contracts/${encodeURIComponent(contractId)}/sign`,
      {
        method: "POST",
        body: JSON.stringify({ content_digest: contentDigest, signature: bytesToBase64(signature) }),
      },
    );
    return (await res.json()) as ContractDTO;
  }

  async remove(recordId: string, contractId: string): Promise<void> {
    await this.fetch(
      `/api/hr/employees/${encodeURIComponent(recordId)}/contracts/${encodeURIComponent(contractId)}`,
      { method: "DELETE" },
    );
  }
}
