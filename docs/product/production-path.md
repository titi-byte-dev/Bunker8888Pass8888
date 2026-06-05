---
title: Caminho para produção
slug: production-path
category: product
order: 7
audience: [admin, developer]
layer: [infra]
feature: production
level: 2
in_app: true
summary: Do docker compose local à VPS com WireGuard, e-mail real e integrações externas.
related: [journey-vps-deploy, fin]
---

:::summary
O desenvolvimento local cumpre os **Definition of Done** das Fases 2 e 3 com
*stubs* (Mailpit, Open Banking mock, Google dev). A produção troca esses stubs
por serviços reais na VPS — sem mudar a arquitetura Zero-Knowledge.
:::

:::concept{id="dev-stub" title="Stub de desenvolvimento" level=1}
Componente que imita o comportamento de um serviço externo (banco, Google, SMTP)
para testar fluxos sem custos nem credenciais de produção.
:::

:::level{level=2 title="Ordem recomendada"}
1. **INFRA-001** — VPS + SSH só chaves (`infra/scripts/01-ssh-hardening.sh`)
2. **INFRA-002** — WireGuard + UFW (`02-`, `03-`)
3. **Deploy** — `docker-compose.prod.yml` + `.env.production.example`
4. **MAIL-002/003** — Postfix + SPF/DKIM/DMARC
5. **FIN-003** — `AEGIS_OB_*` mTLS
6. **GOOGLE-001** — Service Account Workspace
:::

:::level{level=2 title="Checklist"}
Seguir o journey [`journey-vps-deploy`](../roadmap/04-user-journeys/journey-vps-deploy.md)
ou a pasta [`docs/roadmap/10-production/`](../roadmap/10-production/README.md).
:::

:::level{level=3 title="Backups"}
`scripts/backup-postgres.sh` na VPS com `AEGIS_BACKUP_KEY` — ficheiros `.enc`
irrecuperáveis sem a chave offline (INFRA-004).
:::
