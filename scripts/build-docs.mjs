/**
 * Gera manifest + JSON por página a partir de Markdown em docs/.
 *
 * Didático: uma única fonte (SSOT) no repositório; o frontend só lê JSON
 * gerado — evita duplicar texto nem manter HTML à mão.
 *
 * Convenções nos .md:
 *   - Frontmatter YAML entre --- ... ---
 *   - :::summary ... :::           → resumo sempre visível
 *   - :::concept{id title level}   → cartão expansível de conceito
 *   - :::level{level title} ... ::: → secção por nível de complexidade
 *
 * Uso: node scripts/build-docs.mjs
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { marked } from "marked";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const OUT_DIR = path.join(ROOT, "frontend", "src", "lib", "docs", "generated");

/** Pastas lidas como fontes de documentação */
const SOURCE_DIRS = [
  { dir: "docs/concepts", category: "concepts" },
  { dir: "docs/product", category: "product" },
  { dir: "docs/developer", category: "developer" },
  { dir: "docs/competitive", category: "competitive" },
  { dir: "docs/roadmap/04-user-journeys", category: "journeys" },
];

const LEVEL_LABELS = {
  1: "Essencial",
  2: "Intermédio",
  3: "Técnico",
};

const CATEGORY_LABELS = {
  concepts: "Conceitos",
  product: "Funcionalidades",
  developer: "Programador",
  competitive: "Concorrência",
  journeys: "Percursos",
};

marked.setOptions({ gfm: true, breaks: false });

function parseFrontmatter(raw) {
  const trimmed = raw.replace(/^\uFEFF/, "").trimStart();
  if (!trimmed.startsWith("---")) {
    return { meta: {}, body: trimmed };
  }
  // Segundo --- fecha o bloco YAML (índice após o primeiro)
  const close = trimmed.indexOf("\n---", 3);
  if (close === -1) {
    return { meta: {}, body: trimmed };
  }
  const yamlBlock = trimmed.slice(4, close).trim();
  const body = trimmed.slice(close + 4).replace(/^\r?\n/, "").trim();
  const meta = {};
  for (const line of yamlBlock.split(/\r?\n/)) {
    const kv = line.match(/^([a-zA-Z0-9_-]+):\s*(.*)$/);
    if (!kv) continue;
    const key = kv[1];
    let val = kv[2].trim();
    if (val.startsWith("[") && val.endsWith("]")) {
      meta[key] = val
        .slice(1, -1)
        .split(",")
        .map((s) => s.trim().replace(/^['"]|['"]$/g, ""))
        .filter(Boolean);
    } else if (val === "true" || val === "false") {
      meta[key] = val === "true";
    } else if (/^\d+$/.test(val)) {
      meta[key] = Number(val);
    } else {
      meta[key] = val.replace(/^['"]|['"]$/g, "");
    }
  }
  return { meta, body };
}

/** Extrai blocos :::tipo{attrs} ... ::: do corpo Markdown */
function extractBlocks(body) {
  const blockRe = /:::(\w+)(?:\{([^}]*)\})?\s*\n([\s\S]*?)\n:::/g;
  let summary = "";
  const concepts = [];
  const sections = [];
  let lastIndex = 0;
  let match;

  while ((match = blockRe.exec(body)) !== null) {
    const before = body.slice(lastIndex, match.index).trim();
    if (before) {
      sections.push({ level: 1, title: "", html: marked.parse(before), collapsed: false });
    }
    lastIndex = blockRe.lastIndex;

    const kind = match[1];
    const attrs = parseAttrs(match[2] ?? "");
    const inner = match[3].trim();

    if (kind === "summary") {
      summary = inner;
    } else if (kind === "concept") {
      concepts.push({
        id: attrs.id ?? slugify(attrs.title ?? "concept"),
        title: attrs.title ?? "Conceito",
        level: Number(attrs.level ?? 1),
        html: marked.parse(inner),
      });
    } else if (kind === "level") {
      sections.push({
        level: Number(attrs.level ?? 2),
        title: attrs.title ?? LEVEL_LABELS[attrs.level] ?? "Secção",
        html: marked.parse(inner),
        collapsed: Number(attrs.level ?? 2) > 1,
      });
    }
  }

  const tail = body.slice(lastIndex).trim();
  if (tail) {
    // Journeys sem blocos :::level — corpo inteiro como nível 1 + resto colapsável
    const parts = tail.split(/\n(?=## )/);
    if (parts.length > 1) {
      sections.push({
        level: 1,
        title: "Visão geral",
        html: marked.parse(parts[0]),
        collapsed: false,
      });
      for (let i = 1; i < parts.length; i++) {
        const heading = parts[i].match(/^## (.+)/);
        sections.push({
          level: 2,
          title: heading?.[1] ?? `Secção ${i}`,
          html: marked.parse(parts[i]),
          collapsed: true,
        });
      }
    } else {
      sections.push({
        level: 1,
        title: "",
        html: marked.parse(tail),
        collapsed: false,
      });
    }
  }

  return { summary, concepts, sections };
}

function parseAttrs(raw) {
  const out = {};
  const re = /(\w+)=("([^"]*)"|'([^']*)'|(\S+))/g;
  let m;
  while ((m = re.exec(raw)) !== null) {
    out[m[1]] = m[3] ?? m[4] ?? m[5];
  }
  return out;
}

function slugify(s) {
  return s
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-|-$/g, "");
}

function inferJourneyMeta(filePath, body) {
  const base = path.basename(filePath, ".md");
  const titleMatch = body.match(/^#\s+(.+)/m);
  const actorMatch = body.match(/>\s*\*\*Ator:\*\*\s*([^·]+)/);
  return {
    slug: base.replace(/^journey-/, "journey-"),
    title: titleMatch?.[1] ?? base,
    audience: ["user", "admin"],
    feature: base.replace(/^journey-/, "").replace(/-/g, "_"),
    level: 1,
    in_app: true,
    actor: actorMatch?.[1]?.trim(),
  };
}

function collectMarkdownFiles() {
  const files = [];
  for (const { dir, category } of SOURCE_DIRS) {
    const abs = path.join(ROOT, dir);
    if (!fs.existsSync(abs)) continue;
    for (const name of fs.readdirSync(abs)) {
      if (!name.endsWith(".md")) continue;
      files.push({ path: path.join(abs, name), category, dir });
    }
  }
  return files;
}

function buildDoc(filePath, category) {
  const raw = fs.readFileSync(filePath, "utf8");
  const { meta, body } = parseFrontmatter(raw);
  const inferred =
    category === "journeys" ? inferJourneyMeta(filePath, body) : {};
  const slug =
    meta.slug ??
    inferred.slug ??
    path.basename(filePath, ".md").replace(/^journey-/, "journey-");

  const { summary, concepts, sections } = extractBlocks(body);
  const maxLevel = Math.max(
    1,
    ...sections.map((s) => s.level),
    ...concepts.map((c) => c.level),
    meta.level ?? 1,
  );

  return {
    slug,
    title: meta.title ?? inferred.title ?? slug,
    category,
    categoryLabel: CATEGORY_LABELS[category] ?? category,
    audience: meta.audience ?? inferred.audience ?? ["user"],
    layer: meta.layer ?? ["product"],
    feature: meta.feature ?? inferred.feature,
    level: meta.level ?? inferred.level ?? 1,
    maxLevel,
    in_app: meta.in_app ?? inferred.in_app ?? true,
    summary: meta.summary ?? summary ?? "",
    concepts,
    sections,
    related: meta.related ?? [],
    actor: meta.actor ?? inferred.actor,
    order: meta.order ?? 99,
  };
}

function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const files = collectMarkdownFiles();
  const docs = files.map((f) => buildDoc(f.path, f.category));
  docs.sort((a, b) => a.category.localeCompare(b.category) || a.order - b.order);

  const manifest = {
    generatedAt: new Date().toISOString(),
    categories: Object.entries(CATEGORY_LABELS).map(([id, label]) => ({ id, label })),
    levelLabels: LEVEL_LABELS,
    docs: docs.map((d) => ({
      slug: d.slug,
      title: d.title,
      category: d.category,
      categoryLabel: d.categoryLabel,
      audience: d.audience,
      layer: d.layer,
      feature: d.feature,
      level: d.level,
      maxLevel: d.maxLevel,
      in_app: d.in_app,
      summary: d.summary,
      related: d.related,
      actor: d.actor,
      order: d.order,
    })),
  };

  fs.writeFileSync(
    path.join(OUT_DIR, "manifest.json"),
    JSON.stringify(manifest, null, 2),
    "utf8",
  );

  for (const doc of docs) {
  const outPath = path.join(OUT_DIR, `${doc.slug}.json`);
  // Remove JSON antigo se o slug mudou (ex.: renomear ficheiro)
  for (const name of fs.readdirSync(OUT_DIR)) {
    if (!name.endsWith(".json") || name === "manifest.json") continue;
    try {
      const prev = JSON.parse(fs.readFileSync(path.join(OUT_DIR, name), "utf8"));
      if (prev.slug === doc.slug && name !== `${doc.slug}.json`) {
        fs.unlinkSync(path.join(OUT_DIR, name));
      }
    } catch {
      /* ignorar */
    }
  }
  fs.writeFileSync(outPath, JSON.stringify(doc, null, 2), "utf8");
  }

  // Apaga JSON órfãos (slugs antigos após renomear ficheiros)
  const validNames = new Set(docs.map((d) => `${d.slug}.json`));
  validNames.add("manifest.json");
  for (const name of fs.readdirSync(OUT_DIR)) {
    if (name.endsWith(".json") && !validNames.has(name)) {
      fs.unlinkSync(path.join(OUT_DIR, name));
    }
  }

  console.log(`docs: ${docs.length} páginas → ${OUT_DIR}`);
}

main();
