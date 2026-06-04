import { encrypt, decrypt, toBytes, fromBytes, type Bytes } from "../crypto";
import type { VaultPayload, VaultItemType } from "./types";

const PAYLOAD_VERSION = 1;

interface PayloadEnvelope {
  v: number;
  kind: VaultItemType;
  data: VaultPayload;
}

export async function sealItem(key: CryptoKey, payload: VaultPayload): Promise<Bytes> {
  const envelope: PayloadEnvelope = { v: PAYLOAD_VERSION, kind: payload.kind, data: payload };
  return encrypt(key, toBytes(JSON.stringify(envelope)));
}

export async function openItem(key: CryptoKey, blob: Bytes): Promise<VaultPayload> {
  const envelope = JSON.parse(fromBytes(await decrypt(key, blob))) as PayloadEnvelope;
  if (envelope.v !== PAYLOAD_VERSION) throw new Error(`versão não suportada: ${envelope.v}`);
  return envelope.data;
}

export function blobFromBase64(b64: string): Bytes {
  const bin = atob(b64);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out as Bytes;
}

export function blobToBase64(blob: Bytes): string {
  let s = "";
  for (let i = 0; i < blob.length; i++) s += String.fromCharCode(blob[i]!);
  return btoa(s);
}
