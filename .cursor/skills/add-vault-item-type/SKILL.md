---
name: add-vault-item-type
description: Adiciona um novo tipo de item ao cofre (Vault) do AegisPass mantendo o modelo Zero-Knowledge — cifragem client-side, servidor só guarda blobs. Usar ao criar tipos como login, nota, cartão, chave SSH, etc.
disable-model-invocation: true
---

# Adicionar tipo de item ao Vault

Garante que qualquer novo tipo de item respeita o Zero-Knowledge: o servidor
**nunca** vê o conteúdo em claro.

## Princípio

O servidor só conhece **metadados não sensíveis** (id, tipo, tenant, datas) e um
**blob cifrado** com o conteúdo. Toda a estrutura sensível é cifrada no cliente.

## Passos

1. **Cliente (Svelte/TS):** define o schema do novo tipo (ex: `SSHKey { host,
   user, privateKey }`). Serializa → cifra com AES-GCM-256 (WebCrypto) → envia
   só o ciphertext.
2. **API (Go):** aceita e guarda o blob; valida só metadados. Nunca tenta
   decifrar.
3. **Store/SQL:** reutiliza a tabela `vault_items` (coluna `type` + `ciphertext`);
   não criar colunas para campos sensíveis.
4. **Decifragem:** acontece no cliente ao abrir o item.
5. **Testes:** round-trip cifrar→decifrar + verificação de que o servidor não
   regista o conteúdo (ver `06-testing/rgpd-compliance-tests.md`, CT-RGPD-01).

## Checklist

- [ ] Conteúdo sensível cifrado **antes** de sair do cliente.
- [ ] Servidor guarda só blob + metadados não sensíveis.
- [ ] Nonce único; chave nunca enviada.
- [ ] Sem campos sensíveis em colunas/índices/logs em claro.
- [ ] Teste de round-trip + caso de borda.

> ⚠️ Se precisares de pesquisar por um campo sensível, usar *blind indexing*
> (índice sobre um HMAC determinístico do valor), nunca guardar o valor em claro.
