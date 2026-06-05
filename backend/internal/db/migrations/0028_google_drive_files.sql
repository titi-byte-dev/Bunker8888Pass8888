-- Migração 0028 — Ficheiros Drive ZK (GOOGLE-002)
--
-- Didático: o cliente cifra nome+conteúdo com a Master Key e envia um blob opaco.
-- O servidor (ou a Google Drive em produção) só vê bytes sem significado.
-- external_id fica NULL até GOOGLE-002 ligar à API Drive real.

CREATE TABLE IF NOT EXISTS google_drive_files (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blob          BYTEA NOT NULL,
    external_id   TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_google_drive_files_owner ON google_drive_files (owner_id);
