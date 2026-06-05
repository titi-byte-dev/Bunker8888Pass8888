/**
 * Cifragem CAMPO-A-CAMPO das Fichas de Empregado (HR-001).
 *
 * Cada campo é cifrado de forma INDEPENDENTE, com a SUA própria chave aleatória.
 * Essa chave de campo é, por sua vez, cifrada ("embrulhada") com a Master Key do
 * utilizador. O servidor guarda dois blobs opacos por campo:
 *
 *   value_blob  = AES-GCM(chave_campo, valor)
 *   wrapped_key = AES-GCM(master_key, chave_campo)
 *
 * Porque é que cada campo tem a sua chave? Para o crypto-shredding (HR-003):
 * apagar o wrapped_key de UM campo torna esse valor eternamente indecifrável,
 * sem tocar nos restantes. É o "direito ao esquecimento" (RGPD Art. 17) cirúrgico.
 *
 *   ┌─────────────┐  gera   ┌────────────┐
 *   │ valor "5000"│────────▶│ chave_campo│ (32 bytes aleatorios)
 *   └─────┬───────┘         └─────┬──────┘
 *         │ AES-GCM(chave_campo)  │ AES-GCM(master_key)
 *         ▼                       ▼
 *     value_blob              wrapped_key
 *         └──────── servidor (opacos) ───────┘
 */
import { decrypt, encrypt, fromBytes, randomBytes, toBytes, type Bytes } from "$lib/crypto";

/** Importa 32 bytes como CryptoKey AES-GCM para cifrar/decifrar um campo. */
async function importFieldKey(raw: Bytes): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM", length: 256 }, false, [
    "encrypt",
    "decrypt",
  ]);
}

/** Resultado da cifra de um campo: os dois blobs prontos a enviar (bytes). */
export interface EncryptedField {
  valueBlob: Bytes;
  wrappedKey: Bytes;
}

/**
 * Cifra um valor de campo: gera uma chave de campo nova, cifra o valor com ela e
 * embrulha essa chave com a Master Key. Nunca reutiliza chaves entre campos.
 */
export async function encryptField(masterKey: CryptoKey, value: string): Promise<EncryptedField> {
  const rawFieldKey = randomBytes(32);
  const fieldKey = await importFieldKey(rawFieldKey);
  const valueBlob = await encrypt(fieldKey, toBytes(value));
  const wrappedKey = await encrypt(masterKey, rawFieldKey);
  return { valueBlob, wrappedKey };
}

/**
 * Decifra um campo: desembrulha a chave de campo com a Master Key e usa-a para
 * abrir o valor. Lança se o campo tiver sido crypto-shredded (sem wrapped_key)
 * ou se a Master Key estiver errada.
 */
export async function decryptField(
  masterKey: CryptoKey,
  valueBlob: Bytes,
  wrappedKey: Bytes,
): Promise<string> {
  const rawFieldKey = (await decrypt(masterKey, wrappedKey)) as Bytes;
  const fieldKey = await importFieldKey(rawFieldKey);
  return fromBytes(await decrypt(fieldKey, valueBlob));
}

/**
 * Campos sugeridos numa ficha de empregado. São apenas chaves de ESQUEMA
 * (visiveis ao servidor); o utilizador pode acrescentar campos personalizados.
 */
export const SUGGESTED_FIELDS = [
  "full_name",
  "email",
  "phone",
  "nif",
  "iban",
  "address",
  "role",
  "salary",
  "start_date",
  "notes",
] as const;

/** Etiquetas legíveis para os campos sugeridos (PT). */
export const FIELD_LABELS: Record<string, string> = {
  full_name: "Nome completo",
  email: "E-mail",
  phone: "Telefone",
  nif: "NIF",
  iban: "IBAN",
  address: "Morada",
  role: "Função",
  salary: "Salário",
  start_date: "Data de início",
  notes: "Notas",
};

/** Devolve a etiqueta de um campo, ou o próprio nome se for personalizado. */
export function fieldLabel(name: string): string {
  return FIELD_LABELS[name] ?? name;
}
