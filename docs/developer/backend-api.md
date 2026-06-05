---
title: Backend e API
slug: developer-api
category: developer
order: 3
audience: [developer]
layer: [backend]
level: 3
in_app: true
summary: Servidor Go, autenticação, multi-tenancy e módulos internos.
related: [developer-crypto, glossary]
---

:::summary
A API Go valida identidade e contexto (turno, geofence, tenant) mas **não decifra**
conteúdo do cofre. PostgreSQL com RLS isola tenants.
:::

:::concept{id="auth-hash" title="Auth Hash vs Master Key" level=2}
No registo o cliente envia o **Auth Hash** (derivado distinto da Master Key). O login
compara com bcrypt/argon guardado — prova identidade sem revelar material de cifragem.
:::

:::level{level=3 title="Módulos internos"}
| Pacote | Função |
|---|---|
| `internal/auth` | Sessões, JWT, passkeys |
| `internal/vault` | CRUD blobs, WebSocket sync |
| `internal/users` | Perfis por tenant |
| `internal/shifts` | Turnos + NTP |
| `internal/hr` | RH cifrado |
| `internal/sharing` | Cofres partilhados |
| `pkg/crypto` | AEAD, helpers servidor |

Entrypoint: `backend/cmd/server/main.go`.
:::

:::level{level=3 title="Multi-tenancy"}
Cada request autenticado define `app.tenant_id` na sessão SQL (ver
`docs/roadmap/01-architecture/multi-tenancy.md`). Queries sem filtro de tenant são
bug de segurança — escrever testes negativos de isolamento.
:::

:::level{level=3 title="Endpoints de referência"}
- `POST /api/auth/register`, `POST /api/auth/login`
- `GET/POST /api/vault`, `GET/PUT/DELETE /api/vault/:id`
- `GET /healthz` — health check
- WebSocket em `/api/vault/ws` (sync)
:::
