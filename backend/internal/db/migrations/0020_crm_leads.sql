-- Migração 0020 — Leads CRM cifrados (CRM-001)
--
-- Didático: cada lead é um blob AES-GCM com a Master Key (nome, email, estágio
-- do funil, notas...). O servidor não decifra — só armazena bytes opacos.

CREATE TABLE IF NOT EXISTS crm_leads (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blob        BYTEA NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_crm_leads_owner ON crm_leads (owner_id);
