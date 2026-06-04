import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

// vitePreprocess permite usar TypeScript dentro dos componentes .svelte.
export default {
  preprocess: vitePreprocess(),
};
