import { describe, expect, it } from "vitest";
import { decryptField, encryptField, fieldLabel } from "./employees";
import { fromBytes } from "$lib/crypto";

/** Master Key de teste (AES-GCM 256, com encrypt/decrypt — como a derivada). */
async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, [
    "encrypt",
    "decrypt",
  ]);
}

describe("HR-001 fichas de empregado (cripto campo-a-campo)", () => {
  it("cifra e decifra um campo (round-trip)", async () => {
    const mk = await makeMasterKey();
    const { valueBlob, wrappedKey } = await encryptField(mk, "Joana Silva");
    expect(fromBytes(valueBlob)).not.toContain("Joana"); // esta cifrado
    expect(await decryptField(mk, valueBlob, wrappedKey)).toBe("Joana Silva");
  });

  it("cada campo recebe uma chave embrulhada diferente", async () => {
    const mk = await makeMasterKey();
    const a = await encryptField(mk, "mesmo valor");
    const b = await encryptField(mk, "mesmo valor");
    // Chaves de campo distintas → wrapped_key distintos mesmo com valor igual.
    expect(Array.from(a.wrappedKey)).not.toEqual(Array.from(b.wrappedKey));
    expect(Array.from(a.valueBlob)).not.toEqual(Array.from(b.valueBlob));
  });

  it("nao decifra com a Master Key errada", async () => {
    const mk = await makeMasterKey();
    const other = await makeMasterKey();
    const { valueBlob, wrappedKey } = await encryptField(mk, "salario secreto");
    await expect(decryptField(other, valueBlob, wrappedKey)).rejects.toThrow();
  });

  it("simula crypto-shredding: sem a chave embrulhada o valor e irrecuperavel", async () => {
    const mk = await makeMasterKey();
    const { valueBlob, wrappedKey } = await encryptField(mk, "dado RGPD");
    // "Destruir" a chave do campo (HR-003) deixa o value_blob como lixo eterno.
    const shreddedKey = new Uint8Array(wrappedKey.length); // zeros = chave perdida
    await expect(decryptField(mk, valueBlob, shreddedKey as typeof wrappedKey)).rejects.toThrow();
    // O value_blob original mantem-se, mas ja nao abre.
    expect(valueBlob.length).toBeGreaterThan(0);
  });

  it("etiqueta campos conhecidos e devolve o nome cru nos personalizados", () => {
    expect(fieldLabel("salary")).toBe("Salário");
    expect(fieldLabel("campo_custom")).toBe("campo_custom");
  });
});
