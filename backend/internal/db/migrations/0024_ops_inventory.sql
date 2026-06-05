-- Migração 0024 — Inventário operacional (AGENT-008)
--
-- Didático: quantidades e SKUs são metadados operacionais (não PII sensível).
-- O agente de operações reage quando quantity <= reorder_level.

CREATE TABLE IF NOT EXISTS ops_inventory_items (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    sku           TEXT NOT NULL DEFAULT '',
    quantity      INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reorder_level INT NOT NULL DEFAULT 5 CHECK (reorder_level >= 0),
    unit          TEXT NOT NULL DEFAULT 'un',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ops_inventory_owner ON ops_inventory_items (owner_id);
