import { defineConfig } from "astro/config";

/** Site institucional estático — aegispass.com (independente da app SvelteKit). */
export default defineConfig({
  site: "https://aegispass.com",
  output: "static",
  i18n: {
    defaultLocale: "pt",
    locales: ["pt", "fr", "es", "de"],
    routing: {
      prefixDefaultLocale: false,
    },
  },
});
