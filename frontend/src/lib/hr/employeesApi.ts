/**
 * Cliente da API das Fichas de Empregado (HR-001).
 *
 * Só transporta bytes (base64). Toda a cifra/decifra acontece no cliente; o
 * servidor encaminha blobs opacos e faz cumprir a posse da ficha.
 */
import { bytesToBase64, base64ToBytes } from "$lib/auth";
import type { Bytes } from "$lib/crypto";

export interface EmployeeRecordDTO {
  id: string;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

export interface EmployeeFieldDTO {
  id: string;
  field_name: string;
  value_blob: string;
  wrapped_key?: string;
  shredded: boolean;
  created_at: string;
  updated_at: string;
}

export interface EmployeeRecordFullDTO extends EmployeeRecordDTO {
  fields: EmployeeFieldDTO[];
}

export class EmployeesAPI {
  constructor(
    private baseURL: string,
    private token: string,
  ) {}

  private async fetch(path: string, init: RequestInit = {}): Promise<Response> {
    const res = await fetch(this.baseURL + path, {
      ...init,
      headers: {
        Authorization: `Bearer ${this.token}`,
        "Content-Type": "application/json",
        ...(init.headers ?? {}),
      },
    });
    if (!res.ok) {
      const err = (await res.json().catch(() => ({}))) as { error?: string };
      throw new Error(err.error ?? `HTTP ${res.status}`);
    }
    return res;
  }

  /** Cria uma ficha vazia; os campos entram depois, um a um. */
  async create(): Promise<EmployeeRecordDTO> {
    const res = await this.fetch("/api/hr/employees", { method: "POST" });
    return (await res.json()) as EmployeeRecordDTO;
  }

  /** Lista as fichas do utilizador (sem campos). */
  async list(): Promise<EmployeeRecordDTO[]> {
    const res = await this.fetch("/api/hr/employees");
    return ((await res.json()) as { records: EmployeeRecordDTO[] }).records ?? [];
  }

  /** Devolve uma ficha com todos os seus campos cifrados. */
  async get(id: string): Promise<EmployeeRecordFullDTO> {
    const res = await this.fetch(`/api/hr/employees/${encodeURIComponent(id)}`);
    return (await res.json()) as EmployeeRecordFullDTO;
  }

  /** Apaga a ficha inteira. */
  async remove(id: string): Promise<void> {
    await this.fetch(`/api/hr/employees/${encodeURIComponent(id)}`, { method: "DELETE" });
  }

  /** Cria ou actualiza um campo (upsert por nome). */
  async putField(
    id: string,
    field: string,
    valueBlob: Bytes,
    wrappedKey: Bytes,
  ): Promise<EmployeeFieldDTO> {
    const res = await this.fetch(
      `/api/hr/employees/${encodeURIComponent(id)}/fields/${encodeURIComponent(field)}`,
      {
        method: "PUT",
        body: JSON.stringify({
          value_blob: bytesToBase64(valueBlob),
          wrapped_key: bytesToBase64(wrappedKey),
        }),
      },
    );
    return (await res.json()) as EmployeeFieldDTO;
  }

  /** Remove um campo da ficha. */
  async removeField(id: string, field: string): Promise<void> {
    await this.fetch(
      `/api/hr/employees/${encodeURIComponent(id)}/fields/${encodeURIComponent(field)}`,
      { method: "DELETE" },
    );
  }
}

/** Converte um blob DTO (base64) em bytes. Atalho para o serviço. */
export function fieldBlob(b64: string): Bytes {
  return base64ToBytes(b64);
}
