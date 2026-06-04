/** Tipos de item do cofre (VAULT-005). */
export type VaultItemType = "login" | "note" | "card";

export interface LoginItem {
  kind: "login";
  title: string;
  username: string;
  password: string;
  url?: string;
  notes?: string;
}

export interface NoteItem {
  kind: "note";
  title: string;
  body: string;
}

export interface CardItem {
  kind: "card";
  title: string;
  number: string;
  expiry: string;
  cvv: string;
  holder?: string;
}

export type VaultPayload = LoginItem | NoteItem | CardItem;

export interface VaultItemMeta {
  id: string;
  type: VaultItemType;
  blob?: string;
  created_at: string;
  updated_at: string;
}

export interface VaultItemInput {
  type: VaultItemType;
  blob: string;
}

/** Eventos push do WebSocket (VAULT-006) — só metadados, nunca conteúdo em claro. */
export type VaultSyncEventType =
  | "vault.item.created"
  | "vault.item.updated"
  | "vault.item.deleted";

export interface VaultSyncEvent {
  type: VaultSyncEventType;
  item_id: string;
  item_type?: VaultItemType;
  updated_at?: string;
}
