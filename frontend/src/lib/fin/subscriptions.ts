/**
 * Subscrições SaaS (FIN-001) — modelo + cifragem.
 *
 * Cada subscrição é um objeto cifrado com a Master Key (AES-GCM), tal como um
 * item de cofre. O servidor só vê o blob opaco; os custos e alertas (FIN-002)
 * são calculados no cliente depois de decifrar.
 */
import { decrypt, encrypt, fromBytes, toBytes, type Bytes } from "$lib/crypto";

export type BillingCycle = "monthly" | "yearly";

/** Conteúdo cifrado de uma subscrição (tudo menos o id do servidor). */
export interface SubscriptionPayload {
  name: string;
  cost: number; // por ciclo, na moeda indicada
  currency: string; // ex.: "EUR"
  cycle: BillingCycle;
  category?: string;
  /** Liga ao login correspondente no cofre ("cruza com vault"). */
  vaultItemId?: string;
  vaultItemTitle?: string;
  /** Última utilização conhecida (ISO) — base dos alertas de licença sem uso. */
  lastUsedAt?: string;
  active: boolean;
}

/** Subscrição completa (id do servidor + conteúdo decifrado). */
export interface Subscription extends SubscriptionPayload {
  id: string;
}

/** Cifra o payload de uma subscrição com a Master Key. */
export async function encryptSubscription(
  masterKey: CryptoKey,
  payload: SubscriptionPayload,
): Promise<Bytes> {
  return encrypt(masterKey, toBytes(JSON.stringify(payload)));
}

/** Decifra o blob de uma subscrição. */
export async function decryptSubscription(
  masterKey: CryptoKey,
  blob: Bytes,
): Promise<SubscriptionPayload> {
  return JSON.parse(fromBytes(await decrypt(masterKey, blob))) as SubscriptionPayload;
}
