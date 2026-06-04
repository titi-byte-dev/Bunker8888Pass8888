import type { VaultItemInput, VaultItemMeta, VaultItemType } from "./types";
import type { GeoPosition } from "./geofence";
import { geoHeaders } from "./geofence";

export class VaultAPI {
  constructor(
    private baseURL: string,
    private token: string,
    private geo: GeoPosition | null = null,
  ) {}

  /** Actualiza coordenadas GPS enviadas em cada pedido. */
  setGeo(pos: GeoPosition | null): void {
    this.geo = pos;
  }

  async list(type?: VaultItemType): Promise<VaultItemMeta[]> {
    const q = type ? `?type=${encodeURIComponent(type)}` : "";
    const res = await this.fetch(`/api/vault${q}`);
    return ((await res.json()) as { items: VaultItemMeta[] }).items;
  }

  async get(id: string): Promise<VaultItemMeta> {
    return (await this.fetch(`/api/vault/${id}`)).json() as Promise<VaultItemMeta>;
  }

  async create(input: VaultItemInput): Promise<VaultItemMeta> {
    const res = await this.fetch("/api/vault", { method: "POST", body: JSON.stringify(input) });
    return res.json() as Promise<VaultItemMeta>;
  }

  async update(id: string, input: VaultItemInput): Promise<VaultItemMeta> {
    const res = await this.fetch(`/api/vault/${id}`, { method: "PUT", body: JSON.stringify(input) });
    return res.json() as Promise<VaultItemMeta>;
  }

  async delete(id: string): Promise<void> {
    const res = await this.fetch(`/api/vault/${id}`, { method: "DELETE" });
    if (res.status !== 204) throw new Error(`delete falhou: ${res.status}`);
  }

  private async fetch(path: string, init: RequestInit = {}): Promise<Response> {
    const res = await globalThis.fetch(`${this.baseURL}${path}`, {
      ...init,
      headers: {
        Authorization: `Bearer ${this.token}`,
        "Content-Type": "application/json",
        ...geoHeaders(this.geo),
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
