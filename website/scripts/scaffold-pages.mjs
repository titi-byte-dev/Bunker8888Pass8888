/**
 * Gera wrappers finos de páginas por locale (evita duplicar dezenas de ficheiros à mão).
 * Didático: cada .astro só passa locale + slug ao template — copy vive nos locales TS.
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const pagesRoot = path.join(root, "src/pages");
const PRODUCT_SLUGS = ["vault", "security", "team", "workspace"];
const SERVICE_SLUGS = ["agents", "compliance", "deployment"];
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

  const productsDir = path.join(root, base, "products");
  const tProducts = templateImport(productsDir);
  write(
    `${base}/products/index.astro`,
    `---
import CatalogPage from "${tProducts}CatalogPage.astro";
---
<CatalogPage locale="${loc}" kind="products" />
`,
  );

  for (const slug of PRODUCT_SLUGS) {
    write(
      `${base}/products/${slug}.astro`,
      `---
import ProductPage from "${tProducts}ProductPage.astro";
---
<ProductPage locale="${loc}" slug="${slug}" />
`,
    );
  }

  const servicesDir = path.join(root, base, "services");
  const tServices = templateImport(servicesDir);
  write(
    `${base}/services/index.astro`,
    `---
import CatalogPage from "${tServices}CatalogPage.astro";
---
<CatalogPage locale="${loc}" kind="services" />
`,
  );

  for (const slug of SERVICE_SLUGS) {
    write(
      `${base}/services/${slug}.astro`,
      `---
import ServicePage from "${tServices}ServicePage.astro";
---
<ServicePage locale="${loc}" slug="${slug}" />
`,
    );
  }
}

console.log("pages: geradas para pt, fr, es, de (produtos + serviços)");
