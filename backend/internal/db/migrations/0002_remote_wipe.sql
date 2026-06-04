-- Migração 0002 — auditoria de remote wipe (VAULT-012)
-- Registo imutável (append-only) até HR-002 implementar hashing encadeado.

CREATE TABLE IF NOT EXISTS remote_wipe_events (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    initiated_by    TEXT NOT NULL,
    reason          TEXT NOT NULL DEFAULT '',
    devices_notified INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_remote_wipe_target ON remote_wipe_events (target_user_id);
