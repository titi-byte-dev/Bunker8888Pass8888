/**
 * Orquestração das subscrições SaaS (FIN-001) + cruzamento com o cofre.
 * Exige sessão + Master Key desbloqueada.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { base64ToBytes, bytesToBase64 } from "$lib/auth";
import { VaultAPI } from "$lib/vault/api";
import { openItem } from "$lib/vault/items";
import {
  decryptSubscription,
  encryptSubscription,
  type Subscription,
  type SubscriptionPayload,
} from "./subscriptions";

interface SubscriptionDTO {
  id: string;
  blob: string;
  created_at: string;
  updated_at: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessao expirada — inicia sessao de novo.");
  return t;
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para ver custos.");
  return mk;
}

async function authedFetch(path: string, init: RequestInit = {}): Promise<Response> {
  const res = await fetch(path, {
    ...init,
    headers: {
      Authorization: `Bearer ${token()}`,
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

/** Lista e decifra todas as subscrições. */
export async function listSubscriptions(): Promise<Subscription[]> {
  const mk = requireMasterKey();
  const j = (await (await authedFetch("/api/fin/subscriptions")).json()) as {
    subscriptions: SubscriptionDTO[];
  };
  const out: Subscription[] = [];
  for (const dto of j.subscriptions ?? []) {
    try {
      const p = await decryptSubscription(mk, base64ToBytes(dto.blob));
      out.push({ ...p, id: dto.id });
    } catch {
      /* Master Key errada — ignora a entrada ilegível. */
    }
  }
  return out;
}

export async function createSubscription(payload: SubscriptionPayload): Promise<void> {
  const mk = requireMasterKey();
  const blob = await encryptSubscription(mk, payload);
  await authedFetch("/api/fin/subscriptions", {
    method: "POST",
    body: JSON.stringify({ blob: bytesToBase64(blob) }),
  });
}

export async function updateSubscription(id: string, payload: SubscriptionPayload): Promise<void> {
  const mk = requireMasterKey();
  const blob = await encryptSubscription(mk, payload);
  await authedFetch(`/api/fin/subscriptions/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify({ blob: bytesToBase64(blob) }),
  });
}

export async function deleteSubscription(id: string): Promise<void> {
  await authedFetch(`/api/fin/subscriptions/${encodeURIComponent(id)}`, { method: "DELETE" });
}

export interface VaultLoginRef {
  id: string;
  title: string;
}

/**
 * Lista os logins do cofre (id + título) para associar a subscrições
 * ("cruza com vault"). Decifra cada item localmente.
 */
export async function listVaultLogins(): Promise<VaultLoginRef[]> {
  const mk = requireMasterKey();
  const api = new VaultAPI("", token());
  const metas = await api.list("login");
  const out: VaultLoginRef[] = [];
  for (const m of metas) {
    if (!m.blob) continue;
    try {
      const payload = await openItem(mk, base64ToBytes(m.blob));
      out.push({ id: m.id, title: payload.title });
    } catch {
      /* item ilegível — ignora. */
    }
  }
  return out;
}
