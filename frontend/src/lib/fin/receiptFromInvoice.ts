/**
 * FIN-006 — emite recibo (RC) a partir de fatura paga.
 */
import type { InvoiceDocument, InvoicePayload } from "./invoices";
import { invoiceTotals } from "./invoices";

/** Constrói payload de recibo ligado à fatura FT paga. */
export function receiptFromInvoice(invoice: InvoiceDocument): InvoicePayload {
  const t = invoiceTotals(invoice);
  return {
    client: { ...invoice.client },
    lines: [
      {
        description: `Pagamento referente a ${invoice.number}`,
        quantity: 1,
        unitPrice: t.gross,
        vatRate: 0,
      },
    ],
    currency: invoice.currency,
    issuedAt: new Date().toISOString(),
    notes: `Recibo da fatura ${invoice.number}.`,
  };
}

export function canIssueReceipt(doc: InvoiceDocument): boolean {
  return doc.docType === "invoice" && doc.status === "paid";
}
