/**
 * Geofencing no cliente (VAULT-011).
 *
 * Didático: o browser obtém GPS via Geolocation API (requer permissão do
 * utilizador). Enviamos lat/lon nos headers X-Geo-* para o servidor validar
 * contra a zona permitida (círculo em torno do escritório).
 *
 * ⚠️ Segurança: GPS no browser é spoofável; combinamos com IP/CIDR no servidor
 * (defesa em profundidade). Nunca confiar só no cliente.
 */
export interface GeoPosition {
  lat: number;
  lon: number;
  accuracyM?: number;
}

export interface GeofencePolicy {
  enabled: boolean;
  allowed_cidrs: string[];
  gps_enabled: boolean;
  gps_lat?: number | null;
  gps_lon?: number | null;
  gps_radius_m: number;
}

/** Obtém posição actual (timeout 10s). Rejeita se permissão negada. */
export function getCurrentPosition(timeoutMs = 10_000): Promise<GeoPosition> {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error("Geolocation não suportada"));
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) =>
        resolve({
          lat: pos.coords.latitude,
          lon: pos.coords.longitude,
          accuracyM: pos.coords.accuracy,
        }),
      (err) => reject(new Error(err.message)),
      { enableHighAccuracy: true, timeout: timeoutMs, maximumAge: 60_000 },
    );
  });
}

/** Headers para anexar a pedidos HTTP autenticados. */
export function geoHeaders(pos: GeoPosition | null): Record<string, string> {
  if (!pos) return {};
  return {
    "X-Geo-Latitude": String(pos.lat),
    "X-Geo-Longitude": String(pos.lon),
  };
}

/**
 * Distância Haversine em metros (espelha lógica do servidor para UI/preview).
 */
export function distanceMeters(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const rad = Math.PI / 180;
  const dLat = (lat2 - lat1) * rad;
  const dLon = (lon2 - lon1) * rad;
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(lat1 * rad) * Math.cos(lat2 * rad) * Math.sin(dLon / 2) ** 2;
  return 6_371_000 * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

/** Preview local: está dentro do círculo GPS configurado? */
export function isWithinGPSFence(pos: GeoPosition, policy: GeofencePolicy): boolean {
  if (!policy.gps_enabled || policy.gps_lat == null || policy.gps_lon == null) return true;
  return distanceMeters(pos.lat, pos.lon, policy.gps_lat, policy.gps_lon) <= policy.gps_radius_m;
}
