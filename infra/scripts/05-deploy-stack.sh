#!/usr/bin/env bash
# INFRA-003 — Sobe stack Docker de produção em /opt/aegispass
# Executar como utilizador deploy (membro do grupo docker).
set -euo pipefail

REPO_DIR="${REPO_DIR:-/opt/aegispass}"
COMPOSE="docker compose -f docker-compose.yml -f docker-compose.prod.yml"

cd "${REPO_DIR}"

if [[ ! -f .env ]]; then
  echo "ERRO: copia .env.production.example → .env e preenche segredos" >&2
  exit 1
fi

# Valida campos críticos vazios
for key in POSTGRES_PASSWORD AEGIS_ADMIN_KEY AEGIS_MAIL_WEBHOOK_SECRET; do
  val="$(grep -E "^${key}=" .env | cut -d= -f2- | tr -d '\r' || true)"
  if [[ -z "${val}" ]]; then
    echo "ERRO: ${key} vazio no .env" >&2
    exit 1
  fi
done

echo "→ Build e arranque (${COMPOSE})..."
${COMPOSE} up -d --build

echo "→ Aguardar healthcheck..."
for i in $(seq 1 30); do
  if curl -sf "http://127.0.0.1:${BACKEND_PORT:-8080}/healthz" >/dev/null 2>&1; then
    echo "OK: backend healthy em http://127.0.0.1:${BACKEND_PORT:-8080}/healthz"
    ${COMPOSE} ps
    exit 0
  fi
  sleep 2
done

echo "ERRO: backend não respondeu a /healthz em 60s" >&2
${COMPOSE} logs --tail=40 backend
exit 1
