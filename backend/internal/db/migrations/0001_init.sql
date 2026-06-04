-- ============================================================================
-- Migração 0001 — esquema inicial (VAULT-003)
-- ----------------------------------------------------------------------------
-- Nota: gen_random_uuid() é nativo no PostgreSQL 13+. As colunas BYTEA guardam
-- bytes binários (chaves, hashes, blobs cifrados) sem codificação de texto.
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,

    -- verifier = Argon2id(authHash_do_cliente, verifier_salt). O servidor NUNCA
    -- guarda o authHash cru: guarda um hash dele (defesa em profundidade).
    verifier      BYTEA NOT NULL,
    verifier_salt BYTEA NOT NULL,

    -- Parâmetros/salt da KDF do CLIENTE. São necessários para o cliente voltar a
    -- derivar a Master Key noutro dispositivo. Não são segredos.
    kdf_salt      BYTEA NOT NULL,
    kdf_time      INT   NOT NULL,
    kdf_memory    INT   NOT NULL,
    kdf_threads   INT   NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
    -- Guardamos o SHA-256 do token, nunca o token em claro. Se a BD vazar, os
    -- tokens não são utilizáveis (não se consegue inverter o hash).
    token_hash  BYTEA PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions (user_id);

CREATE TABLE IF NOT EXISTS vault_items (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_type  TEXT NOT NULL,

    -- blob = nonce || ciphertext || tag, cifrado NO CLIENTE (Zero-Knowledge).
    -- O servidor armazena bytes opacos: nunca os decifra nem os interpreta.
    blob       BYTEA NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_vault_items_user ON vault_items (user_id);
