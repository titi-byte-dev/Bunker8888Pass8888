/**
 * Remote wipe local (VAULT-012).
 *
 * Didático: o servidor não consegue apagar ficheiros no telemóvel do utilizador
 * (sandbox do browser/OS). O push WebSocket é a ordem; esta função executa a
 * limpeza no lado cliente — cache cifrada + tokens locais.
 *
 * ⚠️ Segurança: mesmo offline, a revogação de sessões no servidor impede
 * novos acessos; a Master Key em memória deve ser descartada pelo callback.
 */
import { purgeMasterKey } from "./masterKeyStore";

export const AEGIS_STORAGE_PREFIX = "aegis:";

/** Apaga entradas localStorage/sessionStorage com prefixo AegisPass. */
export function clearLocalVaultCache(): void {
  for (const storage of [localStorage, sessionStorage]) {
    const keys: string[] = [];
    for (let i = 0; i < storage.length; i++) {
      const key = storage.key(i);
      if (key?.startsWith(AEGIS_STORAGE_PREFIX)) keys.push(key);
    }
    for (const key of keys) storage.removeItem(key);
  }
}

export interface RemoteWipeResult {
  devices_notified: number;
  sessions_revoked: number;
}

/** Pedido autenticado de wipe nos próprios dispositivos. */
export async function requestSelfRemoteWipe(
  baseURL: string,
  token: string,
  reason = "",
): Promise<RemoteWipeResult> {
  const res = await globalThis.fetch(`${baseURL}/api/security/remote-wipe/self`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ reason }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return res.json() as Promise<RemoteWipeResult>;
}

/** Executa a limpeza local após receber evento WebSocket ou acção manual. */
export function executeLocalWipe(onMasterKeyDiscarded?: () => void): void {
  clearLocalVaultCache();
  purgeMasterKey();
  onMasterKeyDiscarded?.();
}
