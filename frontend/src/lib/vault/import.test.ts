import { describe, expect, it } from "vitest";
import { deriveMasterKey, randomBytes } from "../crypto";
import {
  importCsvToVaultInputs,
  parseCsvRows,
  parsePasswordCsv,
  sanitizeField,
  toLoginItems,
} from "./import";

const BITWARDEN_SAMPLE = `folder,favorite,type,name,notes,fields,reprompt,login_uri,login_username,login_password
,Social,,GitHub,conta dev,,,https://github.com,devuser,gh-secret-1
,,,Inject,,,,,injectuser,gh-secret-2`;

const GENERIC_SAMPLE = `title,url,username,password,notes
Gmail,https://mail.google.com,me@gmail.com,mailpass,principal
@formula,https://x.com,u,+secret`;

describe("parseCsvRows", () => {
  it("campos com vírgulas entre aspas", () => {
    const rows = parseCsvRows('a,b\n"hello, world",2');
    expect(rows[1]).toEqual(["hello, world", "2"]);
  });
});

describe("sanitizeField", () => {
  it("remove prefixo de CSV injection", () => {
    expect(sanitizeField("=cmd|'/c calc'!A0")).toBe("cmd|'/c calc'!A0");
    expect(sanitizeField("+123")).toBe("123");
  });
});

describe("parsePasswordCsv", () => {
  it("formato Bitwarden", () => {
    const r = parsePasswordCsv(BITWARDEN_SAMPLE);
    expect(r.format).toBe("bitwarden");
    expect(r.rows).toHaveLength(2);
    expect(r.rows[0]).toMatchObject({
      title: "GitHub",
      username: "devuser",
      password: "gh-secret-1",
      url: "https://github.com",
    });
    expect(r.rows[1]).toMatchObject({
      title: "Inject",
      username: "injectuser",
      password: "gh-secret-2",
    });
  });

  it("formato genérico", () => {
    const r = parsePasswordCsv(GENERIC_SAMPLE);
    expect(r.format).toBe("generic");
    expect(r.rows[0]!.title).toBe("Gmail");
    expect(r.rows[1]!.password).toBe("secret");
  });

  it("rejeita cabeçalho desconhecido", () => {
    const r = parsePasswordCsv("foo,bar\n1,2");
    expect(r.format).toBe("unknown");
    expect(r.rows).toHaveLength(0);
  });
});

describe("importCsvToVaultInputs", () => {
  it("cifra logins importados", async () => {
    const key = await deriveMasterKey("pw", randomBytes(16), { iterations: 1_000, hash: "SHA-256" });
    const { parse, inputs } = await importCsvToVaultInputs(key, BITWARDEN_SAMPLE);
    expect(parse.rows.length).toBe(2);
    expect(inputs).toHaveLength(2);
    expect(inputs[0]!.type).toBe("login");
    expect(inputs[0]!.blob.length).toBeGreaterThan(10);
  });
});

describe("toLoginItems", () => {
  it("mapeia preview para LoginItem", () => {
    const parsed = parsePasswordCsv(BITWARDEN_SAMPLE);
    const items = toLoginItems(parsed.rows);
    expect(items[0]!.kind).toBe("login");
    expect(items[0]!.url).toBe("https://github.com");
  });
});
