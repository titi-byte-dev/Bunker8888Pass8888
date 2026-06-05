-- Migração 0027 — Comissoes sobre faturas pagas (FIN-007)
--
-- Didatico: a comissao fecha o fluxo ERP. Quando uma fatura passa a "paid"
-- (FIN-006), o vendedor que fechou o deal (CRM-003) tem direito a uma comissao.
--
-- Tal como na faturacao, separamos DUAS naturezas:
--
--   1) CONTEUDO (beneficiario, percentagem, valor, moeda, notas) -> PRIVADO.
--      Vai cifrado num blob AES-GCM com a Master Key. O servidor nunca soma nem
--      decifra nada (ZK). Quem fechou o negocio e quanto recebe e sensivel.
--
--   2) LIGACAO + ESTADO -> em claro. A comissao aponta para a fatura de origem
--      (invoice_id) e tem um estado de pagamento ao beneficiario:
--
--        pending  -> ainda nao paga ao vendedor
--        paid     -> ja liquidada
--        void     -> anulada (ex.: a fatura foi anulada)
--
-- Diferenca face a invoices: a comissao NAO e um documento fiscal, logo NAO
-- precisa de numeracao legal sequencial. E um registo interno de acrescimo.
--
-- invoice_id e nullable: se um dia a fatura for removida em cascata (o dono
-- apaga a conta), a comissao perde a referencia mas o historico cifrado fica.

CREATE TABLE IF NOT EXISTS commissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Fatura de origem. SET NULL preserva o registo se a fatura desaparecer.
    invoice_id  UUID NULL REFERENCES invoices(id) ON DELETE SET NULL,

    -- Estado de liquidacao ao beneficiario.
    status      TEXT NOT NULL DEFAULT 'pending'
                     CHECK (status IN ('pending', 'paid', 'void')),

    -- Conteudo cifrado (beneficiario, percentagem, valor...). Opaco ao servidor.
    blob        BYTEA NOT NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_commissions_owner ON commissions (owner_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_commissions_invoice ON commissions (invoice_id);
