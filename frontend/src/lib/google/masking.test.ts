import { describe, expect, it } from "vitest";
import { maskSensitiveText, unmaskText } from "./masking";

describe("masking (GOOGLE-003 dev)", () => {
  it("mascara NIF e permite desmascarar", () => {
    const src = "Cliente NIF 123456789 pagou IBAN PT500026000099876543210 10";
    const { masked, tokens } = maskSensitiveText(src);
    expect(masked).not.toContain("123456789");
    expect(Object.keys(tokens).length).toBeGreaterThan(0);
    expect(unmaskText(masked, tokens)).toContain("123456789");
  });
});
