/**
 * Sentinel Mode — step-up e histórico (DW-004).
 */
import { loadSessionToken } from "$lib/session";

const API = "";

export type LoginEvent = {
  id: string;
  user_id: string;
  email: string;
  client_ip: string;
  geo_lat?: number | null;
  geo_lon?: number | null;
  success: boolean;
  suspicious: boolean;
  reason: string;
  step_up_required: boolean;
  created_at: string;
};

export type SentinelStepUp = {
  challengeId: string;
  reason: string;
  detail: string;
};

export class SentinelStepUpRequired extends Error {
  readonly stepUp: SentinelStepUp;

  constructor(stepUp: SentinelStepUp) {
    super("Verificação Sentinel necessária");
    this.name = "SentinelStepUpRequired";
    this.stepUp = stepUp;
  }
}

export function isSentinelStepUp(err: unknown): err is SentinelStepUpRequired {
  return err instanceof SentinelStepUpRequired;
}

/** Parse resposta 403 do login com code sentinel_step_up. */
export function parseSentinelResponse(status: number, body: unknown): SentinelStepUpRequired | null {
  if (status !== 403 || !body || typeof body !== "object") return null;
  const b = body as { code?: string; challenge_id?: string; reason?: string; detail?: string };
  if (b.code !== "sentinel_step_up" || !b.challenge_id) return null;
  return new SentinelStepUpRequired({
    challengeId: b.challenge_id,
    reason: b.reason ?? "impossible_travel",
    detail: b.detail ?? "",
  });
}

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

function decodeRequestOptions(options: PublicKeyCredentialRequestOptionsJSON): PublicKeyCredentialRequestOptions {
  return {
    ...options,
    challenge: base64URLToBuffer(options.challenge),
    allowCredentials: options.allowCredentials?.map((c) => ({
      type: "public-key" as const,
      id: base64URLToBuffer(c.id),
      transports: c.transports as AuthenticatorTransport[] | undefined,
    })),
  } as PublicKeyCredentialRequestOptions;
}

function encodeAssertion(cred: PublicKeyCredential): Record<string, unknown> {
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

/** Completa step-up Sentinel com passkey — devolve token de sessão. */
export async function completeSentinelStepUp(
  challengeId: string,
  geoHeaders: Record<string, string> = {},
): Promise<string> {
  const begin = await fetch(`${API}/api/auth/sentinel/step-up/begin`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...geoHeaders },
    body: JSON.stringify({ challenge_id: challengeId }),
  });
  if (!begin.ok) throw new Error("Falha ao iniciar verificação Sentinel");
  const beginData = (await begin.json()) as {
    options: PublicKeyCredentialRequestOptionsJSON;
    session_id: string;
    challenge_id: string;
  };

  const cred = (await navigator.credentials.get({
    publicKey: decodeRequestOptions(beginData.options),
  })) as PublicKeyCredential | null;
  if (!cred) throw new Error("Verificação cancelada");

  const finish = await fetch(`${API}/api/auth/sentinel/step-up/finish`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...geoHeaders },
    body: JSON.stringify({
      challenge_id: beginData.challenge_id,
      session_id: beginData.session_id,
      credential: encodeAssertion(cred),
    }),
  });
  if (!finish.ok) throw new Error("Verificação passkey falhou");
  const { token } = (await finish.json()) as { token: string };
  return token;
}

export async function listSentinelEvents(): Promise<{
  events: LoginEvent[];
  suspiciousLast24h: number;
}> {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessão inválida");
  const res = await fetch(`${API}/api/security/sentinel/events`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error("Falha ao carregar eventos Sentinel");
  const data = (await res.json()) as { events: LoginEvent[]; suspicious_last_24h: number };
  return { events: data.events ?? [], suspiciousLast24h: data.suspicious_last_24h ?? 0 };
}

export function reasonLabel(reason: string): string {
  switch (reason) {
    case "impossible_travel":
      return "Viagem impossível";
    default:
      return reason || "Suspeito";
  }
}
