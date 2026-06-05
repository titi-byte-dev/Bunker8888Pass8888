-- Migração 0019 — Auditoria de execuções de tools (AGENT-002)
--
-- Didático: cada chamada a uma tool fica registada sem o conteúdo sensível —
-- só quem, qual agente, qual tool e se correu bem. Zero-Trust auditável.

CREATE TABLE IF NOT EXISTS agent_tool_audit (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    agent_id    TEXT NOT NULL,
    tool_name   TEXT NOT NULL,
    success     BOOLEAN NOT NULL,
    error_msg   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_agent_tool_audit_user ON agent_tool_audit (user_id, created_at DESC);
