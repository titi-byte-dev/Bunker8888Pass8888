/**
 * Score de higiene de passwords (VAULT-008).
 *
 * Didático: toda a análise corre NO CLIENTE, depois de decifrar os logins.
 * O servidor nunca vê passwords em claro — no futuro só pode receber agregados
 * (ex: score global 0–100) para dashboards, nunca os segredos.
 */
import { estimateEntropyBits } from "./password";

/** Problemas detectados num login. */
export type HygieneIssue = "weak" | "reused" | "common" | "short";

/** Entrada mínima: um login já decifrado. */
export interface LoginCredential {
  itemId: string;
  title: string;
  password: string;
}

export interface ItemHygieneResult {
  itemId: string;
  title: string;
  /** 0 (péssima) a 100 (forte). */
  score: number;
  issues: HygieneIssue[];
  /** IDs de outros itens que partilham a mesma password. */
  reusedWith: string[];
}

/** Resumo agregado — é ISTO que poderia ir ao servidor (sem passwords). */
export interface HygieneSummary {
  overallScore: number;
  totalLogins: number;
  weakCount: number;
  reusedCount: number;
  items: ItemHygieneResult[];
}

const WEAK_THRESHOLD = 50;
const SHORT_THRESHOLD = 10;

/** Lista curta de passwords triviais (amostra; em produção expandir ou usar API local). */
const COMMON_PASSWORDS = new Set(
  [
    "password",
    "123456",
    "12345678",
    "123456789",
    "qwerty",
    "abc123",
    "password1",
    "admin",
    "letmein",
    "welcome",
    "monkey",
    "dragon",
    "master",
    "login",
    "passw0rd",
    "iloveyou",
    "sunshine",
    "princess",
    "football",
    "shadow",
    "senha",
    "palavrapasse",
    "alterar",
    "admin123",
  ].map((s) => s.toLowerCase()),
);

/**
 * Avalia a força de uma password (0–100).
 *
 * Combina comprimento, diversidade de charset, entropia estimada e penalizações
 * por padrões fracos (comuns, repetições, sequências óbvias).
 */
export function passwordStrengthScore(password: string): number {
  if (!password) return 0;

  let score = 0;

  // Comprimento: até 30 pontos (15 chars = teto).
  score += Math.min(30, password.length * 2);

  // Diversidade de charset: até 30 pontos.
  if (/[a-z]/.test(password)) score += 8;
  if (/[A-Z]/.test(password)) score += 8;
  if (/[0-9]/.test(password)) score += 7;
  if (/[^a-zA-Z0-9]/.test(password)) score += 7;

  // Bónus de entropia (charset detectado automaticamente).
  const entropy = estimateEntropyBits({
    length: password.length,
    uppercase: /[A-Z]/.test(password),
    lowercase: /[a-z]/.test(password),
    digits: /[0-9]/.test(password),
    symbols: /[^a-zA-Z0-9]/.test(password),
    excludeAmbiguous: false,
  });
  score += Math.min(25, Math.floor(entropy / 5));

  // Penalizações.
  const lower = password.toLowerCase();
  if (COMMON_PASSWORDS.has(lower)) score -= 40;
  if (/^(.)\1+$/.test(password)) score -= 30; // aaaaaa
  if (/(.)\1{2,}/.test(password)) score -= 15; // aaabbb
  if (/01234|12345|23456|abcde|qwert|asdfg/i.test(password)) score -= 20;
  if (/password|senha|admin|login/i.test(password)) score -= 15;

  return clamp(Math.round(score), 0, 100);
}

/** Analisa todos os logins do cofre (já decifrados). */
export async function analyzeVaultHygiene(
  logins: LoginCredential[],
): Promise<HygieneSummary> {
  const reuseMap = await buildReuseMap(logins);

  const items: ItemHygieneResult[] = logins.map((login) => {
    const issues: HygieneIssue[] = [];
    const score = passwordStrengthScore(login.password);

    if (login.password.length < SHORT_THRESHOLD) issues.push("short");
    if (score < WEAK_THRESHOLD) issues.push("weak");
    if (COMMON_PASSWORDS.has(login.password.toLowerCase())) issues.push("common");

    const reusedWith = (reuseMap.get(login.itemId) ?? []).filter(
      (id) => id !== login.itemId,
    );
    if (reusedWith.length > 0) issues.push("reused");

    return { itemId: login.itemId, title: login.title, score, issues, reusedWith };
  });

  const weakCount = items.filter((i) => i.issues.includes("weak")).length;
  const reusedCount = items.filter((i) => i.issues.includes("reused")).length;

  const overallScore =
    items.length === 0
      ? 100
      : Math.round(items.reduce((sum, i) => sum + i.score, 0) / items.length);

  return {
    overallScore: penalizeOverall(overallScore, weakCount, reusedCount, items.length),
    totalLogins: items.length,
    weakCount,
    reusedCount,
    items,
  };
}

/** Agrupa logins pela mesma password via SHA-256 (não guardamos plaintext no Map). */
async function buildReuseMap(
  logins: LoginCredential[],
): Promise<Map<string, string[]>> {
  // fingerprint -> itemIds
  const groups = new Map<string, string[]>();

  for (const login of logins) {
    const fp = await passwordFingerprint(login.password);
    const list = groups.get(fp) ?? [];
    list.push(login.itemId);
    groups.set(fp, list);
  }

  // itemId -> outros itemIds com a mesma fingerprint
  const byItem = new Map<string, string[]>();
  for (const ids of groups.values()) {
    if (ids.length < 2) continue;
    for (const id of ids) {
      byItem.set(id, ids);
    }
  }
  return byItem;
}

/** SHA-256 da password — fingerprint para detetar reutilização sem comparar strings longas. */
async function passwordFingerprint(password: string): Promise<string> {
  const hash = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(password));
  return [...new Uint8Array(hash)].map((b) => b.toString(16).padStart(2, "0")).join("");
}

/** Penaliza o score global se houver muitas fracas/reutilizadas. */
function penalizeOverall(
  avg: number,
  weak: number,
  reused: number,
  total: number,
): number {
  if (total === 0) return 100;
  const weakPct = weak / total;
  const reusedPct = reused / total;
  const penalty = Math.round(weakPct * 30 + reusedPct * 40);
  return clamp(avg - penalty, 0, 100);
}

function clamp(n: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, n));
}

/** Extrai só o agregado para enviar ao servidor (Zero-Knowledge). */
export function toServerSummary(summary: HygieneSummary): Omit<HygieneSummary, "items"> {
  return {
    overallScore: summary.overallScore,
    totalLogins: summary.totalLogins,
    weakCount: summary.weakCount,
    reusedCount: summary.reusedCount,
  };
}
