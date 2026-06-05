-- Migração 0016 — Gestao de contratos (HR-005) + assinatura digital (HR-006)
--
-- Didático: os contratos sao ficheiros (PDF, .docx...) cifrados FICHEIRO-A-
-- FICHEIRO, cada um com a SUA chave aleatoria, embrulhada com a Master Key —
-- exactamente o padrao da HR-001, mas para bytes de ficheiro. Guardamos o
-- ciphertext na BD (object storage didactico); num sistema real o data_blob
-- iria para um bucket S3/minio com a mesma chave por ficheiro.
--
--   ficheiro ──┬─ AES-GCM(file_key, metadados) ─▶ meta_blob  (listagem)
--              ├─ AES-GCM(file_key, bytes)      ─▶ data_blob  (download)
--              └─ AES-GCM(master_key, file_key) ─▶ wrapped_key
--
-- HR-006 acrescenta ASSINATURA DIGITAL (ECDSA P-256). Assina-se o digest do
-- CIPHERTEXT (sha256(data_blob)) — assim a verificacao nao precisa de decifrar
-- (mantem-se Zero-Knowledge) e prova que o signatario se comprometeu com
-- exactamente aqueles bytes.
--
--   sha256(data_blob) ──ECDSA-sign(priv)──▶ signature
--   verify: ECDSA-verify(pub, signature, sha256(data_blob)) == OK ?
--
-- A identidade de assinatura (par de chaves ECDSA) vive em hr_signing_identities:
-- chave publica em claro (qualquer um verifica) e chave privada cifrada com a
-- Master Key (so o dono assina). O servidor nunca assina nada.

CREATE TABLE IF NOT EXISTS hr_signing_identities (
    owner_id            UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    -- Chave publica ECDSA P-256 em formato SPKI (DER). Publica de proposito.
    public_key          BYTEA NOT NULL,
    -- Chave privada PKCS8 cifrada com a Master Key (nonce||ct). So o dono a abre.
    wrapped_private_key BYTEA NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS employee_contracts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id   UUID NOT NULL REFERENCES employee_records(id) ON DELETE CASCADE,

    meta_blob   BYTEA NOT NULL, -- AES-GCM(file_key, JSON{name,mime,size})
    data_blob   BYTEA NOT NULL, -- AES-GCM(file_key, bytes_do_ficheiro)
    wrapped_key BYTEA NOT NULL, -- AES-GCM(master_key, file_key)

    -- Tamanho do ciphertext (visivel ao servidor para impor limites).
    byte_size   BIGINT NOT NULL,

    -- Assinatura digital (HR-006). NULL ate ser assinado.
    content_digest TEXT,                                  -- sha256(data_blob) hex
    signature      BYTEA,                                 -- ECDSA P-256
    signed_by      UUID REFERENCES users(id) ON DELETE SET NULL,
    signed_at      TIMESTAMPTZ,

    created_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_employee_contracts_record ON employee_contracts (record_id);
