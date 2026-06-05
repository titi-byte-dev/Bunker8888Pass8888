-- Migração 0023 — Event Bus de agentes (AGENT-004)
--
-- Didático: log append-only de eventos publicados no bus in-process.
-- Permite auditoria, feed de acções e futura orquestração (AGENT-005).
-- Payload é JSON opaco — sem PII em claro de leads CRM.

CREATE TABLE IF NOT EXISTS agent_events (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    source     TEXT NOT NULL DEFAULT '',
    payload    JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_agent_events_user_time
    ON agent_events (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_events_type_time
    ON agent_events (event_type, created_at DESC);
