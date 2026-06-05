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
 *   - ```mermaid ... ```          → fluxo interactivo (FlowPlayer na app)
 *   - :::flow{id title type} ... ::: → fluxo com título explícito
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

/** Tipo de diagrama Mermaid detectado a partir da primeira linha */
function detectMermaidType(source) {
  const first = source.trim().split("\n")[0]?.trim() ?? "";
  if (first.startsWith("sequenceDiagram")) return "sequence";
  if (first.startsWith("flowchart") || first.startsWith("graph ")) return "flowchart";
  return "diagram";
}

/**
 * Extrai passos legíveis de um sequenceDiagram para o FlowPlayer animado.
 * Didático: cada seta vira um «momento» que o utilizador percorre passo a passo.
 */
function parseSequenceSteps(source) {
  const steps = [];
  for (const line of source.split("\n")) {
    const t = line.trim();
    if (!t || t.startsWith("sequenceDiagram") || t.startsWith("participant")) continue;
    if (t === "alt" || t.startsWith("alt ") || t === "else" || t === "end") {
      steps.push({ kind: "branch", label: t });
      continue;
    }
    // -->> antes de ->> para não partir G-->>App em tokens errados
    const m = t.match(/^([A-Za-z0-9]+)\s*(-->>|->>|->)\s*([A-Za-z0-9]+)\s*:\s*(.+)$/);
    if (m) {
      steps.push({
        kind: "message",
        from: m[1],
        to: m[3],
        arrow: m[2],
        label: m[4].trim(),
      });
    }
  }
  return steps;
}

/** Participantes declarados no sequenceDiagram (participant X as Label). */
function parseParticipants(source) {
  const out = [];
  for (const line of source.split("\n")) {
    const m = line.trim().match(/^participant\s+(\S+)\s+as\s+(.+)$/);
    if (m) out.push({ id: m[1], label: m[2].trim() });
  }
  return out;
}

function normalizeParticipantId(id) {
  return id.replace(/-+$/, "");
}

/**
 * Converte sequenceDiagram em grafo para Svelte Flow (DOC-011/012).
 * Layout horizontal: actores em fila; arestas = mensagens.
 */
function sequenceToGraph(source, steps) {
  const participants = parseParticipants(source);
  const messages = steps.filter((s) => s.kind === "message");
  const labels = new Map(participants.map((p) => [p.id, p.label]));

  for (const msg of messages) {
    const from = normalizeParticipantId(msg.from);
    const to = normalizeParticipantId(msg.to);
    if (!labels.has(from)) labels.set(from, from);
    if (!labels.has(to)) labels.set(to, to);
  }

  const order = participants.length
    ? participants.map((p) => p.id)
    : [...labels.keys()];
  for (const id of labels.keys()) {
    if (!order.includes(id)) order.push(id);
  }

  const nodes = order.map((id, i) => ({
    id,
    label: labels.get(id) ?? id,
    x: i * 200,
  }));

  const edges = messages.map((msg, i) => ({
    id: `step-${i}`,
    source: normalizeParticipantId(msg.from),
    target: normalizeParticipantId(msg.to),
    label: msg.label,
    dashed: msg.arrow.includes("--"),
  }));

  return { nodes, edges };
}

/** Acrescenta renderer (mermaid | svelteflow) e grafo pré-calculado. */
function enrichFlow(flow, attrs = {}) {
  const out = { ...flow };
  const wantsSvelte =
    attrs.renderer === "svelteflow" ||
    (attrs.renderer !== "mermaid" && out.type === "sequence" && out.steps?.length > 0);

  if (wantsSvelte && out.type === "sequence") {
    out.renderer = "svelteflow";
    out.graph = sequenceToGraph(out.source, out.steps);
  } else {
    out.renderer = "mermaid";
  }
  return out;
}

/** Substitui blocos mermaid por comentários HTML (placeholders para o Svelte). */
function injectFlowPlaceholders(text, flows, flowCounter) {
  let counter = flowCounter;
  const replaced = text.replace(/```mermaid\s*\n([\s\S]*?)\n```/g, (_, raw) => {
    counter += 1;
    const source = raw.trim();
    const id = `flow-${counter}`;
    const type = detectMermaidType(source);
    flows.push(
      enrichFlow({
        id,
        title: "",
        type,
        source,
        steps: type === "sequence" ? parseSequenceSteps(source) : [],
      }),
    );
    return `\n<!--DOC_FLOW:${id}-->\n`;
  });
  return { text: replaced, counter };
}

/** Remove <h2> duplicado quando o título da secção já vem do ## Markdown. */
function stripDuplicateHeading(html, title) {
  if (!title) return html;
  const esc = title.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  return html.replace(new RegExp(`<h2[^>]*>\\s*${esc}\\s*</h2>\\s*`, "i"), "");
}

/** Converte Markdown em HTML e recolhe fluxos Mermaid embutidos. */
function markdownToHtmlWithFlows(md) {
  const flows = [];
  const { text } = injectFlowPlaceholders(md, flows, 0);
  return { html: marked.parse(text), flows };
}

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
      const parsed = markdownToHtmlWithFlows(before);
      sections.push({
        level: 1,
        title: "",
        html: parsed.html,
        flows: parsed.flows,
        collapsed: false,
      });
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
      const parsed = markdownToHtmlWithFlows(inner);
      sections.push({
        level: Number(attrs.level ?? 2),
        title: attrs.title ?? LEVEL_LABELS[attrs.level] ?? "Secção",
        html: parsed.html,
        flows: parsed.flows,
        collapsed: Number(attrs.level ?? 2) > 1,
      });
    } else if (kind === "flow") {
      const source = inner.replace(/^```mermaid\s*\n?|\n?```$/g, "").trim();
      const id = attrs.id ?? `flow-${sections.length + 1}`;
      const type = attrs.type ?? detectMermaidType(source);
      const flow = enrichFlow(
        {
          id,
          title: attrs.title ?? "Fluxo",
          type,
          source,
          steps: type === "sequence" ? parseSequenceSteps(source) : [],
        },
        attrs,
      );
      sections.push({
        level: Number(attrs.level ?? 2),
        title: flow.title,
        html: `<!--DOC_FLOW:${id}-->`,
        flows: [flow],
        collapsed: false,
      });
    }
  }

  const tail = body.slice(lastIndex).trim();
  if (tail) {
    // Journeys sem blocos :::level — corpo inteiro como nível 1 + resto colapsável
    const parts = tail.split(/\n(?=## )/);
    if (parts.length > 1) {
      const intro = markdownToHtmlWithFlows(parts[0]);
      sections.push({
        level: 1,
        title: "Visão geral",
        html: intro.html,
        flows: intro.flows,
        collapsed: false,
      });
      for (let i = 1; i < parts.length; i++) {
        const heading = parts[i].match(/^## (.+)/);
        const sectionTitle = heading?.[1] ?? `Secção ${i}`;
        const parsed = markdownToHtmlWithFlows(parts[i]);
        sections.push({
          level: 2,
          title: sectionTitle,
          html: stripDuplicateHeading(parsed.html, sectionTitle),
          flows: parsed.flows,
          collapsed: true,
        });
      }
    } else {
      const parsed = markdownToHtmlWithFlows(tail);
      sections.push({
        level: 1,
        title: "",
        html: parsed.html,
        flows: parsed.flows,
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
