-- Migração 0015 — Logs imutaveis com hashing encadeado (HR-002)
--
-- Didático: um registo de auditoria so vale se nao puder ser adulterado a
-- posteriori. Em vez de confiar nas permissoes da BD, encadeamos cada entrada
-- ao hash da anterior (estilo "blockchain" minimalista). Alterar UMA entrada
-- antiga muda o seu entry_hash e parte a cadeia de todas as seguintes — a
-- adulteracao fica matematicamente visivel.
--
--   entrada N-1            entrada N               entrada N+1
--   ┌──────────┐  prev_hash ┌──────────┐  prev_hash ┌──────────┐
--   │entry_hash│───────────▶│entry_hash│───────────▶│entry_hash│
--   └──────────┘            └────┬─────┘            └──────────┘
--                                │ entry_hash = sha256(
--                                │   v1|seq|owner|action|detail|
--                                │   occurred_at|prev_hash )
--
-- A primeira entrada de cada dono encadeia ao prev_hash = "GENESIS".
-- occurred_at e guardado como string RFC3339Nano — a MESMA usada no hash —
-- para o cliente recalcular a cadeia byte-a-byte.
--
-- Convencao: a tabela e APPEND-ONLY. O codigo nunca faz UPDATE nem DELETE de
-- entradas; a unica operacao e inserir no fim da cadeia do dono.

CREATE TABLE IF NOT EXISTS audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Dono da cadeia (cada utilizador tem a sua sequencia independente).
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Posicao na cadeia do dono (1, 2, 3, ...). UNIQUE evita buracos/duplicados.
    seq         BIGINT NOT NULL,

    -- O que aconteceu (ex.: "field.shred") e um detalhe NAO-secreto (ex.: nome
    -- do campo ou id da ficha). Nunca contem o conteudo cifrado.
    action      TEXT NOT NULL,
    detail      TEXT NOT NULL,

    -- Momento (string RFC3339Nano, a mesma usada no calculo do entry_hash).
    occurred_at TEXT NOT NULL,

    -- Encadeamento: hash da entrada anterior ("GENESIS" na primeira).
    prev_hash   TEXT NOT NULL,

    -- sha256(canonico) desta entrada — o elo que a proxima vai referenciar.
    entry_hash  TEXT NOT NULL,

    UNIQUE (owner_id, seq)
);

CREATE INDEX IF NOT EXISTS idx_audit_owner_seq ON audit_log (owner_id, seq);
