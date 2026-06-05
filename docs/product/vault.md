---
title: O Cofre (Vault)
slug: vault
category: product
order: 2
audience: [user, developer]
layer: [frontend, backend]
feature: vault
level: 1
in_app: true
summary: Guardar e sincronizar credenciais com cifragem Zero-Knowledge.
related: [glossary, security, developer-crypto, journey-employee-byod]
---

:::summary
O cofre guarda **logins**, **notas** e **cartões** como blobs cifrados. O servidor
armazena metadados (título, tipo, datas) e o blob — nunca o conteúdo em claro.
:::

:::concept{id="vault-item" title="Item do cofre" level=1}
Cada item tem um **tipo** (`login`, `note`, `card`) e um **blob** AES-GCM-256 gerado no
cliente. A lista no ecrã mostra títulos; o conteúdo só aparece após desbloqueio.
:::

:::concept{id="vault-sync" title="Sincronização em tempo real" level=2}
Alterações propagam-se via **WebSockets** entre dispositivos do mesmo utilizador.
O servidor reencaminha eventos; não decifra payloads.
:::

:::level{level=1 title="O que podes fazer"}
- Criar, editar e apagar itens em `/vault`
- Gerar palavras-passe fortes no formulário
- Ver códigos **TOTP** inline (2FA)
- Importar credenciais de CSV
- Pesquisar localmente (títulos já decifrados em memória)
:::

:::level{level=2 title="Aprofundar: higiene de passwords"}
O **score de higiene** (fraca, reutilizada, exposta) calcula-se no cliente em
`frontend/src/lib/vault/hygiene.ts`. O servidor pode receber apenas o número agregado
para relatórios — nunca a password em claro.

Abre **Segurança → Saúde** para acções recomendadas.
:::

:::level{level=3 title="Técnico: API e módulos"}
| Camada | Caminho | Responsabilidade |
|---|---|---|
| UI | `frontend/src/routes/(app)/vault/` | Listagem, CRUD, unlock gate |
| Cliente API | `frontend/src/lib/vault/api.ts` | REST + headers geofence |
| Cifragem | `frontend/src/lib/vault/items.ts` | Serializar → cifrar blob |
| Backend | `backend/internal/vault/` | CRUD blobs, RLS, WebSocket |
| Persistência | migrações `vault_items` | `tenant_id`, `user_id`, blob |

Endpoints principais: `GET/POST /api/vault`, `GET/PUT/DELETE /api/vault/:id`.
:::
