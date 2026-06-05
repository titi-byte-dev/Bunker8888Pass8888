# Scripts de infraestrutura — AegisPass

Ordem recomendada na VPS nova:

| # | Script | Task | Notas |
|---|---|---|---|
| 01 | [`01-ssh-hardening.sh`](01-ssh-hardening.sh) | INFRA-001 | `DEPLOY_USER`, `DEPLOY_PUBKEY` |
| 02 | [`02-wireguard-server.sh`](02-wireguard-server.sh) | INFRA-002 | Gera configs cliente |
| 03 | [`03-ufw-firewall.sh`](03-ufw-firewall.sh) | INFRA-002 | `--allow-smtp` quando mail activo |
| 05 | [`05-deploy-stack.sh`](05-deploy-stack.sh) | INFRA-003 | Requer `.env` em `/opt/aegispass` |
| 04 | [`04-postfix-install.sh`](04-postfix-install.sh) | MAIL-002 | Depois do backend healthy |

> 💡 **Ordem:** deploy Docker **antes** do Postfix — o pipe chama `127.0.0.1:8080`.

Dev local — testar ingest sem Postfix:

```powershell
.\scripts\smoke-mail-ingest.ps1 -AliasAddress "SEU_ALIAS@aegis.email"
```

Ver [`docs/roadmap/10-production/demo-client-checklist.md`](../../docs/roadmap/10-production/demo-client-checklist.md).
