import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// Configuração do Vite (bundler/dev server) com o plugin de Svelte.
// O bloco `test` configura o Vitest para correr em ambiente Node, onde a
// WebCrypto API (globalThis.crypto.subtle) já está disponível no Node 22.
export default defineConfig({
  plugins: [svelte()],
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
  },
});
