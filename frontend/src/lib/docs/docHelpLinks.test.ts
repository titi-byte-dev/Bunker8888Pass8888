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
    expect(resolveDocHelp("/fin")?.slug).toBe("fin");
    expect(resolveDocHelp("/security/sentinel")?.slug).toBe("journey-sentinel");
  });

  it("aceita slug explícito", () => {
    expect(resolveDocHelp("/unknown", "glossary")?.slug).toBe("glossary");
  });

  it("devolve null para rotas sem doc", () => {
    expect(resolveDocHelp("/settings")).toBeNull();
  });
});
