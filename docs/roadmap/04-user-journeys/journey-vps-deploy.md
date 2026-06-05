# Journey — Deploy VPS (INFRA-001 → produção)

> **Ator:** titi (admin) · **Pré-requisito:** VPS nova + domínio para e-mail

## Checklist

| # | Passo | Doc | Verificação |
|---|---|---|---|
| 1 | Gerar chave SSH Ed25519 | [infra-001](../10-production/infra-001-vps-hardening.md) | `ssh titi@IP` funciona |
| 2 | `01-ssh-hardening.sh` | infra-001 | Password SSH desactivada |
| 3 | `02-wireguard-server.sh` | [infra-002](../10-production/infra-002-wireguard-firewall.md) | `ping 10.8.0.1` |
| 4 | Importar `titi.conf` no portátil | infra-002 | Túnel activo |
| 5 | `03-ufw-firewall.sh` | infra-002 | `nmap` só mostra 51820/25 |
| 6 | `.env` + compose prod | [production-deploy](../10-production/production-deploy.md) | `/healthz` OK |
| 7 | Postfix + pipe ingest | [mail-002](../10-production/mail-002-003-production.md) | E-mail na inbox |
| 8 | SPF/DKIM/DMARC | mail-002 | mail-tester.com ≥ 8/10 |
| 9 | `AEGIS_OB_*` (opcional) | [fin-003](../10-production/fin-003-open-banking.md) | provider ≠ mock |
| 10 | Google SA (opcional) | [google-001](../10-production/google-001-oauth.md) | Drive API OK |

## Rollback

- Consola do provider (KVM) se SSH falhar após hardening
- `sudo ufw disable` só em emergência (restaura regras depois)
- Dev local continua disponível: [`development-without-vps.md`](../08-dev-environment/development-without-vps.md)
