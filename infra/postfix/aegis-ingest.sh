#!/usr/bin/env bash
# MAIL-002 — Pipe Postfix → POST /api/mail/ingest
# Instalar em /usr/local/bin/aegis-ingest.sh (ver mail-002-003-production.md)
set -euo pipefail

AEGIS_INGEST_URL="${AEGIS_INGEST_URL:-http://127.0.0.1:8080/api/mail/ingest}"
AEGIS_MAIL_WEBHOOK_SECRET="${AEGIS_MAIL_WEBHOOK_SECRET:-}"

if [[ -z "${AEGIS_MAIL_WEBHOOK_SECRET}" ]]; then
  echo "aegis-ingest: AEGIS_MAIL_WEBHOOK_SECRET não definido" >&2
  exit 75
fi

RCPT="${1:-}"
RAW=$(cat)
export RCPT RAW
MSG_ID=$(printf '%s' "${RAW}" | sha256sum | awk '{print $1}')
FROM=$(printf '%s' "${RAW}" | grep -m1 -i '^From:' | sed 's/^[Ff]rom:[[:space:]]*//;s/\r$//')
SUBJECT=$(printf '%s' "${RAW}" | grep -m1 -i '^Subject:' | sed 's/^[Ss]ubject:[[:space:]]*//;s/\r$//')
BODY=$(printf '%s' "${RAW}" | awk 'BEGIN{body=0} /^$/ {body=1; next} body {print}')
export MSG_ID FROM SUBJECT BODY

python3 -c '
import json, os
print(json.dumps({
    "message_id": os.environ.get("MSG_ID", ""),
    "from": os.environ.get("FROM", ""),
    "to": [t for t in [os.environ.get("RCPT", "")] if t],
    "subject": os.environ.get("SUBJECT", ""),
    "body": os.environ.get("BODY") or os.environ.get("RAW", ""),
}))
' | curl -sf -X POST \
  "${AEGIS_INGEST_URL}?secret=${AEGIS_MAIL_WEBHOOK_SECRET}" \
  -H 'Content-Type: application/json' \
  -d @- >/dev/null

exit 0
