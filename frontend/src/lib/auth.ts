/**
 * Autenticação Zero-Knowledge no cliente (VAULT-001 + playground dev).
 *
 * Didático: espelhamos o backend Go (`pkg/crypto/kdf.go`):
 *   1. Master Key = Argon2id(password, kdf_salt, params)
 *   2. Auth Hash  = Argon2id(masterKey, password, params com time=1)
 * Só o Auth Hash viaja para o servidor; a Master Key importa-se como AES-GCM.
 */
import { argon2id } from "hash-wasm";
import { randomBytes, type Bytes } from "./crypto";

/** Parâmetros KDF do cliente (guardados na BD no registo). */
export interface ClientKdfParams {
  time: number;
  memory: number; // KiB (como Argon2 MemoryKiB no Go)
  threads: number;
}

/**
 * Parâmetros de dev (~8 MiB) — derivação rápida o suficiente para testar.
 * ⚠️ Produção: time=3, memory=65536, threads=4 (64 MiB, ~0.5–1s).
 */
export const DEV_CLIENT_KDF: ClientKdfParams = {
  time: 1,
  memory: 8192,
  threads: 1,
};

function toBytes(data: string | Uint8Array): Uint8Array {
  return typeof data === "string" ? new TextEncoder().encode(data) : data;
}

export async function deriveMasterKeyBytes(
  password: string,
  salt: Bytes,
  params: ClientKdfParams,
): Promise<Uint8Array> {
  const hash = await argon2id({
    password,
    salt: toBytes(salt),
    parallelism: params.threads,
    iterations: params.time,
    memorySize: params.memory,
    hashLength: 32,
    outputType: "binary",
  });
  return hash as Uint8Array;
}

/** Segunda derivação — independente da Master Key (Auth Hash). */
export async function deriveAuthHashBytes(
  masterKey: Uint8Array,
  password: string,
  params: ClientKdfParams,
): Promise<Uint8Array> {
  const hash = await argon2id({
    password: masterKey,
    salt: toBytes(password),
    parallelism: params.threads,
    iterations: 1,
    memorySize: params.memory,
    hashLength: 32,
    outputType: "binary",
  });
  return hash as Uint8Array;
}

/** Importa 32 bytes como CryptoKey AES-GCM (não exportável). */
export async function importAesKeyFromBytes(raw: Uint8Array): Promise<CryptoKey> {
  const keyBytes = new Uint8Array(raw) as Bytes;
  return crypto.subtle.importKey("raw", keyBytes, "AES-GCM", false, ["encrypt", "decrypt"]);
}

export function bytesToBase64(b: Uint8Array): string {
  let s = "";
  for (let i = 0; i < b.length; i++) s += String.fromCharCode(b[i]!);
  return btoa(s);
}

export function base64ToBytes(b64: string): Bytes {
  const bin = atob(b64);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out as Bytes;
}

export interface KdfResponse {
  salt: string;
  time: number;
  memory: number;
  threads: number;
}

export async function fetchKdfParams(baseURL: string, email: string): Promise<ClientKdfParams & { salt: Bytes }> {
  const res = await fetch(`${baseURL}/api/auth/kdf?email=${encodeURIComponent(email)}`);
  if (!res.ok) throw new Error("Utilizador não encontrado");
  const j = (await res.json()) as KdfResponse;
  return {
    salt: base64ToBytes(j.salt),
    time: j.time,
    memory: j.memory,
    threads: j.threads,
  };
}

export async function registerUser(
  baseURL: string,
  email: string,
  masterPassword: string,
  kdf: ClientKdfParams = DEV_CLIENT_KDF,
): Promise<{ masterKey: CryptoKey; token: null }> {
  const salt = randomBytes(16);
  const mk = await deriveMasterKeyBytes(masterPassword, salt, kdf);
  const authHash = await deriveAuthHashBytes(mk, masterPassword, kdf);
  const res = await fetch(`${baseURL}/api/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email,
      auth_hash: bytesToBase64(authHash),
      kdf: {
        salt: bytesToBase64(salt),
        time: kdf.time,
        memory: kdf.memory,
        threads: kdf.threads,
      },
    }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `Registo falhou (${res.status})`);
  }
  const masterKey = await importAesKeyFromBytes(mk);
  return { masterKey, token: null };
}

export async function loginUser(
  baseURL: string,
  email: string,
  masterPassword: string,
): Promise<{ masterKey: CryptoKey; token: string }> {
  const kdf = await fetchKdfParams(baseURL, email);
  const mk = await deriveMasterKeyBytes(masterPassword, kdf.salt, kdf);
  const authHash = await deriveAuthHashBytes(mk, masterPassword, kdf);
  const res = await fetch(`${baseURL}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      email,
      auth_hash: bytesToBase64(authHash),
    }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `Login falhou (${res.status})`);
  }
  const { token } = (await res.json()) as { token: string };
  const masterKey = await importAesKeyFromBytes(mk);
  return { masterKey, token };
}

/** Deriva chaves após registo (login imediato). */
export async function loginAfterRegister(
  baseURL: string,
  email: string,
  masterPassword: string,
): Promise<{ masterKey: CryptoKey; token: string }> {
  return loginUser(baseURL, email, masterPassword);
}
