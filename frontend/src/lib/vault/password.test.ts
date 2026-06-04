import { describe, it, expect } from "vitest";
import { generatePassword, estimateEntropyBits } from "./password";

describe("generatePassword (VAULT-007)", () => {
  it("gera password com o comprimento pedido", () => {
    expect(generatePassword({ length: 32 }).length).toBe(32);
  });

  it("inclui maiúsculas, minúsculas, dígitos e símbolos por omissão", () => {
    const pw = generatePassword({ length: 24 });
    expect(/[A-Z]/.test(pw)).toBe(true);
    expect(/[a-z]/.test(pw)).toBe(true);
    expect(/[0-9]/.test(pw)).toBe(true);
    expect(/[!@#$%^&*()\-_=+[\]{}]/.test(pw)).toBe(true);
  });

  it("gera passwords diferentes em chamadas sucessivas", () => {
    const a = generatePassword();
    const b = generatePassword();
    expect(a).not.toBe(b);
  });

  it("rejeita length menor que nº de pools activos", () => {
    expect(() =>
      generatePassword({ length: 2, uppercase: true, lowercase: true, digits: true, symbols: true }),
    ).toThrow();
  });

  it("estimateEntropyBits devolve valor positivo", () => {
    expect(estimateEntropyBits({ length: 20 })).toBeGreaterThan(100);
  });
});
