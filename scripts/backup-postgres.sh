#!/usr/bin/env bash
# INFRA-004 — Backup PostgreSQL cifrado (Linux / VPS)
# Uso: AEGIS_BACKUP_KEY=<base64> ./scripts/backup-postgres.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$ROOT/backups}"
COMPOSE="docker compose -f $ROOT/docker-compose.yml -f $ROOT/docker-compose.prod.yml"

if [[ -z "${AEGIS_BACKUP_KEY:-}" ]]; then
  echo "ERRO: define AEGIS_BACKUP_KEY (go run ./cmd/backup gen-key)" >&2
  exit 1
fi

mkdir -p "$BACKUP_DIR"
STAMP=$(date -u +%Y%m%dT%H%M%SZ)
RAW="$BACKUP_DIR/pg-$STAMP.sql"
ENC="$RAW.enc"

$COMPOSE exec -T db pg_dump -U "${POSTGRES_USER:-aegis}" "${POSTGRES_DB:-aegis}" > "$RAW"

cd "$ROOT/backend"
go run ./cmd/backup encrypt --in "$RAW" --out "$ENC" --key "$AEGIS_BACKUP_KEY"
rm -f "$RAW"

echo "OK: $ENC"
