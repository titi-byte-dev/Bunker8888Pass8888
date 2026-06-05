/**
 * Orquestracao das comissoes (FIN-007). Exige sessao + Master Key desbloqueada.
 * O cliente cifra o conteudo; o servidor guarda a ligacao e o estado.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { base64ToBytes, bytesToBase64 } from "$lib/auth";
import {
  decryptCommission,
  encryptCommission,
  type CommissionDocument,
  type CommissionPayload,
  type CommissionStatus,
} from "./commissions";

interface CommissionDTO {
  id: string;
  invoice_id: string;
  status: CommissionStatus;
  blob: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessao expirada — inicia sessao de novo.");
  return t;
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para gerir comissoes.");
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

/** Lista e decifra todas as comissoes. */
export async function listCommissions(): Promise<CommissionDocument[]> {
  const mk = requireMasterKey();
  const j = (await (await authedFetch("/api/fin/commissions")).json()) as { commissions: CommissionDTO[] };
  const out: CommissionDocument[] = [];
  for (const dto of j.commissions ?? []) {
    try {
      const p = await decryptCommission(mk, base64ToBytes(dto.blob));
      out.push({ ...p, id: dto.id, invoiceId: dto.invoice_id, status: dto.status });
    } catch {
      /* Master Key errada — ignora a entrada ilegivel. */
    }
  }
  return out;
}

/** Regista uma comissao nova (estado inicial pending). */
export async function createCommission(invoiceId: string, payload: CommissionPayload): Promise<void> {
  const mk = requireMasterKey();
  const blob = await encryptCommission(mk, payload);
  await authedFetch("/api/fin/commissions", {
    method: "POST",
    body: JSON.stringify({ invoice_id: invoiceId, blob: bytesToBase64(blob) }),
  });
}

/** Muda o estado de liquidacao (pending -> paid | void). */
export async function updateCommissionStatus(id: string, status: CommissionStatus): Promise<void> {
  await authedFetch(`/api/fin/commissions/${encodeURIComponent(id)}/status`, {
    method: "PUT",
    body: JSON.stringify({ status }),
  });
}
