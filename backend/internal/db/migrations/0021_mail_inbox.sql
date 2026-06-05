-- Migração 0021 — Caixa de entrada simulada (AGENT-003 / MAIL-002 stub)
--
-- Didático: até o servidor SMTP real (MAIL-002), guardamos mensagens recebidas
-- como metadados + corpo em texto. É uma excepção consciente ao Zero-Knowledge
-- (como o destino dos aliases em 0017) — o agente de prospeção lê só o
-- necessário via tool com permissão mail:read_metadata.
--
-- processed_at NULL = ainda não convertida em lead pelo fluxo de prospeção.

CREATE TABLE IF NOT EXISTS mail_inbox_messages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_email   TEXT NOT NULL,
    subject      TEXT NOT NULL,
    body         TEXT NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_mail_inbox_owner ON mail_inbox_messages (owner_id);
CREATE INDEX IF NOT EXISTS idx_mail_inbox_unprocessed
    ON mail_inbox_messages (owner_id) WHERE processed_at IS NULL;
