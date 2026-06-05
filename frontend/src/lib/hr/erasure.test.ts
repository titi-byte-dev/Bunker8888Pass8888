import { describe, expect, it } from "vitest";
import {
  certificateToJSON,
  verifyCertificate,
  type ErasureCertificate,
} from "./erasure";
import { toBytes } from "$lib/crypto";

/** Calcula o cert_hash como o servidor faria, para montar certificados de teste. */
async function sha256hex(text: string): Promise<string> {
  const d = await crypto.subtle.digest("SHA-256", toBytes(text));
  return Array.from(new Uint8Array(d))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

async function makeCert(over: Partial<ErasureCertificate> = {}): Promise<ErasureCertificate> {
  const base: Omit<ErasureCertificate, "certHash"> = {
    id: "cert-1",
    ownerId: "owner-uuid",
    recordId: "rec-uuid",
    fieldName: "salary",
    valueDigest: "abcd1234",
    shreddedAt: "2026-06-05T10:00:00Z",
    issuedAt: "2026-06-05T10:00:00Z",
    ...over,
  };
  const canonical = `v1|${base.recordId}|${base.fieldName}|${base.valueDigest}|${base.shreddedAt}|${base.ownerId}`;
  return { ...base, certHash: await sha256hex(canonical) };
}

describe("HR-004 certificados de eliminação (verificação)", () => {
  it("verifica um certificado íntegro", async () => {
    const cert = await makeCert();
    expect(await verifyCertificate(cert)).toBe(true);
  });

  it("rejeita um certificado adulterado (campo trocado)", async () => {
    const cert = await makeCert();
    expect(await verifyCertificate({ ...cert, fieldName: "nif" })).toBe(false);
  });

  it("rejeita se o value_digest não corresponder", async () => {
    const cert = await makeCert();
    expect(await verifyCertificate({ ...cert, valueDigest: "deadbeef" })).toBe(false);
  });

  it("rejeita se o cert_hash for falsificado", async () => {
    const cert = await makeCert();
    expect(await verifyCertificate({ ...cert, certHash: "0".repeat(64) })).toBe(false);
  });

  it("serializa para JSON auditável com os campos da prova", async () => {
    const cert = await makeCert();
    const json = JSON.parse(certificateToJSON(cert));
    expect(json.kind).toBe("aegispass.erasure-certificate");
    expect(json.cert_hash).toBe(cert.certHash);
    expect(json.field_name).toBe("salary");
  });
});
