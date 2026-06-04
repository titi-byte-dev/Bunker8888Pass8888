import { describe, expect, it, vi, afterEach } from "vitest";
import { distanceMeters, geoHeaders, isWithinGPSFence } from "./geofence";
import type { GeoPosition, GeofencePolicy } from "./geofence";

describe("distanceMeters", () => {
  it("Lisboa–Porto ~275km", () => {
    const d = distanceMeters(38.7223, -9.1393, 41.1579, -8.6291);
    expect(d).toBeGreaterThan(250_000);
    expect(d).toBeLessThan(320_000);
  });
});

describe("geoHeaders", () => {
  it("devolve headers quando há posição", () => {
    const h = geoHeaders({ lat: 38.7, lon: -9.1 });
    expect(h["X-Geo-Latitude"]).toBe("38.7");
    expect(h["X-Geo-Longitude"]).toBe("-9.1");
  });

  it("vazio sem posição", () => {
    expect(geoHeaders(null)).toEqual({});
  });
});

describe("isWithinGPSFence", () => {
  const center: GeoPosition = { lat: 38.7223, lon: -9.1393 };
  const policy: GeofencePolicy = {
    enabled: true,
    allowed_cidrs: [],
    gps_enabled: true,
    gps_lat: center.lat,
    gps_lon: center.lon,
    gps_radius_m: 500,
  };

  it("dentro do raio", () => {
    expect(isWithinGPSFence(center, policy)).toBe(true);
  });

  it("fora do raio", () => {
    expect(isWithinGPSFence({ lat: 41, lon: -8 }, policy)).toBe(false);
  });
});

describe("getCurrentPosition", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("rejeita se geolocation indisponível", async () => {
    vi.stubGlobal("navigator", { geolocation: undefined });
    const { getCurrentPosition } = await import("./geofence");
    await expect(getCurrentPosition()).rejects.toThrow(/não suportada/i);
  });
});
