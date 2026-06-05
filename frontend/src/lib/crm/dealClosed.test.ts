import { describe, expect, it } from "vitest";
import { isDealClosed, proformaFromLead } from "./dealClosed";
import type { Lead } from "./leads";

const wonLead: Lead = {
  id: "lead-1",
  name: "Ana Silva",
  email: "ana@acme.pt",
  company: "ACME Lda",
  stage: "won",
};

describe("CRM-003 deal_closed", () => {
  it("deteta negocio fechado apenas no estagio won", () => {
    expect(isDealClosed({ stage: "won" })).toBe(true);
    expect(isDealClosed({ stage: "proposal" })).toBe(false);
  });

  it("gera pro-forma com cliente da empresa e liga o lead", () => {
    const p = proformaFromLead(wonLead);
    expect(p.client.name).toBe("ACME Lda");
    expect(p.client.email).toBe("ana@acme.pt");
    expect(p.sourceLeadId).toBe("lead-1");
    expect(p.currency).toBe("EUR");
    expect(p.lines).toHaveLength(1);
  });

  it("usa o nome quando nao ha empresa", () => {
    const p = proformaFromLead({ ...wonLead, company: undefined });
    expect(p.client.name).toBe("Ana Silva");
  });

  it("respeita linhas e moeda fornecidas", () => {
    const lines = [{ description: "Setup", quantity: 1, unitPrice: 1000, vatRate: 23 }];
    const p = proformaFromLead(wonLead, { lines, currency: "USD" });
    expect(p.lines).toEqual(lines);
    expect(p.currency).toBe("USD");
  });
});
