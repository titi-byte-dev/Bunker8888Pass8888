/**
 * Gerador de palavras-passe (VAULT-007).
 *
 * Didático: usamos SEMPRE `crypto.getRandomValues` (CSPRNG do browser) — nunca
 * `Math.random()`, que é previsível e inseguro para passwords.
 *
 * ⚠️ Segurança: a password é gerada e exibida só no cliente; nunca passa pelo
 * servidor até o utilizador a guardar num item cifrado do cofre.
 */

export interface PasswordOptions {
  /** Comprimento final. Por omissão 20 (≥128 bits de entropia com charset misto). */
  length?: number;
  uppercase?: boolean;
  lowercase?: boolean;
  digits?: boolean;
  symbols?: boolean;
  /** Exclui caracteres ambíguos (0/O, 1/l/I) — útil ao ditar ou imprimir. */
  excludeAmbiguous?: boolean;
}

const DEFAULTS: Required<PasswordOptions> = {
  length: 20,
  uppercase: true,
  lowercase: true,
  digits: true,
  symbols: true,
  excludeAmbiguous: true,
};

const UPPER = "ABCDEFGHJKLMNPQRSTUVWXYZ";
const LOWER = "abcdefghijkmnopqrstuvwxyz";
const DIGITS = "23456789";
const SYMBOLS = "!@#$%^&*()-_=+[]{}";

/** Gera uma password aleatória conforme as opções. */
export function generatePassword(opts: PasswordOptions = {}): string {
  const o = { ...DEFAULTS, ...opts };

  let upper = o.uppercase ? UPPER : "";
  let lower = o.lowercase ? LOWER : "";
  let digits = o.digits ? DIGITS : "";
  let symbols = o.symbols ? SYMBOLS : "";

  if (o.excludeAmbiguous) {
    // Os charsets acima já excluem ambíguos por desenho; mantemos a flag para
    // futuras extensões (ex: charset custom).
  }

  const pools = [
    { enabled: o.uppercase, chars: upper, name: "uppercase" },
    { enabled: o.lowercase, chars: lower, name: "lowercase" },
    { enabled: o.digits, chars: digits, name: "digits" },
    { enabled: o.symbols, chars: symbols, name: "symbols" },
  ].filter((p) => p.enabled && p.chars.length > 0);

  if (pools.length === 0) {
    throw new Error("pelo menos um conjunto de caracteres tem de estar activo");
  }
  if (o.length < pools.length) {
    throw new Error(
      `length (${o.length}) tem de ser >= nº de conjuntos activos (${pools.length})`,
    );
  }

  let all = "";
  for (const p of pools) all += p.chars;

  // Garantimos pelo menos um char de cada pool activo (requisito de complexidade).
  const chars: string[] = [];
  for (const p of pools) {
    chars.push(p.chars[randomIndex(p.chars.length)]!);
  }

  // Preenchemos o resto com chars aleatórios do charset completo.
  while (chars.length < o.length) {
    chars.push(all[randomIndex(all.length)]!);
  }

  // Fisher-Yates shuffle com CSPRNG — evita padrão previsível nas posições
  // obrigatórias (que ficaram no início antes do shuffle).
  for (let i = chars.length - 1; i > 0; i--) {
    const j = randomIndex(i + 1);
    [chars[i], chars[j]] = [chars[j]!, chars[i]!];
  }

  return chars.join("");
}

/** Índice aleatório em [0, max) usando crypto.getRandomValues. */
function randomIndex(max: number): number {
  // Rejection sampling simples para evitar modulo bias quando max não divide 256.
  const limit = Math.floor(256 / max) * max;
  const buf = new Uint8Array(1);
  let x: number;
  do {
    crypto.getRandomValues(buf);
    x = buf[0]!;
  } while (x >= limit);
  return x % max;
}

/** Estima bits de entropia (log2(charset^length)) — valor aproximado para UI. */
export function estimateEntropyBits(opts: PasswordOptions = {}): number {
  const o = { ...DEFAULTS, ...opts };
  let size = 0;
  if (o.uppercase) size += o.excludeAmbiguous ? UPPER.length : 26;
  if (o.lowercase) size += o.excludeAmbiguous ? LOWER.length : 26;
  if (o.digits) size += o.excludeAmbiguous ? DIGITS.length : 10;
  if (o.symbols) size += SYMBOLS.length;
  if (size === 0) return 0;
  return Math.floor(o.length * Math.log2(size));
}
