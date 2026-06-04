-- ============================================================================
-- Migração 0008 — Acesso de emergência / herdeiro digital (VAULT-016)
-- ----------------------------------------------------------------------------
-- Didático: o servidor orquestra o período de espera e aprovações; o blob
-- cifrado (Master Key) é preparado no cliente (Zero-Knowledge).
-- ============================================================================

CREATE TABLE IF NOT EXISTS emergency_configs (
    owner_user_id   UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    heir_email      TEXT NOT NULL,
    wait_days       INT NOT NULL DEFAULT 7 CHECK (wait_days >= 1 AND wait_days <= 90),
    encrypted_blob  BYTEA,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS emergency_requests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    heir_user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'waiting'
                    CHECK (status IN ('waiting', 'rejected', 'ready', 'consumed')),
    requested_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    unlocks_at      TIMESTAMPTZ NOT NULL,
    rejected_at     TIMESTAMPTZ,
    consumed_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_emergency_requests_owner ON emergency_requests (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_emergency_requests_heir ON emergency_requests (heir_user_id);

-- Uma solicitação activa (waiting/ready) por par titular–herdeiro.
CREATE UNIQUE INDEX IF NOT EXISTS idx_emergency_request_active
    ON emergency_requests (owner_user_id, heir_user_id)
    WHERE status IN ('waiting', 'ready');
