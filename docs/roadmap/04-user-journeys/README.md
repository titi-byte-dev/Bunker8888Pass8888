# User Journeys — índice

> Percursos passo-a-passo para validar funcionalidades. Expostos na app em
> **Definições → Documentação** (`npm run docs:build` gera JSON em `frontend/src/lib/docs/generated/`).

## Por módulo

| Módulo | Journey | Rota na app |
|---|---|---|
| **Onboarding** | [admin-onboarding](journey-admin-onboarding.md) | `/admin` |
| | [employee-byod](journey-employee-byod.md) | `/vault` |
| **Vault** | [passkey](journey-passkey.md) | `/security/devices` |
| | [remote-wipe](journey-remote-wipe.md) | `/security` |
| | [emergency-access](journey-emergency-access.md) | `/security/emergency` |
| **RH / RGPD** | [rgpd-erasure](journey-rgpd-erasure.md) | `/hr` |
| | [hr-agent-recruitment](journey-hr-agent-recruitment.md) | `/hr/recruitment` |
| | [hr-agent-onboarding](journey-hr-agent-onboarding.md) | `/hr` |
| **Partilha** | [shared-vault](journey-shared-vault.md) | `/team/vaults` |
| | [secret-link](journey-secret-link.md) | `/team/links` |
| **Segurança** | [sentinel](journey-sentinel.md) | `/security/sentinel` |
| | [guardian-audit](journey-guardian-audit.md) | `/security/guardian` |
| **Mail** | [mail-alias-relay](journey-mail-alias-relay.md) | `/mail` |
| **CRM** | [crm-prospection](journey-crm-prospection.md) | `/crm` |
| **Finanças** | [saas-costs](journey-saas-costs.md) | `/fin` |
| | [fiscal-categorization](journey-fiscal-categorization.md) | `/fin/fiscal` |
| | [finance-agent-saas](journey-finance-agent-saas.md) | `/fin` |
| | [finance-agent-reconcile](journey-finance-agent-reconcile.md) | `/fin/banking` |
| | [erp-flow-dev](journey-erp-flow-dev.md) | `/fin/invoices` |
| **Google** | [google-dev-stub](journey-google-dev-stub.md) | `/work/google-dev` |
| **Agentes** | [orchestrator](journey-orchestrator.md) | feed em `/crm`, `/fin` |
| | [human-in-the-loop](journey-human-in-the-loop.md) | aprovações no feed |
| | [agent-event-feed](journey-agent-event-feed.md) | `/crm` |
| | [ops-agent-inventory](journey-ops-agent-inventory.md) | `/work/inventory` |
| **Infra** | [vps-deploy](journey-vps-deploy.md) | — (produção) |

## Dev vs produção

| Ambiente | Guia |
|---|---|
| Local sem VPS | [`08-dev-environment/development-without-vps.md`](../08-dev-environment/development-without-vps.md) |
| VPS real | [`10-production/README.md`](../10-production/README.md) |
