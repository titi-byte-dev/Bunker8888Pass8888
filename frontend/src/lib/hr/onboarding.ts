/**
 * Onboarding de empregado em 1 clique (HR-007).
 *
 * Orquestra, num so fluxo, peças que ja existem:
 *   1. cria a ficha de empregado (HR-001),
 *   2. cifra os campos iniciais campo-a-campo (nome, e-mail, função),
 *   3. gera um alias de e-mail (MAIL-001) que reencaminha para o e-mail real.
 *
 * Cada passo e auditado na cadeia imutavel (HR-002) pelos próprios módulos.
 */
import { createRecord, saveField } from "./employeesService";
import { createAlias } from "$lib/mail/aliases";

export interface OnboardingInput {
  fullName: string;
  email: string;
  role?: string;
  /** Ficha já criada (ex.: após aprovar sugestão AGENT-007). */
  recordId?: string;
}

export interface OnboardingStep {
  label: string;
  done: boolean;
}

export interface OnboardingResult {
  recordId: string;
  alias: string;
  steps: OnboardingStep[];
}

/**
 * Executa o onboarding completo. Reporta progresso através de onStep (opcional)
 * para a UI mostrar a checklist a preencher-se em tempo real.
 */
export async function onboardEmployee(
  input: OnboardingInput,
  onStep?: (steps: OnboardingStep[]) => void,
): Promise<OnboardingResult> {
  const steps: OnboardingStep[] = [
    { label: "Criar ficha de empregado", done: false },
    { label: "Cifrar nome", done: false },
    { label: "Cifrar e-mail", done: false },
    { label: "Gerar alias de e-mail", done: false },
  ];
  const tick = (i: number) => {
    steps[i].done = true;
    onStep?.([...steps]);
  };

  const recordId = input.recordId ?? (await createRecord());
  if (!input.recordId) {
    tick(0);
  } else {
    steps[0].done = true;
    onStep?.([...steps]);
  }

  await saveField(recordId, "full_name", input.fullName);
  tick(1);

  await saveField(recordId, "email", input.email);
  tick(2);

  if (input.role && input.role.trim()) {
    await saveField(recordId, "role", input.role.trim());
  }

  const alias = await createAlias(input.email, `${input.fullName} — onboarding`);
  tick(3);

  return { recordId, alias: alias.aliasAddress, steps };
}
