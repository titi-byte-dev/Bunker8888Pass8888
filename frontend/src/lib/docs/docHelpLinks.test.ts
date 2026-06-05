import { describe, expect, it } from "vitest";
import { resolveDocHelp } from "./docHelpLinks";

describe("doc help links (DOC-013)", () => {
  it("resolve rotas específicas antes das genéricas", () => {
    expect(resolveDocHelp("/team/links")?.slug).toBe("journey-secret-link");
    expect(resolveDocHelp("/team/vaults")?.slug).toBe("journey-shared-vault");
    expect(resolveDocHelp("/team")?.slug).toBe("team-sharing");
  });

  it("resolve cofre, finanças e segurança", () => {
    expect(resolveDocHelp("/vault")?.slug).toBe("vault");
    expect(resolveDocHelp("/fin")?.slug).toBe("journey-finance-agent-saas");
    expect(resolveDocHelp("/security/sentinel")?.slug).toBe("journey-sentinel");
  });

  it("aceita slug explícito", () => {
    expect(resolveDocHelp("/unknown", "glossary")?.slug).toBe("glossary");
  });

  it("resolve onboarding RH com journey do agente", () => {
    expect(resolveDocHelp("/hr/onboarding")?.slug).toBe("journey-hr-agent-onboarding");
  });

  it("resolve inventário com journey do agente de operações", () => {
    expect(resolveDocHelp("/work/inventory")?.slug).toBe("journey-ops-agent-inventory");
  });

  it("resolve recrutamento RH com journey de triagem às cegas", () => {
    expect(resolveDocHelp("/hr/recruitment")?.slug).toBe("journey-hr-agent-recruitment");
  });

  it("resolve Open Banking com journey de reconciliação", () => {
    expect(resolveDocHelp("/fin/banking")?.slug).toBe("journey-finance-agent-reconcile");
  });

  it("devolve null para rotas sem doc", () => {
    expect(resolveDocHelp("/settings")).toBeNull();
  });
});
