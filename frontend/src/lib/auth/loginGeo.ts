/**
 * Headers GPS opcionais para login (DW-004 Sentinel).
 *
 * Didático: sem GPS o Sentinel não consegue calcular viagem impossível —
 * o login continua, mas com menos protecção.
 */
import { geoHeaders, getCurrentPosition } from "$lib/vault/geofence";

/** Obtém headers X-Geo-* se o utilizador permitir (timeout curto). */
export async function loginGeoHeaders(): Promise<Record<string, string>> {
  try {
    const pos = await getCurrentPosition(5_000);
    return geoHeaders(pos);
  } catch {
    return {};
  }
}

export function mergeHeaders(
  base: Record<string, string>,
  extra: Record<string, string>,
): Record<string, string> {
  return { ...base, ...extra };
}
