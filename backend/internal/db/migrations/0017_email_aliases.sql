-- Migração 0017 — Aliases de e-mail + reencaminhamento (MAIL-001)
--
-- Didático: um ALIAS e um endereco descartavel (ex.: a7f3k29b@aegis.email) que
-- reencaminha para o e-mail REAL do utilizador. Da-se o alias a cada servico;
-- se um deles vazar ou fizer spam, desactiva-se o alias sem mexer na caixa real.
--
--   servico X ──escreve para──▶ a7f3k29b@aegis.email ──reencaminha──▶ destino
--                                        │ (active?)
--                                        └─ inactivo => correio descartado
--
-- NOTA de arquitectura: ao contrario do resto do AegisPass (Zero-Knowledge), o
-- destino TEM de ser visivel ao servidor — e ele que, no futuro (MAIL-002, com
-- SMTP real), faz o relay. Aqui guardamos so a CONFIGURACAO dos aliases; o envio
-- efectivo fica para o servidor de e-mail dedicado.

CREATE TABLE IF NOT EXISTS email_aliases (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Endereco gerado, unico globalmente (e o que se da aos servicos).
    alias_address TEXT NOT NULL UNIQUE,

    -- Para onde reencaminha (e-mail real). Visivel ao servidor (relay futuro).
    destination   TEXT NOT NULL,

    -- Rotulo livre para o utilizador saber a que servico deu este alias.
    label         TEXT NOT NULL DEFAULT '',

    -- Desligar corta o reencaminhamento sem apagar o historico do alias.
    active        BOOLEAN NOT NULL DEFAULT true,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_email_aliases_owner ON email_aliases (owner_id);
