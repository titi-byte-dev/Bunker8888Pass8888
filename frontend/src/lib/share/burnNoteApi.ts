/**
 * Cliente da API das Notas Auto-Destrutivas (SHARE-005).
 *
 * - Criar exige sessão (Bearer token).
 * - Ler (queima após leitura) e queimar manualmente são PÚBLICOS: a chave de
 *   cifra vive no fragmento do URL.
 */

export interface CreatedNote {
  id: string;
  burn_token: string;
  expires_at: string;
}

/** Cria uma nota efémera a partir de um ciphertext (base64). Requer token. */
export async function createBurnNote(
  token: string,
  ciphertextB64: string,
  ttlSeconds: number,
): Promise<CreatedNote> {
  const res = await fetch("/api/share/notes", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ ciphertext: ciphertextB64, ttl_seconds: ttlSeconds }),
  });
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return (await res.json()) as CreatedNote;
}

/** Erro lançado quando a nota não existe, expirou ou já foi lida/queimada. */
export class BurnNoteGoneError extends Error {
  constructor() {
    super("Esta nota nao existe, expirou ou ja foi destruida.");
    this.name = "BurnNoteGoneError";
  }
}

/** Lê a nota (queima após leitura) e devolve o ciphertext (base64). Público. */
export async function consumeBurnNote(id: string): Promise<string> {
  const res = await fetch(`/api/share/notes/${encodeURIComponent(id)}`, { method: "POST" });
  if (res.status === 404) {
    throw new BurnNoteGoneError();
  }
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
  return ((await res.json()) as { ciphertext: string }).ciphertext;
}

/** Destroi a nota antes de ser lida, mediante o burn_token (capacidade). Público. */
export async function burnNoteManually(id: string, burnToken: string): Promise<void> {
  const res = await fetch(`/api/share/notes/${encodeURIComponent(id)}/burn`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ burn_token: burnToken }),
  });
  if (res.status === 404) {
    throw new BurnNoteGoneError();
  }
  if (!res.ok) {
    const err = (await res.json().catch(() => ({}))) as { error?: string };
    throw new Error(err.error ?? `HTTP ${res.status}`);
  }
}
