import { afterEach, describe, expect, it, vi } from "vitest";
import { checkPasswordBreached, sha1HexUpper } from "./breach";

describe("sha1HexUpper (DW-001)", () => {
  it("devolve SHA-1 uppercase conhecido", async () => {
    // SHA-1 de "password" — vector público HIBP
    const hash = await sha1HexUpper("password");
    expect(hash).toBe("5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8");
  });
});

describe("checkPasswordBreached", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("deteta sufixo na resposta range (k-anonymity local)", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (url: string) => {
        expect(url).toBe("https://api.pwnedpasswords.com/range/5BAA6");
        return {
          ok: true,
          text: async () =>
            "0018A45C4D1DEF81644B54AB7F9695E1CAD1:4\n1E4C9B93F3F0682250B6CF8331B7EE68FD8:3730473\n",
        };
      }),
    );

    const result = await checkPasswordBreached("password");
    expect(result.breached).toBe(true);
    expect(result.exposureCount).toBe(3730473);
  });

  it("password limpa quando sufixo ausente", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        text: async () => "AAAAA:1\n",
      })),
    );

    const result = await checkPasswordBreached("UnicaForte#2026$xyz");
    expect(result.breached).toBe(false);
    expect(result.exposureCount).toBe(0);
  });
});
