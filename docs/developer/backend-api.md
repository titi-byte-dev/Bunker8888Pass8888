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
related: [developer-crypto, glossary, production-path]
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
| `internal/fin` | Subscrições SaaS (blobs ZK) |
| `internal/openbanking` | PSD2 mock / mTLS |
| `internal/invoicing` | Faturas, numeração legal |
| `internal/commissions` | Comissões ligadas a faturas |
| `internal/googleworkspace` | Provider Google (mock / SA) |
| `internal/mail` | Aliases, inbox, ingest SMTP |
| `internal/orchestrator` | Workers multi-agente |
| `internal/agent` | Tools, prospeção, Guardião |
| `internal/crm` | Leads cifrados |
| `internal/hr` | RH + RGPD |

Entrypoint: `backend/cmd/server/main.go`.
:::

:::level{level=3 title="Multi-tenancy"}
Cada request autenticado define `app.tenant_id` na sessão SQL (ver
`docs/roadmap/01-architecture/multi-tenancy.md`). Queries sem filtro de tenant são
bug de segurança — escrever testes negativos de isolamento.
:::

:::level{level=3 title="Endpoints principais"}
**Auth & vault**
- `POST /api/auth/register`, `POST /api/auth/login`
- `GET/POST /api/vault`, `GET/PUT/DELETE /api/vault/{id}`
- `GET /healthz`

**Finanças**
- `GET/POST/PUT/DELETE /api/fin/subscriptions`
- `GET/POST /api/fin/invoices`, `PUT /api/fin/invoices/{id}/status`
- `GET/POST /api/fin/commissions`, `PUT /api/fin/commissions/{id}/status`
- `GET/POST /api/fin/banking/{status,connect,sync}`

**Mail (webhooks sem Bearer — segredo na query)**
- `POST /api/mail/webhook/mailpit` — dev Mailpit
- `POST /api/mail/ingest` — Postfix produção (JSON pipe)

**Google & agentes**
- `GET /api/work/google/status`
- `GET/POST /api/work/google/drive/files`, `DELETE /api/work/google/drive/files/{id}`
- `GET /api/agent/events`, `POST /api/agent/prospection/run`
- `GET /api/agent/audit`, `GET /api/agent/orchestrator/status`

**CRM**
- `GET/POST/PATCH/DELETE /api/crm/leads`
:::

:::level{level=3 title="Variáveis de ambiente (produção)"}
Ver `.env.production.example` e `docs/roadmap/10-production/`. Destaques:
`AEGIS_OB_*` (Open Banking), `AEGIS_GOOGLE_*` (Workspace),
`AEGIS_MAIL_WEBHOOK_SECRET`, `AEGIS_SMTP_RELAY_HOST`.
:::
