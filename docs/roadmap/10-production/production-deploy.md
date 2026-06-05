# Deploy em produção (INFRA-003/004)

> **Objetivo:** correr AegisPass na VPS com Docker Compose, segredos fora do git e
> backups PostgreSQL cifrados.

## Ficheiros

| Ficheiro | Papel |
|---|---|
| `docker-compose.prod.yml` | Serviços sem portas públicas desnecessárias |
| `.env.production.example` | Modelo de variáveis — copiar para `.env` na VPS |
| `scripts/backup-postgres.ps1` / `backend/cmd/backup` | Backups `.enc` (INFRA-004) |

## 1. Preparar `.env` na VPS

```bash
cd /opt/aegispass
cp .env.production.example .env
nano .env   # preencher segredos fortes — NUNCA commitar
```

Gera chaves:

```bash
# Backup AES-256
cd backend && go run ./cmd/backup gen-key

# Segredos webhook / admin (exemplo com openssl)
openssl rand -base64 32
```

## 2. Subir stack

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
docker compose ps
curl -s http://127.0.0.1:8080/healthz
```

## 3. Diferenças face ao dev

| Aspecto | Dev (`docker-compose.yml`) | Prod (`docker-compose.prod.yml`) |
|---|---|---|
| PostgreSQL | Porta `5432` no host | Só rede interna Docker |
| Mailpit | SMTP + UI expostos | **Removido** — Postfix no host |
| `AEGIS_MTLS_AUTO_DEV` | `true` | `false` — certificados reais |
| `AEGIS_ADMIN_KEY` | placeholder | valor longo aleatório |
| Backups | opcional | cron diário + chave offline |

## 4. Migrações e backups

```bash
# Migrações correm no arranque do backend (automático)
# Backup manual:
AEGIS_BACKUP_KEY=<base64> ./scripts/backup-postgres.sh   # ver nota Linux abaixo
```

> No Linux, adapta `scripts/backup-postgres.ps1` ou usa `docker exec` + `pg_dump`
> conforme [`INFRA-004`](../05-tasks/backlog.md).

## 5. Actualizações

```bash
git pull
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

Ordem segura: `db` healthy → `backend` → verificar `/healthz`.

## 6. Integrações seguintes

| Variável | Guia |
|---|---|
| Postfix / relay | [`mail-002-003-production.md`](mail-002-003-production.md) |
| `AEGIS_OB_*` | [`fin-003-open-banking.md`](fin-003-open-banking.md) |
| Google OAuth | [`google-001-oauth.md`](google-001-oauth.md) |
