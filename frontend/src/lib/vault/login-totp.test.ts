import { describe, it, expect } from "vitest";
import { loginTotpCode } from "./login-totp";
import type { LoginItem } from "./types";

describe("loginTotpCode (VAULT-009)", () => {
  it("devolve null sem segredo TOTP", async () => {
    const login: LoginItem = {
      kind: "login",
      title: "App",
      username: "u",
      password: "p",
    };
    expect(await loginTotpCode(login)).toBeNull();
  });

  it("gera código quando há segredo Base32", async () => {
    const login: LoginItem = {
      kind: "login",
      title: "App",
      username: "u",
      password: "p",
      totpSecretBase32: "JBSWY3DPEHPK3PXP",
    };
    const code = await loginTotpCode(login, 0);
    expect(code).toBe("282760");
  });
});
