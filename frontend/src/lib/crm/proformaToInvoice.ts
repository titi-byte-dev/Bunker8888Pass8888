/**
 * Conversao pro-forma -> fatura (CRM-004).
 *
 * Didatico: uma pro-forma (PF) e so um orcamento. Quando o cliente aceita, a
 * empresa emite a fatura definitiva (FT). Esta ponte e PURA e corre no cliente:
 * pega no conteudo ja decifrado da pro-forma e produz um InvoicePayload pronto a
 * ser cifrado e emitido como "invoice". O servidor so vera o novo blob cifrado e
 * atribuira o proximo numero FT — nunca liga uma serie a outra em claro.
 *
 * A ligacao logica preserva-se de duas formas:
 *   - sourceLeadId herda do documento de origem (mantem o rasto ao deal);
 *   - uma nota regista de que pro-forma nasceu a fatura.
 */
import type { InvoiceDocument, InvoicePayload } from "$lib/fin/invoices";

export interface ConvertOptions {
  /** Data de emissao da fatura (ISO). Por omissao, agora. */
  issuedAt?: string;
}

/** True se o documento pode ser convertido em fatura (e uma pro-forma viva). */
export function canConvertToInvoice(doc: InvoiceDocument): boolean {
  return doc.docType === "proforma" && doc.status !== "void";
}

/** Constroi o payload da fatura a partir de uma pro-forma. */
export function invoiceFromProforma(doc: InvoiceDocument, opts: ConvertOptions = {}): InvoicePayload {
  const origin = doc.number ? `Convertida da pro-forma ${doc.number}.` : "Convertida de pro-forma.";
  const notes = doc.notes ? `${doc.notes}\n${origin}` : origin;
  return {
    client: { ...doc.client },
    lines: doc.lines.map((l) => ({ ...l })),
    currency: doc.currency,
    issuedAt: opts.issuedAt ?? new Date().toISOString(),
    dueAt: doc.dueAt,
    notes,
    sourceLeadId: doc.sourceLeadId,
  };
}
