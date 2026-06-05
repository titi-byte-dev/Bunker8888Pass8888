/**
 * Orquestracao da faturacao (FIN-006). Exige sessao + Master Key desbloqueada.
 * O cliente cifra o conteudo; o servidor devolve o numero legal e o estado.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { base64ToBytes, bytesToBase64 } from "$lib/auth";
import {
  decryptInvoice,
  encryptInvoice,
  type DocStatus,
  type DocType,
  type InvoiceDocument,
  type InvoicePayload,
} from "./invoices";

interface InvoiceDTO {
  id: string;
  doc_type: DocType;
  number: string;
  status: DocStatus;
  year: number;
  seq: number;
  source_lead_id: string;
  blob: string;
}

function token(): string {
  const t = loadSessionToken();
  if (!t) throw new Error("Sessao expirada — inicia sessao de novo.");
  return t;
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para gerir faturas.");
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

/** Lista e decifra todos os documentos. */
export async function listInvoices(): Promise<InvoiceDocument[]> {
  const mk = requireMasterKey();
  const j = (await (await authedFetch("/api/fin/invoices")).json()) as { invoices: InvoiceDTO[] };
  const out: InvoiceDocument[] = [];
  for (const dto of j.invoices ?? []) {
    try {
      const p = await decryptInvoice(mk, base64ToBytes(dto.blob));
      out.push({
        ...p,
        id: dto.id,
        docType: dto.doc_type,
        number: dto.number,
        status: dto.status,
        year: dto.year,
        seq: dto.seq,
      });
    } catch {
      /* Master Key errada — ignora a entrada ilegivel. */
    }
  }
  return out;
}

/** Emite um documento novo; o servidor atribui o numero legal. */
export async function issueInvoice(docType: DocType, payload: InvoicePayload): Promise<void> {
  const mk = requireMasterKey();
  const blob = await encryptInvoice(mk, payload);
  await authedFetch("/api/fin/invoices", {
    method: "POST",
    body: JSON.stringify({
      doc_type: docType,
      source_lead_id: payload.sourceLeadId ?? "",
      blob: bytesToBase64(blob),
    }),
  });
}

/** Muda o estado do documento (issued -> paid | void). */
export async function updateInvoiceStatus(id: string, status: DocStatus): Promise<void> {
  await authedFetch(`/api/fin/invoices/${encodeURIComponent(id)}/status`, {
    method: "PUT",
    body: JSON.stringify({ status }),
  });
}
