import { describe, expect, it } from "vitest";
import { canConvertToInvoice, invoiceFromProforma } from "./proformaToInvoice";
import type { InvoiceDocument } from "$lib/fin/invoices";

const proforma: InvoiceDocument = {
  id: "pf-1",
  docType: "proforma",
  number: "PF 2026/0003",
  status: "issued",
  year: 2026,
  seq: 3,
  client: { name: "ACME Lda", taxId: "PT500000000" },
  lines: [{ description: "Consultoria", quantity: 10, unitPrice: 50, vatRate: 23 }],
  currency: "EUR",
  issuedAt: "2026-06-05T00:00:00Z",
  sourceLeadId: "lead-7",
  notes: "Orcamento inicial",
};

describe("CRM-004 conversao pro-forma -> fatura", () => {
  it("so converte pro-formas nao anuladas", () => {
    expect(canConvertToInvoice(proforma)).toBe(true);
    expect(canConvertToInvoice({ ...proforma, status: "void" })).toBe(false);
    expect(canConvertToInvoice({ ...proforma, docType: "invoice" })).toBe(false);
  });

  it("herda cliente, linhas, moeda e o lead de origem", () => {
    const inv = invoiceFromProforma(proforma);
    expect(inv.client).toEqual(proforma.client);
    expect(inv.lines).toEqual(proforma.lines);
    expect(inv.currency).toBe("EUR");
    expect(inv.sourceLeadId).toBe("lead-7");
  });

  it("regista de que pro-forma nasceu, preservando notas", () => {
    const inv = invoiceFromProforma(proforma);
    expect(inv.notes).toContain("Orcamento inicial");
    expect(inv.notes).toContain("PF 2026/0003");
  });

  it("copia as linhas sem partilhar referencia (imutavel)", () => {
    const inv = invoiceFromProforma(proforma);
    inv.lines[0].quantity = 999;
    expect(proforma.lines[0].quantity).toBe(10);
  });
});
