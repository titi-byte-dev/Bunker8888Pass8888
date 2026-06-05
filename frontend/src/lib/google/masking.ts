/**
 * Mascaramento dinâmico NIF/IBAN (GOOGLE-003 dev stub).
 * Simula o proxy Sheets — tokens opacos em vez de PII em claro.
 */
const NIF_RE = /\b\d{9}\b/g;
const IBAN_RE = /\bPT\d{2}\s?\d{4}\s?\d{4}\s?\d{11,13}\b/gi;

let tokenSeq = 0;

function nextToken(kind: string): string {
  tokenSeq += 1;
  return `TOKEN_${kind}_${tokenSeq.toString(36).toUpperCase()}`;
}

export interface MaskResult {
  masked: string;
  tokens: Record<string, string>;
}

/** Substitui NIF e IBAN por tokens reversíveis (só no cliente, dev). */
export function maskSensitiveText(text: string): MaskResult {
  const tokens: Record<string, string> = {};
  let masked = text;

  masked = masked.replace(NIF_RE, (nif) => {
    const t = nextToken("NIF");
    tokens[t] = nif;
    return t;
  });
  masked = masked.replace(IBAN_RE, (iban) => {
    const t = nextToken("IBAN");
    tokens[t] = iban.replace(/\s/g, "");
    return t;
  });

  return { masked, tokens };
}

/** Reinjeta valores reais a partir do mapa de tokens (vista AegisPass). */
export function unmaskText(masked: string, tokens: Record<string, string>): string {
  let out = masked;
  for (const [token, value] of Object.entries(tokens)) {
    out = out.split(token).join(value);
  }
  return out;
}
