import { describe, expect, it } from "vitest";
import {
  buildNoteLink,
  decryptNoteContent,
  encryptNoteContent,
  generateNoteKey,
  parseNoteFragment,
} from "./burnNote";

describe("SHARE-005 notas auto-destrutivas (cripto)", () => {
  it("gera uma chave de nota de 32 bytes (AES-256)", () => {
    expect(generateNoteKey()).toHaveLength(32);
  });

  it("cifra e decifra uma nota sem passphrase (round-trip)", async () => {
    const key = generateNoteKey();
    const { ciphertext, salt } = await encryptNoteContent(key, "a password e tr0ub4dor");
    expect(salt).toBeNull();
    expect(ciphertext).not.toContain("tr0ub4dor"); // esta cifrado
    expect(await decryptNoteContent(key, ciphertext, null, salt)).toBe("a password e tr0ub4dor");
  });

  it("cifra e decifra uma nota com passphrase (2.a camada)", async () => {
    const key = generateNoteKey();
    const { ciphertext, salt } = await encryptNoteContent(key, "segredo duplo", "abracadabra");
    expect(salt).not.toBeNull();
    expect(await decryptNoteContent(key, ciphertext, "abracadabra", salt)).toBe("segredo duplo");
  });

  it("nao abre com a passphrase errada (auth GCM)", async () => {
    const key = generateNoteKey();
    const { ciphertext, salt } = await encryptNoteContent(key, "x", "certa");
    await expect(decryptNoteContent(key, ciphertext, "errada", salt)).rejects.toThrow();
  });

  it("nao abre sem a chave correta do fragmento", async () => {
    const { ciphertext, salt } = await encryptNoteContent(generateNoteKey(), "y");
    await expect(decryptNoteContent(generateNoteKey(), ciphertext, null, salt)).rejects.toThrow();
  });

  it("constroi e analisa o fragmento do link (sem passphrase)", () => {
    const key = generateNoteKey();
    const link = buildNoteLink("https://host", "abc123", key, null);
    expect(link).toContain("/n/abc123#k=");
    const frag = parseNoteFragment(new URL(link).hash);
    expect(frag).not.toBeNull();
    expect([...frag!.key]).toEqual([...key]);
    expect(frag!.requiresPassphrase).toBe(false);
    expect(frag!.salt).toBeNull();
  });

  it("constroi e analisa o fragmento do link (com passphrase)", () => {
    const key = generateNoteKey();
    const salt = new Uint8Array(16).fill(7);
    const link = buildNoteLink("https://host", "abc", key, salt as never);
    expect(link).toContain("&p=1&s=");
    const frag = parseNoteFragment(new URL(link).hash);
    expect(frag).not.toBeNull();
    expect(frag!.requiresPassphrase).toBe(true);
    expect(frag!.salt).not.toBeNull();
    expect([...frag!.salt!]).toEqual([...salt]);
  });

  it("devolve null se faltar a chave no fragmento", () => {
    expect(parseNoteFragment("#nada=1")).toBeNull();
  });
});
