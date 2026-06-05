import { describe, expect, it } from "vitest";
import { buildActionCommands, filterCommands, groupCommands } from "./commands";

describe("command palette (UI-006)", () => {
  it("filterCommands pesquisa por label e keywords", () => {
    const cmds = buildActionCommands();
    const hits = filterCommands(cmds, "sandbox");
    expect(hits.some((c) => c.id === "action-sandbox")).toBe(true);
    expect(filterCommands(cmds, "zzzzz")).toHaveLength(0);
  });

  it("groupCommands agrupa por tipo", () => {
    const cmds = buildActionCommands();
    const groups = groupCommands(cmds);
    expect(groups.get("action")?.length).toBe(cmds.length);
    expect(cmds.some((c) => c.id === "action-settings")).toBe(true);
    expect(cmds.some((c) => c.id === "action-shifts")).toBe(true);
  });
});
