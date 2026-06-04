import { describe, it, expect } from "vitest";
import { deriveMasterKey, randomBytes, type Bytes } from "../crypto";
import { sealItem, openItem, blobToBase64, blobFromBase64 } from "./items";
import { buildWsURL } from "./sync";
import type { LoginItem } from "./types";

const fastKdf = { iterations: 1_000, hash: "SHA-256" } as const;

describe("vault items", () => {
  it("round-trip login cifrado", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const login: LoginItem = {
      kind: "login",
      title: "App",
      username: "u",
      password: "p",
    };
    expect(await openItem(key, await sealItem(key, login))).toEqual(login);
  });

  it("base64 reversível", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), fastKdf);
    const blob = await sealItem(key, { kind: "note", title: "t", body: "b" });
    expect(await openItem(key, blobFromBase64(blobToBase64(blob)) as Bytes)).toMatchObject({
      kind: "note",
    });
  });
});

describe("vault sync URL", () => {
  it("converte http para ws com token", () => {
    const url = buildWsURL("http://localhost:8080", "abc123");
    expect(url).toBe("ws://localhost:8080/api/ws/vault?token=abc123");
  });

  it("converte https para wss", () => {
    const url = buildWsURL("https://api.example.com", "tok");
    expect(url.startsWith("wss://api.example.com/api/ws/vault?token=")).toBe(true);
  });
});
