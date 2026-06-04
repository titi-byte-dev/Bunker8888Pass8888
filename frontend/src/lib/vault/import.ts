/**
 * Importação de palavras-passe (VAULT-015).
 *
 * Didático: a importação corre NO CLIENTE — o ficheiro CSV nunca é enviado em
 * claro ao servidor. Parse → LoginItem → sealItem (AES-GCM) → POST /api/vault.
 *
 * ⚠️ Segurança: sanitizamos campos que começam por =+-@ (CSV injection em Excel).
 * Limitamos o nº de linhas para evitar DoS no browser.
 */
import { sealItem, blobToBase64 } from "./items";
import type { LoginItem, VaultItemInput } from "./types";

/** Máximo de credenciais por importação (evita bloquear a UI). */
export const MAX_IMPORT_ROWS = 2_000;

export type ImportFormat = "bitwarden" | "generic" | "unknown";

export interface ImportPreviewRow {
  line: number;
  title: string;
  username: string;
  password: string;
  url?: string;
  notes?: string;
}

export interface ImportParseResult {
  format: ImportFormat;
  rows: ImportPreviewRow[];
  skipped: number;
  warnings: string[];
}

export interface ImportUploadResult {
  created: number;
  failed: number;
  errors: string[];
}

/** Alias de colunas conhecidos (case-insensitive). */
const COLUMN_ALIASES: Record<string, keyof Omit<ImportPreviewRow, "line">> = {
  name: "title",
  title: "title",
  login_username: "username",
  username: "username",
  user: "username",
  login_password: "password",
  password: "password",
  login_uri: "url",
  url: "url",
  uri: "url",
  website: "url",
  notes: "notes",
  note: "notes",
};

/**
 * Parseia CSV de export (Bitwarden, Chrome, Firefox, etc.).
 * A primeira linha deve ser o cabeçalho.
 */
export function parsePasswordCsv(text: string): ImportParseResult {
  const warnings: string[] = [];
  const rows = parseCsvRows(text);
  if (rows.length === 0) {
    return { format: "unknown", rows: [], skipped: 0, warnings: ["ficheiro vazio"] };
  }

  const [headerRow, ...dataRows] = rows;
  const headers = headerRow.map(normalizeHeader);
  const format = detectFormat(headers);
  if (format === "unknown") {
    return {
      format,
      rows: [],
      skipped: dataRows.length,
      warnings: ["cabeçalho não reconhecido; esperado name/username/password ou export Bitwarden"],
    };
  }

  if (dataRows.length > MAX_IMPORT_ROWS) {
    warnings.push(`limite de ${MAX_IMPORT_ROWS} linhas; linhas extra ignoradas`);
  }

  const preview: ImportPreviewRow[] = [];
  let skipped = 0;

  for (let i = 0; i < Math.min(dataRows.length, MAX_IMPORT_ROWS); i++) {
    const lineNum = i + 2;
    const cells = dataRows[i]!;
    const mapped = mapRow(headers, cells);
    if (!mapped.password && !mapped.title) {
      skipped++;
      continue;
    }
    if (!mapped.title) {
      warnings.push(`linha ${lineNum}: título em falta — ignorada`);
      skipped++;
      continue;
    }
    if (!mapped.password) {
      warnings.push(`linha ${lineNum}: password em falta — ignorada`);
      skipped++;
      continue;
    }
    preview.push({
      line: lineNum,
      title: mapped.title,
      username: mapped.username ?? "",
      password: mapped.password,
      url: mapped.url,
      notes: mapped.notes,
    });
  }

  return { format, rows: preview, skipped, warnings };
}

/** Converte linhas parseadas em LoginItem prontos a cifrar. */
export function toLoginItems(parsed: ImportPreviewRow[]): LoginItem[] {
  return parsed.map((r) => ({
    kind: "login" as const,
    title: r.title,
    username: r.username,
    password: r.password,
    ...(r.url ? { url: r.url } : {}),
    ...(r.notes ? { notes: r.notes } : {}),
  }));
}

/** Cifra cada login e devolve inputs para VaultAPI.create. */
export async function sealImportItems(
  key: CryptoKey,
  items: LoginItem[],
): Promise<VaultItemInput[]> {
  const out: VaultItemInput[] = [];
  for (const item of items) {
    const blob = await sealItem(key, item);
    out.push({ type: "login", blob: blobToBase64(blob) });
  }
  return out;
}

/** Envia itens cifrados para o servidor (um POST por item). */
export async function uploadImportItems(
  create: (input: VaultItemInput) => Promise<unknown>,
  inputs: VaultItemInput[],
  onProgress?: (done: number, total: number) => void,
): Promise<ImportUploadResult> {
  let created = 0;
  let failed = 0;
  const errors: string[] = [];
  for (let i = 0; i < inputs.length; i++) {
    try {
      await create(inputs[i]!);
      created++;
    } catch (e) {
      failed++;
      errors.push(e instanceof Error ? e.message : "erro desconhecido");
    }
    onProgress?.(i + 1, inputs.length);
  }
  return { created, failed, errors };
}

/** Fluxo completo: CSV → cifrar → pronto a enviar. */
export async function importCsvToVaultInputs(
  key: CryptoKey,
  csvText: string,
): Promise<{ parse: ImportParseResult; inputs: VaultItemInput[] }> {
  const parse = parsePasswordCsv(csvText);
  const items = toLoginItems(parse.rows);
  const inputs = await sealImportItems(key, items);
  return { parse, inputs };
}

function detectFormat(headers: string[]): ImportFormat {
  const set = new Set(headers);
  if (set.has("login_username") || set.has("login_password")) return "bitwarden";
  const hasTitle = set.has("name") || set.has("title");
  const hasUser = set.has("username") || set.has("user");
  const hasPass = set.has("password");
  if (hasTitle && hasPass) return "generic";
  if (hasUser && hasPass && (hasTitle || set.has("url"))) return "generic";
  return "unknown";
}

function normalizeHeader(h: string): string {
  return sanitizeField(h).toLowerCase().replace(/\s+/g, "_");
}

interface MappedRow {
  title?: string;
  username?: string;
  password?: string;
  url?: string;
  notes?: string;
}

function mapRow(headers: string[], cells: string[]): MappedRow {
  const out: MappedRow = {};
  for (let i = 0; i < headers.length; i++) {
    const key = COLUMN_ALIASES[headers[i]!];
    if (!key) continue;
    const val = sanitizeField(cells[i] ?? "");
    if (val === "") continue;
    out[key] = val;
  }
  return out;
}

/**
 * Remove prefixos perigosos de CSV injection (=cmd|'/etc/passwd).
 * Didático: ao abrir CSV no Excel, células que começam por = são fórmulas.
 */
export function sanitizeField(raw: string): string {
  let s = raw.trim();
  if (s.startsWith('"') && s.endsWith('"') && s.length >= 2) {
    s = s.slice(1, -1).replace(/""/g, '"');
  }
  while (/^[=+\-@]/.test(s)) {
    s = s.slice(1).trimStart();
  }
  return s;
}

/**
 * Parser CSV mínimo (RFC 4180 simplificado) — suporta campos entre aspas.
 */
export function parseCsvRows(text: string): string[][] {
  const normalized = text.replace(/^\uFEFF/, "").replace(/\r\n/g, "\n").replace(/\r/g, "\n");
  const rows: string[][] = [];
  let row: string[] = [];
  let field = "";
  let inQuotes = false;

  for (let i = 0; i < normalized.length; i++) {
    const c = normalized[i]!;
    if (inQuotes) {
      if (c === '"') {
        if (normalized[i + 1] === '"') {
          field += '"';
          i++;
        } else {
          inQuotes = false;
        }
      } else {
        field += c;
      }
      continue;
    }
    if (c === '"') {
      inQuotes = true;
    } else if (c === ",") {
      row.push(field);
      field = "";
    } else if (c === "\n") {
      row.push(field);
      field = "";
      if (row.some((cell) => cell.trim() !== "")) rows.push(row);
      row = [];
    } else {
      field += c;
    }
  }
  row.push(field);
  if (row.some((cell) => cell.trim() !== "")) rows.push(row);
  return rows;
}
