/**
 * Leads CRM (CRM-001) — modelo + cifragem Zero-Knowledge.
 */
import { decrypt, encrypt, fromBytes, toBytes, type Bytes } from "$lib/crypto";

export type LeadStage = "new" | "contacted" | "qualified" | "proposal" | "won" | "lost";

export const LEAD_STAGES: { id: LeadStage; label: string }[] = [
  { id: "new", label: "Novo" },
  { id: "contacted", label: "Contactado" },
  { id: "qualified", label: "Qualificado" },
  { id: "proposal", label: "Proposta" },
  { id: "won", label: "Ganho" },
  { id: "lost", label: "Perdido" },
];

export interface LeadPayload {
  name: string;
  email: string;
  company?: string;
  stage: LeadStage;
  notes?: string;
  source?: string;
}

export interface Lead extends LeadPayload {
  id: string;
}

export async function encryptLead(masterKey: CryptoKey, payload: LeadPayload): Promise<Bytes> {
  return encrypt(masterKey, toBytes(JSON.stringify(payload)));
}

export async function decryptLead(masterKey: CryptoKey, blob: Bytes): Promise<LeadPayload> {
  return JSON.parse(fromBytes(await decrypt(masterKey, blob))) as LeadPayload;
}
