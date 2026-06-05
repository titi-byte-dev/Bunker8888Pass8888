import { describe, expect, it } from "vitest";
import { decryptSubscription, encryptSubscription, type SubscriptionPayload } from "./subscriptions";

async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]);
}

const payload: SubscriptionPayload = {
  name: "Figma",
  cost: 144,
  currency: "EUR",
  cycle: "yearly",
  category: "Design",
  vaultItemId: "v-123",
  vaultItemTitle: "Figma login",
  lastUsedAt: "2026-01-01T00:00:00Z",
  active: true,
};

describe("FIN-001 subscrições (cifragem)", () => {
  it("cifra e decifra uma subscrição (round-trip)", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptSubscription(mk, payload);
    expect(await decryptSubscription(mk, blob)).toEqual(payload);
  });

  it("o blob não contém o nome em claro", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptSubscription(mk, payload);
    const asText = new TextDecoder().decode(blob);
    expect(asText).not.toContain("Figma");
  });

  it("não decifra com outra Master Key", async () => {
    const mk = await makeMasterKey();
    const other = await makeMasterKey();
    const blob = await encryptSubscription(mk, payload);
    await expect(decryptSubscription(other, blob)).rejects.toThrow();
  });
});
