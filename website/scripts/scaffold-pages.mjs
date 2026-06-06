/**
 * Gera wrappers finos de páginas por locale (evita duplicar 28 ficheiros à mão).
 * Didático: cada .astro só passa locale + slug ao template — copy vive nos locales TS.
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const pagesRoot = path.join(root, "src/pages");
const SLUGS = ["vault", "security", "team", "workspace"];

const LOCALES = ["pt", "fr", "es", "de"];

function templateImport(pageDir) {
  const rel = path.relative(pagesRoot, pageDir);
  const depth = rel.split(/[/\\]/).filter(Boolean).length + 1;
  return `${"../".repeat(depth)}templates/`;
}

function write(rel, content) {
  const file = path.join(root, rel);
  fs.mkdirSync(path.dirname(file), { recursive: true });
  fs.writeFileSync(file, content, "utf8");
}

for (const loc of LOCALES) {
  const base = loc === "pt" ? "src/pages" : `src/pages/${loc}`;

  const homeDir = path.join(root, base);
  const tHome = templateImport(homeDir);
  write(
    `${base}/index.astro`,
    `---
import HomePage from "${tHome}HomePage.astro";
---
<HomePage locale="${loc}" />
`,
  );

  for (const sub of ["platform", "partners"]) {
    const dir = path.join(root, base, sub);
    const t = templateImport(dir);
    const comp = sub === "platform" ? "PlatformPage" : "PartnersPage";
    write(
      `${base}/${sub}/index.astro`,
      `---
import ${comp} from "${t}${comp}.astro";
---
<${comp} locale="${loc}" />
`,
    );
  }

  for (const slug of SLUGS) {
    const dir = path.join(root, base, "products");
    const t = templateImport(dir);
    write(
      `${base}/products/${slug}.astro`,
      `---
import ProductPage from "${t}ProductPage.astro";
---
<ProductPage locale="${loc}" slug="${slug}" />
`,
    );
  }
}

console.log("pages: geradas para pt, fr, es, de");
