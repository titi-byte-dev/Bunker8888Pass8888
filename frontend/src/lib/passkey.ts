/**
 * Passkeys / WebAuthn (VAULT-014) — autenticação ao servidor.
 *
 * Didático: a passkey prova posse do dispositivo (chave privada no Secure Enclave).
 * Obtemos um token de sessão, mas a Master Key continua a vir da Master Password
 * (Zero-Knowledge — o servidor nunca a vê).
 */
import { deriveMasterKeyBytes, fetchKdfParams, importAesKeyFromBytes } from "./auth";

function base64URLToBuffer(b64: string): ArrayBuffer {
  const pad = "=".repeat((4 - (b64.length % 4)) % 4);
  const bin = atob((b64 + pad).replace(/-/g, "+").replace(/_/g, "/"));
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out.buffer;
}

function bufferToBase64URL(buf: ArrayBuffer): string {
  const bytes = new Uint8Array(buf);
  let s = "";
  for (let i = 0; i < bytes.length; i++) s += String.fromCharCode(bytes[i]!);
  return btoa(s).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function decodeCreationOptions(options: PublicKeyCredentialCreationOptionsJSON): PublicKeyCredentialCreationOptions {
  const decoded = {
    ...options,
    challenge: base64URLToBuffer(options.challenge),
    user: {
      ...options.user,
      id: base64URLToBuffer(options.user.id),
    },
    excludeCredentials: options.excludeCredentials?.map((c) => ({
      type: "public-key" as const,
      id: base64URLToBuffer(c.id),
      transports: c.transports as AuthenticatorTransport[] | undefined,
    })),
  };
  return decoded as PublicKeyCredentialCreationOptions;
}

function decodeRequestOptions(options: PublicKeyCredentialRequestOptionsJSON): PublicKeyCredentialRequestOptions {
  const decoded = {
    ...options,
    challenge: base64URLToBuffer(options.challenge),
    allowCredentials: options.allowCredentials?.map((c) => ({
      type: "public-key" as const,
      id: base64URLToBuffer(c.id),
      transports: c.transports as AuthenticatorTransport[] | undefined,
    })),
  };
  return decoded as PublicKeyCredentialRequestOptions;
}

function encodeCredential(cred: PublicKeyCredential): Record<string, unknown> {
  if (cred.response instanceof AuthenticatorAttestationResponse) {
    return {
      id: cred.id,
      rawId: bufferToBase64URL(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: bufferToBase64URL(cred.response.clientDataJSON),
        attestationObject: bufferToBase64URL(cred.response.attestationObject),
        transports: cred.response.getTransports?.() ?? [],
      },
    };
  }
  const ar = cred.response as AuthenticatorAssertionResponse;
  return {
    id: cred.id,
    rawId: bufferToBase64URL(cred.rawId),
    type: cred.type,
    response: {
      clientDataJSON: bufferToBase64URL(ar.clientDataJSON),
      authenticatorData: bufferToBase64URL(ar.authenticatorData),
      signature: bufferToBase64URL(ar.signature),
      userHandle: ar.userHandle ? bufferToBase64URL(ar.userHandle) : null,
    },
  };
}

export interface PasskeyMeta {
  id: string;
  name: string;
  created_at: string;
}

export function passkeysSupported(): boolean {
  return typeof window !== "undefined" && !!window.PublicKeyCredential;
}

/** Regista uma passkey (requer sessão Bearer activa). */
export async function registerPasskey(baseURL: string, token: string, name: string): Promise<void> {
  if (!passkeysSupported()) throw new Error("WebAuthn não suportado neste browser");

  const begin = await fetch(`${baseURL}/api/auth/passkey/register/begin`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: "{}",
  });
  if (!begin.ok) throw new Error(`Registo passkey falhou (${begin.status})`);
  const { options, session_id } = (await begin.json()) as {
    options: PublicKeyCredentialCreationOptionsJSON;
    session_id: string;
  };

  const cred = (await navigator.credentials.create({
    publicKey: decodeCreationOptions(options),
  })) as PublicKeyCredential | null;
  if (!cred) throw new Error("Registo cancelado");

  const finish = await fetch(`${baseURL}/api/auth/passkey/register/finish`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: JSON.stringify({
      session_id,
      name,
      credential: encodeCredential(cred),
    }),
  });
  if (!finish.ok) {
    const err = (await finish.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `Registo falhou (${finish.status})`);
  }
}

/** Login com passkey — devolve token; Master Key requer password separada. */
export async function loginWithPasskey(baseURL: string, email: string): Promise<string> {
  if (!passkeysSupported()) throw new Error("WebAuthn não suportado neste browser");

  const begin = await fetch(`${baseURL}/api/auth/passkey/login/begin`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
  if (!begin.ok) throw new Error("Passkey não disponível para este email");
  const { options, session_id } = (await begin.json()) as {
    options: PublicKeyCredentialRequestOptionsJSON;
    session_id: string;
  };

  const cred = (await navigator.credentials.get({
    publicKey: decodeRequestOptions(options),
  })) as PublicKeyCredential | null;
  if (!cred) throw new Error("Login cancelado");

  const finish = await fetch(`${baseURL}/api/auth/passkey/login/finish`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      session_id,
      credential: encodeCredential(cred),
    }),
  });
  if (!finish.ok) throw new Error("Autenticação passkey falhou");
  const { token } = (await finish.json()) as { token: string };
  return token;
}

/** Login completo: passkey (sessão) + Master Password (desbloqueio ZK). */
export async function unlockWithPasskeyAndPassword(
  baseURL: string,
  email: string,
  masterPassword: string,
): Promise<{ masterKey: CryptoKey; token: string }> {
  const token = await loginWithPasskey(baseURL, email);
  const kdf = await fetchKdfParams(baseURL, email);
  const mk = await deriveMasterKeyBytes(masterPassword, kdf.salt, kdf);
  const masterKey = await importAesKeyFromBytes(mk);
  return { masterKey, token };
}

export async function listPasskeys(baseURL: string, token: string): Promise<PasskeyMeta[]> {
  const res = await fetch(`${baseURL}/api/auth/passkey`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error("Falha ao listar passkeys");
  const j = (await res.json()) as { passkeys: PasskeyMeta[] };
  return j.passkeys;
}

export const _test = { base64URLToBuffer, bufferToBase64URL };
