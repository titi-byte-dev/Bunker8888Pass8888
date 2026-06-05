import { describe, expect, it } from "vitest";
import { buildPurchaseOrderDraft } from "./purchaseOrder";

describe("purchase order draft (AGENT-008)", () => {
  it("sugere quantidade para repor acima do mínimo", () => {
    const d = buildPurchaseOrderDraft({
      item_id: "x",
      name: "Toner",
      sku: "TN-1",
      quantity: 2,
      reorder_level: 5,
      unit: "un",
    });
    expect(d.orderQty).toBe(8);
    expect(d.itemName).toBe("Toner");
    expect(d.note).toContain("Toner");
  });
});
