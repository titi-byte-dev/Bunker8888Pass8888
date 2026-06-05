# Desenvolvimento sem VPS real

> Guia para cumprir os **Definition of Done** das Fases 2 e 3 em `docker compose`
> local, **antes** de INFRA-001/002 (VPS + WireGuard) e integrações externas.

## Pré-requisitos

```bash
docker compose up --build
cd frontend && npm run dev
```

| Serviço | URL / porta | Simula |
|---|---|---|
| API Go | `http://localhost:8080` | Produção |
| PostgreSQL | `localhost:5432` | BD multi-tenant |
| Mailpit SMTP | `localhost:1025` | E-mail real (MAIL-002) |
| Mailpit UI | `http://localhost:8025` | Inbox de teste |
| Open Banking | mock provider | FIN-003 / TPP real |
| Google Workspace | `/work/google-dev` | GOOGLE-001–003 |

> ⚠️ **Segurança:** credenciais de dev estão em `.env.example` — nunca usar em produção.

---

## DoD Fase 2 — como validar em dev

| Critério | Onde | Substituto de produção |
|---|---|---|
| Funil + métricas | `/crm` | Painel conversão/ganhos (cliente) |
| Google Drive ZK | `/work/google-dev` → Drive | `localStorage` + AES-GCM |
| Sheets masking | `/work/google-dev` → Sheets | Tokens NIF/IBAN no cliente |
| Agente + e-mail | `/crm` → simular inbox | Mailpit + webhook |
| Secret links | `/team/links` | Já em produção dev (SHARE-003) |

### E-mail (prospeção)

1. Em `/crm`, preenche **Simular e-mail** e submete.
2. Mailpit dispara webhook → evento `mail.inbox.received`.
3. Aprova **Correr prospeção** no feed → leads no funil.

Ver [`journey-crm-prospection.md`](../04-user-journeys/journey-crm-prospection.md).

### Google (stub)

Ver [`journey-google-dev-stub.md`](../04-user-journeys/journey-google-dev-stub.md).

---

## DoD Fase 3 — como validar em dev

| Critério | Onde | Substituto |
|---|---|---|
| PF → FT → RC | `/fin/invoices` | Numeração legal no servidor; blobs ZK |
| Reconciliação | `/fin/banking` | `MockProvider` + `reconcileTransactions()` |
| Fluxo ERP | `/crm` → `/fin/*` | Orquestrador + HITL |
| Guardião | `/security/guardian` | `GET /api/agent/audit` |
| RGPD PDF | `/hr/compliance` | Print-to-PDF após sugestão do orquestrador |

### Fluxo ERP ponta-a-ponta (dev)

Ver [`journey-erp-flow-dev.md`](../04-user-journeys/journey-erp-flow-dev.md).

Resumo:

1. `/crm` — move lead para **Ganho** → `crm.deal_closed` → aprovar **pro-forma**.
2. `/fin/invoices` — converter PF→FT → **Marcar pago** → aprovar **comissão**.
3. `/fin/commissions` — comissão criada → aprovar **relatório RGPD**.
4. `/hr/compliance` — **Descarregar PDF**.

Opcional: `/fin/banking` — sync mock → reconciliar débitos com subscrições SaaS.

### Guardião

Ver [`journey-guardian-audit.md`](../04-user-journeys/journey-guardian-audit.md).

---

## O que continua a exigir produção

| Task | Porquê | Guia |
|---|---|---|
| INFRA-001/002 | VPS, WireGuard, firewall | [`10-production/`](../10-production/README.md) |
| MAIL-002/003 | Postfix + SPF/DKIM/DMARC | [`mail-002-003-production.md`](../10-production/mail-002-003-production.md) |
| FIN-003 (real) | Certificados mTLS + TPP/banco | [`fin-003-open-banking.md`](../10-production/fin-003-open-banking.md) |
| GOOGLE-001–004 | OAuth2 Service Account Google | [`google-001-oauth.md`](../10-production/google-001-oauth.md) |

Quando a VPS estiver pronta, trocar stubs por variáveis `AEGIS_OB_*` e OAuth — a arquitetura
já está preparada. Checklist: [`journey-vps-deploy.md`](../04-user-journeys/journey-vps-deploy.md).
