import { describe, expect, it } from "vitest";
import { filterLogins, type DecodedLogin } from "./ui";

describe("vault ui (UI-004)", () => {
  const sample: DecodedLogin[] = [
    {
      meta: { id: "1", type: "login", created_at: "", updated_at: "" },
      login: { kind: "login", title: "GitHub", username: "dev", password: "x" },
    },
    {
      meta: { id: "2", type: "login", created_at: "", updated_at: "" },
      login: { kind: "login", title: "Email", username: "a@b.c", password: "y", url: "https://mail.test" },
    },
  ];

  it("filterLogins pesquisa por título e URL", () => {
    expect(filterLogins(sample, "git")).toHaveLength(1);
    expect(filterLogins(sample, "mail.test")).toHaveLength(1);
    expect(filterLogins(sample, "")).toHaveLength(2);
  });
});
