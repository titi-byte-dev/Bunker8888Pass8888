/**
 * Abstração de plataforma nativa (UI-009 Capacitor).
 * Em browser devolve stubs seguros; com Capacitor injecta APIs reais.
 */

export type BiometricResult = "success" | "unavailable" | "failed";

/** Verdadeiro quando a app corre dentro de Capacitor (iOS/Android). */
export function isNativeApp(): boolean {
  if (typeof window === "undefined") return false;
  const cap = (window as Window & { Capacitor?: { isNativePlatform?: () => boolean } }).Capacitor;
  return cap?.isNativePlatform?.() === true;
}

/**
 * Desbloqueio biométrico — stub em web; em mobile liga ao plugin nativo (futuro).
 * ⚠️ Segurança: nunca substitui a Master Key — só re-autentica o utilizador.
 */
export async function unlockWithBiometric(): Promise<BiometricResult> {
  if (!isNativeApp()) return "unavailable";
  // Plugin @capacitor-community/biometric-auth será ligado em UI-009 fase 2.
  return "unavailable";
}
