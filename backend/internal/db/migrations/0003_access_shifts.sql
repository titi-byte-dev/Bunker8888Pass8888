-- Migração 0003 — horários de acesso por turnos (VAULT-010)
--
-- Didático: o turno é uma política por utilizador. Quando enabled=false, não
-- há restrição (compatível com contas pessoais / dev). O schedule é JSONB
-- com janelas por dia da semana no fuso horário indicado.

CREATE TABLE IF NOT EXISTS access_shifts (
    user_id               UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    timezone              TEXT NOT NULL DEFAULT 'UTC',
    schedule              JSONB NOT NULL DEFAULT '{}',
    enabled               BOOLEAN NOT NULL DEFAULT false,
    max_clock_skew_seconds INT NOT NULL DEFAULT 300,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);
