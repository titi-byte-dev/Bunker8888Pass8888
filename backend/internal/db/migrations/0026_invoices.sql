-- Migração 0026 — Faturacao com numeracao legal sequencial (FIN-006)
--
-- Didatico: uma fatura tem DUAS naturezas distintas que esta tabela separa de
-- forma deliberada:
--
--   1) CONTEUDO (cliente, linhas, valores, NIF...) -> e PRIVADO. Vai cifrado
--      num blob AES-GCM com a Master Key. O servidor nunca o ve em claro (ZK).
--
--   2) NUMERO LEGAL (ex.: "FT 2026/0001") -> NAO e secreto, mas TEM de ser
--      sequencial, sem buracos e imutavel por serie/ano (exigencia fiscal). Por
--      isso o numero/seq vivem em COLUNAS em claro, geridas pelo servidor com um
--      advisory lock por (dono, tipo, ano) — tal como a cadeia de auditoria.
--
--   doc_type   prefixo   natureza
--   --------   -------   --------------------------------------------
--   proforma   PF        orcamento (nao fiscal)
--   invoice    FT        fatura (documento fiscal)
--   receipt    RC        recibo (confirma pagamento)
--
--   numero = prefixo + " " + ano + "/" + seq(4 digitos)
--   seq reinicia a 1 por (dono, doc_type, ano).
--
-- Convencao: documentos fiscais NAO se apagam. So muda o status
-- (issued -> paid | void). A unica escrita destrutiva e impossivel por design.

CREATE TABLE IF NOT EXISTS invoices (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Natureza do documento (ver tabela acima).
    doc_type       TEXT NOT NULL CHECK (doc_type IN ('proforma', 'invoice', 'receipt')),

    -- Numeracao legal (em claro, gerida pelo servidor).
    year           INT  NOT NULL,
    seq            BIGINT NOT NULL,
    number         TEXT NOT NULL,

    -- Estado do documento.
    status         TEXT NOT NULL DEFAULT 'issued'
                        CHECK (status IN ('issued', 'paid', 'void')),

    -- Liga ao lead de origem quando nasce de um "deal_closed" (CRM-003).
    -- E apenas um id opaco; nao revela PII.
    source_lead_id UUID NULL,

    -- Conteudo cifrado (cliente, linhas, totais). Opaco para o servidor.
    blob           BYTEA NOT NULL,

    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Garante sequencia sem buracos/duplicados por dono+tipo+ano.
    UNIQUE (owner_id, doc_type, year, seq)
);

CREATE INDEX IF NOT EXISTS idx_invoices_owner ON invoices (owner_id, created_at DESC);
