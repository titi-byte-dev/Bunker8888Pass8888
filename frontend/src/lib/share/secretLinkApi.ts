/**
 * Cliente da API dos Secret Links (SHARE-003).
 *
 * - Criar exige sessão (Bearer token).
 * - Consumir é PÚBLICO (sem token): a chave de cifra vive no fragmento do URL.
 */

export interface CreatedLink {
  id: string;
  expires_at: string;
}

/** Cria um link efémero a partir de um ciphertext (base64). Requer token. */
export async function createSecretLink(
  token: string,
  ciphertextB64: string,
  ttlSeconds: number,
  maxViews: number,
): Promise<CreatedLink> {
  const res = await fetch("/api/share/links", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      ciphertext: ciphertextB64,
      ttl_seconds: ttlSeconds,
      max_views: maxViews,
    }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return (await res.json()) as CreatedLink;
}

/** Erro lançado quando o link não existe, expirou ou já foi usado. */
export class SecretLinkGoneError extends Error {
  constructor() {
    super("Este link nao existe, expirou ou ja foi utilizado.");
    this.name = "SecretLinkGoneError";
  }
}

/** Consome o link (uso único) e devolve o ciphertext (base64). Público. */
export async function consumeSecretLink(id: string): Promise<string> {
  const res = await fetch(`/api/share/links/${encodeURIComponent(id)}`, { method: "POST" });
  if (res.status === 404) {
    throw new SecretLinkGoneError();
  }
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return ((await res.json()) as { ciphertext: string }).ciphertext;
}
