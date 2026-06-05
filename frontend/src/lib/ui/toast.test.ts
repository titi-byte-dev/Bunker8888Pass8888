import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { dismissToast, pushToast, toast, toastQueue, toastStore } from "./toast";

describe("toast store (UI-017)", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    toastStore.set([]);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("adiciona toast à fila", () => {
    pushToast("Guardado", { variant: "success" });
    expect(toastQueue()).toHaveLength(1);
    expect(toastQueue()[0].message).toBe("Guardado");
    expect(toastQueue()[0].variant).toBe("success");
  });

  it("remove toast manualmente", () => {
    const id = pushToast("Temp");
    dismissToast(id);
    expect(toastQueue()).toHaveLength(0);
  });

  it("auto-dismiss após durationMs", () => {
    pushToast("Efémero", { durationMs: 3000 });
    expect(toastQueue()).toHaveLength(1);
    vi.advanceTimersByTime(3000);
    expect(toastQueue()).toHaveLength(0);
  });

  it("atalhos toast.* usam variantes correctas", () => {
    toast.error("Falhou");
    expect(toastQueue()[0].variant).toBe("error");
    dismissToast(toastQueue()[0].id);
    toast.warning("Cuidado");
    expect(toastQueue()[0].variant).toBe("warning");
  });
});
