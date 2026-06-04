-- Migração 0004 — geofencing de acesso (VAULT-011)
--
-- Didático: política por utilizador. IP via CIDRs (rede corporativa/VPN);
-- GPS via círculo (lat/lon + raio em metros). enabled=false = sem restrição.

CREATE TABLE IF NOT EXISTS access_geofence (
    user_id        UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    enabled        BOOLEAN NOT NULL DEFAULT false,
    allowed_cidrs  JSONB NOT NULL DEFAULT '[]',
    gps_enabled    BOOLEAN NOT NULL DEFAULT false,
    gps_lat        DOUBLE PRECISION,
    gps_lon        DOUBLE PRECISION,
    gps_radius_m   DOUBLE PRECISION NOT NULL DEFAULT 500,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
