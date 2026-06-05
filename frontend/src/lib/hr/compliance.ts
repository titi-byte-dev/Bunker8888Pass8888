/**
 * Relatório de conformidade RGPD (HR-008) + acesso à cadeia de auditoria.
 *
 * O servidor só devolve metadados não-secretos (contagens + estado da cadeia);
 * o cliente formata-os num relatório imprimível (PDF via impressão do browser)
 * e reverifica a cadeia de auditoria localmente.
 */
import { loadSessionToken } from "$lib/session";
import { verifyAuditChain, type AuditEntry, type ChainResult } from "./audit";

export interface ComplianceReport {
  generatedAt: string;
  recordCount: number;
  activeFieldCount: number;
  shreddedFieldCount: number;
  certificateCount: number;
  auditEntryCount: number;
  auditChainValid: boolean;
  auditBrokenSeq: number;
}

interface ReportDTO {
  generated_at: string;
  record_count: number;
  active_field_count: number;
  shredded_field_count: number;
  certificate_count: number;
  audit_entry_count: number;
  audit_chain_valid: boolean;
  audit_broken_seq: number;
}

interface AuditEntryDTO {
  owner_id: string;
  seq: number;
  action: string;
  detail: string;
  occurred_at: string;
  prev_hash: string;
  entry_hash: string;
}

async function authedFetch(path: string): Promise<Response> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessao expirada — inicia sessao de novo.");
  const res = await fetch(path, { headers: { Authorization: `Bearer ${token}` } });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return res;
}

/** Obtém o relatório de conformidade (metadados + estado da cadeia). */
export async function fetchComplianceReport(): Promise<ComplianceReport> {
  const dto = (await (await authedFetch("/api/hr/compliance-report")).json()) as ReportDTO;
  return {
    generatedAt: dto.generated_at,
    recordCount: dto.record_count,
    activeFieldCount: dto.active_field_count,
    shreddedFieldCount: dto.shredded_field_count,
    certificateCount: dto.certificate_count,
    auditEntryCount: dto.audit_entry_count,
    auditChainValid: dto.audit_chain_valid,
    auditBrokenSeq: dto.audit_broken_seq,
  };
}

/** Obtém a cadeia de auditoria do utilizador. */
export async function fetchAuditLog(): Promise<AuditEntry[]> {
  const dto = (await (await authedFetch("/api/hr/audit")).json()) as { entries: AuditEntryDTO[] };
  return (dto.entries ?? []).map((e) => ({
    ownerId: e.owner_id,
    seq: e.seq,
    action: e.action,
    detail: e.detail,
    occurredAt: e.occurred_at,
    prevHash: e.prev_hash,
    entryHash: e.entry_hash,
  }));
}

/** Reverifica a cadeia localmente (não confia no servidor). */
export async function verifyLocally(entries: AuditEntry[]): Promise<ChainResult> {
  return verifyAuditChain(entries);
}
