/**
 * Inventário operacional (AGENT-008) — API cliente.
 */
import { loadSessionToken } from "$lib/session";

export type InventoryItem = {
  id: string;
  name: string;
  sku: string;
  quantity: number;
  reorderLevel: number;
  unit: string;
  lowStock: boolean;
  createdAt: string;
  updatedAt: string;
};

type ItemDTO = {
  id: string;
  name: string;
  sku: string;
  quantity: number;
  reorder_level: number;
  unit: string;
  low_stock: boolean;
  created_at: string;
  updated_at: string;
};

function mapItem(d: ItemDTO): InventoryItem {
  return {
    id: d.id,
    name: d.name,
    sku: d.sku,
    quantity: d.quantity,
    reorderLevel: d.reorder_level,
    unit: d.unit,
    lowStock: d.low_stock,
    createdAt: d.created_at,
    updatedAt: d.updated_at,
  };
}

async function authedFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão expirada — inicia sessão de novo.");
  return fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
      ...(init.headers ?? {}),
    },
  });
}

export async function listInventory(): Promise<InventoryItem[]> {
  const res = await authedFetch("/api/ops/inventory");
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  const j = (await res.json()) as { items: ItemDTO[] };
  return (j.items ?? []).map(mapItem);
}

export type InventoryInput = {
  name: string;
  sku?: string;
  quantity: number;
  reorderLevel: number;
  unit?: string;
};

export async function createInventoryItem(input: InventoryInput): Promise<InventoryItem> {
  const res = await authedFetch("/api/ops/inventory", {
    method: "POST",
    body: JSON.stringify({
      name: input.name,
      sku: input.sku ?? "",
      quantity: input.quantity,
      reorder_level: input.reorderLevel,
      unit: input.unit ?? "un",
    }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return mapItem((await res.json()) as ItemDTO);
}

export async function adjustInventory(id: string, delta: number): Promise<InventoryItem> {
  const res = await authedFetch(`/api/ops/inventory/${encodeURIComponent(id)}/adjust`, {
    method: "POST",
    body: JSON.stringify({ delta }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return mapItem((await res.json()) as ItemDTO);
}

export async function deleteInventoryItem(id: string): Promise<void> {
  const res = await authedFetch(`/api/ops/inventory/${encodeURIComponent(id)}`, { method: "DELETE" });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
}
