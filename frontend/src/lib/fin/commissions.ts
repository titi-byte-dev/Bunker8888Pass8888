/**
 * Comissoes (FIN-007) — modelo + cifragem Zero-Knowledge + calculo do valor.
 *
 * Didatico: o CONTEUDO (quem recebe, que percentagem, sobre que base) e PRIVADO
 * e vai cifrado num blob com a Master Key. O servidor so guarda a ligacao a
 * fatura de origem e o estado de liquidacao (pending/paid/void). O valor da
 * comissao calcula-se sempre no cliente — o servidor nunca multiplica nada.
 */
import { decrypt, encrypt, fromBytes, toBytes, type Bytes } from "$lib/crypto";
import { invoiceTotals, type InvoiceDocument } from "./invoices";

export type CommissionStatus = "pending" | "paid" | "void";

export const COMMISSION_STATUSES: { id: CommissionStatus; label: string }[] = [
  { id: "pending", label: "Pendente" },
  { id: "paid", label: "Liquidada" },
  { id: "void", label: "Anulada" },
];

export interface CommissionPayload {
  /** Beneficiario da comissao (ex.: nome do vendedor). */
  beneficiary: string;
  /** Percentagem aplicada sobre a base (ex.: 10 para 10%). */
  ratePct: number;
  /** Base de calculo (tipicamente o liquido da fatura). */
  baseAmount: number;
  currency: string;
  /** Numero legal da fatura de origem, para leitura humana. */
  invoiceNumber?: string;
  notes?: string;
}

/** Comissao ja registada (id/estado vindos do servidor). */
export interface CommissionDocument extends CommissionPayload {
  id: string;
  invoiceId: string;
  status: CommissionStatus;
}

export async function encryptCommission(masterKey: CryptoKey, payload: CommissionPayload): Promise<Bytes> {
  return encrypt(masterKey, toBytes(JSON.stringify(payload)));
}

export async function decryptCommission(masterKey: CryptoKey, blob: Bytes): Promise<CommissionPayload> {
  return JSON.parse(fromBytes(await decrypt(masterKey, blob))) as CommissionPayload;
}

function round2(n: number): number {
  return Math.round((n + Number.EPSILON) * 100) / 100;
}

/** Valor da comissao = base * percentagem. Funcao pura, testavel. */
export function commissionAmount(payload: CommissionPayload): number {
  return round2(payload.baseAmount * (payload.ratePct / 100));
}

/** Rotulo legivel do estado. */
export function commissionStatusLabel(s: CommissionStatus): string {
  return COMMISSION_STATUSES.find((x) => x.id === s)?.label ?? s;
}

/**
 * Ponte FIN-007: a partir de uma fatura paga, propoe a comissao do vendedor.
 * A base e o liquido da fatura (sem IVA), calculado no cliente. Funcao pura —
 * nao toca na rede; quem chama cifra e envia.
 */
export function commissionFromInvoice(
  doc: InvoiceDocument,
  ratePct: number,
  beneficiary: string,
): CommissionPayload {
  return {
    beneficiary,
    ratePct,
    baseAmount: invoiceTotals(doc).net,
    currency: doc.currency,
    invoiceNumber: doc.number,
  };
}
