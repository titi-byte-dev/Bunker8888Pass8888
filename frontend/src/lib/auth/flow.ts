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
  const { masterKey, token } = await loginUser(API, email, masterPassword);
  persistSession(email, token);
  setMasterKey(masterKey);
  return { token, masterKey };
}

/** Passkey autentica o servidor; unlock fica para a página seguinte. */
export async function loginWithPasskeyOnly(email: string): Promise<string> {
  const token = await loginWithPasskey(API, email);
  persistSession(email, token);
  return token;
}

export async function registerAndLogin(
  email: string,
  masterPassword: string,
): Promise<{ token: string; masterKey: CryptoKey }> {
  await registerUser(API, email, masterPassword);
  const { masterKey, token } = await loginAfterRegister(API, email, masterPassword);
  persistSession(email, token);
  setMasterKey(masterKey);
  return { token, masterKey };
}

export async function navigateAfterAuth(redirectTo: string): Promise<void> {
  await goto(redirectTo);
}
