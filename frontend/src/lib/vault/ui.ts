/**
 * Acesso ao cofre na UI (UI-004) — requer sessão + Master Key desbloqueada.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "./masterKeyStore";
import { VaultAPI } from "./api";
import { blobFromBase64, openItem } from "./items";
import type { LoginItem, VaultItemMeta } from "./types";

export type DecodedLogin = {
  meta: VaultItemMeta;
  login: LoginItem;
};

export function requireVaultAccess(): { api: VaultAPI; key: CryptoKey; token: string } {
  const token = loadSessionToken();
  const key = getMasterKey();
  if (!token || !key) {
    throw new Error("Cofre bloqueado — inicia sessão e desbloqueia em /auth/unlock");
  }
  return { api: new VaultAPI("", token), key, token };
}

/** Lista logins e descifra blobs localmente (Zero-Knowledge). */
export async function loadDecodedLogins(): Promise<DecodedLogin[]> {
  const { api, key } = requireVaultAccess();
  const metas = await api.list("login");
  const out: DecodedLogin[] = [];

  for (const meta of metas) {
    if (!meta.blob) continue;
    const login = (await openItem(key, blobFromBase64(meta.blob))) as LoginItem;
    out.push({ meta, login });
  }

  return out.sort((a, b) => a.login.title.localeCompare(b.login.title, "pt"));
}

/** Descifra um login pelo id. */
export async function loadDecodedLogin(id: string): Promise<DecodedLogin> {
  const { api, key } = requireVaultAccess();
  const meta = await api.get(id);
  if (!meta.blob) throw new Error("Item sem conteúdo cifrado");
  const login = (await openItem(key, blobFromBase64(meta.blob))) as LoginItem;
  return { meta, login };
}

/** Pesquisa local por título, utilizador ou URL. */
export function filterLogins(items: DecodedLogin[], query: string): DecodedLogin[] {
  const q = query.trim().toLowerCase();
  if (!q) return items;
  return items.filter(({ login }) => {
    const hay = `${login.title} ${login.username} ${login.url ?? ""}`.toLowerCase();
    return hay.includes(q);
  });
}
