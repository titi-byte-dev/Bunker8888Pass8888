/**
 * Acções partilhadas dos fluxos Auth (UI-003).
 */
import { goto } from "$app/navigation";
import {
  deriveMasterKeyBytes,
  fetchKdfParams,
  importAesKeyFromBytes,
  loginAfterRegister,
  loginUser,
  registerUser,
} from "$lib/auth";
import { loginWithPasskey } from "$lib/passkey";
import { normalizeEmail } from "$lib/auth/http";
import { saveSessionToken, saveUserEmail } from "$lib/session";
import { setMasterKey } from "$lib/vault/masterKeyStore";

const API = "";

export async function unlockWithPassword(email: string, masterPassword: string): Promise<CryptoKey> {
  const kdf = await fetchKdfParams(API, email);
  const mk = await deriveMasterKeyBytes(masterPassword, kdf.salt, kdf);
  const masterKey = await importAesKeyFromBytes(mk);
  setMasterKey(masterKey);
  return masterKey;
}

export function persistSession(email: string, token: string): void {
  saveSessionToken(token);
  saveUserEmail(email);
}

/** Login com Master Password — autentica servidor e desbloqueia cofre (ZK). */
export async function loginWithPassword(
  email: string,
  masterPassword: string,
): Promise<{ token: string; masterKey: CryptoKey }> {
  const normalized = normalizeEmail(email);
  const { masterKey, token } = await loginUser(API, normalized, masterPassword);
  persistSession(normalized, token);
  setMasterKey(masterKey);
  return { token, masterKey };
}

/** Passkey autentica o servidor; unlock fica para a página seguinte. */
export async function loginWithPasskeyOnly(email: string): Promise<string> {
  const normalized = normalizeEmail(email);
  const token = await loginWithPasskey(API, normalized);
  persistSession(normalized, token);
  return token;
}

export async function registerAndLogin(
  email: string,
  masterPassword: string,
): Promise<{ token: string; masterKey: CryptoKey }> {
  const normalized = normalizeEmail(email);
  await registerUser(API, normalized, masterPassword);
  const { masterKey, token } = await loginAfterRegister(API, normalized, masterPassword);
  persistSession(normalized, token);
  setMasterKey(masterKey);
  return { token, masterKey };
}

export async function navigateAfterAuth(redirectTo: string): Promise<void> {
  await goto(redirectTo);
}
