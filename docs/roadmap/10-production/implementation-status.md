# Estado de implementação — dev vs produção

> Matriz viva: o que já corre em `docker compose` local vs o que exige VPS/credenciais.

## Legenda

| Símbolo | Significado |
|---|---|
| 🟢 | Funcional em dev (e/ou prod-ready) |
| 🟡 | Scaffold / guias prontos; falta executar ou credenciais |
| ⚪ | Por implementar |

## Infraestrutura

| Task | Dev | Produção | Notas |
|---|---|---|---|
| INFRA-003/004 | 🟢 | 🟡 | `docker-compose.prod.yml`, `backup-postgres.sh` |
| INFRA-001 | — | 🟡 | Scripts em `infra/scripts/` |
| INFRA-002 | — | 🟡 | WireGuard + UFW documentados |

## E-mail

| Task | Dev | Produção |
|---|---|---|
| MAIL-001/004/005 | 🟢 | 🟢 |
| MAIL-002 | 🟢 Mailpit | 🟡 Postfix + `POST /api/mail/ingest` |
| MAIL-003 | — | 🟡 DNS SPF/DKIM/DMARC |

## Finanças

| Task | Dev | Produção |
|---|---|---|
| FIN-001/002/005/006/007 | 🟢 | 🟢 (ZK no cliente) |
| FIN-003 | 🟢 mock | 🟡 mTLS TPP |
| FIN-004 | ⚪ | ⚪ cartões efémeros |

## Google

| Task | Dev | Produção |
|---|---|---|
| GOOGLE-001 | 🟢 mock + `/work/google` | 🟡 SA JSON na VPS |
| GOOGLE-002/003 | 🟢 `/work/google-dev` + API blobs | 🟡 API Google Drive real |
| GOOGLE-004 | ⚪ | ⚪ Gmail relay |

## Agentes & ERP

| Task | Dev |
|---|---|
| AGENT-003–010 | 🟢 |
| Orquestrador ERP | 🟢 journey `erp-flow-dev` |
| AGENT-006 (banco real) | 🟡 depende FIN-003 prod |

## UI

| Task | Estado |
|---|---|
| UI-001–008, UI-010 | 🟢 |
| UI-009 Capacitor mobile | 🟡 scaffold + guia |

## Próximos passos (prioridade)

1. Executar checklist [`journey-vps-deploy.md`](../04-user-journeys/journey-vps-deploy.md)
2. GOOGLE-002 — upload Drive com blobs ZK via API
3. UI-009 — shell Capacitor + biometria
