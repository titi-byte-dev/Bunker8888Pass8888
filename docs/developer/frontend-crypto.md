---
title: Criptografia no frontend
slug: developer-crypto
category: developer
order: 2
audience: [developer]
layer: [frontend]
feature: vault
level: 3
in_app: true
summary: Argon2id, AES-GCM, Master Key e fluxos no browser.
related: [glossary, vault, developer-api]
---

:::summary
Toda a cifragem sensível corre no cliente via **WebCrypto** e **hash-wasm** (Argon2id).
O bundle nunca inclui segredos — só algoritmos e parâmetros públicos.
:::

:::concept{id="argon2id-client" title="Argon2id no browser" level=2}
`hash-wasm` deriva a Master Key e o Auth Hash a partir da Master Password + salt.
Parâmetros de custo elevados mitigam brute-force offline sobre o Auth Hash armazenado.
:::

:::level{level=3 title="Ficheiros-chave"}
| Ficheiro | Papel |
|---|---|
| `frontend/src/lib/crypto.ts` | AES-GCM, bytes, helpers |
| `frontend/src/lib/vault/password.ts` | Derivação Argon2id |
| `frontend/src/lib/vault/items.ts` | Serialização + blob |
| `frontend/src/lib/vault/masterKeyStore.ts` | Chave em memória volátil |
| `frontend/src/lib/passkey.ts` | WebAuthn (sessão HTTP) |

> ⚠️ **Segurança:** `CryptoKey` com `extractable: false` onde possível; nunca
> `localStorage` para Master Key.
:::

:::level{level=3 title="Fluxo de um item"}
1. Utilizador preenche formulário (login/note/card)
2. JSON serializado → UTF-8 bytes
3. `encrypt()` gera nonce único + ciphertext AES-GCM-256
4. Blob base64 enviado na API; servidor guarda sem decifrar
5. Leitura inversa com Master Key em memória após unlock
:::
