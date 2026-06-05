-- Migração 0012 — Anexos cifrados por ficheiro (SHARE-004)
--
-- Didático: estende os Cofres Partilhados (SHARE-002) com ANEXOS — contratos,
-- chaves .pem, PDFs — cifrados ficheiro-a-ficheiro ANTES do upload. Reutiliza a
-- chave do cofre (AES-256) que cada membro já possui: o servidor continua a ver
-- apenas bytes opacos e a fazer cumprir o acesso por membro.
--
-- Cada anexo guarda DOIS blobs cifrados, separados de propósito:
--   • meta_blob → nome/tipo/tamanho do ficheiro (pequeno; vem na listagem);
--   • data_blob → os bytes do ficheiro (grande; só desce ao descarregar).
-- Separar evita arrastar megabytes só para mostrar a lista de anexos.
--
--   ficheiro.pem ──┬─ AES-GCM(chave_cofre, metadados) ─▶ meta_blob  (listagem)
--                  └─ AES-GCM(chave_cofre, bytes)      ─▶ data_blob  (download)
--                              ▲
--                              └── a mesma chave do cofre da SHARE-002;
--                                  o servidor nunca a vê em claro.
--
-- byte_size é o tamanho do CIPHERTEXT (visível ao servidor) — serve para impor
-- um limite de tamanho sem nunca decifrar nada.

CREATE TABLE IF NOT EXISTS shared_vault_attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id    UUID NOT NULL REFERENCES shared_vaults(id) ON DELETE CASCADE,

    -- meta_blob = AES-GCM(chave_cofre, JSON{name, mime, size}) nonce||ct.
    -- Mesmo o nome do ficheiro é Zero-Knowledge.
    meta_blob   BYTEA NOT NULL,

    -- data_blob = AES-GCM(chave_cofre, bytes_do_ficheiro) nonce||ct.
    data_blob   BYTEA NOT NULL,

    -- Tamanho do ciphertext do ficheiro (bytes). Visível ao servidor só para
    -- impor quota/limites; não revela o conteúdo.
    byte_size   BIGINT NOT NULL,

    -- Quem carregou (auditoria). SET NULL se o utilizador sair do sistema.
    created_by  UUID REFERENCES users(id) ON DELETE SET NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- A listagem de anexos filtra por cofre — daí o índice dedicado.
CREATE INDEX IF NOT EXISTS idx_sva_vault ON shared_vault_attachments (vault_id);
