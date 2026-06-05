/**
 * Protocolo postMessage do browser sandbox (VAULT-013).
 *
 * > ⚠️ **Segurança:** só aceitamos mensagens da mesma origem. A password
 * viaja para o iframe de destino mas NUNCA é renderizada no painel pai.
 */

export const SANDBOX_FILL_MESSAGE = "aegis:sandbox:fill" as const;
export const SANDBOX_READY_MESSAGE = "aegis:sandbox:ready" as const;

export type SandboxFillPayload = {
  type: typeof SANDBOX_FILL_MESSAGE;
  username: string;
  password: string;
};

export type SandboxReadyPayload = {
  type: typeof SANDBOX_READY_MESSAGE;
};

export function isSandboxFillPayload(data: unknown): data is SandboxFillPayload {
  if (!data || typeof data !== "object") return false;
  const d = data as Record<string, unknown>;
  return (
    d.type === SANDBOX_FILL_MESSAGE &&
    typeof d.username === "string" &&
    typeof d.password === "string"
  );
}

export function isSandboxReadyPayload(data: unknown): data is SandboxReadyPayload {
  return !!data && typeof data === "object" && (data as { type?: string }).type === SANDBOX_READY_MESSAGE;
}

/** Envia credenciais para o iframe sandbox (mesma origem). */
export function postSandboxFill(
  iframe: HTMLIFrameElement,
  creds: { username: string; password: string },
  targetOrigin: string,
): void {
  const win = iframe.contentWindow;
  if (!win) throw new Error("Iframe ainda não carregou");
  const msg: SandboxFillPayload = {
    type: SANDBOX_FILL_MESSAGE,
    username: creds.username,
    password: creds.password,
  };
  win.postMessage(msg, targetOrigin);
}

/** URL do alvo demo (same-origin, layout bare) onde a injecção funciona. */
export const SANDBOX_DEMO_TARGET_PATH = "/sandbox/login-target";
