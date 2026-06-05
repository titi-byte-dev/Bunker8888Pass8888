-- Migração 0010 — par de chaves assimétricas por utilizador (SHARE-001)
--
-- Didático: a partilha em Zero-Knowledge faz-se re-cifrando a chave de um item
-- para a CHAVE PÚBLICA do destinatário. Por isso cada utilizador tem um par:
--
--   • public_key          — SPKI DER, em CLARO (qualquer colega a pode pedir
--                            para partilhar segredos comigo).
--   • wrapped_private_key  — PKCS8 cifrado com a Master Key do dono (no cliente,
--                            AES-GCM: nonce||ciphertext). O servidor nunca vê a
--                            chave privada em claro — só o dono a consegue abrir.
--
-- A chave privada decifra o que me foi partilhado; o servidor encaminha sem
-- nunca conseguir ler (modelo Zero-Knowledge).

CREATE TABLE IF NOT EXISTS user_keypairs (
    user_id             UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    public_key          BYTEA NOT NULL,
    wrapped_private_key BYTEA NOT NULL,
    algorithm           TEXT  NOT NULL DEFAULT 'RSA-OAEP-3072-SHA256',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
