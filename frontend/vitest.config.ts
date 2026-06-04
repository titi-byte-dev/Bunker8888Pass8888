import path from "node:path";
import { defineConfig } from "vitest/config";

export default defineConfig({
  resolve: {
    alias: {
      // Espelha o alias $lib do SvelteKit para os testes Vitest.
      $lib: path.resolve("./src/lib"),
    },
  },
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
  },
});
