import { describe, expect, it, vi, beforeEach } from "vitest";
import { AEGIS_STORAGE_PREFIX, clearLocalVaultCache, executeLocalWipe } from "./wipe";

/** Mock mínimo de Storage para correr em ambiente Node (vitest). */
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

describe("clearLocalVaultCache", () => {
  beforeEach(() => {
    vi.stubGlobal("localStorage", mockStorage());
    vi.stubGlobal("sessionStorage", mockStorage());
  });

  it("remove apenas chaves com prefixo aegis:", () => {
    localStorage.setItem(`${AEGIS_STORAGE_PREFIX}token`, "x");
    localStorage.setItem("other-app", "y");
    sessionStorage.setItem(`${AEGIS_STORAGE_PREFIX}cache`, "z");

    clearLocalVaultCache();

    expect(localStorage.getItem(`${AEGIS_STORAGE_PREFIX}token`)).toBeNull();
    expect(localStorage.getItem("other-app")).toBe("y");
    expect(sessionStorage.getItem(`${AEGIS_STORAGE_PREFIX}cache`)).toBeNull();
  });
});

describe("executeLocalWipe", () => {
  beforeEach(() => {
    vi.stubGlobal("localStorage", mockStorage());
    vi.stubGlobal("sessionStorage", mockStorage());
  });

  it("limpa cache e chama callback da Master Key", () => {
    localStorage.setItem(`${AEGIS_STORAGE_PREFIX}mk`, "secret");
    const discard = vi.fn();

    executeLocalWipe(discard);

    expect(localStorage.getItem(`${AEGIS_STORAGE_PREFIX}mk`)).toBeNull();
    expect(discard).toHaveBeenCalledOnce();
  });
});
