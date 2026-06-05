import { describe, expect, it } from "vitest";
import { decryptDriveFile, encryptDriveFile, opaquePreviewB64 } from "./driveCrypto";

async function makeMasterKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey({ name: "AES-GCM", length: 256 }, false, ["encrypt", "decrypt"]);
}

describe("GOOGLE-002 driveCrypto", () => {
  it("round-trip ZK blob", async () => {
    const mk = await makeMasterKey();
    const blob = await encryptDriveFile(mk, "contrato.txt", "secreto");
    const payload = await decryptDriveFile(mk, blob);
    expect(payload.name).toBe("contrato.txt");
    expect(payload.content).toBe("secreto");
  });

  it("opaque preview truncates", () => {
    expect(opaquePreviewB64("abcdef", 3)).toBe("abc…");
  });
});
