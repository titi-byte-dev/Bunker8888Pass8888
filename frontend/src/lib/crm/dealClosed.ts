/**
 * Fluxo deal_closed (CRM-003) — ponte CRM -> Faturacao.
 *
 * Didatico: quando um lead chega ao estagio "won" (negocio fechado), o sistema
 * sugere emitir uma PRO-FORMA pre-preenchida com os dados do cliente. Tudo
 * acontece no cliente, com dados ja decifrados — o servidor so vera o blob
 * cifrado da pro-forma e atribuira o numero legal.
 *
 *   ┌────────────┐  stage=won   ┌────────────────┐  issue   ┌──────────────┐
 *   │  Lead CRM  │─────────────▶│ proformaFromLead│────────▶│  PF 2026/000n │
 *   │ (cifrado)  │  decifrado   │  (InvoicePayload)│  cifrado└──────────────┘
 *   └────────────┘              └────────────────┘
 */
import type { Lead } from "./leads";
import type { InvoiceLine, InvoicePayload } from "$lib/fin/invoices";

/** Um lead esta "fechado" (ganho) quando chega ao estagio won. */
export function isDealClosed(lead: Pick<Lead, "stage">): boolean {
  return lead.stage === "won";
}

export interface ProformaOptions {
  currency?: string;
  /** Linhas iniciais; por omissao uma linha-rascunho a 0. */
  lines?: InvoiceLine[];
  /** Momento de emissao (ISO). Por omissao, agora. */
  issuedAt?: string;
}

const draftLine: InvoiceLine = {
  description: "A descrever",
  quantity: 1,
  unitPrice: 0,
  vatRate: 23,
};

/**
 * Constroi um InvoicePayload de pro-forma a partir de um lead ganho.
 * Pre-preenche o cliente e liga sourceLeadId para rastreabilidade.
 */
export function proformaFromLead(lead: Lead, opts: ProformaOptions = {}): InvoicePayload {
  return {
    client: {
      name: lead.company?.trim() ? lead.company : lead.name,
      email: lead.email || undefined,
    },
    lines: opts.lines && opts.lines.length > 0 ? opts.lines : [{ ...draftLine }],
    currency: opts.currency ?? "EUR",
    issuedAt: opts.issuedAt ?? new Date().toISOString(),
    notes: `Pro-forma gerada do lead "${lead.name}".`,
    sourceLeadId: lead.id,
  };
}
