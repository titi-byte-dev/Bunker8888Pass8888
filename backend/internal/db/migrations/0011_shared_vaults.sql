-- Migração 0011 — Cofres Partilhados (SHARE-002)
--
-- Didático: um cofre partilhado é uma coleção de itens cifrados sob UMA chave
-- simétrica própria (a "chave do cofre", AES-256). Cada membro recebe acesso
-- tendo essa chave do cofre re-cifrada para a sua CHAVE PÚBLICA (RSA-OAEP, da
-- SHARE-001). Assim:
--
--   • O servidor nunca vê a chave do cofre em claro — só guarda, por membro,
--     uma cópia cifrada (`wrapped_vault_key`) que apenas a chave privada desse
--     membro consegue abrir.
--   • Revogar = apagar a linha de membro: o servidor deixa de entregar a cópia
--     cifrada da chave do cofre a esse utilizador (a RLS por membro impede o
--     acesso). Rotação da chave para sigilo perfeito fica como melhoria futura.
--
--                 cria cofre              convida colega
--   ┌────────┐  gera chave K   ┌────────┐  procura PK_colega   ┌────────┐
--   │ Dona   │ ─────────────▶  │ K cifr.│ ───────────────────▶ │ Colega │
--   │ (owner)│  K cifrada p/   │ p/ PK  │   K cifrada p/ PK    │(member)│
--   └────────┘  PK_própria     └────────┘   do colega          └────────┘
--      cada membro guarda a SUA cópia de K cifrada para a sua chave pública.

-- Cofre: dono + nome cifrado com a própria chave do cofre (Zero-Knowledge).
CREATE TABLE IF NOT EXISTS shared_vaults (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- name_blob = AES-GCM(chave_do_cofre, nome) no formato nonce||ciphertext.
    -- O servidor nunca decifra: o nome do cofre também é Zero-Knowledge.
    name_blob   BYTEA NOT NULL,
    algorithm   TEXT  NOT NULL DEFAULT 'AES-GCM-256+RSA-OAEP-3072',

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_shared_vaults_owner ON shared_vaults (owner_id);

-- Membros: papel (permissões) + chave do cofre cifrada para a chave pública
-- deste membro. É esta linha que dá (ou retira, ao ser apagada) o acesso.
CREATE TABLE IF NOT EXISTS shared_vault_members (
    vault_id          UUID NOT NULL REFERENCES shared_vaults(id) ON DELETE CASCADE,
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- owner   → controlo total, apaga o cofre, gere todos os membros
    -- admin   → gere membros (exceto o owner), lê/escreve itens
    -- member  → lê/escreve itens
    -- viewer  → só lê itens
    role              TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),

    -- wrapped_vault_key = RSA-OAEP(PK_deste_membro, chave_do_cofre).
    -- Só a chave privada do membro a abre — o servidor encaminha bytes opacos.
    wrapped_vault_key BYTEA NOT NULL,

    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (vault_id, user_id)
);

-- Procura "os meus cofres" filtra por user_id — daí o índice dedicado.
CREATE INDEX IF NOT EXISTS idx_svm_user ON shared_vault_members (user_id);

-- Itens do cofre: cada blob cifrado com a chave do cofre (não com a Master Key
-- de um utilizador), para que qualquer membro com a chave o consiga abrir.
CREATE TABLE IF NOT EXISTS shared_vault_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id    UUID NOT NULL REFERENCES shared_vaults(id) ON DELETE CASCADE,
    item_type   TEXT NOT NULL,

    -- blob = AES-GCM(chave_do_cofre, payload) no formato nonce||ciphertext.
    blob        BYTEA NOT NULL,

    -- Quem criou (auditoria). SET NULL se o utilizador for removido do sistema.
    created_by  UUID REFERENCES users(id) ON DELETE SET NULL,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_svi_vault ON shared_vault_items (vault_id);
