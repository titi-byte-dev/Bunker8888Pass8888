-- ============================================================================
-- Migração 0007 — credenciais WebAuthn / passkeys (VAULT-014)
-- ----------------------------------------------------------------------------
-- Guardamos a chave pública e metadados — NUNCA a chave privada (fica no
-- dispositivo: Secure Enclave, TPM, etc.).
-- ============================================================================

CREATE TABLE IF NOT EXISTS webauthn_credentials (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id  BYTEA NOT NULL UNIQUE,
    public_key     BYTEA NOT NULL,
    sign_count     BIGINT NOT NULL DEFAULT 0,
    backup_eligible BOOLEAN NOT NULL DEFAULT false,
    backup_state    BOOLEAN NOT NULL DEFAULT false,
    transports     TEXT,
    name           TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_user ON webauthn_credentials (user_id);
