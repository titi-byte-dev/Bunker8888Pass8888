/**
 * Orquestração das Fichas de Empregado (HR-001).
 *
 * Junta a cripto campo-a-campo (employees.ts) com o cliente HTTP (employeesApi.ts).
 * Exige sessão + Master Key desbloqueada — a mesma chave que protege o cofre
 * pessoal embrulha cada chave de campo.
 */
import { loadSessionToken } from "$lib/session";
import { getMasterKey } from "$lib/vault/masterKeyStore";
import { decryptField, encryptField } from "./employees";
import { EmployeesAPI, fieldBlob, type EmployeeRecordDTO } from "./employeesApi";

/** Um campo já decifrado (ou marcado como inacessível). */
export interface DecryptedField {
  name: string;
  value: string;
  /** true quando o campo foi crypto-shredded (sem chave) — valor irrecuperável. */
  shredded: boolean;
}

/** Uma ficha aberta: metadados + campos decifrados. */
export interface OpenRecord {
  id: string;
  createdAt: string;
  updatedAt: string;
  fields: DecryptedField[];
}

export interface RecordSummary {
  id: string;
  createdAt: string;
  updatedAt: string;
}

function api(): EmployeesAPI {
  const token = loadSessionToken();
  if (!token) throw new Error("Sessao expirada — inicia sessao de novo.");
  return new EmployeesAPI("", token);
}

function requireMasterKey(): CryptoKey {
  const mk = getMasterKey();
  if (!mk) throw new Error("Cofre bloqueado — desbloqueia para gerir fichas.");
  return mk;
}

/** Lista as fichas (sem decifrar nada — não têm PII fora dos campos). */
export async function listRecords(): Promise<RecordSummary[]> {
  const recs = await api().list();
  return recs.map((r: EmployeeRecordDTO) => ({
    id: r.id,
    createdAt: r.created_at,
    updatedAt: r.updated_at,
  }));
}

/** Cria uma ficha vazia e devolve o seu id. */
export async function createRecord(): Promise<string> {
  return (await api().create()).id;
}

/** Abre uma ficha e decifra todos os campos acessíveis. */
export async function openRecord(id: string): Promise<OpenRecord> {
  const mk = requireMasterKey();
  const full = await api().get(id);
  const fields: DecryptedField[] = [];
  for (const f of full.fields ?? []) {
    if (f.shredded || !f.wrapped_key) {
      fields.push({ name: f.field_name, value: "", shredded: true });
      continue;
    }
    try {
      const value = await decryptField(mk, fieldBlob(f.value_blob), fieldBlob(f.wrapped_key));
      fields.push({ name: f.field_name, value, shredded: false });
    } catch {
      // Master Key errada ou blob corrompido — não revela nada, marca inacessível.
      fields.push({ name: f.field_name, value: "", shredded: true });
    }
  }
  return { id: full.id, createdAt: full.created_at, updatedAt: full.updated_at, fields };
}

/** Cifra e grava um campo (cria ou actualiza). */
export async function saveField(id: string, name: string, value: string): Promise<void> {
  const mk = requireMasterKey();
  const { valueBlob, wrappedKey } = await encryptField(mk, value);
  await api().putField(id, name, valueBlob, wrappedKey);
}

/** Remove um campo da ficha. */
export async function removeField(id: string, name: string): Promise<void> {
  await api().removeField(id, name);
}

/** Apaga a ficha inteira. */
export async function deleteRecord(id: string): Promise<void> {
  await api().remove(id);
}
