/**
 * sync-tokens.mjs — fonte unica de verdade dos design tokens (SITE-001).
 *
 * Didatico: o site institucional NAO redefine cores. Copia o tokens.css da app
 * (frontend/src/lib/design/tokens.css) para src/lib/tokens.generated.css antes
 * de cada dev/build. Assim a marca do site segue SEMPRE a da app — um so sitio
 * a editar. O ficheiro gerado e .gitignored (artefacto, nao fonte).
 *
 *   frontend/src/lib/design/tokens.css   (FONTE)
 *            |  copia em predev/prebuild
 *            v
 *   site/src/lib/tokens.generated.css    (ARTEFACTO)
 */
import { readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const SRC = resolve(here, "../../frontend/src/lib/design/tokens.css");
const OUT = resolve(here, "../src/lib/tokens.generated.css");

const banner =
  "/* GERADO por scripts/sync-tokens.mjs — NAO editar. Fonte: frontend/src/lib/design/tokens.css */\n";

try {
  const css = readFileSync(SRC, "utf8");
  mkdirSync(dirname(OUT), { recursive: true });
  writeFileSync(OUT, banner + css);
  console.log(`[sync-tokens] OK ${SRC} -> ${OUT}`);
} catch (err) {
  console.error(`[sync-tokens] FALHOU: ${err.message}`);
  process.exit(1);
}
