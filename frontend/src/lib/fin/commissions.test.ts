import { describe, expect, it } from "vitest";
import {
  commissionAmount,
  commissionFromInvoice,
  decryptCommission,
  encryptCommission,
  type CommissionPayload,
} from "./commissions";
import type { InvoiceDocument } from "./invoices";

async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]);
}

const payload: CommissionPayload = {
  beneficiary: "Maria Vendas",
  ratePct: 10,
  baseAmount: 700,
  currency: "EUR",
  invoiceNumber: "FT 2026/0001",
};

const paidInvoice: InvoiceDocument = {
  id: "inv-1",
  docType: "invoice",
  number: "FT 2026/0001",
  status: "paid",
  year: 2026,
  seq: 1,
  client: { name: "ACME Lda" },
  lines: [{ description: "Consultoria", quantity: 10, unitPrice: 50, vatRate: 23 }],
  currency: "EUR",
  issuedAt: "2026-06-05T00:00:00Z",
};

describe("FIN-007 comissoes (cifragem)", () => {
  it("cifra e decifra (round-trip)", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptCommission(mk, payload);
    expect(await decryptCommission(mk, blob)).toEqual(payload);
  });

  it("o blob nao contem o beneficiario em claro", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptCommission(mk, payload);
    expect(new TextDecoder().decode(blob)).not.toContain("Maria Vendas");
  });

  it("nao decifra com outra Master Key", async () => {
    const mk = await makeMasterKey();
    const other = await makeMasterKey();
    const blob = await encryptCommission(mk, payload);
    await expect(decryptCommission(other, blob)).rejects.toThrow();
  });
});

describe("FIN-007 calculo + ponte", () => {
  it("valor da comissao = base * percentagem", () => {
    expect(commissionAmount(payload)).toBe(70); // 700 * 10%
  });

  it("ponte usa o liquido da fatura como base", () => {
    const c = commissionFromInvoice(paidInvoice, 15, "Maria Vendas");
    expect(c.baseAmount).toBe(500); // liquido = 10 * 50
    expect(c.ratePct).toBe(15);
    expect(c.beneficiary).toBe("Maria Vendas");
    expect(c.invoiceNumber).toBe("FT 2026/0001");
    expect(commissionAmount(c)).toBe(75); // 500 * 15%
  });
});
