import { describe, expect, it } from "vitest";
import { AUDIT_GENESIS, actionLabel, verifyAuditChain, type AuditEntry } from "./audit";
import { toBytes } from "$lib/crypto";

const OWNER = "owner-uuid";

async function sha256hex(text: string): Promise<string> {
  const d = await crypto.subtle.digest("SHA-256", toBytes(text));
  return Array.from(new Uint8Array(d))
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

/** Constrói uma cadeia válida a partir de uma lista de (action, detail). */
async function buildChain(steps: [string, string][]): Promise<AuditEntry[]> {
  const out: AuditEntry[] = [];
  let prev = AUDIT_GENESIS;
  let seq = 1;
  for (const [action, detail] of steps) {
    const occurredAt = `2026-06-05T10:00:0${seq}Z`;
    const entryHash = await sha256hex(
      `v1|${seq}|${OWNER}|${action}|${detail}|${occurredAt}|${prev}`,
    );
    out.push({ ownerId: OWNER, seq, action, detail, occurredAt, prevHash: prev, entryHash });
    prev = entryHash;
    seq += 1;
  }
  return out;
}

describe("HR-002 cadeia de auditoria (verificação encadeada)", () => {
  it("verifica uma cadeia íntegra", async () => {
    const chain = await buildChain([
      ["record.create", "rec1"],
      ["field.put", "salary"],
      ["field.shred", "salary"],
    ]);
    expect(await verifyAuditChain(chain)).toEqual({ valid: true, brokenSeq: 0 });
  });

  it("deteta adulteração do detail de uma entrada antiga", async () => {
    const chain = await buildChain([
      ["record.create", "rec1"],
      ["field.put", "salary"],
    ]);
    chain[0].detail = "forjado"; // muda o conteúdo sem recalcular o hash
    expect(await verifyAuditChain(chain)).toEqual({ valid: false, brokenSeq: 1 });
  });

  it("deteta um elo partido (prev_hash trocado)", async () => {
    const chain = await buildChain([
      ["record.create", "rec1"],
      ["field.put", "salary"],
    ]);
    chain[1].prevHash = "0".repeat(64);
    expect(await verifyAuditChain(chain)).toEqual({ valid: false, brokenSeq: 2 });
  });

  it("deteta um seq fora de ordem", async () => {
    const chain = await buildChain([
      ["record.create", "rec1"],
      ["field.put", "salary"],
    ]);
    chain[1].seq = 5;
    expect(await verifyAuditChain(chain)).toEqual({ valid: false, brokenSeq: 5 });
  });

  it("etiqueta acções conhecidas", () => {
    expect(actionLabel("field.shred")).toBe("Campo eliminado (shred)");
    expect(actionLabel("custom.x")).toBe("custom.x");
  });
});
