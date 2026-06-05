-- Migração 0025 — Ligação Open Banking (FIN-003 scaffold)
--
-- Didático: guardamos apenas METADADOS da ligação (estado, provider, datas).
-- Movimentos bancários reais ficam cifrados no cliente ou chegam via API mTLS
-- sem persistir descrições em claro no servidor (Zero-Knowledge).

CREATE TABLE IF NOT EXISTS bank_connections (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider            TEXT NOT NULL DEFAULT 'mock',
    status              TEXT NOT NULL DEFAULT 'pending',
    consent_expires_at  TIMESTAMPTZ,
    last_sync_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_connections_owner_provider
    ON bank_connections (owner_id, provider);

CREATE INDEX IF NOT EXISTS idx_bank_connections_owner ON bank_connections (owner_id);
