import { describe, expect, it } from "vitest";
import { closeConfirm, confirmDialog, confirmState } from "./confirm";

describe("confirmDialog (UI-017)", () => {
  it("abre com opções e resolve true ao confirmar", async () => {
    const p = confirmDialog({
      title: "Apagar?",
      message: "Isto é irreversível.",
      variant: "danger",
    });
    expect(confirmState().open).toBe(true);
    expect(confirmState().options?.title).toBe("Apagar?");
    closeConfirm(true);
    await expect(p).resolves.toBe(true);
    expect(confirmState().open).toBe(false);
  });

  it("resolve false ao cancelar", async () => {
    const p = confirmDialog({ title: "T", message: "M" });
    closeConfirm(false);
    await expect(p).resolves.toBe(false);
  });

  it("cancela diálogo anterior se abrir outro", async () => {
    const first = confirmDialog({ title: "1", message: "a" });
    const second = confirmDialog({ title: "2", message: "b" });
    await expect(first).resolves.toBe(false);
    closeConfirm(true);
    await expect(second).resolves.toBe(true);
  });
});
