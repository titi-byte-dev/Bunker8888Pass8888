import { describe, expect, it } from "vitest";
import { canIssueReceipt, receiptFromInvoice } from "./receiptFromInvoice";
import type { InvoiceDocument } from "./invoices";

const paid: InvoiceDocument = {
  id: "i1",
  docType: "invoice",
  number: "FT 2026/0001",
  status: "paid",
  year: 2026,
  seq: 1,
  client: { name: "ACME" },
  lines: [{ description: "Serviço", quantity: 1, unitPrice: 100, vatRate: 23 }],
  currency: "EUR",
  issuedAt: "2026-01-01",
};

describe("receiptFromInvoice", () => {
  it("só permite recibo de fatura paga", () => {
    expect(canIssueReceipt(paid)).toBe(true);
    expect(canIssueReceipt({ ...paid, status: "issued" })).toBe(false);
  });

  it("gera linha com total bruto", () => {
    const r = receiptFromInvoice(paid);
    expect(r.lines[0].unitPrice).toBe(123);
  });
});
