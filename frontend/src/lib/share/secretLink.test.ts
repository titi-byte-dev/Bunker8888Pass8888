import { describe, expect, it } from "vitest";
import {
  buildSecretLink,
  decryptSecret,
  encryptSecret,
  fromBase64Url,
  generateLinkKey,
  keyFromFragment,
  toBase64Url,
} from "./secretLink";

describe("SHARE-003 secret links (cripto)", () => {
  it("gera uma chave de 32 bytes (AES-256)", () => {
    expect(generateLinkKey()).toHaveLength(32);
  });

  it("cifra e decifra o segredo (round-trip)", async () => {
    const key = generateLinkKey();
    const ct = await encryptSecret(key, "password-super-secreta");
    expect(ct).not.toContain("password");
    expect(await decryptSecret(key, ct)).toBe("password-super-secreta");
  });

  it("não decifra com a chave errada (auth GCM)", async () => {
    const ct = await encryptSecret(generateLinkKey(), "x");
    await expect(decryptSecret(generateLinkKey(), ct)).rejects.toThrow();
  });

  it("codifica/descodifica base64url sem padding", () => {
    const key = generateLinkKey();
    const s = toBase64Url(key);
    expect(s).not.toMatch(/[+/=]/); // url-safe, sem padding
    expect([...fromBase64Url(s)]).toEqual([...key]);
  });

  it("constrói um link e recupera a chave do fragmento", () => {
    const key = generateLinkKey();
    const link = buildSecretLink("https://host.pt", "abc123", key);
    expect(link).toContain("https://host.pt/s/abc123#k=");

    const hash = link.slice(link.indexOf("#"));
    const recovered = keyFromFragment(hash);
    expect(recovered).not.toBeNull();
    expect([...recovered!]).toEqual([...key]);
  });

  it("devolve null para fragmento ausente ou inválido", () => {
    expect(keyFromFragment("")).toBeNull();
    expect(keyFromFragment("#nada=1")).toBeNull();
    expect(keyFromFragment("#k=demasiado-curto")).toBeNull();
  });

  it("a chave abre o que cifrámos, ponta-a-ponta", async () => {
    const key = generateLinkKey();
    const ct = await encryptSecret(key, "segredo-final");
    const link = buildSecretLink("https://h", "id", key);
    const recovered = keyFromFragment(link.slice(link.indexOf("#")))!;
    expect(await decryptSecret(recovered, ct)).toBe("segredo-final");
  });
});
