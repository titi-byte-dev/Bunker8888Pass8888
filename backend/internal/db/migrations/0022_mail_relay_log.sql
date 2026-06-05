-- Migração 0022 — Auditoria e rate limiting de mail (MAIL-005)
--
-- Didático: cada ingest, relay ou compose fica registado para contar eventos
-- na última hora e impedir abuso (anti open-relay complementar).

CREATE TABLE IF NOT EXISTS mail_relay_log (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alias_id   UUID REFERENCES email_aliases(id) ON DELETE SET NULL,
    direction  TEXT NOT NULL CHECK (direction IN ('inbound', 'outbound_relay', 'compose')),
    from_email TEXT NOT NULL,
    to_email   TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_mail_relay_log_owner_time
    ON mail_relay_log (owner_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_mail_relay_log_alias_time
    ON mail_relay_log (alias_id, created_at DESC) WHERE alias_id IS NOT NULL;
