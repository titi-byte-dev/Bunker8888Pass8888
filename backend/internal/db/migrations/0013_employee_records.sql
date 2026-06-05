-- Migração 0013 — Ficha de empregado com cifragem CAMPO-A-CAMPO (HR-001)
--
-- Didático: ao contrário de um item de cofre (UM blob por item), uma ficha de
-- empregado guarda CADA CAMPO cifrado de forma INDEPENDENTE, cada um com a SUA
-- própria chave de campo. Isto é a fundação do crypto-shredding (HR-003): para
-- apagar um único dado (RGPD Art. 17, "direito ao esquecimento" sobre um campo
-- concreto), basta destruir a chave desse campo — o valor cifrado fica
-- permanentemente indecifrável, mesmo que os bytes continuem em disco/backups.
--
-- Modelo de chaves (tudo feito no cliente; o servidor nunca vê nada em claro):
--
--   campo "salario" ──┬─ chave_campo (AES-256 aleatoria, por campo)
--                     │
--                     ├─ value_blob  = AES-GCM(chave_campo, valor)
--                     └─ wrapped_key = AES-GCM(master_key, chave_campo)
--                                                  ▲
--                                                  └── a Master Key do utilizador
--                                                      (VAULT-001/002); nunca sai
--                                                      do dispositivo.
--
--   Para LER  : decifrar wrapped_key c/ master_key → chave_campo → value_blob.
--   Para SHRED: apagar SO o wrapped_key → o value_blob torna-se lixo eterno.
--
-- field_name (ex.: "full_name", "nif", "iban", "salary") e o ESQUEMA, nao o
-- conteudo: e visivel ao servidor tal como item_type nos cofres. So o VALOR e
-- secreto.

CREATE TABLE IF NOT EXISTS employee_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Quem gere a ficha (gestor de RH / dono da organizacao). Se sair, as suas
    -- fichas vao com ele (cascade) — coerente com o modelo Zero-Knowledge: sem
    -- a Master Key dele ninguem decifra os campos de qualquer forma.
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_employee_records_owner ON employee_records (owner_id);

CREATE TABLE IF NOT EXISTS employee_fields (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    record_id   UUID NOT NULL REFERENCES employee_records(id) ON DELETE CASCADE,

    -- Nome do campo = chave de esquema, visivel ao servidor (NAO e segredo).
    field_name  TEXT NOT NULL,

    -- value_blob = AES-GCM(chave_campo, valor) nonce||ct. Mesmo o nome do
    -- empregado e Zero-Knowledge: vive aqui, cifrado.
    value_blob  BYTEA NOT NULL,

    -- wrapped_key = AES-GCM(master_key, chave_campo) nonce||ct. E o ALVO do
    -- crypto-shredding (HR-003): pode passar a NULL para "queimar" o campo sem
    -- mexer no value_blob. shredded_at regista quando isso aconteceu.
    wrapped_key BYTEA,
    shredded_at TIMESTAMPTZ,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Um campo com o mesmo nome so existe uma vez por ficha (upsert por nome).
    UNIQUE (record_id, field_name)
);

CREATE INDEX IF NOT EXISTS idx_employee_fields_record ON employee_fields (record_id);
