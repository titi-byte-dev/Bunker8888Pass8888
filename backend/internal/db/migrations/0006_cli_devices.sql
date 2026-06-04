-- ============================================================================
-- Migração 0006 — dispositivos CLI com certificado mTLS (VAULT-017)
-- ----------------------------------------------------------------------------
-- Cada dispositivo CLI regista um certificado cliente assinado pela CA interna.
-- O servidor guarda apenas a impressão digital (SHA-256) do certificado — nunca
-- a chave privada (fica só no dispositivo do utilizador).
-- ============================================================================

CREATE TABLE IF NOT EXISTS cli_devices (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    cert_fingerprint BYTEA NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at       TIMESTAMPTZ,
    UNIQUE (cert_fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_cli_devices_user ON cli_devices (user_id);
