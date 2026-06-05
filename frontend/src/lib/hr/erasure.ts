/**
 * Verificação dos certificados de eliminação (HR-004).
 *
 * Um certificado prova que um campo foi crypto-shredded (HR-003): a sua chave
 * foi destruída e o valor é matematicamente irrecuperável. A prova é uma
 * impressão digital (sha256) sobre os factos da eliminação. Qualquer pessoa a
 * recalcula — sem ver o conteúdo original — e confirma que bate certo.
 *
 *   canonico = v1|record_id|field_name|value_digest|shredded_at|owner_id
 *   cert_hash = sha256(canonico)        ← tem de ser igual ao do servidor
 */
import { toBytes } from "$lib/crypto";

export interface ErasureCertificate {
  id: string;
  ownerId: string;
  recordId: string;
  fieldName: string;
  valueDigest: string;
  shreddedAt: string;
  certHash: string;
  issuedAt: string;
}

/** Calcula sha256(texto) em hex (igual ao do servidor). */
async function sha256hex(text: string): Promise<string> {
  const digest = await crypto.subtle.digest("SHA-256", toBytes(text));
  return Array.from(new Uint8Array(digest))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

/** Reconstrói o canónico versionado do certificado. */
function canonical(cert: ErasureCertificate): string {
  return `v1|${cert.recordId}|${cert.fieldName}|${cert.valueDigest}|${cert.shreddedAt}|${cert.ownerId}`;
}

/**
 * Verifica um certificado: recalcula o cert_hash a partir dos seus campos e
 * compara com o emitido. true = prova íntegra.
 */
export async function verifyCertificate(cert: ErasureCertificate): Promise<boolean> {
  const recomputed = await sha256hex(canonical(cert));
  return recomputed === cert.certHash;
}

/** Serializa um certificado para download (JSON legível, auditável). */
export function certificateToJSON(cert: ErasureCertificate): string {
  return JSON.stringify(
    {
      version: "v1",
      kind: "aegispass.erasure-certificate",
      id: cert.id,
      owner_id: cert.ownerId,
      record_id: cert.recordId,
      field_name: cert.fieldName,
      value_digest: cert.valueDigest,
      shredded_at: cert.shreddedAt,
      cert_hash: cert.certHash,
      issued_at: cert.issuedAt,
    },
    null,
    2,
  );
}
