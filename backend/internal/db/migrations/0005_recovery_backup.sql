-- Migração 0005 — backup de recuperação da Master Key (VAULT-018)
--
-- Didático: o blob é a Master Key cifrada com a chave de recuperação (feito no
-- cliente). O servidor nunca vê a chave de recuperação nem a Master Key em claro.

CREATE TABLE IF NOT EXISTS recovery_backups (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    blob        BYTEA NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
