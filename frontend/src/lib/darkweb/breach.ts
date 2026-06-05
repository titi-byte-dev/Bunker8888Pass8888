/**
 * Verificação de fugas via k-anonymity (DW-001).
 *
 * Didático: enviamos só os primeiros 5 caracteres do hash SHA-1 da password
 * para a API públicа de breach data. A comparação do sufixo faz-se LOCALMENTE —
 * o serviço externo nunca sabe qual password verificámos.
 *
 * > ⚠️ **Segurança:** SHA-1 aqui é exigência do protocolo HIBP (Pwned Passwords),
 * não escolha nossa para cifragem. Nunca enviamos a password nem o hash completo.
 */

const HIBP_RANGE_URL = "https://api.pwnedpasswords.com/range/";

export type BreachCheckResult = {
  /** true se o hash completo apareceu na resposta da API */
  breached: boolean;
  /** Número de vezes reportada na base (0 se limpa) */
  exposureCount: number;
};

/** SHA-1 em hexadecimal MAIÚSCULAS — formato HIBP. */
export async function sha1HexUpper(text: string): Promise<string> {
  const digest = await crypto.subtle.digest("SHA-1", new TextEncoder().encode(text));
  return [...new Uint8Array(digest)]
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("")
    .toUpperCase();
}

/**
 * Verifica se uma password apareceu em fugas conhecidas (k-anonymity).
 * Requer rede; falha graciosamente com mensagem legível.
 */
export async function checkPasswordBreached(password: string): Promise<BreachCheckResult> {
  if (!password) {
    return { breached: false, exposureCount: 0 };
  }

  const hash = await sha1HexUpper(password);
  const prefix = hash.slice(0, 5);
  const suffix = hash.slice(5);

  let response: Response;
  try {
    response = await fetch(`${HIBP_RANGE_URL}${prefix}`, {
      headers: { "Add-Padding": "true" },
    });
  } catch {
    throw new Error("Verificação indisponível — confirma ligação à Internet");
  }

  if (!response.ok) {
    throw new Error("Serviço de breach data indisponível");
  }

  const body = await response.text();
  for (const line of body.split("\n")) {
    const [hashSuffix, countStr] = line.trim().split(":");
    if (hashSuffix === suffix) {
      return { breached: true, exposureCount: Number.parseInt(countStr, 10) || 1 };
    }
  }

  return { breached: false, exposureCount: 0 };
}

/** Verifica várias passwords com pausa entre pedidos (rate-limit cortês). */
export async function checkPasswordsBreached(
  passwords: string[],
  delayMs = 350,
): Promise<BreachCheckResult[]> {
  const results: BreachCheckResult[] = [];
  for (let i = 0; i < passwords.length; i++) {
    results.push(await checkPasswordBreached(passwords[i]!));
    if (i < passwords.length - 1 && delayMs > 0) {
      await new Promise((r) => setTimeout(r, delayMs));
    }
  }
  return results;
}
