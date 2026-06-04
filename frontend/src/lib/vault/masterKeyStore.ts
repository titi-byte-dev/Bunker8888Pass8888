/**
 * Armazén volátil da Master Key (VAULT-010).
 *
 * Didático: em JS/TS não há garantia de "apagar memória", mas minimizamos a
 * janela — a chave fica só neste módulo (nunca localStorage) e descartamos a
 * referência ao fim do turno ou no remote wipe.
 *
 * ⚠️ Segurança: `purge()` deve ser chamado fora do turno, no wipe remoto e ao
 * fechar sessão.
 */
import type { ShiftPolicy } from "./shift";
import { isWithinShift, msUntilShiftEnd } from "./shift";

let masterKey: CryptoKey | null = null;
let purgeTimer: ReturnType<typeof setTimeout> | null = null;

export function setMasterKey(key: CryptoKey): void {
  masterKey = key;
}

export function getMasterKey(): CryptoKey | null {
  return masterKey;
}

/** Descarta a referência à Master Key e cancela timers de expurgo. */
export function purgeMasterKey(): void {
  masterKey = null;
  if (purgeTimer) {
    clearTimeout(purgeTimer);
    purgeTimer = null;
  }
}

/**
 * Agenda expurgo automático ao fim da janela de turno actual.
 * Se já estiver fora do turno, expurga imediatamente.
 */
export function scheduleShiftPurge(policy: ShiftPolicy, onPurged?: () => void): void {
  if (purgeTimer) {
    clearTimeout(purgeTimer);
    purgeTimer = null;
  }
  if (!policy.enabled) return;

  const now = new Date();
  if (!isWithinShift(now, policy)) {
    purgeMasterKey();
    onPurged?.();
    return;
  }

  const ms = msUntilShiftEnd(now, policy);
  if (ms === null || ms <= 0) {
    purgeMasterKey();
    onPurged?.();
    return;
  }

  purgeTimer = setTimeout(() => {
    purgeMasterKey();
    onPurged?.();
  }, ms);
}
