-- Migração 0014 — Certificados criptográficos de eliminação (HR-003 + HR-004)
--
-- Didático: o crypto-shredding (HR-003) concretiza o "direito ao esquecimento"
-- (RGPD Art. 17) sobre a ficha campo-a-campo da HR-001. Apagar um campo NÃO é
-- reescrever bytes — é DESTRUIR a chave que o abre:
--
--   employee_fields.wrapped_key := NULL    (a chave do campo desaparece)
--   employee_fields.shredded_at := now()   (carimbo de quando ardeu)
--
-- O value_blob (ciphertext) PODE permanecer em disco/backups, mas sem a chave
-- torna-se lixo matematicamente indecifravel. Isto e a essencia do shredding:
-- não confiamos no "apagar" do storage; confiamos na ausencia da chave.
--
-- HR-004 acrescenta a PROVA: cada shred emite um CERTIFICADO verificavel.
--
--   value_blob ──sha256──▶ value_digest  (que ciphertext ficou orfao)
--                                  │
--   canonico = v1|record|field|value_digest|shredded_at|owner
--                                  │ sha256
--                                  ▼
--                              cert_hash   (impressao digital tamper-evident)
--
-- Qualquer pessoa (cliente incluido) recalcula o cert_hash a partir dos campos
-- do certificado e confirma que bate certo — sem nunca ver o conteudo original.
-- O certificado SOBREVIVE a remocao da ficha (record_id ON DELETE SET NULL):
-- a prova de que algo foi eliminado tem de durar mais do que o proprio dado.

CREATE TABLE IF NOT EXISTS erasure_certificates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Dono da prova (gestor de RH). Mantem-se para listagem/auditoria mesmo
    -- depois de a ficha desaparecer.
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Ficha de origem; fica NULL se a ficha for apagada por completo depois.
    record_id   UUID REFERENCES employee_records(id) ON DELETE SET NULL,

    -- Campo eliminado (nome de esquema, ja visivel ao servidor).
    field_name  TEXT NOT NULL,

    -- sha256(value_blob) em hex no momento do shred: identifica o ciphertext
    -- que ficou orfao sem revelar nada do seu conteudo.
    value_digest TEXT NOT NULL,

    -- Momento exacto do shredding (string RFC3339Nano, a MESMA usada no canonico
    -- do cert_hash — garante que cliente e servidor calculam o hash igual).
    shredded_at  TEXT NOT NULL,

    -- Impressao digital tamper-evident: sha256 do canonico (ver acima).
    cert_hash    TEXT NOT NULL,

    issued_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_erasure_certs_owner ON erasure_certificates (owner_id);
CREATE INDEX IF NOT EXISTS idx_erasure_certs_record ON erasure_certificates (record_id);
