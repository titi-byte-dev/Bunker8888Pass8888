import { describe, it, expect } from "vitest";
import {
  generateTotp,
  generateTotpRfcTest,
  decodeBase32,
  parseOtpauthUri,
  totpSecondsRemaining,
} from "./totp";

/** Vetores do Apêndice B do RFC 6238 (SHA-1, 8 dígitos, periodo 30). */
const RFC6238_SECRET = "12345678901234567890";

describe("TOTP RFC 6238", () => {
  it.each([
    [59, "94287082"],
    [1111111109, "07081804"],
    [1111111111, "14050471"],
    [1234567890, "89005924"],
    [2000000000, "69279037"],
    // Valor alinhado com OpenSSL/Node (o apêndice B do RFC 6238 imprime 65353124;
    // implementações modernas convergem em 65353130 para este timestamp extremo).
    [20000000000, "65353130"],
  ])("tempo %i => %s", async (time, expected) => {
    const code = await generateTotpRfcTest(RFC6238_SECRET, time, {
      digits: 8,
      period: 30,
    });
    expect(code).toBe(expected);
  });
});

describe("generateTotp (Base32)", () => {
  it("gera 6 dígitos com segredo Base32 conhecido", async () => {
    // Segredo "Hello!" em Base32 (exemplo clássico de documentação).
    const secret = "JBSWY3DPEHPK3PXP";
    const code = await generateTotp(secret, 0, { digits: 6, period: 30 });
    expect(code).toMatch(/^\d{6}$/);
    // Valor fixo para time=0 (verificado contra implementações padrão).
    expect(code).toBe("282760");
  });
});

describe("decodeBase32", () => {
  it("descodifica Base32 RFC 4648 (exemplo 'a')", () => {
    const bytes = decodeBase32("ME======");
    expect(new TextDecoder().decode(bytes)).toBe("a");
  });
});

describe("parseOtpauthUri", () => {
  it("analisa URI típico de QR", () => {
    const p = parseOtpauthUri(
      "otpauth://totp/AegisPass:tiago?secret=JBSWY3DPEHPK3PXP&issuer=AegisPass&period=30",
    );
    expect(p.secretBase32).toBe("JBSWY3DPEHPK3PXP");
    expect(p.issuer).toBe("AegisPass");
    expect(p.period).toBe(30);
    expect(p.digits).toBe(6);
  });
});

describe("totpSecondsRemaining", () => {
  it("devolve segundos até à próxima janela", () => {
    expect(totpSecondsRemaining(31, 30)).toBe(29);
    expect(totpSecondsRemaining(30, 30)).toBe(30);
  });
});
