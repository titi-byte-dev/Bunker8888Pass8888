/**
 * Faturacao (FIN-006) — modelo + cifragem Zero-Knowledge + calculo de totais.
 *
 * Didatico: o CONTEUDO da fatura (cliente, linhas, valores) e PRIVADO e vai
 * cifrado num blob com a Master Key. O NUMERO legal ("FT 2026/0001") e o estado
 * sao geridos pelo servidor em claro (exigencia fiscal). Os totais calculam-se
 * sempre no cliente, a partir do blob decifrado — o servidor nunca soma nada.
 */
import { decrypt, encrypt, fromBytes, toBytes, type Bytes } from "$lib/crypto";

/** Naturezas legais de documento. */
export type DocType = "proforma" | "invoice" | "receipt";

export const DOC_TYPES: { id: DocType; label: string; prefix: string }[] = [
  { id: "proforma", label: "Pro-forma", prefix: "PF" },
  { id: "invoice", label: "Fatura", prefix: "FT" },
  { id: "receipt", label: "Recibo", prefix: "RC" },
];

export type DocStatus = "issued" | "paid" | "void";

export const DOC_STATUSES: { id: DocStatus; label: string }[] = [
  { id: "issued", label: "Emitido" },
  { id: "paid", label: "Pago" },
  { id: "void", label: "Anulado" },
];

export interface InvoiceClient {
  name: string;
  taxId?: string;
  email?: string;
  address?: string;
}

export interface InvoiceLine {
  description: string;
  quantity: number;
  unitPrice: number;
  /** Taxa de IVA em percentagem (ex.: 23 para 23%). */
  vatRate: number;
}

export interface InvoicePayload {
  client: InvoiceClient;
  lines: InvoiceLine[];
  currency: string;
  issuedAt: string; // ISO
  dueAt?: string;
  notes?: string;
  /** Lead de origem quando nasce de um "deal_closed" (CRM-003). */
  sourceLeadId?: string;
}

/** Documento ja emitido (numero/estado vindos do servidor). */
export interface InvoiceDocument extends InvoicePayload {
  id: string;
  docType: DocType;
  number: string;
  status: DocStatus;
  year: number;
  seq: number;
}

export async function encryptInvoice(masterKey: CryptoKey, payload: InvoicePayload): Promise<Bytes> {
  return encrypt(masterKey, toBytes(JSON.stringify(payload)));
}

export async function decryptInvoice(masterKey: CryptoKey, blob: Bytes): Promise<InvoicePayload> {
  return JSON.parse(fromBytes(await decrypt(masterKey, blob))) as InvoicePayload;
}

// --- Calculo de totais (funcoes puras, testaveis) ---------------------------

function round2(n: number): number {
  return Math.round((n + Number.EPSILON) * 100) / 100;
}

/** Valor liquido de uma linha (sem IVA). */
export function lineNet(line: InvoiceLine): number {
  return round2(line.quantity * line.unitPrice);
}

/** IVA de uma linha. */
export function lineVat(line: InvoiceLine): number {
  return round2(lineNet(line) * (line.vatRate / 100));
}

/** Valor com IVA de uma linha. */
export function lineGross(line: InvoiceLine): number {
  return round2(lineNet(line) + lineVat(line));
}

export interface InvoiceTotals {
  net: number;
  vat: number;
  gross: number;
}

/** Soma os totais de todas as linhas do documento. */
export function invoiceTotals(payload: InvoicePayload): InvoiceTotals {
  let net = 0;
  let vat = 0;
  for (const line of payload.lines) {
    net += lineNet(line);
    vat += lineVat(line);
  }
  net = round2(net);
  vat = round2(vat);
  return { net, vat, gross: round2(net + vat) };
}

/** Rotulo legivel do tipo de documento. */
export function docTypeLabel(t: DocType): string {
  return DOC_TYPES.find((d) => d.id === t)?.label ?? t;
}
