/**
 * Rascunho de ordem de compra (AGENT-008) — gerado no cliente após aprovação.
 *
 * Didático: a ordem real ainda não vai para fornecedores; este stub documenta
 * a quantidade sugerida com base no nível de reordenação.
 */
export type PurchaseOrderDraft = {
  itemId: string;
  itemName: string;
  sku: string;
  orderQty: number;
  unit: string;
  note: string;
  createdAt: string;
};

export function buildPurchaseOrderDraft(payload: Record<string, unknown>): PurchaseOrderDraft {
  const reorder = Number(payload.reorder_level ?? 0);
  const qty = Number(payload.quantity ?? 0);
  const orderQty = Math.max(reorder * 2 - qty, reorder);
  const name = String(payload.name ?? "Artigo");
  const sku = String(payload.sku ?? "");
  const unit = String(payload.unit ?? "un");
  return {
    itemId: String(payload.item_id ?? ""),
    itemName: name,
    sku,
    orderQty,
    unit,
    note: `Repor stock de «${name}» (${qty} ${unit} em armazém, mínimo ${reorder}).`,
    createdAt: new Date().toISOString(),
  };
}
