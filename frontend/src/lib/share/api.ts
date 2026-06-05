/**
 * Cliente da API de chaves de partilha (SHARE-001).
 *
 * O servidor guarda a chave pública em claro e a privada já cifrada — aqui só
 * transportamos bytes (base64). Nenhuma decifragem acontece no servidor.
 */
import { SHARE_KEY_ALGORITHM } from "./keypair";

export interface StoredKeypair {
  public_key: string;
  wrapped_private_key: string;
  algorithm: string;
}

export interface RecipientPublicKey {
  email: string;
  user_id: string;
  public_key: string;
  algorithm: string;
}

export class ShareKeysAPI {
  constructor(
    private baseURL: string,
    private token: string,
  ) {}

  /** Indica se o utilizador já activou a partilha (tem par de chaves). */
  async status(): Promise<boolean> {
    const res = await this.fetch("/api/share/keypair/status");
    return ((await res.json()) as { configured: boolean }).configured;
  }

  /** Grava (ou roda) o par de chaves do utilizador. */
  async upload(publicKey: string, wrappedPrivateKey: string): Promise<void> {
    await this.fetch("/api/share/keypair", {
      method: "PUT",
      body: JSON.stringify({
        public_key: publicKey,
        wrapped_private_key: wrappedPrivateKey,
        algorithm: SHARE_KEY_ALGORITHM,
      }),
    });
  }

  /** Devolve o par do próprio (inclui a chave privada cifrada). */
  async fetchOwn(): Promise<StoredKeypair> {
    const res = await this.fetch("/api/share/keypair");
    return (await res.json()) as StoredKeypair;
  }

  /** Procura a chave pública de um colega, para lhe partilhar um segredo. */
  async fetchPublicKey(email: string): Promise<RecipientPublicKey> {
    const res = await this.fetch(`/api/share/public-key?email=${encodeURIComponent(email)}`);
    return (await res.json()) as RecipientPublicKey;
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
