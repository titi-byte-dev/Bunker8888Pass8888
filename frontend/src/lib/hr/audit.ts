/**
 * Verificação da cadeia de auditoria imutável (HR-002).
 *
 * Cada entrada encadeia ao hash da anterior. Recalculamos toda a cadeia no
 * dispositivo: se UMA entrada foi adulterada, o seu entry_hash deixa de bater
 * certo e a verificação aponta a primeira posição partida — sem confiar no
 * servidor.
 *
 *   entry_hash = sha256( v1|seq|owner|action|detail|occurred_at|prev_hash )
 *   prev_hash da 1.ª entrada = "GENESIS"
 */
import { toBytes } from "$lib/crypto";

export const AUDIT_GENESIS = "GENESIS";

export interface AuditEntry {
  ownerId: string;
  seq: number;
  action: string;
  detail: string;
  occurredAt: string;
  prevHash: string;
  entryHash: string;
}

async function sha256hex(text: string): Promise<string> {
  const d = await crypto.subtle.digest("SHA-256", toBytes(text));
  return Array.from(new Uint8Array(d))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

function canonical(e: AuditEntry, prevHash: string): string {
  return `v1|${e.seq}|${e.ownerId}|${e.action}|${e.detail}|${e.occurredAt}|${prevHash}`;
}

export interface ChainResult {
  valid: boolean;
  /** seq da primeira entrada partida (0 se a cadeia estiver intacta). */
  brokenSeq: number;
}

/**
 * Recalcula a cadeia: confirma que cada seq incrementa, que o prev_hash bate com
 * o entry_hash anterior e que o entry_hash reproduz o canónico.
 */
export async function verifyAuditChain(entries: AuditEntry[]): Promise<ChainResult> {
  let prev = AUDIT_GENESIS;
  let expectedSeq = 1;
  for (const e of entries) {
    if (e.seq !== expectedSeq || e.prevHash !== prev) {
      return { valid: false, brokenSeq: e.seq };
    }
    const want = await sha256hex(canonical(e, prev));
    if (want !== e.entryHash) {
      return { valid: false, brokenSeq: e.seq };
    }
    prev = e.entryHash;
    expectedSeq += 1;
  }
  return { valid: true, brokenSeq: 0 };
}

/** Etiquetas legíveis para as acções registadas. */
export const ACTION_LABELS: Record<string, string> = {
  "record.create": "Ficha criada",
  "record.delete": "Ficha apagada",
  "field.put": "Campo gravado",
  "field.delete": "Campo removido",
  "field.shred": "Campo eliminado (shred)",
};

export function actionLabel(action: string): string {
  return ACTION_LABELS[action] ?? action;
}
