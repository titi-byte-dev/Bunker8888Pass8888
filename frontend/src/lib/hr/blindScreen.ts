/**
 * Triagem às cegas no cliente (AGENT-007 / CT-RGPD-04).
 * Espelha a lógica do backend — defesa em profundidade antes de mostrar ao recrutador.
 */
const BLIND_PATTERNS: RegExp[] = [
  /(g[eé]nero|gender|sexo|sex)\s*[:=]\s*.+/i,
  /(etnia|ethnicity|ra[cç]a|race|cor\s+da\s+pele)\s*[:=]\s*.+/i,
  /(estado\s+civil|marital\s+status)\s*[:=]\s*.+/i,
  /(idade|age|data\s+de\s+nascimento|date\s+of\s+birth|nascimento)\s*[:=]\s*.+/i,
];

export const BLIND_REDACTED = "[oculto — triagem às cegas]";

export function blindScreenCV(text: string): string {
  return text
    .split("\n")
    .map((line) => {
      let out = line;
      for (const re of BLIND_PATTERNS) {
        out = out.replace(re, BLIND_REDACTED);
      }
      return out;
    })
    .join("\n");
}
