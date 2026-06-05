-- Migração 0018 — Monitorizacao de custos SaaS (FIN-001)
--
-- Didático: cada subscricao SaaS (Netflix, Figma, AWS...) e guardada como um
-- BLOB cifrado com a Master Key, tal como um item de cofre. O servidor ve apenas
-- bytes opacos — o nome do servico, o custo e o ciclo sao Zero-Knowledge.
--
--   { name, cost, currency, cycle, category, vault_item_id?, last_used_at?,
--     active } ──AES-GCM(master_key)──▶ blob
--
-- O dashboard de custos (mensal/anual) e os alertas de licencas sem uso (FIN-002)
-- sao calculados NO CLIENTE, depois de decifrar. O campo vault_item_id liga a
-- subscricao ao login correspondente no cofre ("cruza com vault").

CREATE TABLE IF NOT EXISTS saas_subscriptions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- blob = AES-GCM(master_key, JSON da subscricao) nonce||ct. Opaco ao servidor.
    blob        BYTEA NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_saas_subscriptions_owner ON saas_subscriptions (owner_id);
