---
title: Glossário de conceitos-chave
slug: glossary
category: concepts
order: 1
audience: [user, developer]
layer: [product]
level: 1
in_app: true
summary: Termos fundamentais do AegisPass — a base para perceber cofre, segurança e RGPD.
related: [vault, security, developer-crypto]
---

:::summary
O AegisPass assenta em poucos conceitos de segurança e privacidade. Começa por estes
cinco; os restantes expandem o vocabulário à medida que exploras funcionalidades.
:::

:::concept{id="zero-knowledge" title="Zero-Knowledge" level=1}
O servidor **nunca** vê dados sensíveis em claro nem a **Master Key**. Cifragem e
decifragem acontecem no teu dispositivo (*client-side*). Mesmo com acesso à base de
dados, um atacante só encontra blobs ilegíveis.
:::

:::concept{id="master-key" title="Master Key e Master Password" level=1}
A **Master Password** é o que memorizas. A **Master Key** é derivada dela com
**Argon2id** e vive só em memória no browser. O servidor recebe apenas o **Auth Hash**
— uma derivação distinta que prova identidade sem revelar a chave de cifragem.
:::

:::concept{id="tenant" title="Tenant (empresa isolada)" level=1}
Cada empresa-cliente é um **tenant**: os seus dados nunca se misturam com os de outra.
A API e o PostgreSQL aplicam `tenant_id` em cada operação; a **RLS** (*Row-Level Security*)
é rede de segurança na base de dados.
:::

:::concept{id="nonce" title="Nonce" level=2}
Número usado **uma só vez** por operação de cifragem AES-GCM. Reutilizar um nonce com a
mesma chave destrói a confidencialidade do esquema — por isso geramo-lo sempre com
`crypto.getRandomValues` / `crypto/rand`.
:::

:::concept{id="totp" title="TOTP (2FA)" level=2}
Código de 6 dígitos que muda a cada ~30 segundos, gerado a partir de um segredo partilhado
e da hora actual (RFC 6238). No AegisPass o segredo TOTP fica **cifrado no cofre** — o
servidor nunca o vê em claro.
:::

:::level{level=2 title="Aprofundar: pilares Zero-Trust"}
| Pilar | O que significa na prática |
|---|---|
| Rede | Túnel WireGuard antes de falar com a API |
| Identidade | Auth Hash + sessão + passkeys |
| Contexto | Turno (NTP), geofencing, dispositivo registado |
| Dados | Blobs cifrados; score de higiene calculado no cliente |
:::

:::level{level=3 title="Técnico: onde vive cada segredo"}
| Segredo | Onde | Nunca |
|---|---|---|
| Master Key | Memória do browser (volátil) | Servidor, localStorage |
| Auth Hash | PostgreSQL (hash one-way) | Transmitido em claro após registo |
| Blobs do cofre | PostgreSQL (cifrados) | Decifrados no servidor |
| Segredo TOTP | Dentro do blob cifrado | Logs ou API em claro |

Ver implementação: `frontend/src/lib/crypto.ts`, `backend/pkg/crypto/`.
:::
