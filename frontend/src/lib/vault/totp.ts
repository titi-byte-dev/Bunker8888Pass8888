/**
 * TOTP — Time-based One-Time Password (RFC 6238).
 *
 * Didático: o TOTP é uma extensão do HOTP (RFC 4226) onde o "contador" é o tempo
 * dividido em janelas (tipicamente 30s). Servidor e cliente partilham um segredo
 * (Base32) e calculam o mesmo código sem comunicação — daí funcionar offline.
 *
 * ⚠️ Segurança: usamos WebCrypto para HMAC (nunca HMAC caseiro). O segredo TOTP
 * vive cifrado no blob do login (Zero-Knowledge); só é decifrado no cliente.
 */

export type TotpAlgorithm = "SHA-1" | "SHA-256" | "SHA-512";

export interface TotpOptions {
  /** Janela temporal em segundos (Google Authenticator usa 30). */
  period?: number;
  /** Nº de dígitos do código (quase sempre 6). */
  digits?: number;
  algorithm?: TotpAlgorithm;
}

const DEFAULT_OPTS: Required<TotpOptions> = {
  period: 30,
  digits: 6,
  algorithm: "SHA-1",
};

/** Resultado de parse de um URI otpauth:// (QR code). */
export interface OtpauthParsed {
  label: string;
  secretBase32: string;
  issuer?: string;
  period: number;
  digits: number;
  algorithm: TotpAlgorithm;
}

/**
 * Gera o código TOTP actual a partir de um segredo Base32.
 *
 * @param secretBase32 Segredo como no QR (ex: JBSWY3DPEHPK3PXP).
 * @param unixTime Segundos UNIX; por omissão `Date.now()/1000`.
 */
export async function generateTotp(
  secretBase32: string,
  unixTime: number = Math.floor(Date.now() / 1000),
  opts: TotpOptions = {},
): Promise<string> {
  const o = { ...DEFAULT_OPTS, ...opts };
  const key = decodeBase32(secretBase32);
  const counter = Math.floor(unixTime / o.period);
  return hotp(key, counter, o.digits, o.algorithm);
}

/** Segundos até o próximo código (útil para UI com contagem decrescente). */
export function totpSecondsRemaining(
  unixTime: number = Math.floor(Date.now() / 1000),
  period: number = DEFAULT_OPTS.period,
): number {
  return period - (unixTime % period);
}

/**
 * Analisa URIs otpauth://totp/... geradas por Google Authenticator e similares.
 *
 * Exemplo:
 * otpauth://totp/AegisPass:user?secret=JBSWY3DPEHPK3PXP&issuer=AegisPass
 */
export function parseOtpauthUri(uri: string): OtpauthParsed {
  const url = new URL(uri);
  if (url.protocol !== "otpauth:") {
    throw new Error("URI inválida: esperado otpauth://");
  }
  const type = url.hostname;
  if (type !== "totp" && type !== "hotp") {
    throw new Error(`tipo otpauth não suportado: ${type}`);
  }

  const secret = url.searchParams.get("secret");
  if (!secret) throw new Error("secret em falta no URI");

  const label = decodeURIComponent(url.pathname.replace(/^\//, ""));
  const issuer =
    url.searchParams.get("issuer") ??
    (label.includes(":") ? label.split(":")[0] : undefined);

  const algoParam = (url.searchParams.get("algorithm") ?? "SHA1").toUpperCase();
  const algorithm = parseAlgorithm(algoParam);

  return {
    label,
    secretBase32: secret.replace(/\s+/g, "").toUpperCase(),
    issuer: issuer ?? undefined,
    period: parseInt(url.searchParams.get("period") ?? "30", 10) || 30,
    digits: parseInt(url.searchParams.get("digits") ?? "6", 10) || 6,
    algorithm,
  };
}

/** HOTP — HMAC-based OTP (RFC 4226), base do TOTP. */
export async function hotp(
  secret: Uint8Array,
  counter: number,
  digits: number = 6,
  algorithm: TotpAlgorithm = "SHA-1",
): Promise<string> {
  const msg = counterToBytes(counter);
  const hash = await hmac(secret, msg, algorithm);

  // Truncamento dinâmico (RFC 4226 §5.4).
  const h = new Uint8Array(hash);
  const offset = h[h.length - 1]! & 0x0f;
  const bin =
    ((h[offset]! & 0x7f) << 24) |
    ((h[offset + 1]! & 0xff) << 16) |
    ((h[offset + 2]! & 0xff) << 8) |
    (h[offset + 3]! & 0xff);

  const mod = 10 ** digits;
  return String(bin % mod).padStart(digits, "0");
}

/** Vetor de teste RFC 6238: segredo ASCII de 20 bytes (não Base32). */
export async function generateTotpRfcTest(
  secretAscii: string,
  unixTime: number,
  opts: TotpOptions = {},
): Promise<string> {
  const o = { ...DEFAULT_OPTS, ...opts };
  const key = new TextEncoder().encode(secretAscii);
  const counter = Math.floor(unixTime / o.period);
  return hotp(key, counter, o.digits, o.algorithm);
}

// --- internos ---

async function hmac(
  key: Uint8Array,
  message: Uint8Array,
  algorithm: TotpAlgorithm,
): Promise<ArrayBuffer> {
  const hashName = algorithm === "SHA-1" ? "SHA-1" : algorithm;
  const cryptoKey = await crypto.subtle.importKey(
    "raw",
    key.buffer.slice(key.byteOffset, key.byteOffset + key.byteLength) as ArrayBuffer,
    { name: "HMAC", hash: hashName },
    false,
    ["sign"],
  );
  return crypto.subtle.sign("HMAC", cryptoKey, message.buffer.slice(
    message.byteOffset,
    message.byteOffset + message.byteLength,
  ) as ArrayBuffer);
}

/** Contador de 8 bytes em big-endian (RFC 4226). */
function counterToBytes(counter: number): Uint8Array {
  const buf = new Uint8Array(8);
  let c = counter;
  for (let i = 7; i >= 0; i--) {
    buf[i] = c & 0xff;
    c = Math.floor(c / 256);
  }
  return buf;
}

/** Descodifica Base32 (RFC 4648) — formato dos QR de autenticadores. */
export function decodeBase32(input: string): Uint8Array {
  const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";
  const cleaned = input.replace(/\s+/g, "").toUpperCase();
  const out: number[] = [];
  let buffer = 0;
  let bitsLeft = 0;

  for (const ch of cleaned) {
    if (ch === "=") break;
    const idx = alphabet.indexOf(ch);
    if (idx < 0) throw new Error(`carácter Base32 inválido: ${ch}`);
    buffer = (buffer << 5) | idx;
    bitsLeft += 5;
    if (bitsLeft >= 8) {
      out.push((buffer >> (bitsLeft - 8)) & 0xff);
      bitsLeft -= 8;
    }
  }
  return new Uint8Array(out) as Uint8Array<ArrayBuffer>;
}

function parseAlgorithm(s: string): TotpAlgorithm {
  switch (s) {
    case "SHA1":
    case "SHA-1":
      return "SHA-1";
    case "SHA256":
    case "SHA-256":
      return "SHA-256";
    case "SHA512":
    case "SHA-512":
      return "SHA-512";
    default:
      throw new Error(`algoritmo TOTP não suportado: ${s}`);
  }
}
