import { describe, expect, it } from "vitest";
import {
  decryptInvoice,
  encryptInvoice,
  invoiceTotals,
  lineGross,
  lineNet,
  lineVat,
  type InvoicePayload,
} from "./invoices";

async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]);
}

const payload: InvoicePayload = {
  client: { name: "ACME Lda", taxId: "PT500000000", email: "geral@acme.pt" },
  lines: [
    { description: "Consultoria", quantity: 10, unitPrice: 50, vatRate: 23 },
    { description: "Licenca", quantity: 2, unitPrice: 100, vatRate: 23 },
  ],
  currency: "EUR",
  issuedAt: "2026-06-05T00:00:00Z",
  sourceLeadId: "lead-1",
};

describe("FIN-006 faturacao (cifragem)", () => {
  it("cifra e decifra um documento (round-trip)", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptInvoice(mk, payload);
    expect(await decryptInvoice(mk, blob)).toEqual(payload);
  });

  it("o blob nao contem o NIF em claro", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptInvoice(mk, payload);
    expect(new TextDecoder().decode(blob)).not.toContain("PT500000000");
  });

  it("nao decifra com outra Master Key", async () => {
    const mk = await makeMasterKey();
    const other = await makeMasterKey();
    const blob = await encryptInvoice(mk, payload);
    await expect(decryptInvoice(other, blob)).rejects.toThrow();
  });
});

describe("FIN-006 calculo de totais", () => {
  it("calcula liquido, IVA e total de uma linha", () => {
    const line = { description: "x", quantity: 10, unitPrice: 50, vatRate: 23 };
    expect(lineNet(line)).toBe(500);
    expect(lineVat(line)).toBe(115);
    expect(lineGross(line)).toBe(615);
  });

  it("soma os totais do documento", () => {
    const t = invoiceTotals(payload);
    expect(t.net).toBe(700); // 500 + 200
    expect(t.vat).toBe(161); // 115 + 46
    expect(t.gross).toBe(861);
  });

  it("documento sem linhas tem totais a zero", () => {
    expect(invoiceTotals({ ...payload, lines: [] })).toEqual({ net: 0, vat: 0, gross: 0 });
  });
});
