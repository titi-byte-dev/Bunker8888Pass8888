# Checklist — demo a clientes (INFRA + MAIL)

> Ordem prática para sair de `docker compose` local e demonstrar **aliases reais**,
> **orquestrador com e-mail** e journeys `journey-mail-alias-relay` /
> `journey-orchestrator` com tráfego SMTP verdadeiro.

Branch de trabalho sugerida: `feat/ui-doc-infra-mail`.

## Pré-requisitos

- Domínio ou subdomínio dedicado (ex. `mail.suaempresa.pt`)
- VPS Debian/Ubuntu (2 vCPU, 4 GB RAM mínimo)
- Acesso DNS (SPF, DKIM, DMARC)
- Chave SSH gerada localmente

## Fase A — Infra (INFRA-001 / INFRA-002)

| # | Acção | Guia | Verificação |
|---|---|---|---|
| A1 | Criar VPS + utilizador deploy | [`infra-001-vps-hardening.md`](infra-001-vps-hardening.md) | `ssh deploy@vps` só com chave |
| A2 | Executar `infra/scripts/01-ssh-hardening.sh` | idem | Password SSH desactivado |
| A3 | Executar `02-wireguard-server.sh` | [`infra-002-wireguard-firewall.md`](infra-002-wireguard-firewall.md) | Cliente WG liga |
| A4 | Executar `03-ufw-firewall.sh` | idem | Só 22/tcp (ou WG) + UDP WG |
| A5 | Copiar `.env.production.example` → `.env` e preencher segredos | [`production-deploy.md`](production-deploy.md) | Sem campos vazios |
| A6 | Executar `infra/scripts/05-deploy-stack.sh` | idem | `curl /healthz` OK |

> ⚠️ **Segurança:** confirma acesso SSH **antes** de A2. Guarda backup da config WG.

## Fase B — E-mail (MAIL-002 / MAIL-003)

| # | Acção | Guia | Verificação |
|---|---|---|---|
| B1 | Executar `infra/scripts/04-postfix-install.sh` | [`mail-002-003-production.md`](mail-002-003-production.md) | `postfix status` active |
| B2 | (alternativa manual) Copiar `infra/postfix/aegis-ingest.*` | idem | Pipe chama `/api/mail/ingest` |
| B3 | Definir `AEGIS_MAIL_WEBHOOK_SECRET` (prod) | `.env.production.example` | Ingest 401 sem secret |
| B4 | Registar SPF no DNS | mail-002-003 § DNS | `dig TXT dominio` |
| B5 | Configurar DKIM (OpenDKIM ou rspamd) | idem | selector `_domainkey` |
| B6 | Política DMARC `p=quarantine` inicial | idem | Relatórios para admin@ |
| B7 | Criar alias de teste na app | `/mail` | Alias activo no tenant |

## Fase C — Demo funcional (valor de negócio)

| # | Cenário | O que mostrar |
|---|---|---|
| C1 | E-mail externo → alias | Mail chega ao relay/inbox (journey mail-alias) |
| C2 | Orquestrador | Evento `mail.inbox.received` → sugestão em `/crm` |
| C3 | Human-in-the-loop | Aprovar prospeção (AGENT-009) |
| C4 | Documentação | `/settings/docs/journey-orchestrator` com fluxo animado |

## Rollback / dev local

Se a demo falhar, volta a Mailpit:

```bash
docker compose up   # dev — MAIL-001/004 sem Postfix
```

Simular ingest Postfix em dev (alias já criado em `/mail`):

```powershell
.\scripts\smoke-mail-ingest.ps1 -AliasAddress "SEU_ALIAS@aegis.email"
```

## Próximo após MAIL

- **GOOGLE-001** — [`google-001-oauth.md`](google-001-oauth.md)
- **FIN-003** — Open Banking (depende INFRA-002 mTLS)
