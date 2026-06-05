/**
 * Cliente da API de Cofres Partilhados (SHARE-002).
 *
 * Só transporta bytes (base64). Toda a cifra/decifra acontece no cliente; o
 * servidor encaminha blobs opacos e faz cumprir as permissões.
 */
import { SHARED_VAULT_ALGORITHM } from "./vaults";

export type VaultRole = "owner" | "admin" | "member" | "viewer";

export interface SharedVaultDTO {
  id: string;
  owner_id: string;
  name_blob: string;
  algorithm: string;
  role: VaultRole;
  wrapped_vault_key: string;
  created_at: string;
  updated_at: string;
}

export interface VaultMemberDTO {
  user_id: string;
  email: string;
  role: VaultRole;
  created_at: string;
}

export interface VaultItemDTO {
  id: string;
  type: string;
  blob: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export class SharedVaultsAPI {
  constructor(
    private baseURL: string,
    private token: string,
  ) {}

  /** Cria um cofre; o criador fica owner com a sua cópia da chave. */
  async create(nameBlob: string, ownerWrappedKey: string): Promise<SharedVaultDTO> {
    const res = await this.fetch("/api/share/vaults", {
      method: "POST",
      body: JSON.stringify({
        name_blob: nameBlob,
        algorithm: SHARED_VAULT_ALGORITHM,
        wrapped_vault_key: ownerWrappedKey,
      }),
    });
    return (await res.json()) as SharedVaultDTO;
  }

  /** Lista os cofres de que sou membro. */
  async list(): Promise<SharedVaultDTO[]> {
    const res = await this.fetch("/api/share/vaults");
    return ((await res.json()) as { vaults: SharedVaultDTO[] }).vaults;
  }

  /** Devolve um cofre (com o meu papel e a minha cópia da chave). */
  async get(id: string): Promise<SharedVaultDTO> {
    const res = await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}`);
    return (await res.json()) as SharedVaultDTO;
  }

  /** Apaga um cofre (só owner). */
  async remove(id: string): Promise<void> {
    await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}`, { method: "DELETE" });
  }

  /** Lista os membros do cofre. */
  async members(id: string): Promise<VaultMemberDTO[]> {
    const res = await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}/members`);
    return ((await res.json()) as { members: VaultMemberDTO[] }).members;
  }

  /** Convida (ou actualiza) um membro com a chave do cofre cifrada para ele. */
  async addMember(id: string, userID: string, role: VaultRole, wrappedKey: string): Promise<void> {
    await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}/members`, {
      method: "POST",
      body: JSON.stringify({ user_id: userID, role, wrapped_vault_key: wrappedKey }),
    });
  }

  /** Remove um membro (revogação imediata). */
  async removeMember(id: string, userID: string): Promise<void> {
    await this.fetch(
      `/api/share/vaults/${encodeURIComponent(id)}/members/${encodeURIComponent(userID)}`,
      { method: "DELETE" },
    );
  }

  /** Lista os itens cifrados do cofre. */
  async items(id: string): Promise<VaultItemDTO[]> {
    const res = await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}/items`);
    return ((await res.json()) as { items: VaultItemDTO[] }).items;
  }

  /** Adiciona um item cifrado ao cofre. */
  async addItem(id: string, type: string, blob: string): Promise<VaultItemDTO> {
    const res = await this.fetch(`/api/share/vaults/${encodeURIComponent(id)}/items`, {
      method: "POST",
      body: JSON.stringify({ type, blob }),
    });
    return (await res.json()) as VaultItemDTO;
  }

  /** Remove um item do cofre. */
  async removeItem(id: string, itemID: string): Promise<void> {
    await this.fetch(
      `/api/share/vaults/${encodeURIComponent(id)}/items/${encodeURIComponent(itemID)}`,
      { method: "DELETE" },
    );
  }

  private async fetch(path: string, init: RequestInit = {}): Promise<Response> {
    const res = await globalThis.fetch(`${this.baseURL}${path}`, {
      ...init,
      headers: {
        Authorization: `Bearer ${this.token}`,
        "Content-Type": "application/json",
        ...init.headers,
      },
    });
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { error?: string };
      throw new Error(err.error ?? `HTTP ${res.status}`);
    }
    return res;
  }
}
