-- Migração 0009 — Sentinel Mode (DW-004)
-- Registo de tentativas de login e desafios de step-up (passkey).

CREATE TABLE IF NOT EXISTS login_events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID REFERENCES users(id) ON DELETE CASCADE,
    email            TEXT NOT NULL DEFAULT '',
    client_ip        TEXT NOT NULL DEFAULT '',
    geo_lat          DOUBLE PRECISION,
    geo_lon          DOUBLE PRECISION,
    success          BOOLEAN NOT NULL DEFAULT false,
    suspicious       BOOLEAN NOT NULL DEFAULT false,
    reason           TEXT NOT NULL DEFAULT '',
    step_up_required BOOLEAN NOT NULL DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_login_events_user_time ON login_events (user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS sentinel_challenges (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason      TEXT NOT NULL,
    detail      TEXT NOT NULL DEFAULT '',
    client_ip   TEXT NOT NULL DEFAULT '',
    expires_at  TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sentinel_challenges_user ON sentinel_challenges (user_id, created_at DESC);
