/**
 * Ponte entre LoginItem e TOTP (VAULT-009).
 */
import { generateTotp, totpSecondsRemaining, type TotpOptions } from "./totp";
import type { LoginItem } from "./types";

/** Gera o código 2FA actual para um login, ou null se não tiver TOTP. */
export async function loginTotpCode(
  login: LoginItem,
  unixTime?: number,
  opts?: TotpOptions,
): Promise<string | null> {
  if (!login.totpSecretBase32?.trim()) return null;
  return generateTotp(login.totpSecretBase32, unixTime, opts);
}

/** Segundos até o código TOTP do login expirar. */
export function loginTotpSecondsRemaining(
  login: LoginItem,
  unixTime?: number,
  opts?: TotpOptions,
): number | null {
  if (!login.totpSecretBase32?.trim()) return null;
  return totpSecondsRemaining(unixTime, opts?.period ?? 30);
}
