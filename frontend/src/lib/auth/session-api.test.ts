import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { fetchSessionProfile } from "./session-api";

function mockStorage(): Storage {
  const map = new Map<string, string>();
  return {
    get length() {
      return map.size;
    },
    clear: () => map.clear(),
    getItem: (k) => map.get(k) ?? null,
    key: (i) => [...map.keys()][i] ?? null,
    removeItem: (k) => {
      map.delete(k);
    },
    setItem: (k, v) => {
      map.set(k, v);
    },
  };
}

describe("fetchSessionProfile", () => {
  beforeEach(() => {
    vi.stubGlobal("sessionStorage", mockStorage());
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("devolve email normalizado com token válido", async () => {
    sessionStorage.setItem("aegis:session-token", "abc");
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({
        ok: true,
        json: async () => ({ email: " Dev@Test.COM " }),
      })),
    );

    const profile = await fetchSessionProfile("");
    expect(profile.email).toBe("dev@test.com");
    expect(fetch).toHaveBeenCalledWith("/api/auth/session", {
      headers: { Authorization: "Bearer abc" },
    });
  });

  it("falha sem token em sessionStorage", async () => {
    await expect(fetchSessionProfile()).rejects.toThrow("Sessão inválida");
  });
});
