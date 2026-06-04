import adapter from "@sveltejs/adapter-auto";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    // adapter-auto escolhe o alvo de deploy (Node, static, etc.) em CI/produção.
    adapter: adapter(),
  },
};

export default config;
